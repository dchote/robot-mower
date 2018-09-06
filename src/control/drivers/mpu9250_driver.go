package drivers

//
// MPU9250 driver written by Daniel Chote
//
// https://github.com/dchote/robot-mower/blob/master/src/control/drivers/mpu9250_driver.go
//

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
	//"fmt"
	"log"
	"math"
	"time"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/drivers/i2c"
)

const (
	// MPU9250 Default I2C slave address
	SLAVE_ADDRESS = 0x68
	// AK8963 I2C slave address
	AK8963_SLAVE_ADDRESS = 0x0C
	// Device id
	DEVICE_ID = 0x71

	WHO_AM_I_MPU9250 = 0x75
	AK8963_WHO_AM_I  = 0x00

	// sample rate driver
	SMPLRT_DIV     = 0x19
	CONFIG         = 0x1A
	GYRO_CONFIG    = 0x1B
	ACCEL_CONFIG   = 0x1C
	ACCEL_CONFIG_2 = 0x1D
	LP_ACCEL_ODR   = 0x1E
	WOM_THR        = 0x1F
	FIFO_EN        = 0x23
	I2C_MST_CTRL   = 0x24
	I2C_MST_STATUS = 0x36
	INT_PIN_CFG    = 0x37
	INT_ENABLE     = 0x38
	INT_STATUS     = 0x3A

	ACCEL_XOUT = 0x3B
	ACCEL_YOUT = 0x3D
	ACCEL_ZOUT = 0x3F

	TEMP_OUT = 0x41

	GYRO_XOUT = 0x43
	GYRO_YOUT = 0x45
	GYRO_ZOUT = 0x47

	MAG_XOUT  = 0x03
	MAG_YOUT  = 0x05
	MAG_ZOUT  = 0x07
	MAG_CHECK = 0x09

	I2C_MST_DELAY_CTRL = 0x67
	SIGNAL_PATH_RESET  = 0x68
	MOT_DETECT_CTRL    = 0x69
	USER_CTRL          = 0x6A
	PWR_MGMT_1         = 0x6B
	PWR_MGMT_2         = 0x6C
	FIFO_R_W           = 0x74
	WHO_AM_I           = 0x75

	// Gyro Full Scale Select 250dps
	GFS_250 = 0x00
	// Gyro Full Scale Select 500dps
	GFS_500 = 0x01
	// Gyro Full Scale Select 1000dps
	GFS_1000 = 0x02
	// Gyro Full Scale Select 2000dps
	GFS_2000 = 0x03
	// Accel Full Scale Select 2G
	AFS_2G = 0x00
	// Accel Full Scale Select 4G
	AFS_4G = 0x01
	// Accel Full Scale Select 8G
	AFS_8G = 0x02
	// Accel Full Scale Select 16G
	AFS_16G = 0x03

	// AK8963 Register Addresses
	AK8963_WIA        = 0x00
	AK8963_INFO       = 0x01
	AK8963_ST1        = 0x02
	AK8963_HXL        = 0x03
	AK8963_MAGNET_OUT = 0x03
	AK8963_CNTL1      = 0x0A
	AK8963_CNTL2      = 0x0B
	AK8963_I2CDIS     = 0x0F

	AK8963_ASAX = 0x10
	AK8963_ASAY = 0x11
	AK8963_ASAZ = 0x12

	// CNTL1 Mode select
	// Power down mode
	AK8963_MODE_DOWN = 0x00
	// One shot data output
	AK8963_MODE_ONE = 0x01

	// Continous data output 8Hz
	AK8963_MODE_C8HZ = 0x02
	// Continous data output 100Hz
	AK8963_MODE_C100HZ = 0x06

	// Magneto Scale Select
	// 14bit output
	AK8963_BIT_14 = 0x00
	// 16bit output
	AK8963_BIT_16 = 0x01

	MPUREG_I2C_SLV0_DO        = 0x63
	MPUREG_I2C_SLV1_DO        = 0x64
	MPUREG_I2C_MST_CTRL       = 0x24
	MPUREG_I2C_SLV0_ADDR      = 0x25
	MPUREG_I2C_SLV0_REG       = 0x26
	MPUREG_I2C_SLV0_CTRL      = 0x27
	MPUREG_I2C_SLV1_ADDR      = 0x28
	MPUREG_I2C_SLV1_REG       = 0x29
	MPUREG_I2C_SLV1_CTRL      = 0x2A
	MPUREG_I2C_SLV2_ADDR      = 0x2B
	MPUREG_I2C_SLV2_REG       = 0x2C
	MPUREG_I2C_SLV2_CTRL      = 0x2D
	MPUREG_I2C_SLV3_ADDR      = 0x2E
	MPUREG_I2C_SLV3_REG       = 0x2F
	MPUREG_I2C_SLV3_CTRL      = 0x30
	MPUREG_I2C_MST_DELAY_CTRL = 0x67
	MPUREG_I2C_SLV4_ADDR      = 0x31
	MPUREG_I2C_SLV4_REG       = 0x32
	MPUREG_I2C_SLV4_DO        = 0x33
	MPUREG_I2C_SLV4_CTRL      = 0x34
	MPUREG_I2C_SLV4_DI        = 0x35
	MPUREG_INT_PIN_CFG        = 0x37
	MPUREG_EXT_SENS_DATA_00   = 0x49
	MPUREG_EXT_SENS_DATA_01   = 0x4A
	MPUREG_EXT_SENS_DATA_02   = 0x4B
	MPUREG_EXT_SENS_DATA_03   = 0x4C
	MPUREG_EXT_SENS_DATA_04   = 0x4D
	MPUREG_EXT_SENS_DATA_05   = 0x4E
	MPUREG_EXT_SENS_DATA_06   = 0x4F
	MPUREG_EXT_SENS_DATA_07   = 0x50
	MPUREG_EXT_SENS_DATA_08   = 0x51
	MPUREG_EXT_SENS_DATA_09   = 0x52
	MPUREG_EXT_SENS_DATA_10   = 0x53
	MPUREG_EXT_SENS_DATA_11   = 0x54
	MPUREG_EXT_SENS_DATA_12   = 0x55
	MPUREG_EXT_SENS_DATA_13   = 0x56
	MPUREG_EXT_SENS_DATA_14   = 0x57
	MPUREG_EXT_SENS_DATA_15   = 0x58
	MPUREG_EXT_SENS_DATA_16   = 0x59
	MPUREG_EXT_SENS_DATA_17   = 0x5A
	MPUREG_EXT_SENS_DATA_18   = 0x5B
	MPUREG_EXT_SENS_DATA_19   = 0x5C
	MPUREG_EXT_SENS_DATA_20   = 0x5D
	MPUREG_EXT_SENS_DATA_21   = 0x5E
	MPUREG_EXT_SENS_DATA_22   = 0x5F
	MPUREG_EXT_SENS_DATA_23   = 0x60

	READ_FLAG                    = 0x80
	MPU_BANK_SIZE                = 0xFF
	CFG_MOTION_BIAS              = 0x4B8 // Enable/disable gyro bias compensation
	BIT_FIFO_SIZE_1024           = 0x40  // FIFO buffer size
	BIT_AUX_IF_EN          uint8 = 0x20
	BIT_BYPASS_EN                = 0x02
	AKM_POWER_DOWN               = 0x00
	BIT_I2C_READ                 = 0x80
	BIT_SLAVE_EN                 = 0x80
	AKM_SINGLE_MEASUREMENT       = 0x01
	INV_CLK_PLL                  = 0x01
	AK89xx_FSR                   = 9830
	AKM_DATA_READY               = 0x01
	AKM_DATA_OVERRUN             = 0x02
	AKM_OVERFLOW                 = 0x80
)

// MPUData contains all the values measured by an MPU9250.
type MPUData struct {
	G1, G2, G3 float64
	A1, A2, A3 float64
	M1, M2, M3 float64
	Temp       float64
}

type MPU9250Driver struct {
	name       string
	connector  i2c.Connector
	connection i2c.Connection
	i2c.Config

	gResolution float64
	aResolution float64
	mResolution float64

	magXcoef float64
	magYcoef float64
	magZcoef float64

	Data *MPUData
}

var compassBearing = []struct {
	label        string
	start, ended float64
}{
	{"N", 000.00, 011.25}, {"NNE", 011.25, 033.75}, {"NE", 033.75, 056.25}, {"ENE", 056.25, 078.75},
	{"E", 078.75, 101.25}, {"ESE", 101.25, 123.75}, {"SE", 123.75, 146.25}, {"SSE", 146.25, 168.75},
	{"S", 168.75, 191.25}, {"SSW", 191.25, 213.75}, {"SW", 213.75, 236.25}, {"WSW", 236.25, 258.75},
	{"W", 258.75, 281.25}, {"WNW", 281.25, 303.75}, {"NW", 303.75, 326.25}, {"NNW", 326.25, 348.75},
	{"N", 348.75, 360.00},
}

// NewMPU9250Driver creates a new Gobot Driver for an MPU9250 I2C Accelerometer/Gyroscope.
//
// Params:
//		conn Connector - the Adaptor to use with this Driver
//
// Optional params:
//		i2c.WithBus(int):	bus to use with this driver
//		i2c.WithAddress(int):	address to use with this driver
//
func NewMPU9250Driver(a i2c.Connector, options ...func(i2c.Config)) *MPU9250Driver {
	m := &MPU9250Driver{
		name:      gobot.DefaultName("MPU9250"),
		connector: a,
		Config:    i2c.NewConfig(),
	}

	for _, option := range options {
		option(m)
	}

	// TODO: add commands to API
	return m
}

// Name returns the name of the device.
func (mpu *MPU9250Driver) Name() string { return mpu.name }

// SetName sets the name of the device.
func (mpu *MPU9250Driver) SetName(n string) { mpu.name = n }

// Connection returns the connection for the device.
func (mpu *MPU9250Driver) Connection() gobot.Connection { return mpu.connector.(gobot.Connection) }

// Start writes initialization bytes to sensor
func (mpu *MPU9250Driver) Start() (err error) {
	if err := mpu.initialize(); err != nil {
		return err
	}

	return
}

// Halt returns true if devices is halted successfully
func (mpu *MPU9250Driver) Halt() (err error) {
	return
}

func (mpu *MPU9250Driver) CurrentHeading() (heading float64, headingLabel string, err error) {
	if mpu.Data.M1 == 0 && mpu.Data.M2 == 0 {
		return 0, "", errors.New("MPU9250Driver Warning: no magnetometer data")
	}

	heading = 180 * math.Atan2(mpu.Data.M1, mpu.Data.M2) / math.Pi
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

func (mpu *MPU9250Driver) initialize() (err error) {
	bus := mpu.GetBusOrDefault(mpu.connector.GetDefaultBus())
	address := mpu.GetAddressOrDefault(SLAVE_ADDRESS)

	mpu.connection, err = mpu.connector.GetConnection(address, bus)
	if err != nil {
		return err
	}

	//mpu.connection.WriteByteData(register, value)
	//mpu.connection.WriteBlockData(register data)
	//mpu.connection.ReadByteData(register)
	//mpu.connection.ReadWordData(register) (int16)
	mpuIdentity, _ := mpu.connection.ReadByteData(WHO_AM_I_MPU9250)
	log.Printf("MPU9250Driver mpuIdentity: %v should be %v", mpuIdentity, 0x71)

	// reset and autoselect clock source
	mpu.connection.WriteByteData(PWR_MGMT_1, 0x00)
	time.Sleep(10 * time.Millisecond)

	mpu.connection.WriteByteData(PWR_MGMT_1, 0x01)
	mpu.connection.WriteByteData(PWR_MGMT_2, 0x00)

	time.Sleep(10 * time.Millisecond)

	// configure gyro/accelerometer
	mpu.connection.WriteByteData(CONFIG, 0x03)
	mpu.connection.WriteByteData(SMPLRT_DIV, 0x04)

	mpu.connection.WriteByteData(GYRO_CONFIG, GFS_250<<3) // gres = 250.0/32768.0
	mpu.gResolution = 250.0 / 32768.0

	mpu.connection.WriteByteData(ACCEL_CONFIG, AFS_2G<<3) // ares = 2.0/32768.0
	mpu.aResolution = 2.0 / 32768.0

	mpu.connection.WriteByteData(ACCEL_CONFIG_2, 0x03)

	// configure interrupts and setup i2c master
	mpu.connection.WriteByteData(INT_PIN_CFG, 0x30)
	mpu.connection.WriteByteData(USER_CTRL, 0x20)
	mpu.connection.WriteByteData(MPUREG_I2C_MST_CTRL, 0x0D)

	// setup the i2c slave for the AK8963
	mpu.connection.WriteByteData(MPUREG_I2C_SLV0_ADDR, AK8963_SLAVE_ADDRESS)
	mpu.connection.WriteByteData(MPUREG_I2C_SLV0_REG, AK8963_CNTL2)

	// reset AK8963
	mpu.connection.WriteByteData(MPUREG_I2C_SLV0_DO, 0x01)
	mpu.connection.WriteByteData(MPUREG_I2C_SLV0_CTRL, 0x81)

	mpu.connection.WriteByteData(MPUREG_I2C_SLV0_DO, 0x12)
	mpu.connection.WriteByteData(MPUREG_I2C_SLV0_CTRL, 0x81)

	time.Sleep(10 * time.Millisecond)

	// read mag calibration data
	mpu.connection.WriteByteData(MPUREG_I2C_SLV0_ADDR, AK8963_SLAVE_ADDRESS|READ_FLAG)
	mpu.connection.WriteByteData(MPUREG_I2C_SLV0_REG, AK8963_ASAX)
	mpu.connection.WriteByteData(MPUREG_I2C_SLV0_CTRL, 0x83)

	time.Sleep(10 * time.Millisecond)

	if err = mpu.connection.WriteByte(MPUREG_EXT_SENS_DATA_00); err != nil {
		return errors.New("MPU9250Driver mag coef read error")
	}

	magBuf := []byte{0, 0, 0}
	_, _ = mpu.connection.Read(magBuf)

	mpu.magXcoef = (float64(magBuf[0])-128)/256.0 + 1.0
	mpu.magYcoef = (float64(magBuf[1])-128)/256.0 + 1.0
	mpu.magZcoef = (float64(magBuf[2])-128)/256.0 + 1.0

	log.Printf("mag coef: %v, %v, %v - %v, %v, %v", magBuf[0], magBuf[1], magBuf[2], mpu.magXcoef, mpu.magYcoef, mpu.magZcoef)

	// AK8963 power down
	mpu.connection.WriteByteData(MPUREG_I2C_SLV0_DO, 0x00)
	time.Sleep(10 * time.Millisecond)

	// set scale&continuous mode
	mpu.connection.WriteByteData(MPUREG_I2C_SLV0_CTRL, (AK8963_BIT_16<<4 | AK8963_MODE_C8HZ))
	mpu.mResolution = 4912.0 / 32760.0

	time.Sleep(100 * time.Millisecond)

	// attempt to fetch data for the first time
	go mpu.GetData()

	return nil
}

// readSensors polls the gyro, accelerometer and magnetometer sensors as well as the die temperature.
func (mpu *MPU9250Driver) GetData() (err error) {
	dataReady, _ := mpu.connection.ReadByteData(INT_STATUS)
	if dataReady != 0x01 {
		return errors.New("MPU9250Driver Error: Data not ready")
	}

	var (
		g1, g2, g3, a1, a2, a3, m1, m2, m3, m4 int16
	)

	g1, _ = mpu.readWordFromRegisterLittle(GYRO_XOUT)
	g2, _ = mpu.readWordFromRegisterLittle(GYRO_YOUT)
	g3, _ = mpu.readWordFromRegisterLittle(GYRO_ZOUT)

	a1, _ = mpu.readWordFromRegisterLittle(ACCEL_XOUT)
	a2, _ = mpu.readWordFromRegisterLittle(ACCEL_YOUT)
	a3, _ = mpu.readWordFromRegisterLittle(ACCEL_ZOUT)

	temp, _ := mpu.readWordFromRegisterBig(TEMP_OUT)

	mpu.Data = &MPUData{
		G1:   float64(g1) * mpu.gResolution,
		G2:   float64(g2) * mpu.gResolution,
		G3:   float64(g3) * mpu.gResolution,
		A1:   float64(a1) * mpu.aResolution,
		A2:   float64(a2) * mpu.aResolution,
		A3:   float64(a3) * mpu.aResolution,
		Temp: float64(temp)/333.87 + 21.0,
	}

	mpu.connection.WriteByteData(MPUREG_I2C_SLV0_ADDR, AK8963_SLAVE_ADDRESS|READ_FLAG)
	mpu.connection.WriteByteData(MPUREG_I2C_SLV0_REG, AK8963_HXL)
	mpu.connection.WriteByteData(MPUREG_I2C_SLV0_CTRL, 0x87)

	time.Sleep(10 * time.Millisecond)

	magRegMap := map[*int16]byte{
		&m1: MPUREG_EXT_SENS_DATA_00, &m2: MPUREG_EXT_SENS_DATA_02, &m3: MPUREG_EXT_SENS_DATA_04, &m4: MPUREG_EXT_SENS_DATA_06,
	}
	for p, reg := range magRegMap {
		*p, _ = mpu.readWordFromRegisterLittle(reg)
	}

	// validate mag data
	if (byte(m1&0xFF)&AKM_DATA_READY) == 0x00 && (byte(m1&0xFF)&AKM_DATA_OVERRUN) != 0x00 {
		log.Println("MPU9250Driver Warning: mag data not ready or overflow")
	} else if (byte((m4>>8)&0xFF) & AKM_OVERFLOW) != 0x00 {
		log.Println("MPU9250Driver Warning: mag data overflow")
	} else {
		mpu.Data.M1 = float64(m1) * mpu.mResolution * mpu.magXcoef
		mpu.Data.M2 = float64(m2) * mpu.mResolution * mpu.magYcoef
		mpu.Data.M3 = float64(m3) * mpu.mResolution * mpu.magZcoef
	}

	return nil
}

func (mpu *MPU9250Driver) readWordFromRegisterBig(reg uint8) (val int16, err error) {
	if err = mpu.connection.WriteByte(reg); err != nil {
		return 0, errors.New("MPU9250Driver readWordFromRegister error")
	}

	buf := []byte{0, 0}
	if _, err = mpu.connection.Read(buf); err != nil {
		return 0, errors.New("MPU9250Driver readWordFromRegister error")
	}

	return int16(buf[0])<<8 + int16(buf[1]), nil
}

func (mpu *MPU9250Driver) readWordFromRegisterLittle(reg uint8) (val int16, err error) {
	bigEndian, err := mpu.readWordFromRegisterBig(reg)
	if err != nil {
		return 0, err
	}
	return (bigEndian&0xFF)<<8 + bigEndian>>8, nil // return littleEndian
}
