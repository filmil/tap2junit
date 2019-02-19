// Package main represents the program tap2junit, a utility that converts
// a testanything.org's TAP test format into junit format.
package main

import (
	"os"

	"github.com/filmil/tap2junit/pkg/tap"
	"github.com/golang/glog"
)

func main() {
	_, err := tap.Read(os.Stdin)
	if err != nil {
		glog.Fatalf("unexpected error: %v", err)
	}
}
