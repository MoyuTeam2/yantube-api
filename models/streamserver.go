package models

type StreamServer struct {
	Host      string `json:"host" gorm:"primaryKey"`
	Heartbeat int64  `json:"heartbeat"`
	IsActive  bool   `json:"is_active"`
}

func (StreamServer) TableName() string {
	return "stream_server"
}
