package repository

import (
	"crypto/rand"
	"errors"
	"fmt"
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
	"ia-boilerplate/src/infrastructure"
	"os"
)

type Repository struct {
	DB     *gorm.DB
	Logger *infrastructure.Logger
	Auth   *infrastructure.Auth
}

func NewRepository(db *gorm.DB, logger *infrastructure.Logger, auth *infrastructure.Auth) *Repository {
	return &Repository{
		DB:     db,
		Logger: logger,
		Auth:   auth,
	}
}

func (r *Repository) SetLogger(logger *infrastructure.Logger) {
	r.Logger = logger
}

type Config struct {
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	SSLMode    string
}

func (r *Repository) LoadDBConfig() (Config, error) {
	dbHost, errHost := getEnv("DB_HOST")
	dbPort, errPort := getEnv("DB_PORT")
	dbUser, errUser := getEnv("DB_USER")
	dbPassword, errPassword := getEnv("DB_PASSWORD")
	dbName, errName := getEnv("DB_NAME")
	sslMode, errSSL := getEnv("DB_SSLMODE")

	initErrs := []error{errHost, errPort, errUser, errPassword, errName, errSSL}
	initErrs = r.removeNilErrors(initErrs)
	if len(initErrs) > 0 {
		cfgErr := r.errorsMessagesCollector(initErrs...)
		r.Logger.Error("Error loading database configuration", zap.Error(cfgErr))
		return Config{}, cfgErr
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
	if value, exists := os.LookupEnv(key); exists {
		return value, nil
	}
	return "", errors.New("the environment variable " + key + " is not defined")
}

func (c Config) GetDSN() string {
	return "host=" + c.DBHost +
		" port=" + c.DBPort +
		" user=" + c.DBUser +
		" password=" + c.DBPassword +
		" dbname=" + c.DBName +
		" sslmode=" + c.SSLMode
}

func (r *Repository) InitDatabase() error {
	cfg, err := r.LoadDBConfig()
	if err != nil {
		return err
	}

	gormZap := infrastructure.NewGormLogger(r.Logger.Log).
		LogMode(gormlogger.Warn)

	r.DB, err = gorm.Open(postgres.Open(cfg.GetDSN()), &gorm.Config{
		Logger: gormZap,
	})
	if err != nil {
		r.Logger.Error("Error connecting to the database", zap.Error(err))
		return err
	}

	if err := r.MigrateEntitiesGORM(); err != nil {
		r.Logger.Error("Error migrating the database", zap.Error(err))
		return err
	}

	r.Logger.Info("Database connection and migrations successful")
	return nil
}

func (r *Repository) errorsMessagesCollector(errs ...error) error {
	var result string
	for _, err := range errs {
		result += err.Error() + "; "
	}
	if result != "" {
		return errors.New(result)
	}
	return nil
}

func (r *Repository) removeNilErrors(errs []error) []error {
	var result []error
	for _, err := range errs {
		if err != nil {
			result = append(result, err)
		}
	}
	return result
}

type ErrorType string

const (
	NotFound              ErrorType = "NotFound"
	ValidationError       ErrorType = "ValidationError"
	ResourceAlreadyExists ErrorType = "ResourceAlreadyExists"
	RepositoryError       ErrorType = "RepositoryError"
	NotAuthenticated      ErrorType = "NotAuthenticated"
	TokenGeneratorError   ErrorType = "TokenGeneratorError"
	NotAuthorized         ErrorType = "NotAuthorized"
	UnknownError          ErrorType = "UnknownError"
)

type ErrorTypeMessage string

const (
	notFoundMessage              ErrorTypeMessage = "record not found"
	validationErrorMessage       ErrorTypeMessage = "validation error"
	alreadyExistsErrorMessage    ErrorTypeMessage = "resource already exists"
	repositoryErrorMessage       ErrorTypeMessage = "error in repository operation"
	notAuthenticatedErrorMessage ErrorTypeMessage = "not authenticated"
	tokenGeneratorErrorMessage   ErrorTypeMessage = "token generator error"
	notAuthorizedErrorMessage    ErrorTypeMessage = "not authorized on this action or resource"
	unknownErrorMessage          ErrorTypeMessage = "unknown error, we are working to improve this experience for you"
)

type AppError struct {
	Err  error
	Type ErrorType
}

func NewAppError(err error, errType ErrorType) *AppError {
	return &AppError{
		Err:  err,
		Type: errType,
	}
}

func NewAppErrorWithType(errType ErrorType) *AppError {
	var err error
	switch errType {
	case NotFound:
		err = errors.New(string(notFoundMessage))
	case ValidationError:
		err = errors.New(string(validationErrorMessage))
	case ResourceAlreadyExists:
		err = errors.New(string(alreadyExistsErrorMessage))
	case RepositoryError:
		err = errors.New(string(repositoryErrorMessage))
	case NotAuthenticated:
		err = errors.New(string(notAuthenticatedErrorMessage))
	case NotAuthorized:
		err = errors.New(string(notAuthorizedErrorMessage))
	case TokenGeneratorError:
		err = errors.New(string(tokenGeneratorErrorMessage))
	case UnknownError:
		err = errors.New(string(unknownErrorMessage))
	default:
		err = errors.New(string(unknownErrorMessage))
	}
	return &AppError{
		Err:  err,
		Type: errType,
	}
}

func (appErr *AppError) Error() string {
	return appErr.Err.Error()
}

func GenerateNewUUID() (string, error) {
	uuid := make([]byte, 16)
	if _, err := rand.Read(uuid); err != nil {
		return "", err
	}
	uuid[6] = (uuid[6] & 0x0f) | 0x40
	uuid[8] = (uuid[8] & 0x3f) | 0x80

	return fmt.Sprintf("%x-%x-%x-%x-%x", uuid[0:4], uuid[4:6], uuid[6:8], uuid[8:10], uuid[10:]), nil
}
