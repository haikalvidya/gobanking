package mysql

import (
	"context"
	"fmt"
	"time"

	"github.com/pkg/errors"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Config struct {
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	User     string `yaml:"user"`
	DBName   string `yaml:"dbName"`
	SSLMode  bool   `yaml:"sslMode"`
	Password string `yaml:"password"`
}

const (
	maxConn           = 50
	healthCheckPeriod = 1 * time.Minute
	maxConnIdleTime   = 1 * time.Minute
	maxConnLifetime   = 3 * time.Minute
	minConns          = 1
	lazyConnect       = false
)

// New returns a new MySQL connection pool using gorm
func NewMysqlConn(cfg *Config) (*gorm.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local&tls=%t",
		cfg.User,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.DBName,
		cfg.SSLMode,
	)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, errors.Wrap(err, "failed to connect to database")
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, errors.Wrap(err, "failed to get sql.DB from gorm.DB")
	}

	sqlDB.SetMaxOpenConns(maxConn)
	sqlDB.SetMaxIdleConns(maxConn)
	sqlDB.SetConnMaxIdleTime(maxConnIdleTime)
	sqlDB.SetConnMaxLifetime(maxConnLifetime)
	sqlDB.SetConnMaxIdleTime(maxConnIdleTime)

	ctx, cancel := context.WithTimeout(context.Background(), healthCheckPeriod)
	defer cancel()

	if err := sqlDB.PingContext(ctx); err != nil {
		return nil, errors.Wrap(err, "failed to ping database")
	}

	return db, nil
}
