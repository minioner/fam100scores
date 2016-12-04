package main

import (
	"encoding/json"
	"flag"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
)

var (
	delete   = false
	keepDays = 60

	deleted = 0
)

func main() {
	searchDir := flag.String("path", "scores", "path contains json file to clean up")
	flag.IntVar(&keepDays, "days", 60, "days to keep the file")
	flag.BoolVar(&delete, "f", false, "actually delete the file")
	flag.Parse()

	log.Printf("processing %q\n", *searchDir)
	err := filepath.Walk(*searchDir, func(path string, f os.FileInfo, err error) error {
		if f.IsDir() || !strings.HasSuffix(f.Name(), ".json") {
			return nil
		}

		processFile(path)
		return nil
	})

	if err != nil {
		log.Fatal("failed processing directory", err)
	}

	log.Printf("deleted %d files", deleted)

}

func processFile(path string) {
	f, err := os.Open(path)
	if err != nil {
		log.Fatalf("error opening %s: %s", path, err)
	}
	dec := json.NewDecoder(f)

	var s struct {
		LastUpdated time.Time `json:"lastUpdated"`
	}
	if err := dec.Decode(&s); err != nil {
		log.Printf("failed processing %q, %s", path, err)
		return
	}
	f.Close()
	age := time.Now().Sub(s.LastUpdated)

	if age > time.Duration(keepDays)*24*time.Hour {
		log.Printf("deleting %s (age %f)", path, age.Hours()/24)
		if delete {
			if err := os.Remove(path); err != nil {
				log.Printf("failed removing %q, %s", path, err)
			}
		}
		deleted++
	}
}
