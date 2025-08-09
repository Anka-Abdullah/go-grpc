package database

import (
	"fmt"
	"go-grpc-crud/app/model"
	"go-grpc-crud/config"
	"log"
	"net/url"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func Connect() {
	config := config.LoadConfig()

	// Step 1: Buat base URL koneksi
	baseURL := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s",
		config.DBUser,
		config.DBPassword,
		config.DBHost,
		config.DBPort,
		config.DBName,
	)

	// Step 2: Parse URL
	conn, err := url.Parse(baseURL)
	if err != nil {
		log.Fatalf("[DATABASE]::URI_PARSE_ERROR: %v", err)
	}

	// Step 3: Tambahkan query parameter SSL
	query := conn.Query()
	query.Set("sslmode", config.DBSSLMode)

	// Tambahkan sslrootcert hanya jika disediakan
	if config.DBSSLRootCert != "" {
		query.Set("sslrootcert", config.DBSSLRootCert)
	}

	conn.RawQuery = query.Encode()

	// Step 4: Buka koneksi GORM
	db, err := gorm.Open(postgres.Open(conn.String()), &gorm.Config{
		NowFunc: func() time.Time { return time.Now().Local() },
		Logger:  logger.Default.LogMode(logger.Error),
	})

	if err != nil {
		log.Fatalf("[DATABASE]::CONNECTION_ERROR: %v", err)
	}

	// Step 5: Migrasi model Book
	err = db.AutoMigrate(&model.Book{})
	if err != nil {
		log.Fatalf("[DATABASE]::MIGRATION_ERROR: %v", err)
	}

	DB = db
	fmt.Println("[DATABASE]::CONNECTED & MIGRATED")
}
