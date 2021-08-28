package main

import (
	"fmt"
	//"html "
	"log"
	"net/http"

	"gobot.io/x/gobot/platforms/raspi"
)

var pi *raspi.Adaptor

func main() {

	pi = raspi.NewAdaptor()

	http.HandleFunc("/sprinkler/on", sprinklerOn)

	http.HandleFunc("/sprinkler/off", sprinklerOff)

	http.HandleFunc("/GetMoistureLevel", GetOnMoistureLevel)

	static := http.FileServer(http.Dir("../webcontent"))
	http.Handle("/", static)
	go readSensor("A")
	go readSensor("B")


	log.Printf("About to listen on 8081 to http://localhost:8081/")
	log.Fatal(http.ListenAndServe(":8081", nil))

}
func readSensor(sensor string) {
	reading := 0

	for {
		time.Sleep(3 * time.Second)

		switch sensor {
			case "A":
				reading = pi.DigitalRead("35")
				TurnSprinkler(sensor,reading,"31")
			case "B":
				reading = pi.DigitalRead("37")
				TurnSprinkler(sensor,reading,"33")
			default:
				fmt.Fprintf(w, "Invalid name for Sensor  %v", sensor)
				return		reading = pi.DigitalRead("35")
		
		}

		fmt.Printf("Reading at sensor %v: %v\n", sensor, reading)

	}
}
func TurnSprinkler(ensor string,reading int,pin string) {
	if reading > 0 {
		fmt.Printf("Sprinkler %v will be turned OFF\n", sensor)
		pi.DigitalWrite(pin, 1)
	}
	else {
		fmt.Printf("Sprinkler %v will be turned ON\n", sensor)
		pi.DigitalWrite(pin, 0)
	}
}

func sprinklerOn(w http.ResponseWriter, r *http.Request) {
	whichSprinkler := r.URL.Query()["which"]
	fmt.Fprintf(w, "Will attempt to turn on  %v", whichSprinkler)

	switch whichSprinkler[0] {
	case "A":
		pi.DigitalWrite("31", 0)
	case "B":
		pi.DigitalWrite("33", 0)
	case "C":
		pi.DigitalWrite("35", 0)
	case "D":
		pi.DigitalWrite("37", 0)
	default:
		fmt.Fprintf(w, "Invalid name for sprinkler  %v", whichSprinkler)
		return
	}

	fmt.Fprintf(w, "Sprinkler %v is turned on", whichSprinkler)
}

func sprinklerOff(w http.ResponseWriter, r *http.Request) {
	whichSprinkler := r.URL.Query()["which"]
	fmt.Fprintf(w, "Will attempt to turn off  %v", whichSprinkler)

	switch whichSprinkler[0] {
	case "A":
		pi.DigitalWrite("31", 1)
	case "B":
		pi.DigitalWrite("33", 1)
	case "C":
		pi.DigitalWrite("35", 1)
	case "D":
		pi.DigitalWrite("37", 1)
	default:
		fmt.Fprintf(w, "Invalid name for sprinkler  %v", whichSprinkler)
		return
	}
}
func GetMoistureReading(w http.ResponseWriter, r *http.Request) {
	whichMoisture := r.URL.Query()["which"]
	fmt.Fprintf(w, "Will attempt to find moisture level  of  %v", whichMoisture)
	var moiread int = 0

	switch whichMoisture[0] {
	case "A":
		moiread = pi.DigitalRead("31", 0)
	case "B":
		moiread = pi.DigitalWrite("33", 0)
	default:
		fmt.Fprintf(w, "Invalid name for sprinkler  %v", whichSprinkler)
		return
	}
	if (moiread > 50)
		fmt.Fprintf(w, "Turn ON moisture of  %v ", whichSprinkler)
	else
		fmt.Fprintf(w, "Turn OFF moisture of  %v ", whichSprinkler)

}
	