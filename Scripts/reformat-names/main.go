package main

import (
	"bufio"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io/ioutil"
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

	formattedNames := []FormattedName{}
	symbols := plistDict.GetDict("symbols")
	for symbolIndex, symbol := range symbols.Keys {
		releaseYear := symbols.Strings[symbolIndex]
		formattedName := FormattedName{
			Name:        symbol,
			ReleaseYear: releaseYear,
		}
		formattedNames = append(formattedNames, formattedName)
	}

	formattedNamesBytes, err := json.MarshalIndent(formattedNames, "", "  ")
	if err != nil {
		log.Fatalln(err)
	}

	err = ioutil.WriteFile("../../Shared/Resources/Names/names.json", formattedNamesBytes, 0644)
	if err != nil {
		log.Fatalln(err)
	}

	elapsed := time.Since(start)
	fmt.Printf("done reformatting names in %s\n", elapsed)
}

type FormattedName struct {
	Name        string `json:"name"`
	ReleaseYear string `json:"release_year"`
}

type SupportedVersions struct {
	IOS string `json:"iOS"`
}

type PlistDict struct {
	XMLName xml.Name    `xml:"dict"`
	Keys    []string    `xml:"key"`
	Dicts   []PlistDict `xml:"dict"`
	Strings []string    `xml:"string"`
}

func (plistDict PlistDict) GetDict(key string) *PlistDict {
	for dictKeyIndex, dictKey := range plistDict.Keys {
		if dictKey == key {
			return &plistDict.Dicts[dictKeyIndex]
		}
	}
	return nil
}
