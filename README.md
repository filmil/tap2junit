# TAP-to-jUnit converter

This is a go implementation of a converter from the TAP test format (from
www.testanything.org) to the jUnit format.

As both formats are somewhat loosely specified, the conversion is somewhat
of an interpretative dance.

# Installation

go get github.com/filmil/tap2junit
go install github.com/filmil/tap2junit/...

# Features

- Support for Version 12 of the TAP specification.
- Support for a custom extension to measure and report test duration.


