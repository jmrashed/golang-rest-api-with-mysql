package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

type DB struct {
	*sql.DB
}

// Config holds database configuration
type Config struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
}

// NewConnection creates a new database connection
func NewConnection(config Config) (*DB, error) {
	// First connect without database to create it if needed
	var dsn string
	if config.Password == "" {
		dsn = fmt.Sprintf("%s@tcp(%s:%s)/?parseTime=true&multiStatements=true",
			config.User, config.Host, config.Port)
	} else {
		dsn = fmt.Sprintf("%s:%s@tcp(%s:%s)/?parseTime=true&multiStatements=true",
			config.User, config.Password, config.Host, config.Port)
	}

	log.Printf("Connecting to MySQL with DSN: %s", dsn)
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	// Create database if it doesn't exist
	_, err = db.Exec(fmt.Sprintf("CREATE DATABASE IF NOT EXISTS %s", config.DBName))
	if err != nil {
		return nil, fmt.Errorf("failed to create database: %w", err)
	}
	db.Close()

	// Now connect to the specific database
	if config.Password == "" {
		dsn = fmt.Sprintf("%s@tcp(%s:%s)/%s?parseTime=true&multiStatements=true",
			config.User, config.Host, config.Port, config.DBName)
	} else {
		dsn = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true&multiStatements=true",
			config.User, config.Password, config.Host, config.Port, config.DBName)
	}

	db, err = sql.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	// Set connection pool settings
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)

	return &DB{db}, nil
}

// InitializeSchema creates tables and inserts default data
func (db *DB) InitializeSchema() error {
	schemaFile := "schema/schema.sql"
	
	// Read schema file
	schema, err := os.ReadFile(schemaFile)
	if err != nil {
		return fmt.Errorf("failed to read schema file: %w", err)
	}

	// Execute schema
	if _, err := db.Exec(string(schema)); err != nil {
		return fmt.Errorf("failed to execute schema: %w", err)
	}

	log.Println("Database schema initialized successfully")
	return nil
}

// GetDefaultConfig returns default database configuration
func GetDefaultConfig() Config {
	return Config{
		Host:     getEnv("DB_HOST", "127.0.0.1"),
		Port:     getEnv("DB_PORT", "3306"),
		User:     getEnv("DB_USER", "root"),
		Password: getEnv("DB_PASSWORD", "password"),
		DBName:   getEnv("DB_NAME", "goblog"),
	}
}

func getEnv(key, defaultValue string) string {
	value, exists := os.LookupEnv(key)
	if exists {
		return value
	}
	return defaultValue
}