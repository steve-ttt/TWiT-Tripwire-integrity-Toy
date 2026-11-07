package tripwire

import (
	"encoding/json"
	"log"
)

type Config struct {
	Database     string   `json:"database"`
	ScanDirs     []string `json:"scan_dirs"`
	ScanFiles    []string `json:"scan_files"`
	ScanInterval int      `json:"scan_interval"`
}

// os.ReadFile("config.json") returns ([]byte, error)
func readConfig(data []byte) (Config, error) {
	var config Config
	err := json.Unmarshal(data, &config)
	if err != nil {
		log.Fatal(err)
		panic(err)
	}
	return config, nil
}
