package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq" // Driver PostgreSQL
)

// DB adalah instance database global
var DB *sql.DB

// ConnectDB menginisialisasi koneksi ke database menggunakan DSN dari .env
func ConnectDB() {
	dsn := os.Getenv("DB_DSN")
	if dsn == "" {
		log.Fatal("DB_DSN not set in .env file")
	}

	var err error
	DB, err = sql.Open("postgres", dsn)
	if err != nil {
		log.Fatal("Gagal koneksi ke database:", err)
	}

	if err = DB.Ping(); err != nil {
		log.Fatal("Gagal ping database:", err)
	}

	fmt.Println("Berhasil terhubung ke database PostgreSQL 'project_uas'")
}