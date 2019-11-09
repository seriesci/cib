package cov

import "encoding/xml"

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
