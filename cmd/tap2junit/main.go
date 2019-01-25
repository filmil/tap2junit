// Package main represents the program tap2junit, a utility that converts
// a testanything.org's TAP test format into junit format.
package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"regexp"
	"strconv"
	"time"

	"github.com/davecgh/go-spew/spew"
)

type Status int

const (
	// UNKNOWN status represents a test of unknown status.
	UNKNOWN Status = iota
	OK
	NOT_OK
	SKIPPED
	TODO
)

// Result is the result of a single TAP test.
type TAPResult struct {
	// Status shows the status of this test.
	Status Status
	// Duration is how long the test took.
	Duration time.Duration
	// Raw is the raw content of the test result dump.
	Raw string
}

// TAPCase is the result of running a TAP test suite.
type TAPCase struct {
	// Version is the TAP test specification.  If unset, it is a specification
	// earlier than version 13.
	Version int
	// First is the first test that was executed. Optional.
	First *int
	// Last is the last test. Optional.
	Last *int
	// Test results, if any.
	Results []TAPResult
	// Raw contents of the TAPSpec
	Raw string
	// Duration is how long the test took, if known.
	Duration time.Duration
}

var (
	// Spec is the TAP line representing the version.
	// "TAP version 13"
	Spec = regexp.MustCompile(`TAP version (\d+)`)
	// Range is the range of the tests to run.
	// "1..42"
	Range = regexp.MustCompile(`(\d+)\.\.(\d+)`)

	// OKTest is an OK test line.
	// "ok 41 some text # TODO some comment"
	OkTest = regexp.MustCompile(`xok( ((\d+)?) ([^#])# (TODO|todo|SKIP|skip)?(\W+))?`)

	// NotOkTest is a failed test line.
	// "not ok 42 some test # SKIP some comment"
	// 1: test number, optional
	// 3: text before #
	NotOkTest = regexp.MustCompile(`not ok( ((\d+)?) ([^#])((# (TODO|todo|SKIP|skip)?(\W+))?))`)
)

// toInt parses a string to int.  The string is known to be parseable to int.
func toInt(s string) int {
	i, _ := strconv.Atoi(s)
	return i
}

func copyResize(dest *[]TAPResult, newSize int) {
	if len(*dest) >= newSize {
		return
	}
	new := make([]TAPResult, newSize)
	copy(*dest, new)
	*dest = new
}

// ReadTAP parses the contents of i into a TAPResult.
func ReadTAP(i io.Reader) (TAPCase, error) {
	var (
		r  TAPCase = TAPCase{Version: 12}
		lt int     // last test read
	)
	s := bufio.NewScanner(i)
	for s.Scan() {
		t := s.Text()
		r.Raw = fmt.Sprintf("%s%s\n", r.Raw, t)

		log.Printf("Text: %q", t)
		if v := Spec.FindStringSubmatch(t); v != nil {
			log.Printf("version: %v", spew.Sdump(v))
			r.Version = toInt(v[1])
		} else if v := Range.FindStringSubmatch(t); v != nil {
			log.Printf("range: %v", spew.Sdump(v))
			f := toInt(v[1])
			r.First = &f
			if lt == 0 {
				// If we haven't yet scanned any tests, update the test counter:
				// next unnumbered test result will be for test lt.
				lt = f
			}
			l := toInt(v[2])
			r.Last = &l
			// Resize the results array to fit.
			copyResize(&r.Results, *r.Last)
		} else if v := OkTest.FindStringSubmatch(t); v != nil {
			log.Printf("ok: %v", spew.Sdump(v))
			if v[1] != "" {
				lt = toInt(v[1])
			}
			copyResize(&r.Results, lt)
			if r.Last == nil || *r.Last < lt {
				l := lt
				r.Last = &l
			}
			r.Results[lt].Status = OK
		} else if v := NotOkTest.FindStringSubmatch(t); v != nil {
			log.Printf("not ok: %v", spew.Sdump(v))
			if v[1] != "" {
				lt = toInt(v[1])
			}
			copyResize(&r.Results, lt)
			if r.Last == nil || *r.Last < lt {
				l := lt
				r.Last = &l
			}
			r.Results[lt].Status = NOT_OK
			lt++
		} else {
			log.Printf("no match: %q", t)
		}
	}
	if s.Err() != nil {
		return r, s.Err()
	}
	return r, nil
}

func main() {
	r, err := ReadTAP(os.Stdin)
	if err != nil {
		log.Fatalf("unexpected error: %v", err)
	}
	fmt.Printf("%v", spew.Sdump(r))
}
