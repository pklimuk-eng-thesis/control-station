package domain

import "time"

type Sensor struct {
	Name    string `json:"name"`
	Address string `json:"address"`
}

type SensorInfo struct {
	Enabled  bool `json:"enabled"`
	Detected bool `json:"detected"`
}

type SensorData struct {
	ID        int       `json:"id" db:"id"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	IsEnabled bool      `json:"is_enabled" db:"is_enabled"`
	Detected  bool      `json:"detected"  db:"detected"`
}
