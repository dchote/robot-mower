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
		CPULoad struct {
			Count int     `json:"count"`
			Total float64 `json:"total"`
			Core1 float64 `json:"core_1"`
			Core2 float64 `json:"core_2"`
			Core3 float64 `json:"core_3"`
			Core4 float64 `json:"core_4"`
			Core5 float64 `json:"core_5"`
			Core6 float64 `json:"core_6"`
			Core7 float64 `json:"core_7"`
			Core8 float64 `json:"core_8"`
		} `json:"cpu"`
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
		VoltageNominal float64 `json:"voltage_nominal"`
		VoltageWarn    float64 `json:"voltage_warn"`
		Voltage        float64 `json:"voltage"`
		Current        float64 `json:"current"`
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
