package driver

import (
	"net/http"
	"testing"

	"github.com/ghazlabs/wa-scheduler/internal/core"
)

func TestNewErrorResp_CannotRetryMapsToConflict(t *testing.T) {
	t.Parallel()

	resp := NewErrorResp(core.NewCannotRetryMessageError(core.MessageStatusSent))

	if resp.StatusCode != http.StatusConflict {
		t.Fatalf("StatusCode = %d, want %d", resp.StatusCode, http.StatusConflict)
	}
	if resp.Err != "ERR_CONFLICT" {
		t.Fatalf("Err = %q, want ERR_CONFLICT", resp.Err)
	}
	wantMsg := "Cannot retry message: current status is 'SENT'. Only 'FAILED' messages can be retried."
	if resp.Message != wantMsg {
		t.Fatalf("Message = %q, want %q", resp.Message, wantMsg)
	}
}

func TestNewErrorResp_MessageNotFoundMapsToNotFound(t *testing.T) {
	t.Parallel()

	resp := NewErrorResp(core.ErrMessageNotFound)

	if resp.StatusCode != http.StatusNotFound {
		t.Fatalf("StatusCode = %d, want %d", resp.StatusCode, http.StatusNotFound)
	}
	if resp.Err != "ERR_MESSAGE_NOT_FOUND" {
		t.Fatalf("Err = %q, want ERR_MESSAGE_NOT_FOUND", resp.Err)
	}
}
