package main

import (
	"flag"
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
		log.Fatal("failed to scan directories and files: " + err.Error())
	}
	for path, hash := range scanResult {
		if currentHash, ok := currentDataBase[path]; !ok {
			integretyFailed = true
			fmt.Printf("New Unknown file: %s\n", path)
		} else if currentHash != hash {
			integretyFailed = true
			fmt.Printf("File changed: %s\n", path)
		}
	}
	// check if there are any files in the database that are not in the current scan result
	for path := range currentDataBase {
		if _, ok := scanResult[path]; !ok {
			integretyFailed = true
			fmt.Printf("Deleted file: %s\n", path)
		}
	}
	return integretyFailed

}

func updateBaseline(config Config) {
	scanResult, err := scanDirsAndFiles(config)
	if err != nil {
		log.Fatal("failed to scan directories and files: " + err.Error())
	}
	// check if there are any files in the database that are not in the current scan result
	// Only compare against previous baseline if the DB file already exists
	if currentDataBase, err := loadBaseline(config.Database); err == nil {
		for path := range currentDataBase {
			if _, ok := scanResult[path]; !ok {
				fmt.Printf("Deleted file: %s\n", path)
			}
		}
	}

	if err = saveBaseline(config.Database, scanResult); err != nil {
		log.Fatalf("failed to save baseline: %v", err)
	}

}

func main() {
	check := flag.Bool("check", false, "Check file integrity against the baseline")
	update := flag.Bool("update", false, "Initialize or update the baseline database")
	configFile := flag.String("config", "/root/tripwire.json", "Path to the configuration file")
	flag.Parse()

	if *check == *update {
		fmt.Fprintln(os.Stderr, "Usage: tripwire --check|--update [--config <file>]")
		flag.PrintDefaults()
		os.Exit(1)
	}

	data, err := os.ReadFile(*configFile)
	if err != nil {
		log.Fatal(err)
	}

	config, err := readConfig(data)
	if err != nil {
		log.Fatal(err)
	}

	if *check {
		integretyFailed := checkBaseline(config)
		fmt.Println("==================================")
		if integretyFailed {
			fmt.Println(RED + "Integrity check failed." + RESET)
			fmt.Println("==================================")
			os.Exit(1)
		}
		fmt.Println(GREEN + "Integrity check passed." + RESET)
		fmt.Println("==================================")
	} else if *update {
		updateBaseline(config)
	}
}
