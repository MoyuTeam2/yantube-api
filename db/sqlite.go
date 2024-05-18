package db

import (
	"api/models"
	"errors"
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
	err = db.AutoMigrate(&models.StreamServer{}, &models.User{})
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

func (s *Sqlite) GetAllActiveStreamServers() ([]*models.StreamServer, error) {
	var servers []*models.StreamServer
	err := s.db.Where("is_active = ?", true).Find(&servers).Error
	if err != nil {
		return nil, err
	}
	return servers, nil
}

func (s *Sqlite) CreateUser(user *models.User) error {
	return s.db.Create(user).Error
}

func (s *Sqlite) GetUserByUsername(username string) (*models.User, error) {
	var user models.User
	err := s.db.Where("username = ?", username).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}
	return &user, nil
}

func (s *Sqlite) RevokeUserStreamCode(username string, streamCode string) error {
	result := s.db.Model(&models.User{}).
		Where("username = ?", username).
		Update("stream_code", streamCode)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrUserNotFound
	}
	return nil
}

func (s *Sqlite) DeleteUserByUsername(username string) error {
	result := s.db.Where("username = ?", username).Delete(&models.User{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrUserNotFound
	}
	return nil
}

func (s *Sqlite) Close() error {
	sqlDB, err := s.db.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}
