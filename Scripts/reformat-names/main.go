package main

import (
	"bufio"
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"time"
)

const SAVE_LOCATION = "../../Shared/Resources/Names"

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

	formattedNamesChannel := make(chan error)
	go makeFormattedNames(plistDict, formattedNamesChannel)

	supportedVersionsChannel := make(chan error)
	go makeSupportedVersions(plistDict, supportedVersionsChannel)

	channels := []chan error{formattedNamesChannel, supportedVersionsChannel}
	for _, channel := range channels {
		err = <-channel
		if err != nil {
			log.Fatalln(err)
		}
	}

	elapsed := time.Since(start)
	fmt.Printf("done reformatting names in %s\n", elapsed)
}

func makeSupportedVersions(plistDict PlistDict, channel chan error) {
	yearToRelease := plistDict.GetDict("year_to_release")
	if yearToRelease == nil {
		channel <- errors.New("year to release not found in plist")
		return
	}

	supportedVersions := make(map[string]SupportedVersions)
	for yearIndex, year := range yearToRelease.Keys {
		releases := yearToRelease.Dicts[yearIndex]
		iOSRelease := releases.Strings[0]
		supportedVersions[year] = SupportedVersions{IOS: iOSRelease}
	}

	supportedVersionsBytes, err := json.MarshalIndent(supportedVersions, "", "  ")
	if err != nil {
		channel <- err
		return
	}

	err = ioutil.WriteFile(fmt.Sprintf("%s/supported_versions.json", SAVE_LOCATION), supportedVersionsBytes, 0644)
	channel <- err
}

func makeFormattedNames(plistDict PlistDict, channel chan error) {
	formattedNames := []FormattedName{}
	symbols := plistDict.GetDict("symbols")
	if symbols == nil {
		channel <- errors.New("symbols not found in plist")
		return
	}

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
		channel <- err
		return
	}

	err = ioutil.WriteFile(fmt.Sprintf("%s/names.json", SAVE_LOCATION), formattedNamesBytes, 0644)
	channel <- err
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
