package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/saiharsha/money-manager/internal/data"
	"github.com/saiharsha/money-manager/internal/mail"
	"github.com/saiharsha/money-manager/pkg/logger"

	_ "github.com/lib/pq"
)

type config struct {
	port        int
	logLevel    string
	host        string
	environment string
	debugLevel  string
	secretKey   string
	db          struct {
		dsn          string
		maxOpenConns int
		maxIdleConns int
		maxIdleTime  string
	}
	smtp struct {
		host     string
		sslport  int
		tlsport  int
		username string
		password string
		from     string
	}
}

type application struct {
	config config
	logger *logger.Logger
	models data.Models
	mailer *mail.Mailer
	wg     sync.WaitGroup
}

func main() {

	var config config
	flag.IntVar(&config.port, "port", 8080, "API server port")
	flag.StringVar(&config.logLevel, "loglevel", "info", "log level for the application can be one of debug | info | error")
	flag.StringVar(&config.environment, "environment", "development", "env type development | production")
	flag.StringVar(&config.host, "host", "localhost", "host for the application")
	flag.StringVar(&config.debugLevel, "debuglevel", "INFO", "Options can be DEBUG, INFO, ERROR, FATAL, OFF level from lowest to highest")
	secretKey := os.Getenv("SECRET_KEY")
	flag.StringVar(&config.secretKey, "SecretKey", secretKey, "Secret Key for generating json tokens")
	pw := os.Getenv("DATABASE_PASSWORD")
	flag.StringVar(&config.db.dsn, "db-dsn", fmt.Sprintf("postgres://moneymanager:%s@localhost:5432/moneymanager?sslmode=disable", pw), "PostgreSQL DSN")
	flag.IntVar(&config.db.maxOpenConns, "db-max-open-conns", 25, "PostgreSQL max open connections")
	flag.IntVar(&config.db.maxIdleConns, "db-max-idle-conns", 25, "PostgreSQL max open idle connections")
	flag.StringVar(&config.db.maxIdleTime, "db-max-idle-time", "15m", "PostgreSQL max connection idle time")

	flag.StringVar(&config.smtp.host, "smtp-host", "smtp.gmail.com", "SMTP host")
	flag.IntVar(&config.smtp.sslport, "smtp-ssl-port", 465, "SMTP SSL port")
	flag.IntVar(&config.smtp.tlsport, "smtp-tls-port", 587, "SMTP TLS port")
	flag.StringVar(&config.smtp.username, "smtp-username", "moneymanager.bot3330@gmail.com", "SMTP username")
	flag.StringVar(&config.smtp.password, "smtp-password", os.Getenv("SMTP_PASSWORD"), "SMTP password")
	flag.StringVar(&config.smtp.from, "smtp-from", "moneymanager.bot3330@gmail.com", "SMTP from")

	flag.Parse()

	level := logger.StringLevel(config.debugLevel)
	logger := logger.NewLogger(os.Stdout, level.GetLevel())

	db, err := openDB(config)
	if err != nil {
		logger.PrintFatal(err, nil)
	}

	logger.PrintInfo("connected to psql database", nil)

	mailer := mail.NewMailer(
		config.smtp.host,
		config.smtp.sslport,
		config.smtp.username,
		config.smtp.password,
		config.smtp.from,
		true,
		false,
	)

	// Test SMTP connection
	logger.PrintInfo("Testing SMTP connection...", nil)
	if err := mailer.TestConnection(); err != nil {
		logger.PrintError(err, map[string]string{
			"operation": "smtp_connection_test",
		})
		logger.PrintInfo("SMTP connection failed, but continuing with server startup", nil)
	} else {
		logger.PrintInfo("SMTP connection successful", nil)
	}

	app := &application{
		config: config,
		logger: logger,
		models: data.NewModels(db),
		mailer: mailer,
		wg:     sync.WaitGroup{},
	}

	message := fmt.Sprintf("starting server at port %d and host %s", config.port, config.host)
	logger.PrintInfo(message, nil)
	err = app.serve()
	if err != nil {
		logger.PrintFatal(err, nil)
	}
}

func openDB(cfg config) (*sql.DB, error) {
	// Use sql.Open() to create an empty connection pool, using the DSN from the config struct.
	db, err := sql.Open("postgres", cfg.db.dsn)
	if err != nil {
		return nil, err
	}

	// Set the maximum number of open (in-use + idle) connections in the pool.
	// Note that passing a value less than or equal to 0 will mean there is no limit.
	db.SetMaxOpenConns(cfg.db.maxOpenConns)

	// Set the maximum number of idle connection in the pool. Again,
	// passing a value less than or equal to 0 will mean there is no limit
	db.SetMaxIdleConns(cfg.db.maxIdleConns)

	// Use the time.ParseDuration() function to convert the idle timeout duration string to a
	// time.Duration type.
	duration, err := time.ParseDuration(cfg.db.maxIdleTime)
	if err != nil {
		return nil, err
	}

	// Set the maximum idle timeout.
	db.SetConnMaxIdleTime(duration)

	// Create a context with a 5-second timeout deadline.
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = db.PingContext(ctx)
	if err != nil {
		return nil, err
	}

	// Return the sql.DB connection pool.
	return db, nil
}
