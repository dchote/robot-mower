package control

import ()

type MowerStateStruct struct {
	Battery struct {
		Status  string  `json:"status"`
		Voltage float32 `json:"voltage"`
		Current float32 `json:"current"`
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
		Speed int `json:"speed"`
	} `json:"drive"`
	Cutter struct {
		Speed int `json:"speed"`
	} `json:"cutter"`
}

var (
	MowerState *MowerStateStruct
)
