package control

//
// Compass bearing logic from https://github.com/pd0mz/go-maidenhead/blob/master/point.go
//
// The MIT License (MIT)
//
// Copyright (c) 2016 pd0mz
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

import (
	"errors"
	"log"
	"math"
	"time"

	"github.com/dchote/robot-mower/src/control/drivers"
	//"github.com/dchote/robot-mower/src/control/filters"
)

const (
	RAD_TO_DEG = math.Pi / 180
)

var (
	// Compass bearing constraints
	compassBearing = []struct {
		label        string
		start, ended float64
	}{
		{"N", 000.00, 011.25}, {"NNE", 011.25, 033.75}, {"NE", 033.75, 056.25}, {"ENE", 056.25, 078.75},
		{"E", 078.75, 101.25}, {"ESE", 101.25, 123.75}, {"SE", 123.75, 146.25}, {"SSE", 146.25, 168.75},
		{"S", 168.75, 191.25}, {"SSW", 191.25, 213.75}, {"SW", 213.75, 236.25}, {"WSW", 236.25, 258.75},
		{"W", 258.75, 281.25}, {"WNW", 281.25, 303.75}, {"NW", 303.75, 326.25}, {"NNW", 326.25, 348.75},
		{"N", 348.75, 360.00},
	}

	IMUDeltaTime time.Time
)

func InitFilters() {

}

func SetIMUValues(data *drivers.MPUData) {
	newTime := time.Now()
	deltaTime := newTime.Sub(IMUDeltaTime)
	IMUDeltaTime = newTime

	dt := float64(deltaTime / time.Millisecond)

	// log the values
	log.Println("----------")
	log.Printf("IMU dt: %v", dt)
	log.Printf("IMU Temperature: %v", data.Temp)
	log.Printf("IMU Accelerometer: %v, %v, %v", data.A1, data.A2, data.A3)
	log.Printf("IMU Gyroscope: %v, %v, %v", data.G1, data.G2, data.G3)
	log.Printf("IMU Magnetometer: %v, %v, %v", data.M1, data.M2, data.M3)

	/*log.Println("----------")
	log.Printf("IMU roll: %v", roll)
	log.Printf("IMU gyroXangle: %v", gyroXangle)
	log.Printf("IMU kalAngle_X: %v", kalAngle_X)

	log.Printf("IMU pitch: %v", pitch)
	log.Printf("IMU gyroYangle: %v", gyroYangle)
	log.Printf("IMU kalAngle_Y: %v", kalAngle_Y)
	log.Println("----------")
	*/

	heading, label, err := CurrentHeading(data.M1, data.M2)
	if err == nil {

		log.Printf("CurrentHeading: %v, %v", heading, label)
		log.Println("----------")
	}

}

// Compass bearing methods

func CurrentHeading(M1 float64, M2 float64) (heading float64, headingLabel string, err error) {
	if M1 == 0 && M2 == 0 {
		return 0, "", errors.New("M1 & M2 must not be 0")
	}

	heading = 180 * math.Atan2(M1, M2) / math.Pi
	if heading < 0 {
		heading += 360
	}

	for _, compass := range compassBearing {
		if heading >= compass.start && heading <= compass.ended {
			return heading, compass.label, nil
		}
	}

	return heading, "", nil
}
