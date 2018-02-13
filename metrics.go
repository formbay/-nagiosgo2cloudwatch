package main

import (
	"log"
	"regexp"
	"strconv"
	"strings"
)

// StripUnits removes trailing units from a metric value
func StripUnits(value string) float64 {
	reg, err := regexp.Compile("[^0-9.]+")
	if err != nil {
		log.Fatal(err)
	}
	processedString := reg.ReplaceAllString(value, "")
	newFloat, _ := strconv.ParseFloat(processedString, 64)
	return newFloat
}

// ProcessOutput parses raw nagios check output into key/value pairs
func ProcessOutput(output string) map[string]float64 {
	outPut := make(map[string]float64)
	a := strings.Split(output, "|")
	b := a[len(a)-1]
	b = strings.TrimSpace(b)
	for _, token := range strings.Split(b, " ") {
		pair := strings.Split(token, ";")
		keypair := strings.Split(pair[0], "=")
		outPut[keypair[0]] = StripUnits(keypair[1])
	}
	return outPut
}
