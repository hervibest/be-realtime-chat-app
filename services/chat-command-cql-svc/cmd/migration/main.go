package main

import (
	"be-realtime-chat-app/services/chat-command-cql-svc/internal/config"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/gocql/gocql"
)

func runCQLFile(session *gocql.Session, path string) error {
	content, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	queries := strings.Split(string(content), ";")

	for _, query := range queries {
		q := strings.TrimSpace(query)
		if q == "" {
			continue
		}
		if err := session.Query(q).Exec(); err != nil {
			return fmt.Errorf("error executing query %q: %w", q, err)
		}
	}
	return nil
}

func main() {
	session, _ := config.NewCQLDB()
	defer session.Close()

	// Load and run .up.cql files
	files, err := filepath.Glob("../../db/migrations/*.up.cql")
	if err != nil {
		log.Fatalf("Failed to read migration files: %v", err)
	}

	for _, file := range files {
		log.Println("Running migration:", file)
		if err := runCQLFile(session, file); err != nil {
			log.Fatalf("Migration failed: %v", err)
		}
	}
	log.Println("Migration finished successfully")
}
