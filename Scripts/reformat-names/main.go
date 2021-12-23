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
	yearToRelease, err := plistDict.GetDict("year_to_release")
	if err != nil {
		channel <- err
		return
	}

	supportedVersions := make(map[string]SupportedVersions)
	for yearIndex, year := range yearToRelease.Keys {
		releases := yearToRelease.Dicts[yearIndex]
		iOSRelease, err := releases.GetString("iOS")
		if err != nil {
			channel <- err
			return
		}

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
	symbols, err := plistDict.GetDict("symbols")
	if err != nil {
		channel <- err
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

func (plistDict PlistDict) GetDict(key string) (PlistDict, error) {
	if len(plistDict.Keys) != len(plistDict.Dicts) {
		return PlistDict{}, errors.New("dict items not formatted correctly")
	}

	for dictKeyIndex, dictKey := range plistDict.Keys {
		if dictKey == key {
			return plistDict.Dicts[dictKeyIndex], nil
		}
	}

	return PlistDict{}, fmt.Errorf("%s key not found in plist", key)
}

func (plistDict PlistDict) GetString(key string) (string, error) {
	if len(plistDict.Keys) != len(plistDict.Strings) {
		return "", errors.New("string items not formatted correctly")
	}

	for dictKeyIndex, dictKey := range plistDict.Keys {
		if dictKey == key {
			return plistDict.Strings[dictKeyIndex], nil
		}
	}

	return "", fmt.Errorf("%s key not found in plist", key)
}
