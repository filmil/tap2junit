# TAP-to-jUnit converter

This is a go implementation of a converter from the TAP test format (from
www.testanything.org) to the jUnit format.

Example:

```console
$ tap2junit -test_name="my_test" <<EOF
1..2
ok 1 This test # comment 1
# TAP2JUNIT: Duration: 10s
not ok 2 That test # comment 2
# TAP2JUNIT: Duration: 10s
EOF
<?xml version="1.0" encoding="UTF-8"?>
   <testsuites tests="2" failures="1" time="20.000">
      <testsuite id="3f1f8851c5e6eee11c8c1a5bb777cc303dbb406e6e8eba79ea1c1477e15781ac" name="my_test" tests="2" failures="1" time="20.000">
         <testcase id="ba9ec74f753775734860835065bd83505683bc030a628250cbf4695e45c80c60" name="ok 1 This test # comment 1" time="10.000"></testcase>
         <testcase id="3a2718564b85de1cafbef2ac551beab4246645a25db12e8ad1fbd7abd709265c" name="not ok 2 That test # comment 2" time="10.000">
            <failure message="not ok 2 That test # comment 2" type="TestFailed"> 2 That test # comment 2&#xA;# TAP2JUNIT: Duration: 10s</failure>
         </testcase>
      </testsuite>
   </testsuites>
$
```

As both formats are somewhat loosely specified, the conversion is somewhat
of an interpretative dance.  We use a comprehensive test suite to guard the
functionality.

# Installation

```
go get github.com/filmil/tap2junit
go install github.com/filmil/tap2junit/...
```

# Testing

```
go test github.com/filmil/tap2junit/...
```

# Features

- Support for Version 12 of the TAP specification.
- Support for a custom extension to measure and report test duration.



