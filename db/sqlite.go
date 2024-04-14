package db

import (
	"api/models"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Sqlite struct {
	filename string
	db       *gorm.DB
}

func newSqlite(filename string) (*Sqlite, error) {
	db, err := gorm.Open(sqlite.Open(filename))
	if err != nil {
		return nil, err
	}
	err = db.AutoMigrate(&models.StreamServer{})
	if err != nil {
		return nil, err
	}
	return &Sqlite{
		filename: filename,
		db:       db,
	}, nil
}

func (s *Sqlite) RegisterStreamServer(host string) error {
	return s.db.Clauses(clause.OnConflict{
		Columns: []clause.Column{{Name: "host"}},
		DoUpdates: clause.Assignments(map[string]interface{}{
			"heartbeat": time.Now().Unix(),
			"is_active": true,
		}),
	}).Create(&models.StreamServer{
		Host:      host,
		Heartbeat: time.Now().Unix(),
		IsActive:  true,
	}).Error
}

func (s *Sqlite) StreamServerKeepLive(host string) error {
	result := s.db.Model(&models.StreamServer{}).
		Where("host = ?", host).
		Where("is_active = ?", true).
		Update("heartbeat", time.Now().Unix())
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrStreamServerUnregistered
	}
	return nil
}

func (s *Sqlite) UnregisterStreamServer(host string) error {
	result := s.db.Model(&models.StreamServer{}).
		Where("host = ?", host).
		Where("is_active = ?", true).
		Update("is_active", false)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrStreamServerUnregistered
	}
	return nil
}
