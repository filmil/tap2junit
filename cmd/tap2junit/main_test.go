package main

import (
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestMain(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
		reorder  bool
	}{
		{
			name: "Basic",
			input: `1..2
ok 1 This test # comment 1
# TAP2JUNIT: Duration: 10s
not ok 2 That test # comment 2
# TAP2JUNIT: Duration: 20s
`,
			expected: `<?xml version="1.0" encoding="UTF-8"?>
   <testsuites tests="2" failures="1" time="30.000">
      <testsuite id="7cc84235ce3aaeab160cebf213fdff2a0d92dcb4e6304dee5fb2762673f107f1" name="named_test" tests="2" failures="1" time="30.000">
         <testcase id="ba9ec74f753775734860835065bd83505683bc030a628250cbf4695e45c80c60" name="ok 1 This test # comment 1" time="10.000"></testcase>
         <testcase id="3a2718564b85de1cafbef2ac551beab4246645a25db12e8ad1fbd7abd709265c" name="not ok 2 That test # comment 2" time="20.000">
            <failure message="not ok 2 That test # comment 2" type="TestFailed"><![CDATA[ 2 That test # comment 2
# TAP2JUNIT: Duration: 20s]]></failure>
         </testcase>
      </testsuite>
   </testsuites>`,
		},
		{
			name:    "Reordered",
			reorder: true,
			input: `1..2
# TAP2JUNIT: Duration: 20s
ok 1 This test # comment 1
# TAP2JUNIT: Duration: 10s
not ok 2 That test # comment 2
`,
			expected: `<?xml version="1.0" encoding="UTF-8"?>
   <testsuites tests="2" failures="1" time="30.000">
      <testsuite id="7cc84235ce3aaeab160cebf213fdff2a0d92dcb4e6304dee5fb2762673f107f1" name="named_test" tests="2" failures="1" time="30.000">
         <testcase id="ba9ec74f753775734860835065bd83505683bc030a628250cbf4695e45c80c60" name="ok 1 This test # comment 1" time="20.000"></testcase>
         <testcase id="3a2718564b85de1cafbef2ac551beab4246645a25db12e8ad1fbd7abd709265c" name="not ok 2 That test # comment 2" time="10.000">
            <failure message="not ok 2 That test # comment 2" type="TestFailed"><![CDATA[# TAP2JUNIT: Duration: 10s
 2 That test # comment 2]]></failure>
         </testcase>
      </testsuite>
   </testsuites>`,
		},
	}
	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			var b strings.Builder
			run(strings.NewReader(test.input), &b, "named_test", test.reorder)
			actual := strings.Split(b.String(), "\n")
			exp := strings.Split(test.expected, "\n")
			if !cmp.Equal(exp, actual) {
				t.Errorf("diff:\n%v\nexpected:\n%v\nactual:\n%v",
					cmp.Diff(exp, actual), exp, actual)
			}
		})
	}
}
