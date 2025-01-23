package main

import (
	"em4/config"
	"fmt"
	"os"

	"em4/internal/storage"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

func run() error {
	err := godotenv.Load()
	if err != nil {
		return fmt.Errorf("error loading environment: %v", err)
	}

	cfg := config.Config{}
	if err := envconfig.Process("", &cfg); err != nil {
		return fmt.Errorf("error loading environment: %v", err)
	}

	classic_storage, err := storage.NewPostgresStorage(cfg.Storage.Path)
	if err != nil {
		return fmt.Errorf("error initializing storage: %v", err)
	}
	defer classic_storage.Close()

	orm_storage, err := storage.NewORMPostgresStorage(cfg.Storage.Path)
	if err != nil {
		return fmt.Errorf("error initializing storage: %v", err)
	}
	defer orm_storage.Close()

	pgx_storage, err := storage.NewPGXStorage(cfg.Storage.Path)
	if err != nil {
		return fmt.Errorf("error initializing storage: %v", err)
	}
	defer pgx_storage.Close()

	// call funcs ...
	return nil
}

func main() {
	if err := run(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
