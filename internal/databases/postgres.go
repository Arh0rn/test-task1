package databases

import (
	"database/sql"
	"fmt"
	"github.com/avast/retry-go"
	"test-task1/pkg/config"
	"time"
)

var (
	dsnTemplate      = "host=%s user=%s password=%s dbname=%s sslmode=disable"
	attempts    uint = 3
	delay            = time.Second
)

func NewPostgresConnection(c *config.Database) (*sql.DB, error) {
	dsn := fmt.Sprintf(dsnTemplate, c.Host, c.User, c.Password, c.DBName)
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	err = retry.Do(
		db.Ping,
		retry.Attempts(attempts),
		retry.Delay(delay),
	)
	if err != nil {
		return nil, err
	}

	return db, nil
}
