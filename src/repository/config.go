package repository

import (
	"errors"
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"ia-boilerplate/src/infrastructure"
	"os"
)

type Repository struct {
	DB             *gorm.DB
	Logger         *zap.Logger
	Infrastructure *infrastructure.Infrastructure
}

func NewRepository(db *gorm.DB, logger *zap.Logger) *Repository {
	return &Repository{
		DB:     db,
		Logger: logger,
	}
}

func (r *Repository) SetLogger(logger *zap.Logger) {
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

	r.DB, err = gorm.Open(postgres.Open(cfg.GetDSN()), &gorm.Config{})
	if err != nil {
		r.Logger.Error("Error connecting to the database", zap.Error(err))
		return err
	}

	err = r.MigrateEntitiesGORM()
	if err != nil {
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

// Error types and AppError struct remain unchanged
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
	notAuthenticatedErrorMessage = "not authenticated"
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
