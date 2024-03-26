package database

import (
	"fmt"
	"strconv"

	"github.com/mnsh5/fiber-jwt-auth/src/config"
	"github.com/mnsh5/fiber-jwt-auth/src/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func ConnectDB() {
	var err error
	p := config.Config("DB_PORT")
	port, err := strconv.ParseUint(p, 10, 32)
	if err != nil {
		panic("Failed to parse database port")
	}

	dsn := fmt.Sprintf("host=%s port=%d dbname=%s user=%s password=%s sslmode=disable",
		config.Config("DB_HOST"),
		port,
		config.Config("DB_NAME"),
		config.Config("DB_USER"),
		config.Config("DB_PASSWORD"))

	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(fmt.Errorf("failed to connect to database: %w", err))
	}
	fmt.Println("Connection Open to database")

	DB.AutoMigrate(&models.User{})
	fmt.Println("Database migrated successfully")
}
