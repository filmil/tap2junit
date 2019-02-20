package tojunit

import (
	"testing"
	"time"

	"github.com/filmil/tap2junit/pkg/junit"
	"github.com/filmil/tap2junit/pkg/tap"
	"github.com/google/go-cmp/cmp"
)

func TestConversion(t *testing.T) {
	tests := []struct {
		name     string
		input    tap.Case
		expected junit.Testsuites
	}{
		{
			name: "Basic",
			input: tap.Case{
				Version: 12,
				Name:    "test_name_here",
				Results: []tap.Result{
					{
						Status:   tap.PASSED,
						Duration: 2 * time.Second,
						Raw:      "Some string here",
						Header:   "Header0",
					},
					{
						Status:   tap.FAILED,
						Duration: 3 * time.Second,
						Raw: `not ok 2 Test failed
# Some failure message
`,
						Header: "Header1",
					},
					{
						Status: tap.SKIPPED,
						Header: "Header2",
					},
					{
						Status: tap.UNKNOWN,
					},
				},
				Raw: "Raw string",
			},
			expected: junit.Testsuites{
				NumTests:    4,
				NumFailures: 1,
				Time:        junit.DurationSec{5 * time.Second},
				Suites: []junit.Suite{
					{
						ID:          strHash("test_name_here"),
						Name:        "test_name_here",
						NumTests:    4,
						NumFailures: 1,
						Time:        junit.DurationSec{5 * time.Second},
						Testcases: []junit.Case{
							{
								ID:   strHash("Header0"),
								Name: "Header0",
								Time: junit.DurationSec{2 * time.Second},
							},
							{
								ID:   strHash("Header1"),
								Name: "Header1",
								Time: junit.DurationSec{3 * time.Second},
								Failures: []junit.Failure{
									{
										Type:    "TestFailed",
										Message: "Header1",
										Text: `not ok 2 Test failed
# Some failure message
`,
									},
								},
							},
							{
								ID:   strHash("Header2"),
								Name: "Header2",
							},
							{
								ID: strHash(""),
							},
						},
					},
				},
			},
		},
	}
	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			actual, err := FromTAP(test.input)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if !cmp.Equal(test.expected, actual) {
				t.Errorf("diff:\n%v\nexpected:\n%+v\nactual:\n%+v",
					cmp.Diff(test.expected, actual), test.expected, actual)
			}
		})
	}
}
