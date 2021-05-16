package morserino_com

/*
Copyright © 2021 Jean-Marc Meessen, ON4KJM <on4kjm@gmail.com>

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/

import (
	"fmt"
	"go.bug.st/serial/enumerator"
	"log"
)

//structure to store port details
type comPortItem struct {
	//Full name of the Com Port
	portName string
	//HEX USB vendor identification (if available)
	usbVendorID string
	//HEX USB product identification (if available)
	usbProductID string
	//port serial number
	serialNumber string
	//is a Morserino device?
	isMorserinoPort bool
}

//Structure to store all the detected ports
type comPortList struct {
	// number of detected ports
	nbrOfPorts int
	//number of Morserino ports found
	nbrOfMorserinoPorts int
	// Name of the detected morserino port. Empty if none was found
	morserinoPortName string
	// Array of port items
	portList []comPortItem
}

// Gets the list of COM devices and displays them on the console
func List_com() {
	comList := Get_com_list()
	buffer := prettyPrint_comList(comList)
	for _, line := range buffer {
		fmt.Println(line)
	}
}

//Gets all the ports on the system , checks whether it is a moreserino,
// and returns an array of port description
func Get_com_list() comPortList {

	var workComPortList comPortList
	workComPortList.nbrOfPorts = 0

	ports, err := enumerator.GetDetailedPortsList()
	if err != nil {
		log.Fatal(err)
	}
	if len(ports) == 0 {
		workComPortList.nbrOfPorts = 0
		return workComPortList
	}

	for _, port := range ports {
		var wrkPortItem comPortItem
		wrkPortItem.portName = port.Name
		if port.IsUSB {
			wrkPortItem.usbVendorID = port.VID
			wrkPortItem.usbProductID = port.PID
			wrkPortItem.serialNumber = port.SerialNumber
			//FIXME: Get the VID and PID of newer models
			if (port.VID == "10C4") && (port.PID == "EA60") {
				wrkPortItem.isMorserinoPort = true
			}
		}
		workComPortList.nbrOfPorts++
		workComPortList.portList = append(workComPortList.portList, wrkPortItem)
		if wrkPortItem.isMorserinoPort {
			workComPortList.nbrOfMorserinoPorts++
			workComPortList.morserinoPortName = wrkPortItem.portName
		}
	}
	return workComPortList
}

func prettyPrint_comList(portList comPortList) []string {
	var buffer []string

	if portList.nbrOfPorts == 0 {
		buffer = append(buffer, "No ports found !")
	}
	if portList.nbrOfMorserinoPorts > 1 {
		buffer = append(buffer, "WARNING: Multiple multiple Morserino devices detected")
	}
	buffer = append(buffer, "")

	for _, portItem := range portList.portList {
		buffer = append(buffer, format_com_item(portItem))
	}

	return buffer
}

// Generates a printable string with the details of a comPort item
func format_com_item(item comPortItem) string {
	var morserinoFlag string
	if item.isMorserinoPort {
		morserinoFlag = "=> "
	} else {
		morserinoFlag = "   "
	}
	buffer := fmt.Sprintf("%s %-30s (USB ID: %4s:%4s, USB Serial: %s)", morserinoFlag, item.portName, item.usbVendorID, item.usbProductID, item.serialNumber)
	return buffer
}
