package main

import (
	"fmt"
	"log"

	"github.com/alecthomas/kingpin"
	mhz16 "github.com/kefi550/mh-z16-go"
)

func main() {
	portName := kingpin.Arg("port", "portName").Required().ExistingFile()
	zero := kingpin.Flag("zero", "zero calibration mode").Bool()

	kingpin.Parse()

	m, err := mhz16.Open(*portName)
	if err != nil {
		log.Fatalln(err)
	}
	defer m.Close()

	if *zero {
		m.ZeroCalibration()
		return
	}
	co2, err := m.GetCo2()
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(co2)
}
