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
         <testcase id="d32c977c8ba0374c3c0e821206cc08d19a041daa9caec8c7373de9175b1189e8" name="This test" time="10.000"></testcase>
         <testcase id="b3b1d666dfa8d2b061fc60641b53d49cd8df01ac940265b168e808c28e66a11e" name="That test" time="20.000">
            <failure message="That test" type="TestFailed"><![CDATA[ 2 That test # comment 2
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
         <testcase id="d32c977c8ba0374c3c0e821206cc08d19a041daa9caec8c7373de9175b1189e8" name="This test" time="20.000"></testcase>
         <testcase id="b3b1d666dfa8d2b061fc60641b53d49cd8df01ac940265b168e808c28e66a11e" name="That test" time="10.000">
            <failure message="That test" type="TestFailed"><![CDATA[# TAP2JUNIT: Duration: 10s
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
