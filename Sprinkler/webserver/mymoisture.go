package main

//package ads1x15_test

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"gobot.io/x/gobot/platforms/raspi"
	"periph.io/x/periph/conn/i2c/i2creg"
	"periph.io/x/periph/conn/physic"
	"periph.io/x/periph/experimental/devices/ads1x15"
	"periph.io/x/periph/host"
)

var pi *raspi.Adaptor

func main() {

	pi = raspi.NewAdaptor()

	http.HandleFunc("/sprinkler/on", sprinklerOn)

	http.HandleFunc("/sprinkler/off", sprinklerOff)

	static := http.FileServer(http.Dir("../webcontent"))
	http.Handle("/", static)

	log.Printf("About to listen on 8081 to http://localhost:8081/")
	log.Fatal(http.ListenAndServe(":8081", nil))

}

func getSensorReading(sprinklerNum string) float64 {
	var valueStr string
	var moistureLevel float64
	var channelNum ads1x15.Channel

	// Make sure periph is initialized.
	if _, err := host.Init(); err != nil {
		log.Fatal(err)
	}
	if sprinklerNum == "A" {
		channelNum = ads1x15.Channel0
	} else if sprinklerNum == "B" {
		channelNum = ads1x15.Channel1
	} else if sprinklerNum == "C" {
		channelNum = ads1x15.Channel2
	} else if sprinklerNum == "D" {
		channelNum = ads1x15.Channel3
	}
	// Open default I2C bus.
	bus, err := i2creg.Open("")
	if err != nil {
		log.Fatalf("failed to open I²C: %v", err)
	}
	defer bus.Close()

	// Create a new ADS1115 ADC.
	adc, err := ads1x15.NewADS1115(bus, &ads1x15.DefaultOpts)
	if err != nil {
		log.Fatalln(err)
	}

	// Obtain an analog pin from the ADC. PinForChannel(c Channel, maxVoltage physic.ElectricPotential, f physic.Frequency, q ConversionQuality)
	pin, err := adc.PinForChannel(channelNum, 5*physic.Volt, 1*physic.Hertz, ads1x15.SaveEnergy)
	fmt.Printf("value of  pin  :%v \n", pin)
	if err != nil {
		log.Fatalln(err)
	}
	defer pin.Halt()
	// Single Reading
	reading, err := pin.Read()
	if err != nil {
		log.Fatalln(err)
	}
	valueStr = reading.V.String()
	valueStr = valueStr[0 : len(valueStr)-2]
	moistureLevel, _ = strconv.ParseFloat(valueStr, 2)
	return moistureLevel
}

/*  // Read values continuously from ADC.
    fmt.Println("Continuous reading")
    c := pin.ReadContinuous()

    for reading := range c {
        valueStr = reading.V.String()
        valueStr = valueStr[0:len(valueStr)-2]
        currentReading,_ = strconv.ParseFloat(valueStr,2)
        mL = (currentReading - ReadingAtNoMoisture) / (ReadingAtFullMoisture - ReadingAtNoMoisture) * 100
        var mL float64
        moistureLevel = math.Trunc(mL)
        fmt.Println("Moisture Level : ",moistureLevel,"%")
        if enoughMoisture(moistureLevel) == false {
            pi.DigitalWrite("35", 0)
        }
        if enoughMoisture(moistureLevel) == true {
            pi.DigitalWrite("35", 1)
        }
    }
*/

func sprinklerOn(w http.ResponseWriter, r *http.Request) {
	var moistureLevel float64

	whichSprinkler := r.URL.Query()["which"]
	fmt.Fprintf(w, "Will attempt to turn on  %v", whichSprinkler)

	moistureLevel = getSensorReading(whichSprinkler[0])
	fmt.Println("Moisture Level : ", moistureLevel)
	if enoughMoisture(moistureLevel) == true {
		fmt.Println("Enough Moisture. Sprinkler will not Turn On : ", moistureLevel)
		return
	}

	switch whichSprinkler[0] {
	case "A":
		pi.DigitalWrite("31", 0)
	case "B":
		pi.DigitalWrite("35", 0)
	default:
		fmt.Fprintf(w, "Invalid name for sprinkler  %v", whichSprinkler)
		return
	}

	fmt.Fprintf(w, "Sprinkler %v is turned on", whichSprinkler)
}

func sprinklerOff(w http.ResponseWriter, r *http.Request) {

	var moistureLevel float64

	whichSprinkler := r.URL.Query()["which"]
	fmt.Fprintf(w, "Will attempt to turn off  %v", whichSprinkler)

	moistureLevel = getSensorReading(whichSprinkler[0])
	fmt.Println("Moisture Level : ", moistureLevel)
	if enoughMoisture(moistureLevel) == false {
		fmt.Println("Not Enough Moisture. Sprinkler will not Turn Off : ", moistureLevel)
		return
	}

	switch whichSprinkler[0] {
	case "A":
		pi.DigitalWrite("31", 1)
	case "B":
		pi.DigitalWrite("35", 1)
	default:
		fmt.Fprintf(w, "Invalid name for sprinkler  %v", whichSprinkler)
		return
	}

	fmt.Fprintf(w, "Sprinkler %v is turned off", whichSprinkler)
}

func enoughMoisture(moistureLevel float64) bool {
	var rValue bool
	if moistureLevel < 1.90 {
		rValue = true
	}
	if moistureLevel > 2.80 {
		rValue = false
	}
	return rValue
}
