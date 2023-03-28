package db

import (
	"database/sql"
	"errors"
	"fmt"
	"log"

	"outbox/config"

	_ "github.com/go-sql-driver/mysql"
)

type Configuration interface {
	Configure() (*sql.DB, error)
}

type mysql struct {
	cfg *config.AppConfig
}

func NewMySql(cfg *config.AppConfig) Configuration {
	return &mysql{cfg: cfg}
}

func (c *mysql) Configure() (*sql.DB, error) {
	if c.cfg.Connection == "" {
		log.Panic("Attempt to connect to MySQL with an unset connection value")
		return nil, errors.New("unset db connection value")
	}

	// Open the connection and write the settings
	client, err := sql.Open("mysql", c.cfg.Connection)
	if err != nil {
		return nil, err
	}
	client.SetConnMaxLifetime(c.cfg.ConnMaxLifetime)
	if int(c.cfg.MaxIdleConns) > 0 {
		client.SetMaxIdleConns(int(c.cfg.MaxIdleConns))
	}
	client.SetMaxOpenConns(int(c.cfg.MaxOpenConns))

	// We ping just to ensure the connection is established
	if err := client.Ping(); err != nil {
		if err := client.Close(); err != nil {
			log.Panicf("Unexpected error while closing MySQL connection after ping failure. %+v", err)
			return nil, err
		}
		fmt.Printf("error while generate ping to mysql server. %+v", err)
		return nil, err
	}

	return client, nil
}
