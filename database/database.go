package database

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/lib/pq" // Driver PostgreSQL
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// DB adalah instance database PostgreSQL global
var DB *sql.DB

// MongoDB adalah instance database MongoDB global
var MongoDB *mongo.Database

// ConnectDB menginisialisasi koneksi ke DUA database (Postgres & Mongo)
func ConnectDB() {
	// --- 1. KONEKSI POSTGRESQL ---
	dsn := os.Getenv("DB_DSN")
	if dsn == "" {
		log.Fatal("DB_DSN not set in .env file")
	}

	var err error
	DB, err = sql.Open("postgres", dsn)
	if err != nil {
		log.Fatal("Gagal koneksi ke PostgreSQL:", err)
	}

	if err = DB.Ping(); err != nil {
		log.Fatal("Gagal ping PostgreSQL:", err)
	}
	fmt.Println("✅ Berhasil terhubung ke database PostgreSQL")

	// --- 2. KONEKSI MONGODB ---
	mongoURI := os.Getenv("MONGO_URI")
	mongoDBName := os.Getenv("MONGO_DB_NAME")

	if mongoURI == "" || mongoDBName == "" {
		log.Fatal("MONGO_URI atau MONGO_DB_NAME belum diset di .env")
	}

	// Buat context dengan timeout 10 detik agar tidak hang jika koneksi lambat
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	clientOptions := options.Client().ApplyURI(mongoURI)
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal("Gagal inisialisasi client MongoDB:", err)
	}

	// Cek koneksi (Ping)
	if err = client.Ping(ctx, nil); err != nil {
		log.Fatal("Gagal ping MongoDB:", err)
	}

	// Set variabel global
	MongoDB = client.Database(mongoDBName)
	fmt.Println("✅ Berhasil terhubung ke database MongoDB:", mongoDBName)
}