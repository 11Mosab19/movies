package database

import (
	"database/sql"
	"os"

	"github.com/joho/godotenv"
)

type Database struct {
	DB *sql.DB
}

func ReadEnv() (string, error) {
	err := godotenv.Load()
	if err != nil {
		return "", err
	}
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	name := os.Getenv("DB_NAME")
	return "postgres://" + user + ":" + password + "@" + host + ":" + port + "/" + name + "?sslmode=disable", nil
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
	DB.DB = connection
	if err := DB.DB.Ping(); err != nil {
		return err
	}
	return nil
}

func (DB *Database) CloseConnection() error {
	return DB.DB.Close()
}
