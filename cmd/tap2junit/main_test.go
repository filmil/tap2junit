package main

import (
	"fmt"
	"runtime/debug"
	"strings"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func ptr(v int) *int {
	return &v
}

func duration(s string) time.Duration {
	d, err := time.ParseDuration(s)
	if err != nil {
		panic(fmt.Sprintf("for: %q: %v", s, err))
	}
	return d
}

func TestOutput(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected TAPCase
	}{
		{
			name: "Basic",
			input: `
`,
			expected: TAPCase{
				Version: 12,
			},
		},
		{
			name: "TAP Version 42",
			input: `TAP version 42
`,
			expected: TAPCase{
				Version: 42,
			},
		},
		{
			name: "One test",
			input: `
1..1
`,
			expected: TAPCase{
				Version: 12,
				First:   ptr(1),
				Last:    ptr(1),
				Results: []TAPResult{
					{},
				},
			},
		},
		{
			name: "One OK test with comment",
			input: `
1..2
ok 2 Hello world # Some comment
`,
			expected: TAPCase{
				Version: 12,
				First:   ptr(1),
				Last:    ptr(2),
				Results: []TAPResult{
					{
						Status: UNKNOWN,
					},
					{
						Status: OK,
						Raw:    " 2 Hello world # Some comment",
					},
				},
			},
		},
		{
			name: "One OK test with TODO",
			input: `
1..2
ok 2 Hello world # TODO not done yet
`,
			expected: TAPCase{
				Version: 12,
				First:   ptr(1),
				Last:    ptr(2),
				Results: []TAPResult{
					{
						Status: UNKNOWN,
					},
					{
						Status: TODO,
						Raw:    " 2 Hello world # TODO not done yet",
					},
				},
			},
		},
		{
			name: "Full test example",
			input: `
1..9
ok 2 Hello world # Some comment
ok 3 Third test # SKIP not implemented yet
ok 4 Fourth test # TODO this is to be done
not ok 5 Fifth test # Failed here
not ok 6 Sixth test # SKIP Failed here
# Some annotation
# TAP2JUNIT: Duration: 10s
not ok 7 Seventh test # TODO Failed here
ok Unnumbered test
`,
			expected: TAPCase{
				Version: 12,
				First:   ptr(1),
				Last:    ptr(9),
				Results: []TAPResult{
					{
						Status: UNKNOWN,
					},
					{
						Status: OK,
						Raw:    " 2 Hello world # Some comment",
					},
					{
						Status: SKIPPED,
						Raw:    " 3 Third test # SKIP not implemented yet",
					},
					{
						Status: TODO,
						Raw:    " 4 Fourth test # TODO this is to be done",
					},
					{
						Status: NOT_OK,
						Raw:    " 5 Fifth test # Failed here",
					},
					{
						// 5
						Status: SKIPPED,
						Raw: " 6 Sixth test # SKIP Failed here\n" +
							"# Some annotation\n" +
							"# TAP2JUNIT: Duration: 10s",
						Duration: duration("10s"),
					},
					{
						Status: TODO,
						Raw:    " 7 Seventh test # TODO Failed here",
					},
					{
						// 7
						Status: OK,
						Raw:    " Unnumbered test",
					},
					{
						// 8, this test was not ran.
						Status: UNKNOWN,
					},
				},
			},
		},
		{
			name: "Bail out",
			input: `
1..5
ok 2 Hello world # Some comment
Bail out! Some justification.
ok 3 Belated result
`,
			expected: TAPCase{
				Version: 12,
				First:   ptr(1),
				Last:    ptr(5),
				Results: []TAPResult{
					{Status: UNKNOWN},
					{Status: OK, Raw: " 2 Hello world # Some comment"},
					{Status: UNKNOWN},
					{Status: UNKNOWN},
					{Status: UNKNOWN},
				},
			},
		},
	}

	opts := cmp.Options{
		cmpopts.IgnoreFields(TAPCase{}, "Raw"),
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			r := strings.NewReader(test.input)
			defer func() {
				if r := recover(); r != nil {
					t.Errorf("recovered: %v", r)
					debug.PrintStack()
				}
			}()
			actual, err := ReadTAP(r)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if !cmp.Equal(test.expected, actual, opts) {
				t.Errorf("diff:\n%v\n, expected:\n%+v\nactual:\n%+v",
					cmp.Diff(test.expected, actual), test.expected, actual)
			}
		})
	}
}
