package domain

type Sensor struct {
	Name    string `json:"name"`
	Address string `json:"address"`
}

type SensorInfo struct {
	Enabled  bool `json:"enabled"`
	Detected bool `json:"detected"`
}
