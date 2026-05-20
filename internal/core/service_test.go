package core

import (
	"context"
	"errors"
	"testing"
)

type mockStorage struct {
	message *Message
	getErr  error
}

func (m *mockStorage) GetAllMessages(ctx context.Context, input GetAllMessagesInput) ([]Message, error) {
	return nil, nil
}

func (m *mockStorage) SaveMessage(ctx context.Context, message Message) error {
	return nil
}

func (m *mockStorage) UpdateMessage(ctx context.Context, message Message) error {
	return nil
}

func (m *mockStorage) GetMessage(ctx context.Context, id string) (*Message, error) {
	if m.getErr != nil {
		return nil, m.getErr
	}
	return m.message, nil
}

type mockScheduler struct {
	retryCalled bool
	retryErr    error
}

func (m *mockScheduler) ScheduleMessage(ctx context.Context, message Message) error {
	return nil
}

func (m *mockScheduler) RetryMessage(ctx context.Context, msg Message) error {
	m.retryCalled = true
	return m.retryErr
}

func newRetryTestService(t *testing.T, status MessageStatus) (Service, *mockScheduler) {
	t.Helper()

	scheduler := &mockScheduler{}
	svc, err := NewService(ServiceConfig{
		Storage: &mockStorage{
			message: &Message{ID: "msg-1", Status: status},
		},
		Scheduler: scheduler,
	})
	if err != nil {
		t.Fatalf("NewService() error = %v", err)
	}
	return svc, scheduler
}

func callRetryMessage(t *testing.T, svc Service) error {
	t.Helper()
	return svc.RetryMessage(context.Background(), RetryMessageInput{
		ID:                 "msg-1",
		ScheduledSendingAt: 12345,
	})
}

func assertRetrySucceeded(t *testing.T, scheduler *mockScheduler, err error) {
	t.Helper()
	if err != nil {
		t.Fatalf("RetryMessage() error = %v, want nil", err)
	}
	if !scheduler.retryCalled {
		t.Fatal("expected scheduler.RetryMessage to be called")
	}
}

func assertRetryRejected(t *testing.T, scheduler *mockScheduler, err error, wantErrIs error, wantErrText string) {
	t.Helper()
	if scheduler.retryCalled {
		t.Fatal("scheduler.RetryMessage should not be called for non-failed messages")
	}
	if !errors.Is(err, wantErrIs) {
		t.Fatalf("RetryMessage() errors.Is = %v, want %v", err, wantErrIs)
	}
	if err.Error() != wantErrText {
		t.Fatalf("RetryMessage() error = %q, want %q", err.Error(), wantErrText)
	}
}

func TestRetryMessage_FailedStatusCanBeRetried(t *testing.T) {
	t.Parallel()

	svc, scheduler := newRetryTestService(t, MessageStatusFailed)
	err := callRetryMessage(t, svc)
	assertRetrySucceeded(t, scheduler, err)
}

func TestRetryMessage_SentStatusCannotBeRetried(t *testing.T) {
	t.Parallel()

	svc, scheduler := newRetryTestService(t, MessageStatusSent)
	err := callRetryMessage(t, svc)
	assertRetryRejected(t, scheduler, err, ErrCannotRetryActiveMessage,
		"Cannot retry message: current status is 'SENT'. Only 'FAILED' messages can be retried.")
}

func TestRetryMessage_ScheduledStatusCannotBeRetried(t *testing.T) {
	t.Parallel()

	svc, scheduler := newRetryTestService(t, MessageStatusScheduled)
	err := callRetryMessage(t, svc)
	assertRetryRejected(t, scheduler, err, ErrCannotRetryActiveMessage,
		"Cannot retry message: current status is 'SCHEDULED'. Only 'FAILED' messages can be retried.")
}

func TestRetryMessage_MessageNotFound(t *testing.T) {
	t.Parallel()

	scheduler := &mockScheduler{}
	svc, err := NewService(ServiceConfig{
		Storage:   &mockStorage{message: nil},
		Scheduler: scheduler,
	})
	if err != nil {
		t.Fatalf("NewService() error = %v", err)
	}

	err = svc.RetryMessage(context.Background(), RetryMessageInput{ID: "missing"})
	if !errors.Is(err, ErrMessageNotFound) {
		t.Fatalf("RetryMessage() errors.Is = %v, want ErrMessageNotFound", err)
	}
	if scheduler.retryCalled {
		t.Fatal("scheduler.RetryMessage should not be called when message is missing")
	}
}
