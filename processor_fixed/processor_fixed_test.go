package processor_fixed

import (
	"io"
	"strings"
	"testing"
)

// --- Mocks ---
type MockReader struct{ Content string }

func (m MockReader) Open(p string) (io.ReadCloser, error) {
	return io.NopCloser(strings.NewReader(m.Content)), nil
}

type MockPayments struct{ Status int }

func (m MockPayments) GetStatus(id string) (int, error) { return m.Status, nil }

func TestProcessUser_BannedUser(t *testing.T) {
	exitCalled := false
	processedCount := 0

	// 1. Setup the Processor with Mocks
	proc := NewUserProcessor(
		MockReader{Content: `{"mode":"strict"}`},
		MockPayments{Status: 200},
		"config.json",
	)

	// 2. Override Function Fields for Sensing
	proc.OnExit = func(code int) { exitCalled = true }
	proc.OnProcess = func() { processedCount++ }

	// 3. Execute
	proc.ProcessUser("banned-user-123")

	// 4. Assertions
	if !exitCalled {
		t.Error("Expected OnExit to be called for banned user")
	}
	if processedCount != 1 {
		t.Errorf("Expected count 1, got %d", processedCount)
	}
}

func TestIsBannedAction(t *testing.T) {
	// Testing the primitivized pure logic
	if !IsBannedAction("strict", "banned-user-123") {
		t.Error("Should have flagged as banned")
	}
}
