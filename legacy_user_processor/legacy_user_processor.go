package legacy_user_processor

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

// This variable is a GLOBAL.
// It keeps track of how many users we've processed in the current run.
var ProcessCount = 0

type LegacyUserProcessor struct {
	ConfigPath string
}

func (lup *LegacyUserProcessor) ProcessUser(userID string) error {
	// 1. Dependency: Global variable access
	ProcessCount++

	// 2. Dependency: Hard-coded File System access
	// This will fail if the file doesn't exist on the machine running the test.
	configFile, err := os.Open(lup.ConfigPath)
	if err != nil {
		return err
	}
	defer configFile.Close()

	// 3. Dependency: Hard-coded JSON decoding logic
	// Imagine this is a very complex calculation based on the config.
	var config map[string]string
	json.NewDecoder(configFile).Decode(&config)

	// 4. Dependency: Global Network Call (The "Toxic" call)
	// This hits a real production API. It's slow and requires an API key.
	resp, err := http.Get("https://api.payments.com/v1/status/" + userID)
	if err != nil || resp.StatusCode != 200 {
		return fmt.Errorf("payment check failed for %s", userID)
	}

	// 5. Dependency: Hard-coded OS signal
	// If the user status is "banned", we kill the whole program!
	if config["mode"] == "strict" && userID == "banned-user-123" {
		os.Exit(1)
	}

	fmt.Printf("User %s processed successfully\n", userID)
	return nil
}
