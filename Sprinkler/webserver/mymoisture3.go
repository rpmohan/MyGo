package main

//package ads1x15_test

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"gobot.io/x/gobot/platforms/raspi"
	"periph.io/x/periph/conn/i2c/i2creg"
	"periph.io/x/periph/conn/physic"
	"periph.io/x/periph/experimental/devices/ads1x15"
	"periph.io/x/periph/host"
)

var pi *raspi.Adaptor
var sprinkerAEnabled bool
var maxSprinklerOnTime int = 10000

func main() {

	pi = raspi.NewAdaptor()

	http.HandleFunc("/sprinkler/on", sprinklerOn)

	http.HandleFunc("/sprinkler/off", sprinklerOff)
	http.HandleFunc("/sprinkler/readsensor", getSensorReadingLevelForWeb)
	go operateSprinklerWithMoisture("A")

	static := http.FileServer(http.Dir("../webcontent"))
	http.Handle("/", static)
	//    defer bus.Close()
	log.Printf("About to listen on 8081 to http://localhost:8081/")
	log.Fatal(http.ListenAndServe(":8081", nil))

}
func getSensorReadingLevelForWeb(w http.ResponseWriter, r *http.Request) {
	/* your code goes here*/
	whichSprinkler := r.URL.Query()["which"]
	fmt.Fprintf(w, "%v", getSensorReading(whichSprinkler[0]))

	return
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

func sprinklerOn(w http.ResponseWriter, r *http.Request) {
	var moistureLevel float64

	whichSprinkler := r.URL.Query()["which"]
	fmt.Fprintf(w, "Will attempt to turn on  %v", whichSprinkler)

	moistureLevel = getSensorReading(whichSprinkler[0])
	fmt.Println("Current Moisture Level : ", moistureLevel)
	if enoughMoisture(moistureLevel) == true {

		fmt.Fprintf(w, "Enough Moisture. Sprinkler will not Turn On : ", moistureLevel)

		fmt.Println("Enough Moisture. Sprinkler will not Turn On : ", moistureLevel)
		return
	}
	switchRelay(whichSprinkler[0], 0)
	if !(whichSprinkler[0] == "A" || whichSprinkler[0] == "B") {
		fmt.Fprintf(w, "Invalid name for sprinkler  %v", whichSprinkler[0])
	}

	sprinkerAEnabled = true

	fmt.Fprintf(w, "Sprinkler %v is turned on", whichSprinkler)
}

func switchRelay(relayName string, mode byte) {
	if mode == 0 {
		DurationOfTime := time.Duration(10) * time.Second
		Timer1 := time.AfterFunc(DurationOfTime, func() { switchRelay(relayName, 1) })
		defer Timer1.Stop()
	}
	switch relayName {
	case "A":
		pi.DigitalWrite("31", mode)
	case "B":
		pi.DigitalWrite("35", mode)
	default:
		return
	}

}

func sprinklerOff(w http.ResponseWriter, r *http.Request) {

	whichSprinkler := r.URL.Query()["which"]
	fmt.Fprintf(w, "Will attempt to turn off  %v", whichSprinkler)

	switchRelay(whichSprinkler[0], 1)
	fmt.Fprintf(w, "Sprinkler %v is turned off", whichSprinkler)
	sprinkerAEnabled = false
}

func operateSprinklerWithMoisture(relayName string) {
	var moistureLevel float64
	for {
		moistureLevel = getSensorReading(relayName)
		fmt.Printf("Realy Name: %s Moisture Level : %d \n ", relayName, moistureLevel)
		if sprinkerAEnabled == true {
			if enoughMoisture(moistureLevel) == false {
				fmt.Println("Not Enough Moisture. Turning ON Sprinkler : ", moistureLevel)
				switchRelay(relayName, 0)
			} else {
				fmt.Println("Enough Moisture. Turning OFF Sprinkler : ", moistureLevel)
				switchRelay(relayName, 1)

			}
		}
		time.Sleep(2 * time.Second)

	}

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
