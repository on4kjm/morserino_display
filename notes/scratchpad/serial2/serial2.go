package main

import (
	"fmt"
	"log"
	"strings"

	"go.bug.st/serial"
	"go.bug.st/serial/enumerator"
)

func main() {
	ports, err := enumerator.GetDetailedPortsList()
	if err != nil {
		log.Fatal(err)
	}
	if len(ports) == 0 {
		fmt.Println("No serial ports found!")
		return
	}
	var morserinoPortName string
	for _, port := range ports {
		fmt.Printf("Found port: %s\n", port.Name)
		if port.IsUSB {
			fmt.Printf("   USB ID     %s:%s\n", port.VID, port.PID)
			fmt.Printf("   USB serial %s\n", port.SerialNumber)
			if (port.VID == "10C4") && (port.PID == "EA60") {
				morserinoPortName = port.Name
			}
		}
	}
	fmt.Printf("===================\n")

	mode := &serial.Mode{
		BaudRate: 115200,
		Parity:   serial.NoParity,
		DataBits: 8,
		StopBits: serial.OneStopBit,
	}

	if morserinoPortName == "" {
		fmt.Println("Did not find a usable port.")
		return
	}
	fmt.Println("Automatically selected Morserino port: " + morserinoPortName + "\n")

	myPort, err := serial.Open(morserinoPortName, mode)
	if err != nil {
		log.Fatal(err)
	}

	buff := make([]byte, 100)
	for {
		// Reads up to 100 bytes
		n, err := myPort.Read(buff)
		if err != nil {
			log.Fatal(err)
		}
		if n == 0 {
			fmt.Println("\nEOF")
			break
		}

		if strings.Contains(string(buff[:n]), "=") {
			fmt.Println("=")
		} else {
			fmt.Printf("%s", string(buff[:n]))
		}

	}

}
