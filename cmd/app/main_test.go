package main

import (
	"database/sql"
	"fmt"
	"os"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/task-manager/app"
	"github.com/task-manager/cache"
	"github.com/task-manager/config"
	"github.com/task-manager/db"
)

func TestMain(m *testing.M) {

	cfg, err := config.LoadTestConfig()

	if err != nil {
		logrus.Fatalf("couldn't load configuration: %v", err)
	}
	postgresDB, err := db.InitDB(*cfg)
	if err != nil {
		logrus.Fatalf("couldn't initialize db: %v", err)
	}
	redis, err := cache.ConnectToRedis(*cfg)
	if err != nil {
		logrus.Fatalf("couldn't initialize db: %v", err)
	}

	testApp := app.BuildApp(cfg, postgresDB, redis)
	// Setup test database
	if err := setupTestDB(testApp); err != nil {
		logrus.Fatalf("Error setting up test database: %v", err)
	}
	defer teardownTestDB(testApp.PostgresDB().Conn)

	// Run tests
	exitVal := m.Run()

	// Teardown (if necessary)
	os.Exit(exitVal)
}

func setupTestDB(testApp *app.App) error {

	if err := resetTestDB(testApp.PostgresDB().Conn); err != nil {
		return err
	}
	if err := seedTestData(testApp.PostgresDB().Conn); err != nil {
		return err
	}

	return nil
}

func resetTestDB(db *sql.DB) error {
	_, err := db.Exec(`
		DROP SCHEMA IF EXISTS public CASCADE;
		CREATE SCHEMA public;
	`)
	if err != nil {
		return fmt.Errorf("error dropping schema: %v", err)
	}
	if err := createTables(db); err != nil {
		return fmt.Errorf("error creating tables: %v", err)
	}

	return nil
}

func createTables(db *sql.DB) error {

	// Read SQL statements from file
	sqlFile, err := os.ReadFile("../db/db.sql")
	if err != nil {
		return fmt.Errorf("error reading SQL file: %v", err)
	}

	// Create necessary tables
	_, err = db.Exec(string(sqlFile))
	if err != nil {
		return fmt.Errorf("error creating tables: %v", err)
	}

	return nil
}

func seedTestData(db *sql.DB) error {
	// Insert initial test data
	_, err := db.Exec(`
		INSERT INTO task (title, description) VALUES
			('Task 1', 'Description for Task 1'),
			('Task 2', 'Description for Task 2');
	`)
	if err != nil {
		return fmt.Errorf("error seeding test data: %v", err)
	}
	return nil
}

func teardownTestDB(db *sql.DB) {
	if db != nil {
		db.Close()
	}
}
