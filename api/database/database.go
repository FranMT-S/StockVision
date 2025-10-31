package database

import (
	"api/config"
	apilogger "api/logger"
	"api/models"
	"api/sanatizer"
	"fmt"
	"log"
	"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

type Database struct {
	DB *gorm.DB
}

var databaseInstance *Database

// NewDatabase creates a new database connection and initializes the schema
func NewDatabase(config config.DatabaseConfig) (*Database, error) {

	// validate config
	if !isValid(config) {
		return nil, fmt.Errorf("invalid database configuration")
	}

	// Connect to the database
	db, err := connectToDatabase(config)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Auto-migrate models
	if err := migrateModels(db); err != nil {
		return nil, fmt.Errorf("failed to migrate models: %w", err)
	}

	databaseInstance = &Database{DB: db}

	return databaseInstance, nil
}

// connectToDatabase establishes a connection to the database
func connectToDatabase(config config.DatabaseConfig) (*gorm.DB, error) {
	dsn := GetDSN(config)
	dbLogger := ConfigLogger()

	apilogger.Logger().Info().Msg("Successfully connected to the database")
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			TablePrefix: config.SchemaName + ".",
		},
		DisableForeignKeyConstraintWhenMigrating: true,
		PrepareStmt:                              true,
		Logger:                                   dbLogger,
	})

	// Create schema if it doesn't exist
	schemaName := sanatizer.SanitizeForSQL(config.SchemaName)
	if err := createSchemaIfNotExists(db, schemaName); err != nil {
		return nil, fmt.Errorf("error creating schema: %w", err)
	}

	// Configure connection pool
	if err = configureConnectionPool(db); err != nil {
		return nil, fmt.Errorf("error configuring connection pool: %w", err)
	}

	return db, nil
}

// configureConnectionPool configures the connection pool
func configureConnectionPool(db *gorm.DB) error {
	// Configure connection pool for CockroachDB
	sqlDB, err := db.DB()
	if err != nil {
		return err
	}

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)
	sqlDB.SetConnMaxIdleTime(10 * time.Minute)

	return nil
}

// createSchemaIfNotExists creates the schema if it doesn't exist
func createSchemaIfNotExists(db *gorm.DB, schema string) error {
	if schema == "public" {
		return nil // public schema always exists
	}

	query := fmt.Sprintf("CREATE SCHEMA IF NOT EXISTS %s", schema)
	return db.Exec(query).Error
}

// migrateModels runs auto-migration for all models
func migrateModels(db *gorm.DB) error {
	if err := db.AutoMigrate(
		&models.Brokerage{},
		&models.Ticker{},
		&models.Recommendation{},
		&models.Onboarding{},
	); err != nil {
		return fmt.Errorf("failed to migrate models: %w", err)
	}

	return nil
}

// Close closes the database connection
func (d *Database) Close() error {
	sqlDB, err := d.DB.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}

// GetDB returns the underlying GORM DB instance
func (d *Database) GetDB() *gorm.DB {
	return d.DB
}

// isValid check if the Host, port, user and db name are not empty
func isValid(cfg config.DatabaseConfig) bool {
	return cfg.Host != "" && cfg.Port != "" && cfg.User != "" && cfg.DBName != ""
}

// GetDSN returns the DSN string for the database connection in gorm
func GetDSN(cfg config.DatabaseConfig) string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s search_path=%s",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.DBName, cfg.SSLMode, cfg.SchemaName)
}

// ConfigLogger returns a new logger instance if config.Log().DB is true
func ConfigLogger() logger.Interface {
	var newLogger logger.Interface = nil
	if config.Log().DB {
		newLogger = logger.New(
			log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
			logger.Config{
				SlowThreshold:             2 * time.Second, // Slow SQL threshold
				LogLevel:                  logger.Info,     // Log level
				IgnoreRecordNotFoundError: true,            // Ignore ErrRecordNotFound error for logger
				ParameterizedQueries:      true,            // Don't include params in the SQL log
				Colorful:                  true,            // Disable color
			},
		)
	}
	return newLogger
}

// GetDB the database instance
// if the instance is not initialized, it will be initialized with the default configuration
func GetDB() (*Database, error) {
	if databaseInstance == nil {
		return NewDatabase(*config.Database())
	}

	return databaseInstance, nil
}
