package junit

import (
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestOutput(t *testing.T) {
	tests := []struct {
		name     string
		input    Testsuites
		expected string
	}{
		{
			name: "Example from https://www.ibm.com/support/knowledgecenter/en/SSUFAU_1.0.0/com.ibm.rsar.analysis.codereview.cobol.doc/topics/cac_useresults_junit.html",
			input: Testsuites{
				ID:           "20140612_170519",
				Name:         "New_configuration (14/06/12 17:05:19)",
				NumTests:     225,
				NumFailiures: 1262,
				Time:         1 * time.Millisecond,
				Suites: []Suite{
					ID:          "codereview.cobol.analysisProvider",
					Name:        "COBOL Code Review",
					NumTests:    45,
					NumFailures: 17,
					Time:        1 * time.Millisecond,
					Testcases: []Case{
						ID:   "codereview.cobol.analysisProvider",
						Name: "COBOL Code Review",
						Time: 1 * time.Millisecond,
						Failures: []Failure{
							Message: "PROGRAM.cbl:2 Use a program name that matches the source file name",
							Text: `WARNING: Use a program name that matches the source file name
Category: COBOL Code Review – Naming Conventions
File: /project/PROGRAM.cbl
Line: 2
`,
						},
					},
				},
			},
			expected: `<?xml version="1.0" encoding="UTF-8" ?> 
   <testsuites id="20140612_170519" name="New_configuration (14/06/12 17:05:19)" tests="225" failures="1262" time="0.001">
      <testsuite id="codereview.cobol.analysisProvider" name="COBOL Code Review" tests="45" failures="17" time="0.001">
         <testcase id="codereview.cobol.rules.ProgramIdRule" name="Use a program name that matches the source file name" time="0.001">
            <failure message="PROGRAM.cbl:2 Use a program name that matches the source file name" type="WARNING">
WARNING: Use a program name that matches the source file name
Category: COBOL Code Review – Naming Conventions
File: /project/PROGRAM.cbl
Line: 2
      </failure>
    </testcase>
  </testsuite>
</testsuites>
`,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var out strings.Builder
			Write(test.input, &out)
			if !cmp.Equal(test.expected, out.String()) {
				t.Errorf("diff:\n%v\nexpected:\n%+v\nactual:\n%+v",
					cmp.Diff(test.expected, out.String()), test.expected, out.String())
			}
		})
	}
}
