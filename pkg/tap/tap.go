// Package tap parses the TAP test protocol file into a Case.  TAP protocol is
// defined by the specification at www.testanything.org.
package tap

import (
	"bufio"
	"fmt"
	"io"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/golang/glog"
)

// Status is the test status type.
type Status int

const (
	// UNKNOWN status represents a test of unknown status.
	UNKNOWN Status = Status(0)
	// OK status represents a passed test.
	PASSED = Status(1)
	// FAILED status represents a failed test.
	FAILED = Status(2)
	// SKIPPED status represents a test that was skipped.
	SKIPPED = Status(3)
	// TODO status represents a test that was marked TODO
	TODO = Status(4)
)

// Result is the result of a single TAP test.
type Result struct {
	// Status shows the status of this test.
	Status Status
	// Duration is how long the test took.
	Duration time.Duration
	// Header is the first line of the test, complete.
	Header string
	// Raw is the raw content of the test result dump.
	Raw string
}

// Case is the result of running a TAP test suite.
type Case struct {
	// Version is the TAP test specification.  If unset, it is a specification
	// earlier than version 13.
	Version int
	// Name is the unique name of the test.
	Name string
	// First is the first test that was executed. Optional.
	First *int
	// Last is the last test. Optional.
	Last *int
	// Test results, if any.
	Results []Result
	// Raw contents of the TAPSpec
	Raw string
	// Duration is how long the test took, if known.
	Duration time.Duration
}

// toInt parses a string to int.  The string is known to be parseable to int.
func toInt(s string) int {
	i, _ := strconv.Atoi(s)
	return i
}

func copyResize(dest *[]Result, newSize int) {
	if len(*dest) >= newSize {
		return
	}
	new := make([]Result, newSize)
	copy(*dest, new)
	*dest = new
}

// StatusFrom returns a Status from a supplied string.
func StatusFrom(str string, def Status) Status {
	switch str {
	case "TODO":
		fallthrough
	case "todo":
		return TODO
	case "SKIP":
		fallthrough
	case "skip":
		return SKIPPED
	}
	return def
}

// parser contains the parser status
type parser struct {
	// last test read
	lt int
}

// Read parses the contents of i into a Result. name is a given test name.
func Read(i io.Reader, name string) (Case, error) {
	var (
		r  Case = Case{Version: 12, Name: name}
		ps parser
	)
	s := bufio.NewScanner(i)
	for s.Scan() {
		t := s.Text()
		r.Raw = fmt.Sprintf("%s%s\n", r.Raw, t)

		glog.V(2).Infof("Text: %q", t)

		// Spec is the TAP line representing the version.
		// "TAP version 13"
		var Spec = regexp.MustCompile(`TAP version (\d+)`)
		if v := Spec.FindStringSubmatch(t); v != nil {
			glog.V(2).Infof("version: %v", spew.Sdump(v))
			r.Version = toInt(v[1])
			continue
		}

		// Range is the range of the tests to run.
		// "1..42"
		var Range = regexp.MustCompile(`(\d+)\.\.(\d+)`)
		if v := Range.FindStringSubmatch(t); v != nil {
			glog.V(2).Infof("range: %v", spew.Sdump(v))
			f := toInt(v[1])
			r.First = &f
			if ps.lt == 0 {
				// If we haven't yet scanned any tests, update the test counter:
				// next unnumbered test result will be for test lt.
				ps.lt = f
			}
			l := toInt(v[2])
			r.Last = &l
			// Resize the results array to fit.
			copyResize(&r.Results, *r.Last)
			continue
		}

		// OKTest is an OK test line.
		// "ok 41 some text # TODO some comment"
		var OKTest = regexp.MustCompile(
			`^ok( (\d+)?(\s+)?(([^#]*))?(#\s+(TODO|todo|SKIP|skip)?(.*))?)`)

		// "42 ok Some comment"
		// Regex analysis using https://regex101.com
		if v := OKTest.FindStringSubmatch(t); v != nil {
			glog.V(2).Infof("ok: %v", spew.Sdump(v))
			tiStr := v[2]
			if tiStr != "" {
				ps.lt = toInt(tiStr)
			} else {
				ps.lt++
			}
			copyResize(&r.Results, ps.lt)
			if r.Last == nil || *r.Last < ps.lt {
				l := ps.lt
				r.Last = &l
			}
			r.Results[ps.lt-1].Status = StatusFrom(v[7], PASSED)
			r.Results[ps.lt-1].Raw = v[1]
			r.Results[ps.lt-1].Header = t
			continue
		}

		// NotOKTest is a failed test line.
		// "not ok 42 some test # SKIP some comment"
		// 1: test number, optional
		// 3: text before #
		var NotOKTest = regexp.MustCompile(
			`^not ok( (\d+)?(\s+)?(([^#]*))?(#\s+(TODO|todo|SKIP|skip)?(.*))?)`)
		if v := NotOKTest.FindStringSubmatch(t); v != nil {
			glog.V(2).Infof("not ok: %v", spew.Sdump(v))
			tiStr := v[2]
			if tiStr != "" {
				ps.lt = toInt(tiStr)
			} else {
				ps.lt++
			}
			copyResize(&r.Results, ps.lt)
			if r.Last == nil || *r.Last < ps.lt {
				l := ps.lt
				r.Last = &l
			}
			r.Results[ps.lt-1].Status = StatusFrom(v[7], FAILED)
			r.Results[ps.lt-1].Raw = v[1]
			r.Results[ps.lt-1].Header = t
			continue
		}

		// An annotation is attached to the "current" test.
		var TestAnnotation = regexp.MustCompile(`^#(\s+)?(.+)?`)
		if v := TestAnnotation.FindStringSubmatch(t); v != nil {
			glog.V(2).Infof("annotation: %v", spew.Sdump(v))
			line := v[0]
			// This is an annotation for the current test.
			r.Results[ps.lt-1].Raw = strings.Join([]string{r.Results[ps.lt-1].Raw, line}, "\n")

			// Extension parsing
			if strings.HasPrefix(line, "# TAP2JUNIT:") {
				line = strings.TrimPrefix(line, "# TAP2JUNIT:")
				line = strings.TrimSpace(line)
			}
			glog.V(2).Infof("extension: %v", line)
			if strings.HasPrefix(line, "Duration:") {
				line = strings.TrimPrefix(line, "Duration:")
				line = strings.TrimSpace(line)

				d, err := time.ParseDuration(line)
				if err != nil {
					glog.Warningf("could not parse duration: %v", line)
				}
				r.Results[ps.lt-1].Duration = d
			}
			continue
		}

		var BailOut = regexp.MustCompile(`Bail out!\s*(.*)`)
		if v := BailOut.FindStringSubmatch(t); v != nil {
			glog.V(3).Infof("Found bail out! text: %q", v[1])
			break
		}
		glog.V(2).Infof("no match: %q", t)
	}
	if s.Err() != nil {
		return r, s.Err()
	}
	return r, nil
}
