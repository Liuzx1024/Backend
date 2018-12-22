package generalgpio

import (
	"errors"
	. "github.com/chfanghr/Backend/hardware"
	"sync"
)

type GeneralGPIOHub struct {
	c Controller
	m sync.Locker
}

func NewGeneralGPIOHub(c Controller) *GeneralGPIOHub {
	return &GeneralGPIOHub{
		c: c,
		m: new(sync.Mutex),
	}
}
func (e *GeneralGPIOHub) withMutex(job func() (interface{}, error)) (res interface{}, err error) {
	e.m.Lock()
	defer e.m.Unlock()
	res, err = job()
	return
}
func (e *GeneralGPIOHub) doCommand(m CommandMethod, p PinNumber, v PinValue) (PinValue, error) {
	res, err := e.withMutex(func() (interface{}, error) {
		req := NewGeneralGpioRequester(m, p, v)
		err := req.Commit(e.c)
		if err != nil {
			return 0, err
		}
		return req.GetRes(), nil
	})
	if err != nil {
		return 0, err
	}
	return res.(PinValue), err
}
func (e *GeneralGPIOHub) PinMode(pin PinNumber, mode PinValue) error {
	switch mode {
	case GpioInput:
		break
	case GpioOutput:
		break
	case GpioInputPullUp:
		break
	case GpioInputPullDown:
		break
	default:
		return errors.New("invalid mode")
	}
	err := e.c.IsValidPin(pin)
	if err != nil {
		return err
	}
	_, err = e.doCommand(GpioPinmode, pin, mode)
	return err
}
func (e *GeneralGPIOHub) DigitalWrite(pin PinNumber, value PinValue) error {
	switch value {
	case GpioHigh:
		break
	case GpioLow:
		break
	default:
		return errors.New("invalid value")
	}
	_, err := e.doCommand(GpioDigitalWrite, pin, value)
	return err
}
func (e *GeneralGPIOHub) DigitalRead(pin PinNumber) (PinValue, error) {
	err := e.c.IsValidPin(pin)
	if err != nil {
		return 0, err
	}
	v, err := e.doCommand(GpioDigitalRead, pin, 0)
	if err != nil {
		return 0, err
	}
	switch v {
	case GpioLow:
		return v, nil
	case GpioHigh:
		return v, nil
	default:
		return 0, errors.New("unknown value")
	}
	return 0, errors.New("unknown error")
}
func (e *GeneralGPIOHub) AnalogWrite(pin PinNumber, value PinValue) error {
	err := e.c.IsValidPin(pin)
	if err != nil {
		return err
	}
	_, err = e.doCommand(GpioAnalogWrite, pin, value)
	return err
}
func (e *GeneralGPIOHub) AnalogRead(pin PinNumber) (PinValue, error) {
	err := e.c.IsValidPin(pin)
	if err != nil {
		return 0, err
	}
	return e.doCommand(GpioAnalogRead, pin, 0)
}

type GeneralGPIO struct {
	h GPIOHub
	p PinNumber
}

func NewGeneralGPIO(hub GPIOHub, pin PinNumber) *GeneralGPIO {
	return &GeneralGPIO{
		h: hub,
		p: pin,
	}
}
func (g GeneralGPIO) PinMode(m PinValue) error {
	return g.h.PinMode(g.p, m)
}
func (g GeneralGPIO) DigitalWrite(v PinValue) error {
	return g.h.DigitalWrite(g.p, v)
}
func (g GeneralGPIO) DigitalRead() (PinValue, error) {
	return g.h.DigitalRead(g.p)
}
func (g GeneralGPIO) AnalogRead() (PinValue, error) {
	return g.h.AnalogRead(g.p)
}
func (g GeneralGPIO) AnalogWrite(v PinValue) error {
	return g.h.AnalogWrite(g.p, v)
}
