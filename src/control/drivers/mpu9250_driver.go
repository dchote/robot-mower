package drivers

//
// MPU9250 driver written by Daniel Chote
//
// https://github.com/dchote/robot-mower/blob/master/src/control/drivers/mpu9250_driver.go
//

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
	MPUREG_INT_STATUS         = 0x3A
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

	MPU9250M_4800uT = 0.6            // 0.6 uT/LSB
	MPU9250T_85degC = 0.002995177763 // 0.002995177763 degC/LSB

	scaleMag = 9830.0 / 65536
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

	a01 float64
	a02 float64
	a03 float64

	g01 float64
	g02 float64
	g03 float64

	magXcoef float64
	magYcoef float64
	magZcoef float64

	Data *MPUData
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

func (mpu *MPU9250Driver) initialize() (err error) {
	bus := mpu.GetBusOrDefault(mpu.connector.GetDefaultBus())
	address := mpu.GetAddressOrDefault(MPU_ADDRESS)

	mpu.connection, err = mpu.connector.GetConnection(address, bus)
	if err != nil {
		return err
	}

	//mpu.connection.WriteByteData(register, value)
	//mpu.connection.WriteBlockData(register data)
	//mpu.connection.ReadByteData(register)
	//mpu.connection.ReadWordData(register) (int16)

	// Full reset
	mpu.connection.WriteByteData(MPUREG_PWR_MGMT_1, 0x80)
	time.Sleep(100 * time.Millisecond)

	mpu.connection.WriteByteData(MPUREG_PWR_MGMT_1, 0x01)
	mpu.connection.WriteByteData(MPUREG_PWR_MGMT_2, 0x00)

	time.Sleep(10 * time.Millisecond)

	mpu.connection.WriteByteData(MPUREG_CONFIG, 0x03)
	mpu.connection.WriteByteData(MPUREG_SMPLRT_DIV, 0x04)

	time.Sleep(10 * time.Millisecond)

	// configure interrupts and setup i2c master
	mpu.connection.WriteByteData(MPUREG_INT_PIN_CFG, 0x30)
	mpu.connection.WriteByteData(MPUREG_USER_CTRL, 0x20)
	mpu.connection.WriteByteData(MPUREG_I2C_MST_CTRL, 0x0D)

	time.Sleep(10 * time.Millisecond)

	// setup the i2c slave for the AK8963
	mpu.connection.WriteByteData(MPUREG_I2C_SLV0_ADDR, AK8963_I2C_ADDR)
	mpu.connection.WriteByteData(MPUREG_I2C_SLV0_REG, AK8963_CNTL2)
	mpu.connection.WriteByteData(MPUREG_I2C_SLV0_DO, 0x01)
	mpu.connection.WriteByteData(MPUREG_I2C_SLV0_CTRL, 0x81)

	time.Sleep(10 * time.Millisecond)

	mpu.connection.WriteByteData(MPUREG_I2C_SLV0_REG, AK8963_CNTL1)
	mpu.connection.WriteByteData(MPUREG_I2C_SLV0_DO, 0x12)
	mpu.connection.WriteByteData(MPUREG_I2C_SLV0_CTRL, 0x81)

	time.Sleep(10 * time.Millisecond)

	//
	// validation and calibration
	//
	mpuIdentity, err := mpu.connection.ReadByteData(MPUREG_WHOAMI)
	if err != nil {
		return errors.New("MPU9250Driver unable to fetch device identity")
	}
	log.Printf("MPU9250Driver mpuIdentity: %v should be %v", mpuIdentity, 0x71)

	// gyro config
	mpu.gResolution = 250.0 / float64(math.MaxInt16)
	mpu.connection.WriteByteData(MPUREG_GYRO_CONFIG, 0x00)
	time.Sleep(25 * time.Millisecond)

	gyroConf, err := mpu.connection.ReadByteData(MPUREG_GYRO_CONFIG)
	if err != nil {
		return errors.New("GYRO_CONFIG read error")
	}
	gyroConf = gyroConf &^ 0x03
	gyroConf = gyroConf &^ 0x18
	gyroConf = gyroConf | BITS_FS_250DPS<<3
	// 0x03      // Clear Fchoice bits [1:0]
	// 0x18      // Clear AFS bits [4:3]
	// GFS_250<<3 // Set full scale range for the gyro
	if err = mpu.connection.WriteByteData(MPUREG_GYRO_CONFIG, gyroConf); err != nil {
		return errors.New("GYRO_CONFIG write error")
	}

	// accel config
	mpu.aResolution = 2.0 / float64(math.MaxInt16)
	mpu.connection.WriteByteData(MPUREG_ACCEL_CONFIG, 0x00) // ares = 2.0/32768.0
	time.Sleep(25 * time.Millisecond)

	accelConf, err := mpu.connection.ReadByteData(MPUREG_ACCEL_CONFIG)
	if err != nil {
		return errors.New("ACCEL_CONFIG read error")
	}
	accelConf = accelConf &^ 0x18
	accelConf = accelConf | BITS_FS_2G<<3
	// 0x18     // Clear AFS bits [4:3]
	// AFS_2G<<3 // Set full scale range for the accelerometer
	if err = mpu.connection.WriteByteData(MPUREG_ACCEL_CONFIG, accelConf); err != nil {
		return errors.New("ACCEL_CONFIG write error")
	}

	/*
		mpu.connection.WriteByteData(ACCEL_CONFIG_2, 0x03)
	*/
	accelConf2, err := mpu.connection.ReadByteData(MPUREG_ACCEL_CONFIG_2)
	if err != nil {
		return errors.New("ACCEL_CONFIG_2 read error")
	}
	accelConf2 = accelConf2 &^ 0x0F
	accelConf2 = accelConf2 | 0x03
	// 0x0F // Clear accel_fchoice_b (bit 3) and A_DLPFG (bits [2:0])
	// 0x03  // Set accelerometer rate to 1 kHz and bandwidth to 41 Hz
	if err = mpu.connection.WriteByteData(MPUREG_ACCEL_CONFIG_2, accelConf2); err != nil {
		return errors.New("ACCEL_CONFIG_2 write error")
	}

	a0x, _ := mpu.i2cRead16(MPUREG_XA_OFFSET_H)
	a0y, _ := mpu.i2cRead16(MPUREG_YA_OFFSET_H)
	a0z, _ := mpu.i2cRead16(MPUREG_ZA_OFFSET_H)

	mpu.a01 = float64(a0x << 2)
	mpu.a02 = float64(a0y << 2)
	mpu.a03 = float64(a0z << 2)

	log.Printf("MPU9250Driver accel hardware bias read: %v, %v, %v", mpu.a01, mpu.a02, mpu.a03)

	g0x, _ := mpu.i2cRead16(MPUREG_XG_OFFS_USRH)
	g0y, _ := mpu.i2cRead16(MPUREG_YG_OFFS_USRH)
	g0z, _ := mpu.i2cRead16(MPUREG_ZG_OFFS_USRH)

	mpu.g01 = float64(g0x << 2)
	mpu.g02 = float64(g0y << 2)
	mpu.g03 = float64(g0z << 2)

	log.Printf("MPU9250Driver gyro hardware bias read: %v, %v, %v", mpu.g01, mpu.g02, mpu.g03)

	time.Sleep(50 * time.Millisecond)

	// AK8963 whoami
	mpu.connection.WriteByteData(MPUREG_I2C_SLV0_ADDR, AK8963_I2C_ADDR|READ_FLAG)
	mpu.connection.WriteByteData(MPUREG_I2C_SLV0_REG, AK8963_WIA)
	mpu.connection.WriteByteData(MPUREG_I2C_SLV0_CTRL, 0x81)

	time.Sleep(10 * time.Millisecond)
	magIdentity, err := mpu.connection.ReadByteData(MPUREG_EXT_SENS_DATA_00)
	if err != nil {
		return errors.New("MPU9250Driver unable to fetch mag identity")
	}

	log.Printf("MPU9250Driver magIdentity: %v should be %v", magIdentity, AK8963_Device_ID)

	if magIdentity != AK8963_Device_ID {
		return errors.New("MPU9250Driver invalid mag identity")
	}

	mpu.connection.WriteByteData(MPUREG_I2C_SLV0_DO, 0x00)
	time.Sleep(10 * time.Millisecond)
	mpu.connection.WriteByteData(MPUREG_I2C_SLV0_DO, 0x0F) // fuse rom access
	time.Sleep(10 * time.Millisecond)

	// read mag calibration data
	mpu.connection.WriteByteData(MPUREG_I2C_SLV0_ADDR, AK8963_I2C_ADDR|READ_FLAG)
	mpu.connection.WriteByteData(MPUREG_I2C_SLV0_REG, AK8963_ASAX)
	mpu.connection.WriteByteData(MPUREG_I2C_SLV0_CTRL, 0x83)

	time.Sleep(20 * time.Millisecond)
	/*
		mpu.connection.WriteByteData(AK8963_CNTL1, 0x00)
		time.Sleep(10 * time.Millisecond)
		mpu.connection.WriteByteData(AK8963_CNTL1, 0x0F) // fuse rom access
		time.Sleep(10 * time.Millisecond)
	*/

	// read the data
	if err = mpu.connection.WriteByte(MPUREG_EXT_SENS_DATA_00); err != nil {
		return errors.New("MPU9250Driver mag coef read error")
	}
	magBuf := []byte{0, 0, 0}
	_, _ = mpu.connection.Read(magBuf)

	mpu.magXcoef = (float64(int16(magBuf[0]))-128)/256.0 + 1.0
	mpu.magYcoef = (float64(int16(magBuf[1]))-128)/256.0 + 1.0
	mpu.magZcoef = (float64(int16(magBuf[2]))-128)/256.0 + 1.0

	log.Printf("MPU9250Driver mag hardware bias: %v, %v, %v coef: %v, %v, %v", magBuf[0], magBuf[1], magBuf[2], mpu.magXcoef, mpu.magYcoef, mpu.magZcoef)

	// AK8963 power down & cleanup
	mpu.connection.WriteByteData(MPUREG_I2C_SLV0_DO, 0x00)
	//mpu.connection.WriteByteData(AK8963_CNTL1, 0x00)
	time.Sleep(50 * time.Millisecond)

	// set scale&continuous mode
	mpu.connection.WriteByteData(MPUREG_I2C_SLV0_DO, (AK8963_BIT_16<<4 | AK8963_MODE_C8HZ))
	//mpu.connection.WriteByteData(AK8963_CNTL1, (AK8963_BIT_16<<4 | AK8963_MODE_C8HZ))
	mpu.mResolution = 4912.0 / float64(math.MaxInt16)

	time.Sleep(10 * time.Millisecond)
	// motion bias
	//enableRegs := []byte{0xb8, 0xaa, 0xb3, 0x8d, 0xb4, 0x98, 0x0d, 0x35, 0x5d}
	disableRegs := []byte{0xb8, 0xaa, 0xaa, 0xaa, 0xb0, 0x88, 0xc3, 0xc5, 0xc7}
	mpu.memWrite(CFG_MOTION_BIAS, &disableRegs)

	time.Sleep(200 * time.Millisecond)

	// attempt to fetch data for the first time
	go mpu.GetData()

	return nil
}

// readSensors polls the gyro, accelerometer and magnetometer sensors as well as the die temperature.
func (mpu *MPU9250Driver) GetData() (err error) {
	dataReady, _ := mpu.connection.ReadByteData(MPUREG_INT_STATUS)
	if dataReady != 0x01 {
		return errors.New("MPU9250Driver Error: Data not ready")
	}

	var (
		g1, g2, g3, a1, a2, a3, m1, m2, m3, m4 int16
	)

	g1, _ = mpu.i2cRead16(MPUREG_GYRO_XOUT_H)
	g2, _ = mpu.i2cRead16(MPUREG_GYRO_YOUT_H)
	g3, _ = mpu.i2cRead16(MPUREG_GYRO_ZOUT_H)

	a1, _ = mpu.i2cRead16(MPUREG_ACCEL_XOUT_H)
	a2, _ = mpu.i2cRead16(MPUREG_ACCEL_YOUT_H)
	a3, _ = mpu.i2cRead16(MPUREG_ACCEL_ZOUT_H)

	temp, _ := mpu.i2cRead16BigEnd(MPUREG_TEMP_OUT_H)

	mpu.Data = &MPUData{
		G1:   (float64(g1) - mpu.g01) * mpu.gResolution,
		G2:   (float64(g2) - mpu.g02) * mpu.gResolution,
		G3:   (float64(g3) - mpu.g03) * mpu.gResolution,
		A1:   (float64(a1) - mpu.a01) * mpu.aResolution,
		A2:   (float64(a2) - mpu.a02) * mpu.aResolution,
		A3:   (float64(a3) - mpu.a03) * mpu.aResolution,
		Temp: float64(temp)/333.87 + 21.0,
	}

	// request to read data from the AK8963
	mpu.connection.WriteByteData(MPUREG_I2C_SLV0_ADDR, AK8963_I2C_ADDR|READ_FLAG)
	mpu.connection.WriteByteData(MPUREG_I2C_SLV0_REG, AK8963_HXL)
	mpu.connection.WriteByteData(MPUREG_I2C_SLV0_CTRL, 0x87) // we want 7 bytes

	time.Sleep(10 * time.Millisecond)

	if err = mpu.connection.WriteByte(MPUREG_EXT_SENS_DATA_00); err != nil {
		return errors.New("MPU9250Driver mag coef read error")
	}

	magBuf := []byte{0, 0, 0, 0, 0, 0, 0}
	_, _ = mpu.connection.Read(magBuf)

	m1 = mpu.bufConvert(magBuf[0], magBuf[1])
	m2 = mpu.bufConvert(magBuf[2], magBuf[3])
	m3 = mpu.bufConvert(magBuf[4], magBuf[5])
	m4 = int16(magBuf[6])

	log.Printf("mag: %v, %v, %v", m1, m2, m3)
	// validate mag data
	if (byte(m1&0xFF)&AKM_DATA_READY) == 0x00 && (byte(m1&0xFF)&AKM_DATA_OVERRUN) != 0x00 {
		log.Println("MPU9250Driver Warning: mag data not ready or overflow")
	} else if (byte((m4>>8)&0xFF) & AKM_OVERFLOW) != 0x00 {
		log.Println("MPU9250Driver Warning: mag data overflow")
	} else {
		mpu.Data.M1 = (float64(m1) * mpu.magXcoef) * mpu.mResolution
		mpu.Data.M2 = (float64(m2) * mpu.magYcoef) * mpu.mResolution
		mpu.Data.M3 = (float64(m3) * mpu.magZcoef) * mpu.mResolution
	}

	return nil
}

func (mpu *MPU9250Driver) calibrate() {
	// reset and autoselect clock source
	mpu.connection.WriteByteData(MPUREG_PWR_MGMT_1, 0x80)
	time.Sleep(50 * time.Millisecond)

	//mpu.connection.WriteByteData(MPUREG_PWR_MGMT_1, 0x00)
	time.Sleep(10 * time.Millisecond)

	mpu.connection.WriteByteData(MPUREG_PWR_MGMT_1, 0x01)
	mpu.connection.WriteByteData(MPUREG_PWR_MGMT_2, 0x00)

	time.Sleep(10 * time.Millisecond)

	// configure device for bias calc
	mpu.connection.WriteByteData(MPUREG_INT_ENABLE, 0x00)
	mpu.connection.WriteByteData(MPUREG_FIFO_EN, 0x00)
	mpu.connection.WriteByteData(MPUREG_PWR_MGMT_1, 0x00)
	mpu.connection.WriteByteData(MPUREG_I2C_MST_CTRL, 0x00)
	mpu.connection.WriteByteData(MPUREG_USER_CTRL, 0x00)
	mpu.connection.WriteByteData(MPUREG_USER_CTRL, 0x0C)

	time.Sleep(10 * time.Millisecond)

	// configure for bias calc
	mpu.connection.WriteByteData(MPUREG_CONFIG, 0x01)
	mpu.connection.WriteByteData(MPUREG_SMPLRT_DIV, 0x00)
	mpu.connection.WriteByteData(MPUREG_GYRO_CONFIG, 0x00)
	mpu.connection.WriteByteData(MPUREG_ACCEL_CONFIG, 0x00)

	//gyrosensitivity := 131
	//accelsensitivity := 16384

	// configure FIFO for bias capture
	mpu.connection.WriteByteData(MPUREG_USER_CTRL, 0x40)
	mpu.connection.WriteByteData(MPUREG_FIFO_EN, 0x78)

	time.Sleep(40 * time.Millisecond)

	mpu.connection.WriteByteData(MPUREG_FIFO_EN, 0x00)
	mpu.connection.WriteByteData(MPUREG_USER_CTRL, 0x40)

	// TODO implement the rest...
}

func (mpu *MPU9250Driver) i2cRead16(reg uint8) (val int16, err error) {
	v, errRead := mpu.connection.ReadWordData(reg)
	if errRead != nil {
		return 0, errors.New("MPU9250Driver i2cRead16 error")
	} else {
		return int16(v), nil
	}
}

func (mpu *MPU9250Driver) i2cRead16BigEnd(reg uint8) (val int16, err error) {
	if err = mpu.connection.WriteByte(reg); err != nil {
		return 0, errors.New("MPU9250Driver readWordFromRegister error")
	}

	buf := []byte{0, 0}
	if _, err = mpu.connection.Read(buf); err != nil {
		return 0, errors.New("MPU9250Driver readWordFromRegister error")
	}

	return int16(buf[0])<<8 + int16(buf[1]), nil
}

func (mpu *MPU9250Driver) bufConvert(data1 uint8, data2 uint8) (val int16) {
	value := int32(data1) | (int32(data2) << 8)

	if value > 0x7FFF {
		value -= 0x10000
	}

	return int16(value)
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
		return errors.New("MPU9250 Error selecting memory bank")
	}

	err = mpu.connection.WriteBlockData(MPUREG_MEM_R_W, *data)
	if err != nil {
		return errors.New("MPU9250 Error writing to the memory bank")
	}

	return nil
}
