package db

import (
	"api/config"
	"errors"
)

type DB interface {
	RegisterStreamServer(host string) error
	StreamServerKeepLive(host string) error
	UnregisterStreamServer(host string) error
}

var db DB

func Get() DB {
	return db
}

func Init(cfg *config.Conf) error {
	var err error
	switch cfg.DB.Driver {
	case "sqlite":
		db, err = newSqlite(cfg.DB.FilePath)
	default:
		err = ErrUnknownDatabaseDriver
	}
	return err
}

var (
	ErrUnknownDatabaseDriver = errors.New("unknown database driver")
)