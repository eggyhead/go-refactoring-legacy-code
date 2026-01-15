package processor_fixed

import (
	"encoding/json"
	"fmt"
	"io"
)

// --- 1. Interface Substitution & Parameterization ---
// We replace the OS File and HTTP Client with interfaces.
type ReaderProvider interface {
	Open(path string) (io.ReadCloser, error)
}

type PaymentClient interface {
	GetStatus(userID string) (int, error)
}

type UserProcessor struct {
	configPath string
	reader     ReaderProvider
	payments   PaymentClient
	// --- 2. Function Fields ---
	// We wrap global side-effects (ProcessCount and os.Exit)
	OnProcess func()    // Replaces global ProcessCount++
	OnExit    func(int) // Replaces os.Exit
}

func NewUserProcessor(rp ReaderProvider, pc PaymentClient, path string) *UserProcessor {
	return &UserProcessor{
		configPath: path,
		reader:     rp,
		payments:   pc,
		OnProcess:  func() {}, // Default: do nothing
		OnExit:     func(int) {},
	}
}

// --- 3. Primitivize Parameter ---
// We pull the "Banned Check" out of the complex I/O method.
// This is now a "Pure Function" we can test with 100 cases.
func IsBannedAction(mode string, userID string) bool {
	return mode == "strict" && userID == "banned-user-123"
}

func (up *UserProcessor) ProcessUser(userID string) error {
	// 1. Handle Global State via Function Field
	up.OnProcess()

	// 2. Handle File System via Interface
	configFile, err := up.reader.Open(up.configPath)
	if err != nil {
		return err
	}
	defer configFile.Close()

	// 3. Simple JSON logic
	var config map[string]string
	if err := json.NewDecoder(configFile).Decode(&config); err != nil {
		return err
	}

	// 4. Handle Network via Interface
	statusCode, err := up.payments.GetStatus(userID)
	if err != nil || statusCode != 200 {
		return fmt.Errorf("payment check failed for %s", userID)
	}

	// 5. Use Primitivized Logic & Function Field for Exit
	if IsBannedAction(config["mode"], userID) {
		up.OnExit(1)
		return fmt.Errorf("banned user encountered")
	}

	return nil
}
