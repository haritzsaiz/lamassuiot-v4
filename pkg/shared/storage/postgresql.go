package storage

import (
	"context"
	"fmt"
	"time"

	"github.com/lamassuiot/lamassuiot/v4/pkg/shared/config"
	"github.com/lamassuiot/lamassuiot/v4/pkg/shared/logger"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
)

func CreatePostgresDBConnection(logger *logger.Logger, cfg config.PostgresConfig, database string) (*gorm.DB, error) {
	dbLogger := &GormLogger{
		logger: logger,
	}

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable", cfg.Hostname, cfg.Username, cfg.Password, database, cfg.Port)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: dbLogger,
	})

	return db, err
}

type GormLogger struct {
	logger *logger.Logger
}

func (l *GormLogger) LogMode(lvl gormlogger.LogLevel) gormlogger.Interface {
	newlogger := *l
	return &newlogger
}

func (l *GormLogger) Info(ctx context.Context, str string, rest ...interface{}) {
	le := logger.ConfigureLogger(ctx, l.logger)
	le.Info(fmt.Sprintf(str, rest...))
}

func (l *GormLogger) Warn(ctx context.Context, str string, rest ...interface{}) {
	le := logger.ConfigureLogger(ctx, l.logger)
	le.Warn(fmt.Sprintf(str, rest...))
}

func (l *GormLogger) Error(ctx context.Context, str string, rest ...interface{}) {
	le := logger.ConfigureLogger(ctx, l.logger)
	le.Error(fmt.Sprintf(str, rest...))
}

func (l *GormLogger) Trace(ctx context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error) {
	le := logger.ConfigureLogger(ctx, l.logger)
	sql, rows := fc()
	if err != nil {
		le.Errorf("Took: %s, Err:%s, SQL: %s, AffectedRows: %d", time.Since(begin).String(), err, sql, rows)
	} else {
		le.Tracef("Took: %s, SQL: %s, AffectedRows: %d", time.Since(begin).String(), sql, rows)
	}

}
