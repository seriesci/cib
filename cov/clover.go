package cov

import (
	"encoding/xml"
	"io/ioutil"
)

// Metrics is clover metrics.
type Metrics struct {
	Statements          int `xml:"statements,attr"`
	CoveredStatements   int `xml:"coveredstatements,attr"`
	Conditionals        int `xml:"conditionals,attr"`
	CoveredConditionals int `xml:"coveredconditionals,attr"`
	Methods             int `xml:"methods,attr"`
	CoveredMethods      int `xml:"coveredmethods,attr"`
	Elements            int `xml:"elements,attr"`
	CoveredElements     int `xml:"coveredelements,attr"`
	Complexity          int `xml:"complexity,attr"`
}

// Project is clover project.
type Project struct {
	XMLName   xml.Name `xml:"project"`
	Timestamp int64    `xml:"timestamp,attr"`
	Name      string   `xml:"name,attr"`
	Metrics   Metrics  `xml:"metrics"`
}

// Coverage is clover coverage.
type Coverage struct {
	XMLName   xml.Name `xml:"coverage"`
	Generated int64    `xml:"generated,attr"`
	Clover    string   `xml:"clover,attr"`
	Project   Project  `xml:"project"`
}

// ReadUnmarshal reads XML file and returns proper structs.
func ReadUnmarshal(path string) (*Coverage, error) {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var c Coverage
	return &c, xml.Unmarshal(b, &c)
}
