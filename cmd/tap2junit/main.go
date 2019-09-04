// Package main represents the program tap2junit, a utility that converts
// a testanything.org's TAP test format into junit format.
package main

import (
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/filmil/tap2junit/pkg/junit"
	"github.com/filmil/tap2junit/pkg/tap"
	"github.com/filmil/tap2junit/pkg/tap/tojunit"
	"github.com/golang/glog"
)

var (
	testName        = flag.String("test_name", "unnamed_test", "Sets the test name to use")
	reorderDuration = flag.Bool("reorder_duration", false, "If set, will reorder durations to work around https://github.com/bats-core/bats-core/issues/187")
	singleSuite     = flag.Bool("single_suite", false, "If set, will output only the <testsuite> as top-level tag; not <testsuites>")
)

func run(r io.Reader, w io.Writer, name string, reorderDuration, singleSuite bool) error {
	t, err := tap.Read(r, name, reorderDuration)
	if err != nil {
		return fmt.Errorf("while reading TAP: %v", err)
	}
	j, err := tojunit.FromTAP(t)
	if err != nil {
		return fmt.Errorf("while converting to jUnit: %v", err)
	}
	if err := junit.Write(j, w, singleSuite); err != nil {
		return fmt.Errorf("while writing jUnit: %v", err)
	}
	return nil
}

func init() {
	flag.Parse()
}

func main() {
	if err := run(os.Stdin, os.Stdout, *testName, *reorderDuration, *singleSuite); err != nil {
		glog.Fatalf("unexpected error: %v", err)
	}
}
