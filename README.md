# bme280
[![GoDoc](https://godoc.org/github.com/quhar/bme280?status.svg)](https://godoc.org/github.com/quhar/bme280)

Golang library to read data from Bosch BME280 sensor. It requires library to talk to the device via following interface:

```go
type bus interface {
	ReadReg(byte, []byte) error
	WriteReg(byte, []byte) error
}
```

One of the libraries which matches this interface is experimental standard I/O [lib](https://godoc.org/golang.org/x/exp/io/i2c).

## Example usage

```go
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
```
