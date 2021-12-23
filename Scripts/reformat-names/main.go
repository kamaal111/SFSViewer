package main

import (
	"bufio"
	"encoding/xml"
	"fmt"
	"log"
	"os"
	"strings"
	"time"
)

func main() {
	start := time.Now()

	nameAvailabilityFile, err := os.Open("name_availability.plist")
	if err != nil {
		log.Fatalln(err)
	}
	defer nameAvailabilityFile.Close()

	scanner := bufio.NewScanner(nameAvailabilityFile)
	scanner.Split(bufio.ScanLines)
	var nameAvailabilityArray []string

	for scanner.Scan() {
		nameAvailabilityArray = append(nameAvailabilityArray, scanner.Text())
	}

	nameAvailabilityArray = nameAvailabilityArray[3 : len(nameAvailabilityArray)-1]
	nameAvailabilityString := strings.Join(nameAvailabilityArray, "\n")
	nameAvailabilityBytes := []byte(nameAvailabilityString)

	var plistDict PlistDict
	err = xml.Unmarshal(nameAvailabilityBytes, &plistDict)
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println(plistDict)

	elapsed := time.Since(start)
	fmt.Printf("done reformatting names in %s\n", elapsed)
}

type PlistDict struct {
	XMLName xml.Name `xml:"dict"`
	Keys    []string `xml:"key"`
}
