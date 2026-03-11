package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/Shubhouy1/asset-management/database"
	"github.com/Shubhouy1/asset-management/router"
	"github.com/Shubhouy1/asset-management/utils"
)

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

func main() {
	err := utils.InitFirebase()
	if err != nil {
		log.Fatal("firebase init failed:", err)
	}
	r := router.SetupRouter()

	dbHost := getEnv("DB_HOST", "localhost")
	dbPort := getEnv("DB_PORT", "5434")
	dbUser := getEnv("DB_USER", "local")
	dbPassword := getEnv("DB_PASSWORD", "local")
	dbName := getEnv("DB_NAME", "storeX")
	sslMode := getEnv("DB_SSLMODE", string(database.SSLModeDisabled))
	serverPort := getEnv("SERVER_PORT", "8080")

	err = database.CreateAndMigrate(
		dbHost,
		dbPort,
		dbUser,
		dbPassword,
		dbName,
		database.SSLMode(sslMode),
	)
	if err != nil {
		panic(err)
	}

	fmt.Println("Server running on port", serverPort)

	if err := http.ListenAndServe(":"+serverPort, r); err != nil {
		panic(err)
	}
}
