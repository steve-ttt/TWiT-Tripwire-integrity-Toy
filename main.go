package main

import (
	"fmt"
	"log"
	"os"
)

const (
	RED   = "\033[31m"
	GREEN = "\033[32m"
	RESET = "\033[0m"
)

func scanDirsAndFiles(config Config) (map[string]string, error) {
	scanResult := make(map[string]string)
	for _, dir := range config.ScanDirs {
		dirScanResult, err := hashTree(dir)
		if err != nil {
			return nil, fmt.Errorf("failed to scan directory %s: %v", dir, err)
		}
		for path, hash := range dirScanResult {
			scanResult[path] = hash
		}
	}

	for _, file := range config.ScanFiles {
		fileHash, err := hashFile(file)
		if err != nil {
			return nil, fmt.Errorf("failed to scan file %s: %v", file, err)
		}
		scanResult[file] = fileHash
	}

	return scanResult, nil
}

func checkBaseline(config Config) bool {
	currentDataBase, err := loadBaseline(config.Database)
	if err != nil {
		log.Fatalf("failed to load baseline (run --update first): %v", err)
	}
	integretyFailed := false
	scanResult, err := scanDirsAndFiles(config)
	if err != nil {
		log.Fatal("failed to scan directories and files")
	}
	for path, hash := range scanResult {
		if currentHash, ok := currentDataBase[path]; !ok {
			integretyFailed = true
			log.Printf("New Unknown file: %s", path)
			fmt.Printf("New Unknown file: %s\n", path)
		} else if currentHash != hash {
			integretyFailed = true
			log.Printf("File changed: %s", path)
			fmt.Printf("File changed: %s\n", path)
		}
	}
	// check if there are any files in the database that are not in the current scan result
	for path := range currentDataBase {
		if _, ok := scanResult[path]; !ok {
			integretyFailed = true
			log.Printf("Deleted file: %s", path)
			fmt.Printf("Deleted file: %s\n", path)
		}
	}
	return integretyFailed

}

func updateBaseline(config Config) {
	scanResult, err := scanDirsAndFiles(config)
	if err != nil {
		log.Fatal("failed to scan directories and files")
	}
	// check if there are any files in the database that are not in the current scan result
	// This is a naive check we will delete the entry from the database if it is not in the current scan result.
	// and log the unknown file to the log file and console.
	// Only compare against previous baseline if the DB file already exists
	if currentDataBase, err := loadBaseline(config.Database); err == nil {
		for path := range currentDataBase {
			if _, ok := scanResult[path]; !ok {
				log.Printf("Deleted file: %s", path)
				fmt.Printf("Deleted file: %s\n", path)
			}
		}
	}

	if err = saveBaseline(config.Database, scanResult); err != nil {
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

	switch os.Args[1] {
	case "--check":
		integretyFailed := checkBaseline(config)
		fmt.Println("==================================")
		if integretyFailed {
			fmt.Println(RED + "Integrity check failed." + RESET)
		} else {
			fmt.Println(GREEN + "Integrity check passed." + RESET)
		}
		fmt.Println("==================================")
	case "--update":
		updateBaseline(config)
	case "--help":
		fmt.Println("Usage: tripwire --check|--update [--config <file>]")
		fmt.Println("Default config file: /root/tripwire.json")
	default:
		log.Fatal("Unknown command. Usage: tripwire --check|--update [--config <file>]")
	}

}
