package db

import (
	"errors"
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"os"
)

type GormErr struct {
	Number  int    `json:"Number"`
	Message string `json:"Message"`
}

const (
	NotFound                     = "NotFound"
	notFoundMessage              = "record not found"
	ValidationError              = "ValidationError"
	validationErrorMessage       = "validation error"
	ResourceAlreadyExists        = "ResourceAlreadyExists"
	alreadyExistsErrorMessage    = "resource already exists"
	RepositoryError              = "RepositoryError"
	repositoryErrorMessage       = "error in repository operation"
	NotAuthenticated             = "NotAuthenticated"
	notAuthenticatedErrorMessage = "not Authenticated"
	TokenGeneratorError          = "TokenGeneratorError"
	tokenGeneratorErrorMessage   = "error in token generation"
	NotAuthorized                = "NotAuthorized"
	notAuthorizedErrorMessage    = "not authorized"
	UnknownError                 = "UnknownError"
	unknownErrorMessage          = "something went wrong"
)

type AppError struct {
	Err  error
	Type string
}

func NewAppError(err error, errType string) *AppError {
	return &AppError{
		Err:  err,
		Type: errType,
	}
}

func NewAppErrorWithType(errType string) *AppError {
	var err error

	switch errType {
	case NotFound:
		err = errors.New(notFoundMessage)
	case ValidationError:
		err = errors.New(validationErrorMessage)
	case ResourceAlreadyExists:
		err = errors.New(alreadyExistsErrorMessage)
	case RepositoryError:
		err = errors.New(repositoryErrorMessage)
	case NotAuthenticated:
		err = errors.New(notAuthenticatedErrorMessage)
	case NotAuthorized:
		err = errors.New(notAuthorizedErrorMessage)
	case TokenGeneratorError:
		err = errors.New(tokenGeneratorErrorMessage)
	default:
		err = errors.New(unknownErrorMessage)
	}

	return &AppError{
		Err:  err,
		Type: errType,
	}
}

func (appErr *AppError) Error() string {
	return appErr.Err.Error()
}

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
		return "", fmt.Errorf("the environment variable %s is not defined", key)
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
