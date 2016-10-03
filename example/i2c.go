// i2c is example usage of BME280 library with I2C bus.
package main

import (
	"fmt"

	"golang.org/x/exp/io/i2c"

	"github.com/quhar/bme280"
)

func main() {

	d, err := i2c.Open(&i2c.Devfs{Dev: "/dev/i2c-1"}, bme280.I2CAddr)
	if err != nil {
		panic(err)
	}

	b := bme280.New(d)
	err = b.Init()

	t, p, h, err := b.EnvData()

	if err != nil {
		panic(err)
	}

	fmt.Printf("Temp: %fC, Press: %fhPa, Hum: %f%%\n", t, p, h)
}
