package drivers

//
// INA219Driver adapted from https://github.com/NeuralSpaz/ti-ina219/blob/master/ina219.go
//
// Original author @NeuralSpaz did not have an attached license.
//

import (
	"gobot.io/x/gobot"
	"gobot.io/x/gobot/drivers/i2c"
)

const (
	INA219_ADDRESS01 = 0x40
	INA219_ADDRESS02 = 0x41

	INA219_CONFIG_RESET                     = 0x8000
	INA219_CONFIG_BVOLTAGERANGE_MASK        = 0x2000
	INA219_CONFIG_BVOLTAGERANGE_16V         = 0x0000
	INA219_CONFIG_BVOLTAGERANGE_32V         = 0x2000
	INA219_CONFIG_GAIN_MASK                 = 0x1800
	INA219_CONFIG_GAIN_1_40MV               = 0x0000
	INA219_CONFIG_GAIN_2_80MV               = 0x0800
	INA219_CONFIG_GAIN_4_160MV              = 0x1000
	INA219_CONFIG_GAIN_8_320MV              = 0x1800
	INA219_CONFIG_BADCRES_MASK              = 0x0780
	INA219_CONFIG_BADCRES_9BIT              = 0x0080
	INA219_CONFIG_BADCRES_10BIT             = 0x0100
	INA219_CONFIG_BADCRES_11BIT             = 0x0200
	INA219_CONFIG_BADCRES_12BIT             = 0x0400
	INA219_CONFIG_SADCRES_MASK              = 0x0078
	INA219_CONFIG_SADCRES_9BIT_1S_84US      = 0x0000
	INA219_CONFIG_SADCRES_10BIT_1S_148US    = 0x0008
	INA219_CONFIG_SADCRES_11BIT_1S_276US    = 0x0010
	INA219_CONFIG_SADCRES_12BIT_1S_532US    = 0x0018
	INA219_CONFIG_SADCRES_12BIT_2S_1060US   = 0x0048
	INA219_CONFIG_SADCRES_12BIT_4S_2130US   = 0x0050
	INA219_CONFIG_SADCRES_12BIT_8S_4260US   = 0x0058
	INA219_CONFIG_SADCRES_12BIT_16S_8510US  = 0x0060
	INA219_CONFIG_SADCRES_12BIT_32S_17MS    = 0x0068
	INA219_CONFIG_SADCRES_12BIT_64S_34MS    = 0x0070
	INA219_CONFIG_SADCRES_12BIT_128S_69MS   = 0x0078
	INA219_CONFIG_MODE_MASK                 = 0x0007
	INA219_CONFIG_MODE_POWERDOWN            = 0x0000
	INA219_CONFIG_MODE_SVOLT_TRIGGERED      = 0x0001
	INA219_CONFIG_MODE_BVOLT_TRIGGERED      = 0x0002
	INA219_CONFIG_MODE_SANDBVOLT_TRIGGERED  = 0x0003
	INA219_CONFIG_MODE_ADCOFF               = 0x0004
	INA219_CONFIG_MODE_SVOLT_CONTINUOUS     = 0x0005
	INA219_CONFIG_MODE_BVOLT_CONTINUOUS     = 0x0006
	INA219_CONFIG_MODE_SANDBVOLT_CONTINUOUS = 0x0007

	INA219_SHUNTRESISTOR_VALUE float64 = 0.1

	INA219_REG_CONFIG       = 0x00
	INA219_REG_SHUNTVOLTAGE = 0x01
	INA219_REG_BUSVOLTAGE   = 0x02
	INA219_REG_POWER        = 0x03
	INA219_REG_CURRENT      = 0x04
	INA219_REG_CALIBRATION  = 0x05
)

// INA219Driver is a driver for the INA219 current and bus voltage monitoring device.
type INA219Driver struct {
	name       string
	connector  i2c.Connector
	connection i2c.Connection
	i2c.Config
	CalibrationValue uint16
	halt             chan bool
}

// NewINA219Driver creates a new driver with the specified i2c interface.
// Params:
//		conn Connector - the Adaptor to use with this Driver
//
// Optional params:
//		i2c.WithBus(int):		bus to use with this driver
//		i2c.WithAddress(int):		address to use with this driver
func NewINA219Driver(c i2c.Connector, options ...func(i2c.Config)) *INA219Driver {
	i := &INA219Driver{
		name:             gobot.DefaultName("INA219"),
		connector:        c,
		Config:           i2c.NewConfig(),
		CalibrationValue: 4027,
	}

	for _, option := range options {
		option(i)
	}

	return i
}

// Name returns the name of the device.
func (i *INA219Driver) Name() string {
	return i.name
}

// SetName sets the name of the device.
func (i *INA219Driver) SetName(name string) {
	i.name = name
}

// Connection returns the connection of the device.
func (i *INA219Driver) Connection() gobot.Connection {
	return i.connector.(gobot.Connection)
}

// Start initializes the INA219
func (i *INA219Driver) Start() error {
	var err error
	bus := i.GetBusOrDefault(i.connector.GetDefaultBus())
	address := i.GetAddressOrDefault(int(INA219_ADDRESS01))

	if i.connection, err = i.connector.GetConnection(address, bus); err != nil {
		return err
	}

	if err := i.initialize(); err != nil {
		return err
	}

	i.connection.WriteBlockData(INA219_REG_CALIBRATION, []byte{byte(i.CalibrationValue >> 8), byte(i.CalibrationValue & 0x00FF)})

	return nil
}

// Halt halts the device.
func (i *INA219Driver) Halt() error {
	return nil
}

// GetBusVoltage gets the bus voltage in Volts
func (i *INA219Driver) GetBusVoltage() (float64, error) {
	value, err := i.getBusVoltageRaw()
	if err != nil {
		return 0, err
	}

	return float64((value>>3)*4) * 0.001, nil
}

// GetShuntVoltage Gets the shunt voltage in mV
func (i *INA219Driver) GetShuntVoltage() (float64, error) {
	value, err := i.getShuntVoltageRaw()
	if err != nil {
		return 0, err
	}

	return float64(value) * 0.00001, nil
}

// GetCurrent gets the current value in mA, taking into account the config settings and current LSB
func (i *INA219Driver) GetCurrent() (float64, error) {
	value, err := i.GetShuntVoltage()
	if err != nil {
		return 0, err
	}

	ma := value / INA219_SHUNTRESISTOR_VALUE
	return ma, nil
}

// GetLoadVoltage gets the load voltage in mV
func (i *INA219Driver) GetLoadVoltage() (float64, error) {
	bv, err := i.GetBusVoltage()
	if err != nil {
		return 0, err
	}

	sv, err := i.GetShuntVoltage()
	if err != nil {
		return 0, err
	}

	return bv + (sv / 1000.0), nil
}

// getBusVoltageRaw gets the raw bus voltage (16-bit signed integer, so +-32767)
func (i *INA219Driver) getBusVoltageRaw() (uint16, error) {
	val, err := i.readWordFromRegister(INA219_REG_BUSVOLTAGE)
	if err != nil {
		return 0, err
	}

	return uint16(val), nil
}

// getShuntVoltageRaw gets the raw shunt voltage (16-bit signed integer, so +-32767)
func (i *INA219Driver) getShuntVoltageRaw() (int16, error) {
	val, err := i.readWordFromRegister(INA219_REG_SHUNTVOLTAGE)
	if err != nil {
		return 0, err
	}

	value := int32(val)
	if value > 0x7FFF {
		value -= 0x10000
	}

	return int16(value), nil
}

// getCurrentRaw
func (i *INA219Driver) getCurrentRaw() (int16, error) {
	val, err := i.readWordFromRegister(INA219_REG_CURRENT)
	if err != nil {
		return 0, err
	}

	value := int32(val)
	if value > 0x7FFF {
		value -= 0x10000
	}

	return int16(value), nil
}

// getPowerRaw
func (i *INA219Driver) getPowerRaw() (int16, error) {
	val, err := i.readWordFromRegister(INA219_REG_POWER)
	if err != nil {
		return 0, err
	}

	value := int32(val)
	if value > 0x7FFF {
		value -= 0x10000
	}

	return int16(value), nil
}

// reads word from supplied register address
func (i *INA219Driver) readWordFromRegister(reg uint8) (uint16, error) {
	val, err := i.connection.ReadWordData(reg)
	if err != nil {
		return 0, err
	}

	return uint16(((val & 0x00FF) << 8) | ((val & 0xFF00) >> 8)), nil
}

// initialize initializes the INA219 device
func (i *INA219Driver) initialize() error {
	config := INA219_CONFIG_BVOLTAGERANGE_32V |
		INA219_CONFIG_GAIN_8_320MV |
		INA219_CONFIG_BADCRES_12BIT |
		INA219_CONFIG_SADCRES_12BIT_1S_532US |
		INA219_CONFIG_MODE_SANDBVOLT_CONTINUOUS

	return i.connection.WriteBlockData(INA219_REG_CONFIG, []byte{byte(config >> 8), byte(config & 0x00FF)})
}
