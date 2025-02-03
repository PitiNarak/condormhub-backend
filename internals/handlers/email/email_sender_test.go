package email

import (
	"testing"
)

func TestSendEmail(t *testing.T) {
	if err := SendEmail("equesdedeus@gmail.com", "test-token"); err != nil {
		t.Errorf("Failed to send email: %v", err)
	}
}
