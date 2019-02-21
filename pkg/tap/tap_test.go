package tap

import (
	"flag"
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
		testName string
		input    string
		reorder  bool
		expected Case
	}{
		{
			name:     "Basic",
			testName: "named_test",
			input: `
`,
			expected: Case{
				Version: 12,
				Name:    "named_test",
			},
		},
		{
			name: "TAP Version 42",
			input: `TAP version 42
`,
			expected: Case{
				Version: 42,
			},
		},
		{
			name: "One test",
			input: `
1..1
`,
			expected: Case{
				Version: 12,
				First:   ptr(1),
				Last:    ptr(1),
				Results: []Result{
					{},
				},
			},
		},
		{
			name:    "One OK test with comment",
			reorder: true,
			input: `
1..2
ok 2 Hello world # Some comment
# This is part of test 2
`,
			expected: Case{
				Version: 12,
				First:   ptr(1),
				Last:    ptr(2),
				Results: []Result{
					{
						Status: UNKNOWN,
					},
					{
						Status: PASSED,
						Header: "Hello world",
						Raw: ` 2 Hello world # Some comment
# This is part of test 2`,
					},
				},
				Raw: `
1..2
ok 2 Hello world # Some comment
# This is part of test 2
`,
			},
		},
		{
			name:    "Reorder timing report",
			reorder: true,
			input: `
1..2
# TAP2JUNIT: Duration: 4.3ms
ok 1 Hello world # Some comment
# This is part of test 1
ok 2 Test 2
# This is part of test 2
`,
			expected: Case{
				Version: 12,
				First:   ptr(1),
				Last:    ptr(2),
				Results: []Result{
					{
						Status: PASSED,
						Header: "Hello world",
						Raw: `# TAP2JUNIT: Duration: 4.3ms
 1 Hello world # Some comment
# This is part of test 1`,
						Duration: 4300 * time.Microsecond,
					},
					{
						Status: PASSED,
						Header: "Test 2",
						Raw: ` 2 Test 2
# This is part of test 2`,
					},
				},
				Raw: `
1..2
# TAP2JUNIT: Duration: 4.3ms
ok 1 Hello world # Some comment
# This is part of test 1
ok 2 Test 2
# This is part of test 2
`,
			},
		},
		{
			name: "One OK test with TODO",
			input: `
1..2
ok 2 Hello world # TODO not done yet
`,
			expected: Case{
				Version: 12,
				First:   ptr(1),
				Last:    ptr(2),
				Results: []Result{
					{
						Status: UNKNOWN,
					},
					{
						Status: TODO,
						Raw:    " 2 Hello world # TODO not done yet",
						Header: "Hello world",
					},
				},
				Raw: `
1..2
ok 2 Hello world # TODO not done yet
`,
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
# Part of fifth test
not ok 6 Sixth test # SKIP Failed here
# Some annotation
# TAP2JUNIT: Duration: 10s
not ok 7 Seventh test # TODO Failed here
ok Unnumbered test
`,
			expected: Case{
				Version: 12,
				First:   ptr(1),
				Last:    ptr(9),
				Results: []Result{
					{
						Status: UNKNOWN,
					},
					{
						Status: PASSED,
						Raw:    " 2 Hello world # Some comment",
						Header: "Hello world",
					},
					{
						Status: SKIPPED,
						Raw:    " 3 Third test # SKIP not implemented yet",
						Header: "Third test",
					},
					{
						Status: TODO,
						Raw:    " 4 Fourth test # TODO this is to be done",
						Header: "Fourth test",
					},
					{
						Status: FAILED,
						Raw: ` 5 Fifth test # Failed here
# Part of fifth test`,
						Header: `Fifth test`,
					},
					{
						// 5
						Status: SKIPPED,
						Raw: " 6 Sixth test # SKIP Failed here\n" +
							"# Some annotation\n" +
							"# TAP2JUNIT: Duration: 10s",
						Header:   "Sixth test",
						Duration: duration("10s"),
					},
					{
						Status: TODO,
						Raw:    " 7 Seventh test # TODO Failed here",
						Header: "Seventh test",
					},
					{
						// 7
						Status: PASSED,
						Raw:    " Unnumbered test",
						Header: "Unnumbered test",
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
			expected: Case{
				Version: 12,
				First:   ptr(1),
				Last:    ptr(5),
				Results: []Result{
					{Status: UNKNOWN},
					{Status: PASSED, Raw: " 2 Hello world # Some comment",
						Header: "Hello world"},
					{Status: UNKNOWN},
					{Status: UNKNOWN},
					{Status: UNKNOWN},
				},
			},
		},
	}
	flag.Parse()

	opts := cmp.Options{
		cmpopts.IgnoreFields(Case{}, "Raw"),
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			r := strings.NewReader(test.input)
			defer func() {
				if r := recover(); r != nil {
					t.Errorf("recovered: %v", r)
					debug.PrintStack()
				}
			}()
			actual, err := Read(r, test.testName, test.reorder)
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
