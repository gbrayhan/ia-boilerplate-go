package db

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"os"
	"strings"
)

type Config struct {
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	SSLMode    string
}

func LoadDBConfig() (Config, error) {
	var dbHost, errHost = getEnv("DB_HOST")
	var dbPort, errPort = getEnv("DB_PORT")
	var dbUser, errUser = getEnv("DB_USER")
	var dbPassword, errPassword = getEnv("DB_PASSWORD")
	var dbName, errName = getEnv("DB_NAME")
	var sslMode, errSSL = getEnv("DB_SSLMODE")

	var errors = []error{errHost, errPort, errUser, errPassword, errName, errSSL}
	errors = removeNilErrors(errors)
	if len(errors) > 0 {
		var message = errorsMessagesCollector(errors...)
		return Config{}, message
	}

	return Config{
		DBHost:     dbHost,
		DBPort:     dbPort,
		DBUser:     dbUser,
		DBPassword: dbPassword,
		DBName:     dbName,
		SSLMode:    sslMode,
	}, nil
}

func getEnv(key string) (string, error) {
	value, exists := os.LookupEnv(key)
	if !exists {
		return "", fmt.Errorf("La variable de entorno %s no está	definida", key)
	}
	return value, nil
}

func (c Config) GetDSN() string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s", c.DBHost, c.DBPort, c.DBUser, c.DBPassword, c.DBName, c.SSLMode)
}

var DB *gorm.DB

func InitDatabase() error {
	cfg, err := LoadDBConfig()
	if err != nil {
		return err
	}

	DB, err = gorm.Open(postgres.Open(cfg.GetDSN()), &gorm.Config{})
	if err != nil {
		fmt.Println("Error connecting to the database:", err)
		return err
	}

	err = MigrateEntitiesGORM()
	if err != nil {
		fmt.Println("Error migrating the database:", err)
		return err
	}

	return nil
}

func errorsMessagesCollector(errors ...error) error {
	var result string
	for _, err := range errors {
		if err != nil {
			result += err.Error() + "\n"
		}
	}
	if result != "" {
		return fmt.Errorf(result)
	}
	return nil
}

func removeNilErrors(errors []error) []error {
	var result []error
	for _, err := range errors {
		if err != nil {
			result = append(result, err)
		}
	}
	return result
}

func joinWithOr(conditions []string) string {
	return strings.Join(conditions, " OR ")
}
