package db

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
	log "github.com/sirupsen/logrus"
	"github.com/task-manager/config"
)

type DB struct {
	Conn *sql.DB
}

func InitDB(cfg config.Config) (*DB, error) {
	var err error
	db := &DB{}
	conn, err := sql.Open(cfg.DB.DriverName,
		cfg.DB.URL)
	if err != nil {
		log.Errorf("couldn't connect to db: %v", err)
		return db, err
	}

	if conn.Ping() != nil {
		log.Errorf("couldn't connect: %v", conn.Ping())
		return db, fmt.Errorf("couldn't connect:%v", conn.Ping())
	}
	return &DB{Conn: conn}, nil
}
