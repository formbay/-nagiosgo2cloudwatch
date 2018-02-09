package main

import (
	"strings"
)

// CW holds all cloudwatch metric data for sending
type CW struct {
	Namespace  string
	BaseName   string
	Dimensions []map[string]string
	Data       []CWData
}

// CWData holds a single metric key/value pair
type CWData struct {
	MetricName string
	Value      float64
	Dimensions []map[string]string
}

// NewCW constructor for CW objects
func NewCW(namespace string, basename string, dimensions string) *CW {
	var dimMaps []map[string]string
	for _, pair := range strings.Split(dimensions, ",") {
		keyVal := strings.Split(pair, "=")
		if len(keyVal) == 1 {
			continue
		}
		dimMaps = append(dimMaps, map[string]string{keyVal[0]: keyVal[1]})
	}
	return &CW{namespace, basename, dimMaps, []CWData{}}
}

// AddData method to create and add CWData to CW object
func (cw *CW) AddData(suffix string, value float64) *CW {
	data := CWData{
		cw.BaseName + "-" + suffix,
		value,
		cw.Dimensions,
	}
	cw.Data = append(cw.Data, data)
	return cw
}
