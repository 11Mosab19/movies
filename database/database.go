package database

import (
	"database/sql"
	"errors"
	"fmt"
	"net"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/joho/godotenv"

	_ "github.com/jackc/pgx/v5/stdlib"
)

type Database struct {
	DB *sql.DB
}

func ReadEnv() (string, error) {
	if err := godotenv.Load(); err != nil && !errors.Is(err, os.ErrNotExist) {
		return "", err
	}

	host := strings.TrimSpace(os.Getenv("DB_HOST"))
	port := strings.TrimSpace(os.Getenv("DB_PORT"))
	user := strings.TrimSpace(os.Getenv("DB_USER"))
	password := os.Getenv("DB_PASSWORD")
	name := strings.TrimSpace(os.Getenv("DB_NAME"))
	if host == "" || port == "" || user == "" || password == "" || name == "" {
		return "", fmt.Errorf("database configuration is incomplete")
	}

	sslMode := strings.TrimSpace(os.Getenv("DB_SSLMODE"))
	if sslMode == "" {
		sslMode = "require"
	}

	connectionURL := &url.URL{
		Scheme: "postgres",
		User:   url.UserPassword(user, password),
		Host:   net.JoinHostPort(host, port),
		Path:   name,
	}
	query := connectionURL.Query()
	query.Set("sslmode", sslMode)
	connectionURL.RawQuery = query.Encode()

	return connectionURL.String(), nil
}

func (DB *Database) ConnectDatabase() error {
	connectString, err := ReadEnv()
	if err != nil {
		return err
	}
	connection, err := sql.Open("pgx", connectString)
	if err != nil {
		return err
	}
	connection.SetMaxOpenConns(25)
	connection.SetMaxIdleConns(25)
	connection.SetConnMaxLifetime(5 * time.Minute)

	DB.DB = connection
	if err := DB.DB.Ping(); err != nil {
		_ = connection.Close()
		DB.DB = nil
		return err
	}
	return nil
}

func (DB *Database) CloseConnection() error {
	if DB.DB == nil {
		return nil
	}
	return DB.DB.Close()
}
