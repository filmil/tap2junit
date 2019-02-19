// Package junit is a data model for a jUnit test.
package junit

import (
	"encoding/xml"
	"fmt"
	"io"
	"time"
)

// DurationSec is a duration, expressed in seconds when marshaling.
type DurationSec struct {
	time.Duration
}

var _ xml.MarshalerAttr = DurationSec{}

// MarshalXML implements xml.MarshalerAttr.
func (d DurationSec) MarshalXMLAttr(name xml.Name) (xml.Attr, error) {
	s := fmt.Sprintf("%.3f", d.Duration.Seconds())
	return xml.Attr{Name: name, Value: s}, nil
}

// Testsuites is a definition of the test suites.
type Testsuites struct {
	XMLName     xml.Name    `xml:"testsuites"`
	ID          string      `xml:"id,attr"`
	Name        string      `xml:"name,attr"`
	NumTests    int         `xml:"tests,attr"`
	NumFailures int         `xml:"failures,attr"`
	Time        DurationSec `xml:"time,attr"`
	Suites      []Suite
}

type Suite struct {
	XMLName     xml.Name    `xml:"testsuite"`
	ID          string      `xml:"id,attr"`
	Name        string      `xml:"name,attr"`
	NumTests    int         `xml:"tests,attr"`
	NumFailures int         `xml:"failures,attr"`
	Time        DurationSec `xml:"time,attr"`
	Testcases   []Case
}

// Case is a description of a single result test case.
type Case struct {
	XMLName  xml.Name    `xml:"testcase"`
	ID       string      `xml:"id,attr"`
	Name     string      `xml:"name,attr"`
	Time     DurationSec `xml:"time,attr"`
	Failures []Failure
}

// Failure is a message about a single test failure.
type Failure struct {
	XMLName xml.Name `xml:"failure"`
	Message string   `xml:"message,attr"`
	Type    string   `xml:"type,attr"`
	Text    string   `xml:",chardata"`
}

// Write writes out the test suites information into the supplied writer.
func Write(suites Testsuites, w io.Writer) error {
	e := xml.NewEncoder(w)
	e.Indent("   ", "   ")
	fmt.Fprintf(w, xml.Header)
	return e.Encode(suites)
}
