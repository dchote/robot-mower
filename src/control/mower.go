package control

import ()

type MowerStateStruct struct {
	Platform struct {
		Hostname        string `json:"hostname"`
		OperatingSystem string `json:"operating_system"`
		Platform        string `json:"platform"`
		LoadAverage     struct {
			Load1  float64 `json:"load1"`
			Load5  float64 `json:"load5"`
			Load15 float64 `json:"load15"`
		} `json:"load_average"`
		MemoryUsage struct {
			Total     uint64 `json:"total"`
			Available uint64 `json:"available"`
		} `json:"memory"`
		DiskUsage struct {
			Total uint64 `json:"total"`
			Free  uint64 `json:"free"`
		} `json:"disk"`
	} `json:"platform"`
	Battery struct {
		Status         string  `json:"status"`
		VoltageNominal float32 `json:"voltage_nominal"`
		VoltageWarn    float32 `json:"voltage_warn"`
		Voltage        float32 `json:"voltage"`
		Current        float32 `json:"current"`
	} `json:"battery"`
	Compass struct {
		Status  string `json:"status"`
		Bearing string `json:"bearing"`
	} `json:"compass"`
	GPS struct {
		Status      string `json:"status"`
		Coordinates string `json:"coordinates"`
	} `json:"gps"`
	Drive struct {
		Speed     int    `json:"speed"`
		Direction string `json:"direction"`
	} `json:"drive"`
	Cutter struct {
		Speed int `json:"speed"`
	} `json:"cutter"`
}

var (
	MowerState *MowerStateStruct
)
