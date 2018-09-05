package drivers

//
// This code is adapted from @westphae's MPU9250 implementation that can be found here: https://github.com/westphae/goflying/blob/master/mpu9250/mpu9250.go
//
// Approach adapted from the InvenSense DMP 6.1 drivers
// Also referenced https://github.com/brianc118/MPU9250/blob/master/MPU9250.cpp
//
// The MIT License (MIT)
// Copyright (c) 2016 Eric Westphal
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
	"fmt"
	"log"
	"math"
	"time"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/drivers/i2c"
)

const (
	MPU_ADDRESS               = 0x68
	MPUREG_XG_OFFS_TC         = 0x00
	MPUREG_YG_OFFS_TC         = 0x01
	MPUREG_ZG_OFFS_TC         = 0x02
	MPUREG_X_FINE_GAIN        = 0x03
	MPUREG_Y_FINE_GAIN        = 0x04
	MPUREG_Z_FINE_GAIN        = 0x05
	MPUREG_XA_OFFS_H          = 0x06
	MPUREG_XA_OFFS_L          = 0x07
	MPUREG_YA_OFFS_H          = 0x08
	MPUREG_YA_OFFS_L          = 0x09
	MPUREG_ZA_OFFS_H          = 0x0A
	MPUREG_ZA_OFFS_L          = 0x0B
	MPUREG_PRODUCT_ID         = 0x0C
	MPUREG_SELF_TEST_X        = 0x0D
	MPUREG_SELF_TEST_Y        = 0x0E
	MPUREG_SELF_TEST_Z        = 0x0F
	MPUREG_SELF_TEST_A        = 0x10
	MPUREG_XG_OFFS_USRH       = 0x13
	MPUREG_XG_OFFS_USRL       = 0x14
	MPUREG_YG_OFFS_USRH       = 0x15
	MPUREG_YG_OFFS_USRL       = 0x16
	MPUREG_ZG_OFFS_USRH       = 0x17
	MPUREG_ZG_OFFS_USRL       = 0x18
	MPUREG_SMPLRT_DIV         = 0x19
	MPUREG_CONFIG             = 0x1A
	MPUREG_GYRO_CONFIG        = 0x1B
	MPUREG_ACCEL_CONFIG       = 0x1C
	MPUREG_ACCEL_CONFIG_2     = 0x1D
	MPUREG_LP_ACCEL_ODR       = 0x1E
	MPUREG_MOT_THR            = 0x1F
	MPUREG_FIFO_EN            = 0x23
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
	MPUREG_I2C_SLV4_ADDR      = 0x31
	MPUREG_I2C_SLV4_REG       = 0x32
	MPUREG_I2C_SLV4_DO        = 0x33
	MPUREG_I2C_SLV4_CTRL      = 0x34
	MPUREG_I2C_SLV4_DI        = 0x35
	MPUREG_I2C_MST_STATUS     = 0x36
	MPUREG_INT_PIN_CFG        = 0x37
	MPUREG_INT_ENABLE         = 0x38
	MPUREG_ACCEL_XOUT_H       = 0x3B
	MPUREG_ACCEL_XOUT_L       = 0x3C
	MPUREG_ACCEL_YOUT_H       = 0x3D
	MPUREG_ACCEL_YOUT_L       = 0x3E
	MPUREG_ACCEL_ZOUT_H       = 0x3F
	MPUREG_ACCEL_ZOUT_L       = 0x40
	MPUREG_TEMP_OUT_H         = 0x41
	MPUREG_TEMP_OUT_L         = 0x42
	MPUREG_GYRO_XOUT_H        = 0x43
	MPUREG_GYRO_XOUT_L        = 0x44
	MPUREG_GYRO_YOUT_H        = 0x45
	MPUREG_GYRO_YOUT_L        = 0x46
	MPUREG_GYRO_ZOUT_H        = 0x47
	MPUREG_GYRO_ZOUT_L        = 0x48
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
	MPUREG_I2C_SLV0_DO        = 0x63
	MPUREG_I2C_SLV1_DO        = 0x64
	MPUREG_I2C_SLV2_DO        = 0x65
	MPUREG_I2C_SLV3_DO        = 0x66
	MPUREG_I2C_MST_DELAY_CTRL = 0x67
	MPUREG_SIGNAL_PATH_RESET  = 0x68
	MPUREG_MOT_DETECT_CTRL    = 0x69
	MPUREG_USER_CTRL          = 0x6A
	MPUREG_PWR_MGMT_1         = 0x6B
	MPUREG_PWR_MGMT_2         = 0x6C
	MPUREG_BANK_SEL           = 0x6D
	MPUREG_MEM_START_ADDR     = 0x6E
	MPUREG_MEM_R_W            = 0x6F
	MPUREG_DMP_CFG_1          = 0x70
	MPUREG_DMP_CFG_2          = 0x71
	MPUREG_FIFO_COUNTH        = 0x72
	MPUREG_FIFO_COUNTL        = 0x73
	MPUREG_FIFO_R_W           = 0x74
	MPUREG_WHOAMI             = 0x75
	MPUREG_XA_OFFSET_H        = 0x77
	MPUREG_XA_OFFSET_L        = 0x78
	MPUREG_YA_OFFSET_H        = 0x7A
	MPUREG_YA_OFFSET_L        = 0x7B
	MPUREG_ZA_OFFSET_H        = 0x7D
	MPUREG_ZA_OFFSET_L        = 0x7E
	/* ---- AK8963 Reg In MPU9250 ----------------------------------------------- */
	AK8963_I2C_ADDR        = 0x0C //0x18
	AK8963_Device_ID       = 0x48
	AK8963_MAX_SAMPLE_RATE = 0x64 // 100 Hz
	// Read-only Reg
	AK8963_WIA  = 0x00
	AK8963_INFO = 0x01
	AK8963_ST1  = 0x02
	AK8963_HXL  = 0x03
	AK8963_HXH  = 0x04
	AK8963_HYL  = 0x05
	AK8963_HYH  = 0x06
	AK8963_HZL  = 0x07
	AK8963_HZH  = 0x08
	AK8963_ST2  = 0x09
	// Write/Read Reg
	AK8963_CNTL1  = 0x0A
	AK8963_CNTL2  = 0x0B
	AK8963_ASTC   = 0x0C
	AK8963_TS1    = 0x0D
	AK8963_TS2    = 0x0E
	AK8963_I2CDIS = 0x0F
	// Read-only Reg ( ROM )
	AK8963_ASAX = 0x10
	AK8963_ASAY = 0x11
	AK8963_ASAZ = 0x12
	// Configuration bits mpu9250
	BIT_SLEEP                  = 0x40
	BIT_H_RESET                = 0x80
	BITS_CLKSEL                = 0x07
	MPU_CLK_SEL_PLLGYROX       = 0x01
	MPU_CLK_SEL_PLLGYROZ       = 0x03
	MPU_EXT_SYNC_GYROX         = 0x02
	BITS_FS_250DPS             = 0x00
	BITS_FS_500DPS             = 0x08
	BITS_FS_1000DPS            = 0x10
	BITS_FS_2000DPS            = 0x18
	BITS_FS_2G                 = 0x00
	BITS_FS_4G                 = 0x08
	BITS_FS_8G                 = 0x10
	BITS_FS_16G                = 0x18
	BITS_FS_MASK               = 0x18
	BITS_DLPF_CFG_256HZ_NOLPF2 = 0x00
	BITS_DLPF_CFG_188HZ        = 0x01
	BITS_DLPF_CFG_98HZ         = 0x02
	BITS_DLPF_CFG_42HZ         = 0x03
	BITS_DLPF_CFG_20HZ         = 0x04
	BITS_DLPF_CFG_10HZ         = 0x05
	BITS_DLPF_CFG_5HZ          = 0x06
	BITS_DLPF_CFG_2100HZ_NOLPF = 0x07
	BITS_DLPF_CFG_MASK         = 0x07
	BIT_INT_ANYRD_2CLEAR       = 0x10
	BIT_RAW_RDY_EN             = 0x01
	BIT_I2C_IF_DIS             = 0x10

	// Misc
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

	/* = ---- Sensitivity --------------------------------------------------------- */

	MPU9250M_4800uT                       = 0.6            // 0.6 uT/LSB
	MPU9250T_85degC                       = 0.002995177763 // 0.002995177763 degC/LSB
	Magnetometer_Sensitivity_Scale_Factor = 0.15

	// scaling/buffers
	bufSize  = 250 // Size of buffer storing instantaneous sensor values
	scaleMag = 9830.0 / 65536
)

// MPUData contains all the values measured by an MPU9250.
type MPUData struct {
	G1, G2, G3        float64
	A1, A2, A3        float64
	M1, M2, M3        float64
	Temp              float64
	GAError, MagError error
	N, NM             int
	T, TM             time.Time
	DT, DTM           time.Duration
}

// MPU9250Driver is adapted from https://github.com/westphae/goflying/blob/master/mpu9250/mpu9250.go
type MPU9250Driver struct {
	name       string
	connector  i2c.Connector
	connection i2c.Connection
	i2c.Config

	scaleGyro, scaleAccel float64         // Max sensor reading for value 2**15-1
	sampleRate            int             // Sample rate for sensor readings, Hz
	enableMag             bool            // Read the magnetometer?
	mcal1, mcal2, mcal3   float64         // Hardware magnetometer calibration values, uT
	a01, a02, a03         float64         // Hardware accelerometer calibration values, G
	g01, g02, g03         float64         // Hardware gyro calibration values, °/s
	C                     <-chan *MPUData // Current instantaneous sensor values
	CAvg                  <-chan *MPUData // Average sensor values (since CAvg last read)
	CBuf                  <-chan *MPUData // Buffer of instantaneous sensor values

	stop chan bool // Turn off MPU polling
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
	mpu.stop <- true

	return
}

func (mpu *MPU9250Driver) initialize() (err error) {
	bus := mpu.GetBusOrDefault(mpu.connector.GetDefaultBus())
	address := mpu.GetAddressOrDefault(MPU_ADDRESS)

	mpu.connection, err = mpu.connector.GetConnection(address, bus)
	if err != nil {
		return err
	}

	// hardcode these variables for now
	sensitivityGyro := 250
	sensitivityAccel := 4

	applyHWOffsets := false

	mpu.sampleRate = 25
	mpu.enableMag = false

	// Initialization of MPU
	// Reset device.
	if err := mpu.i2cWrite(MPUREG_PWR_MGMT_1, BIT_H_RESET); err != nil {
		return errors.New(fmt.Sprintf("Error resetting MPU9250: %s", err))
	}

	// Note: the following is in inv_mpu.c, but doesn't appear to be necessary from the MPU-9250 register map.
	// Wake up chip.
	time.Sleep(100 * time.Millisecond)
	if err := mpu.i2cWrite(MPUREG_PWR_MGMT_1, 0x00); err != nil {
		return errors.New(fmt.Sprintf("Error waking MPU9250: %s", err))
	}

	// Note: inv_mpu.c sets some registers here to allocate 1kB to the FIFO buffer and 3kB to the DMP.
	// It doesn't seem to be supported in the 1.6 version of the register map and we're not using FIFO anyway,
	// so we skip this.
	// Don't let FIFO overwrite DMP data
	if err := mpu.i2cWrite(MPUREG_ACCEL_CONFIG_2, BIT_FIFO_SIZE_1024|0x8); err != nil {
		return errors.New(fmt.Sprintf("Error setting up MPU9250: %s", err))
	}

	// Set Gyro and Accel sensitivities
	if err := mpu.SetGyroSensitivity(sensitivityGyro); err != nil {
		return errors.New(fmt.Sprintf("Error setting MPU9250 gyro sensitivity: %s", err))
	}

	if err := mpu.SetAccelSensitivity(sensitivityAccel); err != nil {
		return errors.New(fmt.Sprintf("Error setting MPU9250 accel sensitivity: %s", err))
	}

	sampRate := byte(1000/mpu.sampleRate - 1)
	// Default: Set Gyro LPF to half of sample rate
	if err := mpu.SetGyroLPF(sampRate >> 1); err != nil {
		return errors.New(fmt.Sprintf("Error setting MPU9250 Gyro LPF: %s", err))
	}

	// Default: Set Accel LPF to half of sample rate
	if err := mpu.SetAccelLPF(sampRate >> 1); err != nil {
		return errors.New(fmt.Sprintf("Error setting MPU9250 Accel LPF: %s", err))
	}

	// Set sample rate to chosen
	if err := mpu.SetSampleRate(sampRate); err != nil {
		return errors.New(fmt.Sprintf("Error setting MPU9250 Sample Rate: %s", err))
	}

	// Turn off FIFO buffer
	if err := mpu.i2cWrite(MPUREG_FIFO_EN, 0x00); err != nil {
		return errors.New(fmt.Sprintf("MPU9250 Error: couldn't disable FIFO: %s", err))
	}

	// Turn off interrupts
	if err := mpu.i2cWrite(MPUREG_INT_ENABLE, 0x00); err != nil {
		return errors.New(fmt.Sprintf("MPU9250 Error: couldn't disable interrupts: %s", err))
	}

	// Set up magnetometer
	if mpu.enableMag {
		if err := mpu.ReadMagCalibration(); err != nil {
			return errors.New(fmt.Sprintf("Error reading calibration from magnetometer: %s", err))
		}

		// Set up AK8963 master mode, master clock and ES bit
		if err := mpu.i2cWrite(MPUREG_I2C_MST_CTRL, 0x40); err != nil {
			return errors.New(fmt.Sprintf("Error setting up AK8963: %s", err))
		}
		// Slave 0 reads from AK8963
		if err := mpu.i2cWrite(MPUREG_I2C_SLV0_ADDR, BIT_I2C_READ|AK8963_I2C_ADDR); err != nil {
			return errors.New(fmt.Sprintf("Error setting up AK8963: %s", err))
		}
		// Compass reads start at this register
		if err := mpu.i2cWrite(MPUREG_I2C_SLV0_REG, AK8963_ST1); err != nil {
			return errors.New(fmt.Sprintf("Error setting up AK8963: %s", err))
		}
		// Enable 8-byte reads on slave 0
		if err := mpu.i2cWrite(MPUREG_I2C_SLV0_CTRL, BIT_SLAVE_EN|8); err != nil {
			return errors.New(fmt.Sprintf("Error setting up AK8963: %s", err))
		}
		// Slave 1 can change AK8963 measurement mode
		if err := mpu.i2cWrite(MPUREG_I2C_SLV1_ADDR, AK8963_I2C_ADDR); err != nil {
			return errors.New(fmt.Sprintf("Error setting up AK8963: %s", err))
		}
		if err := mpu.i2cWrite(MPUREG_I2C_SLV1_REG, AK8963_CNTL1); err != nil {
			return errors.New(fmt.Sprintf("Error setting up AK8963: %s", err))
		}
		// Enable 1-byte reads on slave 1
		if err := mpu.i2cWrite(MPUREG_I2C_SLV1_CTRL, BIT_SLAVE_EN|1); err != nil {
			return errors.New(fmt.Sprintf("Error setting up AK8963: %s", err))
		}
		// Set slave 1 data
		if err := mpu.i2cWrite(MPUREG_I2C_SLV1_DO, AKM_SINGLE_MEASUREMENT); err != nil {
			return errors.New(fmt.Sprintf("Error setting up AK8963: %s", err))
		}
		// Triggers slave 0 and 1 actions at each sample
		if err := mpu.i2cWrite(MPUREG_I2C_MST_DELAY_CTRL, 0x03); err != nil {
			return errors.New(fmt.Sprintf("Error setting up AK8963: %s", err))
		}

		// Set AK8963 sample rate to same as gyro/accel sample rate, up to max
		var ak8963Rate byte
		if mpu.sampleRate < AK8963_MAX_SAMPLE_RATE {
			ak8963Rate = 0
		} else {
			ak8963Rate = byte(mpu.sampleRate/AK8963_MAX_SAMPLE_RATE - 1)
		}

		// Not so sure of this one--I2C Slave 4??!
		if err := mpu.i2cWrite(MPUREG_I2C_SLV4_CTRL, ak8963Rate); err != nil {
			return errors.New(fmt.Sprintf("Error setting up AK8963: %s", err))
		}

		time.Sleep(100 * time.Millisecond) // Make sure mag is ready
	}

	// Set clock source to PLL
	if err := mpu.i2cWrite(MPUREG_PWR_MGMT_1, INV_CLK_PLL); err != nil {
		return errors.New(fmt.Sprintf("Error setting up MPU9250: %s", err))
	}
	// Turn off all sensors -- Not sure if necessary, but it's in the InvenSense DMP driver
	if err := mpu.i2cWrite(MPUREG_PWR_MGMT_2, 0x63); err != nil {
		return errors.New(fmt.Sprintf("Error setting up MPU9250: %s", err))
	}
	time.Sleep(100 * time.Millisecond)
	// Turn on all gyro, all accel
	if err := mpu.i2cWrite(MPUREG_PWR_MGMT_2, 0x00); err != nil {
		return errors.New(fmt.Sprintf("Error setting up MPU9250: %s", err))
	}

	if applyHWOffsets {
		if err := mpu.ReadAccelBias(sensitivityAccel); err != nil {
			return err
		}
		if err := mpu.ReadGyroBias(sensitivityGyro); err != nil {
			return err
		}
	}

	// Usually we don't want the automatic gyro bias compensation - it pollutes the gyro in a non-inertial frame.
	if err := mpu.EnableGyroBiasCal(false); err != nil {
		return err
	}

	go mpu.readSensors()

	// Give the IMU time to fully initialize and then clear out any bad values from the averages.
	//time.Sleep(500 * time.Millisecond) // Make sure it's ready
	//<-mpu.CAvg

	return nil
}

// readSensors polls the gyro, accelerometer and magnetometer sensors as well as the die temperature.
// Communication is via channels.
func (mpu *MPU9250Driver) readSensors() {
	var (
		g1, g2, g3, a1, a2, a3, m1, m2, m3, m4, tmp int16   // Current values
		avg1, avg2, avg3, ava1, ava2, ava3, avtmp   float64 // Accumulators for averages
		avm1, avm2, avm3                            int32
		n, nm                                       float64
		gaError, magError                           error
		t0, t, t0m, tm                              time.Time
		magSampleRate                               int
		curdata                                     *MPUData
	)

	acRegMap := map[*int16]byte{
		&g1: MPUREG_GYRO_XOUT_H, &g2: MPUREG_GYRO_YOUT_H, &g3: MPUREG_GYRO_ZOUT_H,
		&a1: MPUREG_ACCEL_XOUT_H, &a2: MPUREG_ACCEL_YOUT_H, &a3: MPUREG_ACCEL_ZOUT_H,
		&tmp: MPUREG_TEMP_OUT_H,
	}
	magRegMap := map[*int16]byte{
		&m1: MPUREG_EXT_SENS_DATA_00, &m2: MPUREG_EXT_SENS_DATA_02, &m3: MPUREG_EXT_SENS_DATA_04, &m4: MPUREG_EXT_SENS_DATA_06,
	}

	if mpu.sampleRate > 100 {
		magSampleRate = 100
	} else {
		magSampleRate = mpu.sampleRate
	}

	cC := make(chan *MPUData)
	defer close(cC)
	mpu.C = cC
	cAvg := make(chan *MPUData)
	defer close(cAvg)
	mpu.CAvg = cAvg
	cBuf := make(chan *MPUData, bufSize)
	defer close(cBuf)
	mpu.CBuf = cBuf
	mpu.stop = make(chan bool)
	defer close(mpu.stop)

	clock := time.NewTicker(time.Duration(int(1000.0/float32(mpu.sampleRate)+0.5)) * time.Millisecond)
	//TODO westphae: use the clock to record actual time instead of a timer
	defer clock.Stop()

	clockMag := time.NewTicker(time.Duration(int(1000.0/float32(magSampleRate)+0.5)) * time.Millisecond)
	t0 = time.Now()
	t0m = time.Now()

	makeMPUData := func() *MPUData {
		d := MPUData{
			G1:      (float64(g1) - mpu.g01) * mpu.scaleGyro,
			G2:      (float64(g2) - mpu.g02) * mpu.scaleGyro,
			G3:      (float64(g3) - mpu.g03) * mpu.scaleGyro,
			A1:      (float64(a1) - mpu.a01) * mpu.scaleAccel,
			A2:      (float64(a2) - mpu.a02) * mpu.scaleAccel,
			A3:      (float64(a3) - mpu.a03) * mpu.scaleAccel,
			M1:      float64(m1) * mpu.mcal1,
			M2:      float64(m2) * mpu.mcal2,
			M3:      float64(m3) * mpu.mcal3,
			Temp:    float64(tmp)/340 + 36.53,
			GAError: gaError, MagError: magError,
			N: 1, NM: 1,
			T: t, TM: tm,
			DT: time.Duration(0), DTM: time.Duration(0),
		}
		if gaError != nil {
			d.N = 0
		}
		if magError != nil {
			d.NM = 0
		}

		return &d
	}

	makeAvgMPUData := func() *MPUData {
		d := MPUData{}
		if n > 0.5 {
			d.G1 = (avg1/n - mpu.g01) * mpu.scaleGyro
			d.G2 = (avg2/n - mpu.g02) * mpu.scaleGyro
			d.G3 = (avg3/n - mpu.g03) * mpu.scaleGyro
			d.A1 = (ava1/n - mpu.a01) * mpu.scaleAccel
			d.A2 = (ava2/n - mpu.a02) * mpu.scaleAccel
			d.A3 = (ava3/n - mpu.a03) * mpu.scaleAccel
			d.Temp = (float64(avtmp)/n)/340 + 36.53
			d.N = int(n + 0.5)
			d.T = t
			d.DT = t.Sub(t0)
		} else {
			d.GAError = errors.New("MPU9250 Warning: No new accel/gyro values")
		}
		if nm > 0 {
			d.M1 = float64(avm1) * mpu.mcal1 / nm
			d.M2 = float64(avm2) * mpu.mcal2 / nm
			d.M3 = float64(avm3) * mpu.mcal3 / nm
			d.NM = int(nm + 0.5)
			d.TM = tm
			d.DTM = t.Sub(t0m)
		} else {
			d.MagError = errors.New("MPU9250 Warning: No new magnetometer values")
		}
		return &d
	}

	for {
		select {
		case t = <-clock.C: // Read accel/gyro data:
			for p, reg := range acRegMap {
				*p, gaError = mpu.i2cRead2(reg)
				if gaError != nil {
					log.Println("MPU9250 Warning: error reading gyro/accel")
				}
			}
			curdata = makeMPUData()
			// Update accumulated values and increment count of gyro/accel readings
			avg1 += float64(g1)
			avg2 += float64(g2)
			avg3 += float64(g3)
			ava1 += float64(a1)
			ava2 += float64(a2)
			ava3 += float64(a3)
			avtmp += float64(tmp)
			avm1 += int32(m1)
			avm2 += int32(m2)
			avm3 += int32(m3)
			n++
			select {
			case cBuf <- curdata: // We update the buffer every time we read a new value.
			default: // If buffer is full, remove oldest value and put in newest.
				<-cBuf
				cBuf <- curdata
			}
		case tm = <-clockMag.C: // Read magnetometer data:
			if mpu.enableMag {
				// Set AK8963 to slave0 for reading
				if err := mpu.i2cWrite(MPUREG_I2C_SLV0_ADDR, AK8963_I2C_ADDR|READ_FLAG); err != nil {
					log.Printf("MPU9250 Warning: couldn't set AK8963 address for reading: %s", err)
				}
				//I2C slave 0 register address from where to begin data transfer
				if err := mpu.i2cWrite(MPUREG_I2C_SLV0_REG, AK8963_HXL); err != nil {
					log.Printf("MPU9250 Warning: couldn't set AK8963 read register: %s", err)
				}
				//Tell AK8963 that we will read 7 bytes
				if err := mpu.i2cWrite(MPUREG_I2C_SLV0_CTRL, 0x87); err != nil {
					log.Printf("MPU9250 Warning: couldn't communicate with AK8963: %s", err)
				}

				// Read the actual data
				for p, reg := range magRegMap {
					*p, magError = mpu.i2cRead2(reg)
					if magError != nil {
						log.Println("MPU9250 Warning: error reading magnetometer")
					}
				}

				// Test validity of magnetometer data
				if (byte(m1&0xFF)&AKM_DATA_READY) == 0x00 && (byte(m1&0xFF)&AKM_DATA_OVERRUN) != 0x00 {
					log.Println("MPU9250 Warning: mag data not ready or overflow")
					log.Printf("MPU9250 Warning: m1 LSB: %X\n", byte(m1&0xFF))
					continue // Don't update the accumulated values
				}

				if (byte((m4>>8)&0xFF) & AKM_OVERFLOW) != 0x00 {
					log.Println("MPU9250 Warning: mag data overflow")
					log.Printf("MPU9250 Warning: m4 MSB: %X\n", byte((m1>>8)&0xFF))
					continue // Don't update the accumulated values
				}

				// Update values and increment count of magnetometer readings
				avm1 += int32(m1)
				avm2 += int32(m2)
				avm3 += int32(m3)
				nm++
			}
		case cC <- curdata: // Send the latest values
		case cAvg <- makeAvgMPUData(): // Send the averages
			avg1, avg2, avg3 = 0, 0, 0
			ava1, ava2, ava3 = 0, 0, 0
			avm1, avm2, avm3 = 0, 0, 0
			avtmp = 0
			n, nm = 0, 0
			t0, t0m = t, tm
		case <-mpu.stop: // Stop the goroutine, ease up on the CPU
			break
		}
	}
}

// SetSampleRate changes the sampling rate of the MPU.
func (mpu *MPU9250Driver) SetSampleRate(rate byte) (err error) {
	errWrite := mpu.i2cWrite(MPUREG_SMPLRT_DIV, byte(rate)) // Set sample rate to chosen
	if errWrite != nil {
		err = fmt.Errorf("MPU9250 Error: Couldn't set sample rate: %s", errWrite)
	}
	return
}

// SetGyroLPF sets the low pass filter for the gyro.
func (mpu *MPU9250Driver) SetGyroLPF(rate byte) (err error) {
	var r byte
	switch {
	case rate >= 188:
		r = BITS_DLPF_CFG_188HZ
	case rate >= 98:
		r = BITS_DLPF_CFG_98HZ
	case rate >= 42:
		r = BITS_DLPF_CFG_42HZ
	case rate >= 20:
		r = BITS_DLPF_CFG_20HZ
	case rate >= 10:
		r = BITS_DLPF_CFG_10HZ
	default:
		r = BITS_DLPF_CFG_5HZ
	}

	errWrite := mpu.i2cWrite(MPUREG_CONFIG, r)
	if errWrite != nil {
		err = fmt.Errorf("MPU9250 Error: couldn't set Gyro LPF: %s", errWrite)
	}
	return
}

// SetAccelLPF sets the low pass filter for the accelerometer.
func (mpu *MPU9250Driver) SetAccelLPF(rate byte) (err error) {
	var r byte
	switch {
	case rate >= 218:
		r = BITS_DLPF_CFG_188HZ
	case rate >= 99:
		r = BITS_DLPF_CFG_98HZ
	case rate >= 45:
		r = BITS_DLPF_CFG_42HZ
	case rate >= 21:
		r = BITS_DLPF_CFG_20HZ
	case rate >= 10:
		r = BITS_DLPF_CFG_10HZ
	default:
		r = BITS_DLPF_CFG_5HZ
	}

	errWrite := mpu.i2cWrite(MPUREG_ACCEL_CONFIG_2, r)
	if errWrite != nil {
		err = fmt.Errorf("MPU9250 Error: couldn't set Accel LPF: %s", errWrite)
	}
	return
}

// EnableGyroBiasCal enables or disables motion bias compensation for the gyro.
// For flying we generally do not want this!
func (mpu *MPU9250Driver) EnableGyroBiasCal(enable bool) error {
	enableRegs := []byte{0xb8, 0xaa, 0xb3, 0x8d, 0xb4, 0x98, 0x0d, 0x35, 0x5d}
	disableRegs := []byte{0xb8, 0xaa, 0xaa, 0xaa, 0xb0, 0x88, 0xc3, 0xc5, 0xc7}

	if enable {
		if err := mpu.memWrite(CFG_MOTION_BIAS, &enableRegs); err != nil {
			return errors.New("Unable to enable motion bias compensation")
		}
	} else {
		if err := mpu.memWrite(CFG_MOTION_BIAS, &disableRegs); err != nil {
			return errors.New("Unable to disable motion bias compensation")
		}
	}

	return nil
}

// SampleRate returns the current sample rate of the MPU9250, in Hz.
func (mpu *MPU9250Driver) SampleRate() int {
	return mpu.sampleRate
}

// MagEnabled returns whether or not the magnetometer is being read.
func (mpu *MPU9250Driver) MagEnabled() bool {
	return mpu.enableMag
}

// SetGyroSensitivity sets the gyro sensitivity of the MPU9250; it must be one of the following values:
// 250, 500, 1000, 2000 (all in °/s).
func (mpu *MPU9250Driver) SetGyroSensitivity(sensitivityGyro int) (err error) {
	var sensGyro byte

	switch sensitivityGyro {
	case 2000:
		sensGyro = BITS_FS_2000DPS
		mpu.scaleGyro = 2000.0 / float64(math.MaxInt16)
	case 1000:
		sensGyro = BITS_FS_1000DPS
		mpu.scaleGyro = 1000.0 / float64(math.MaxInt16)
	case 500:
		sensGyro = BITS_FS_500DPS
		mpu.scaleGyro = 500.0 / float64(math.MaxInt16)
	case 250:
		sensGyro = BITS_FS_250DPS
		mpu.scaleGyro = 250.0 / float64(math.MaxInt16)
	default:
		err = fmt.Errorf("MPU9250 Error: %d is not a valid gyro sensitivity", sensitivityGyro)
	}

	if errWrite := mpu.i2cWrite(MPUREG_GYRO_CONFIG, sensGyro); errWrite != nil {
		err = errors.New("MPU9250 Error: couldn't set gyro sensitivity")
	}

	return
}

// SetAccelSensitivity sets the accelerometer sensitivity of the MPU9250; it must be one of the following values:
// 2, 4, 8, 16, all in G (gravity).
func (mpu *MPU9250Driver) SetAccelSensitivity(sensitivityAccel int) (err error) {
	var sensAccel byte

	switch sensitivityAccel {
	case 16:
		sensAccel = BITS_FS_16G
		mpu.scaleAccel = 16.0 / float64(math.MaxInt16)
	case 8:
		sensAccel = BITS_FS_8G
		mpu.scaleAccel = 8.0 / float64(math.MaxInt16)
	case 4:
		sensAccel = BITS_FS_4G
		mpu.scaleAccel = 4.0 / float64(math.MaxInt16)
	case 2:
		sensAccel = BITS_FS_2G
		mpu.scaleAccel = 2.0 / float64(math.MaxInt16)
	default:
		err = fmt.Errorf("MPU9250 Error: %d is not a valid accel sensitivity", sensitivityAccel)
	}

	if errWrite := mpu.i2cWrite(MPUREG_ACCEL_CONFIG, sensAccel); errWrite != nil {
		err = errors.New("MPU9250 Error: couldn't set accel sensitivity")
	}

	return
}

// ReadAccelBias reads the bias accelerometer value stored on the chip.
// These values are set at the factory.
func (mpu *MPU9250Driver) ReadAccelBias(sensitivityAccel int) error {
	a0x, err := mpu.i2cRead2(MPUREG_XA_OFFSET_H)
	if err != nil {
		return errors.New("MPU9250 Error: ReadAccelBias error reading chip")
	}
	a0y, err := mpu.i2cRead2(MPUREG_YA_OFFSET_H)
	if err != nil {
		return errors.New("MPU9250 Error: ReadAccelBias error reading chip")
	}
	a0z, err := mpu.i2cRead2(MPUREG_ZA_OFFSET_H)
	if err != nil {
		return errors.New("MPU9250 Error: ReadAccelBias error reading chip")
	}

	switch sensitivityAccel {
	case 16:
		mpu.a01 = float64(a0x >> 1)
		mpu.a02 = float64(a0y >> 1)
		mpu.a03 = float64(a0z >> 1)
	case 8:
		mpu.a01 = float64(a0x)
		mpu.a02 = float64(a0y)
		mpu.a03 = float64(a0z)
	case 4:
		mpu.a01 = float64(a0x << 1)
		mpu.a02 = float64(a0y << 1)
		mpu.a03 = float64(a0z << 1)
	case 2:
		mpu.a01 = float64(a0x << 2)
		mpu.a02 = float64(a0y << 2)
		mpu.a03 = float64(a0z << 2)
	default:
		return fmt.Errorf("MPU9250 Error: %d is not a valid acceleration sensitivity", sensitivityAccel)
	}

	log.Printf("MPU9250 Info: accel hardware bias read: %6f %6f %6f\n", mpu.a01, mpu.a02, mpu.a03)
	return nil
}

// ReadGyroBias reads the bias gyro value stored on the chip.
// These values are set at the factory.
func (mpu *MPU9250Driver) ReadGyroBias(sensitivityGyro int) error {
	g0x, err := mpu.i2cRead2(MPUREG_XG_OFFS_USRH)
	if err != nil {
		return errors.New("MPU9250 Error: ReadGyroBias error reading chip")
	}
	g0y, err := mpu.i2cRead2(MPUREG_YG_OFFS_USRH)
	if err != nil {
		return errors.New("MPU9250 Error: ReadGyroBias error reading chip")
	}
	g0z, err := mpu.i2cRead2(MPUREG_ZG_OFFS_USRH)
	if err != nil {
		return errors.New("MPU9250 Error: ReadGyroBias error reading chip")
	}

	switch sensitivityGyro {
	case 2000:
		mpu.g01 = float64(g0x >> 1)
		mpu.g02 = float64(g0y >> 1)
		mpu.g03 = float64(g0z >> 1)
	case 1000:
		mpu.g01 = float64(g0x)
		mpu.g02 = float64(g0y)
		mpu.g03 = float64(g0z)
	case 500:
		mpu.g01 = float64(g0x << 1)
		mpu.g02 = float64(g0y << 1)
		mpu.g03 = float64(g0z << 1)
	case 250:
		mpu.g01 = float64(g0x << 2)
		mpu.g02 = float64(g0y << 2)
		mpu.g03 = float64(g0z << 2)
	default:
		return fmt.Errorf("MPU9250 Error: %d is not a valid gyro sensitivity", sensitivityGyro)
	}

	log.Printf("MPU9250 Info: Gyro hardware bias read: %6f %6f %6f\n", mpu.g01, mpu.g02, mpu.g03)
	return nil
}

// ReadMagCalibration reads the magnetometer bias values stored on the chpi.
// These values are set at the factory.
func (mpu *MPU9250Driver) ReadMagCalibration() error {
	// Enable bypass mode
	var tmp uint8
	var err error
	tmp, err = mpu.i2cRead(MPUREG_USER_CTRL)
	if err != nil {
		return errors.New("ReadMagCalibration error reading chip")
	}
	if err = mpu.i2cWrite(MPUREG_USER_CTRL, tmp & ^BIT_AUX_IF_EN); err != nil {
		return errors.New("ReadMagCalibration error reading chip")
	}
	time.Sleep(3 * time.Millisecond)
	if err = mpu.i2cWrite(MPUREG_INT_PIN_CFG, BIT_BYPASS_EN); err != nil {
		return errors.New("ReadMagCalibration error reading chip")
	}

	// Prepare for getting sensitivity data from AK8963
	//Set the I2C slave address of AK8963
	if err = mpu.i2cWrite(MPUREG_I2C_SLV0_ADDR, AK8963_I2C_ADDR); err != nil {
		return errors.New("ReadMagCalibration error reading chip")
	}
	// Power down the AK8963
	if err = mpu.i2cWrite(MPUREG_I2C_SLV0_CTRL, AK8963_CNTL1); err != nil {
		return errors.New("ReadMagCalibration error reading chip")
	}
	// Power down the AK8963
	if err = mpu.i2cWrite(MPUREG_I2C_SLV0_DO, AKM_POWER_DOWN); err != nil {
		return errors.New("ReadMagCalibration error reading chip")
	}
	time.Sleep(time.Millisecond)
	// Fuse AK8963 ROM access
	if mpu.i2cWrite(MPUREG_I2C_SLV0_DO, AK8963_I2CDIS); err != nil {
		return errors.New("ReadMagCalibration error reading chip")
	}
	time.Sleep(time.Millisecond)

	// Get sensitivity data from AK8963 fuse ROM
	mcal1, err := mpu.i2cRead(AK8963_ASAX)
	if err != nil {
		return errors.New("ReadMagCalibration error reading chip")
	}
	mcal2, err := mpu.i2cRead(AK8963_ASAY)
	if err != nil {
		return errors.New("ReadMagCalibration error reading chip")
	}
	mcal3, err := mpu.i2cRead(AK8963_ASAZ)
	if err != nil {
		return errors.New("ReadMagCalibration error reading chip")
	}

	log.Printf("MPU9250 Info: Raw mag calibrations: %d %d %d\n", mcal1, mcal2, mcal3)
	mpu.mcal1 = float64(int16(mcal1)+128) / 256 * scaleMag
	mpu.mcal2 = float64(int16(mcal2)+128) / 256 * scaleMag
	mpu.mcal3 = float64(int16(mcal3)+128) / 256 * scaleMag

	// Clean up from getting sensitivity data from AK8963
	// Fuse AK8963 ROM access
	if err = mpu.i2cWrite(MPUREG_I2C_SLV0_DO, AK8963_I2CDIS); err != nil {
		return errors.New("ReadMagCalibration error reading chip")
	}
	time.Sleep(time.Millisecond)

	// Disable bypass mode now that we're done getting sensitivity data
	tmp, err = mpu.i2cRead(MPUREG_USER_CTRL)
	if err != nil {
		return errors.New("ReadMagCalibration error reading chip")
	}
	if err = mpu.i2cWrite(MPUREG_USER_CTRL, tmp|BIT_AUX_IF_EN); err != nil {
		return errors.New("ReadMagCalibration error reading chip")
	}
	time.Sleep(3 * time.Millisecond)
	if err = mpu.i2cWrite(MPUREG_INT_PIN_CFG, 0x00); err != nil {
		return errors.New("ReadMagCalibration error reading chip")
	}
	time.Sleep(3 * time.Millisecond)

	log.Printf("MPU9250 Info: Mag hardware bias: %f %f %f\n", mpu.mcal1, mpu.mcal2, mpu.mcal3)
	return nil
}

func (mpu *MPU9250Driver) i2cWrite(register, value byte) (err error) {
	if errWrite := mpu.connection.WriteByteData(register, value); errWrite != nil {
		err = fmt.Errorf("MPU9250 Error writing %X to %X: %s\n",
			value, register, errWrite)
	} else {
		time.Sleep(time.Millisecond)
	}
	return
}

func (mpu *MPU9250Driver) i2cRead(register byte) (value uint8, err error) {
	value, errWrite := mpu.connection.ReadByteData(register)
	if errWrite != nil {
		err = fmt.Errorf("i2cRead error: %s", errWrite)
	}
	return
}

func (mpu *MPU9250Driver) i2cRead2(register byte) (value int16, err error) {

	v, errWrite := mpu.connection.ReadWordData(register)
	if errWrite != nil {
		err = fmt.Errorf("MPU9250 Error reading %x: %s\n", register, err)
	} else {
		value = int16(v)
	}
	return
}

func (mpu *MPU9250Driver) memWrite(addr uint16, data *[]byte) error {
	var err error
	var tmp = make([]byte, 2)

	tmp[0] = byte(addr >> 8)
	tmp[1] = byte(addr & 0xFF)

	// Check memory bank boundaries
	if tmp[1]+byte(len(*data)) > MPU_BANK_SIZE {
		return errors.New("Bad address: writing outside of memory bank boundaries")
	}

	err = mpu.connection.WriteBlockData(MPUREG_BANK_SEL, tmp)
	if err != nil {
		return fmt.Errorf("MPU9250 Error selecting memory bank: %s\n", err)
	}

	err = mpu.connection.WriteBlockData(MPUREG_MEM_R_W, *data)
	if err != nil {
		return fmt.Errorf("MPU9250 Error writing to the memory bank: %s\n", err)
	}

	return nil
}
