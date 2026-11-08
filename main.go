package tripwire

import (
	"fmt"
	"log"
	"os"
)

func checkBaseline(config Config) {
	// TDD Stub
}

func updateBaseline(config Config, currentDataBase map[string]string) {
	scanResult := make(map[string]string)
	for _, dir := range config.ScanDirs {
		// append the scan result to the scanResult map
		dirScanResult, err := hashTree(dir)
		if err != nil {
			log.Fatalf("failed to scan directory %s: %v", dir, err)
		}
		for path, hash := range dirScanResult {
			scanResult[path] = hash
		}
	}

	for _, file := range config.ScanFiles {
		fileHash, err := hashFile(file)
		scanResult[file] = fileHash
		if err != nil {
			log.Fatalf("failed to scan file %s: %v", file, err)
		}
	}
	// check if there are any unknown files in the database that are not in the current scan result
	// This is a naive check we will delete the entry from the database if it is not in the current scan result.
	// and log the unknown file to the log file and console.
	for path := range currentDataBase {
		if _, ok := scanResult[path]; !ok {
			log.Printf("Unknown file: %s", path)
			fmt.Printf("Unknown file: %s\n", path)
		}
	}

	err := saveBaseline(config.Database, scanResult)
	if err != nil {
		log.Fatalf("failed to save baseline: %v", err)
	}

}

func main() {
	// command line args --check, --update, --config <file>[default: /root/tripwire.json]
	// parse command line args
	if len(os.Args) < 2 {
		log.Fatal("Usage: tripwire --check|--update | --help [--config <file>]")
	}
	var configFile string = "/root/tripwire.json"
	if len(os.Args) >= 4 && os.Args[2] == "--config" {
		configFile = os.Args[3]
	}

	data, err := os.ReadFile(configFile)
	if err != nil {
		log.Fatal(err)
	}

	config, err := readConfig(data)
	if err != nil {
		log.Fatal(err)
	}

	currentDataBase, err := loadBaseline(config.Database)
	if err != nil {
		log.Fatal(err)
	}

	switch os.Args[1] {
	case "--check":
		checkBaseline(config)
	case "--update":
		updateBaseline(config, currentDataBase)
	case "--help":
		fmt.Println("Usage: tripwire --check|--update [--config <file>]")
		fmt.Println("Default config file: /root/tripwire.json")
	default:
		log.Fatal("Unknown command. Usage: tripwire --check|--update [--config <file>]")
	}

}
