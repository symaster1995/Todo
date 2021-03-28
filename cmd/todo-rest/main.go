package main

import (
	"Todo"
	"Todo/http"
	"Todo/store"
	"context"
	"fmt"
	"github.com/spf13/viper"
	"os"
	"os/signal"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())

	c := make(chan os.Signal, 1)

	signal.Notify(c, os.Interrupt)

	go func() { <-c; cancel() }()

	m := NewMain()

	// Execute program.
	if err := m.Run(ctx); err != nil {
		m.Close()
		fmt.Fprintln(os.Stderr, err)
		Todo.ReportError(ctx, err)
		os.Exit(1)
	}

	// Wait for CTRL-C.
	<-ctx.Done()

	// Clean up program.
	if err := m.Close(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

type Main struct {
	HTTPServer *http.Server
	DB         *store.DB
	Config     Config
}

func NewMain() *Main {
	return &Main{
		HTTPServer: http.NewServer(),
		DB:         store.NewDB(),
	}
}

func (m *Main) Close() error {
	if m.HTTPServer != nil {
		if err := m.HTTPServer.Close(); err != nil {
			return err
		}
	}
	return nil
}

func (m *Main) Run(ctx context.Context) (err error) {
	//Set Config From yml file
	config, err := setConfig()
	if err != nil {
		return err
	}

	//Set MySql DB Config
	m.DB.DSN = dsn(config.Database)

	if err := m.DB.Open(); err != nil {
		return fmt.Errorf("cannot open db: %w", err)
	}

	//Instantiate store service
	taskService := store.NewTaskService(m.DB)

	// Copy configuration settings to the HTTP server
	m.HTTPServer.Addr = config.Server.Port

	// Attach task service to the HTTP server
	m.HTTPServer.TaskService = taskService

	if err := m.HTTPServer.Open(); err != nil {
		return err
	}

	return nil
}

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
}

type ServerConfig struct {
	Port   int
}

type DatabaseConfig struct {
	DBName     string
	DBUser     string
	DBPassword string
	DBHost     string
}

func setConfig() (Config, error) {
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	viper.AutomaticEnv()
	viper.SetConfigType("yml")

	var config Config

	if err := viper.ReadInConfig(); err != nil {
		return config, err
	}

	// Set undefined variables
	viper.SetDefault("database.dbname", "todo")

	err := viper.Unmarshal(&config)
	if err != nil {
		return config, err
	}

	return config, nil

}

func dsn(db DatabaseConfig) string {
	return fmt.Sprintf("%s:%s@tcp(%s)/%s?parseTime=true", db.DBUser, db.DBPassword, db.DBHost, db.DBName)
}
