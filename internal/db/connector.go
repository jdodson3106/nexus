package db

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/jdodson3106/nexus/log"
	_ "github.com/lib/pq"
)

type ConnectionConfiguration struct {
	Host               string
	Port               int
	Database           string
	Username           string
	Password           string
	DriverName         string
	ConnectionPoolSize int

	connectionPoolMax int
}

type DB struct {
	conn   *sql.DB
	config *ConnectionConfiguration
}

func (d *DB) Close() error {
	return d.conn.Close()
}

func DefaultConfiguration() *ConnectionConfiguration {
	return &ConnectionConfiguration{
		Host:               "localhost",
		Port:               5432,
		Database:           "postgres",
		Username:           "postgres",
		Password:           "password",
		DriverName:         "postgres",
		ConnectionPoolSize: 10,
		connectionPoolMax:  100,
	}
}

func NewDefaultDbConnection() (*DB, error) {
	return NewDbConnection(DefaultConfiguration())
}

func NewDbConnection(c *ConnectionConfiguration) (*DB, error) {
	connStr := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		c.Host, c.Port, c.Username, c.Password, c.Database,
	)
	conn, err := sql.Open(c.Database, connStr)
	if err != nil {
		e := fmt.Sprintf("error connecting to db: %s", err)
		log.Error(e)
		return nil, errors.New(e)
	}

	conn.SetMaxIdleConns(int(c.ConnectionPoolSize))
	conn.SetMaxOpenConns(int(c.connectionPoolMax))

	db := &DB{conn: conn, config: c}
	return db, nil
}
