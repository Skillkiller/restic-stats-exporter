package main

import (
	"log"
	"os"
)

func main() {
	log.Println("Starting restic statistics exporter...")
	checkEnv("RESTIC_REPOSITORY")
}

func checkEnv(name string) {
	val, ok := os.LookupEnv(name)
	if !ok {
		log.Fatalf("%s environment variable is not set", name)
	}

	if val == "" {
		log.Fatalf("%s environment variable is empty", name)
	}

}
