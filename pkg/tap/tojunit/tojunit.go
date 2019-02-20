// Package tojunint contains code for TAP to junit conversion.
package tojunit

import (
	"crypto/sha256"
	"fmt"
	"time"

	"github.com/filmil/tap2junit/pkg/junit"
	"github.com/filmil/tap2junit/pkg/tap"
)

func strHash(s string) string {
	h := sha256.New()
	h.Write([]byte(s))
	return fmt.Sprintf("%x", h.Sum(nil))
}

// FromTAP converts a TAP test case into a jUnit testsuite.
func FromTAP(c tap.Case) (junit.Testsuites, error) {
	var (
		r      junit.Testsuites
		s      junit.Suite
		nt, nf int
		d      time.Time
	)
	for _, r := range c.Results {
		var c junit.Case
		c.ID = strHash(r.Header)
		c.Name = r.Header
		nt++
		if r.Status == tap.FAILED {
			var f junit.Failure
			nf++
			f.Type = "TestFailed"
			f.Text = r.Raw
			f.Message = r.Header
			// Test message - full first line
			c.Failures = append(c.Failures, f)
		}
		d = d.Add(r.Duration)
		c.Time = junit.DurationSec{r.Duration}
		s.Testcases = append(s.Testcases, c)
	}
	td := d.Sub(time.Time{})
	r.Time = junit.DurationSec{td}
	r.NumTests = nt
	r.NumFailures = nf

	s.Name = c.Name
	s.ID = strHash(s.Name)
	s.Time = junit.DurationSec{td}
	s.NumTests = nt
	s.NumFailures = nf

	r.Suites = append(r.Suites, s)
	return r, nil
}
