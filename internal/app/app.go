package app

import (
	"context"
	"log"
	"test-task1/pkg/config"
)

type App struct {
	cfg *config.Config
	// db  *sql.DB

}

func NewApp(ctx context.Context) (*App, error) {
	cfg, err := config.LoadConfig()
	if err != nil {
		return nil, err
	}
	log.Printf("Config loaded: %+v\n", cfg)

	return nil, nil
}

func (a *App) Run(ctx context.Context) error {
	return nil
}
