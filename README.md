# TAP-to-jUnit converter

[![CircleCI Build Status](https://circleci.com/gh/circleci/githubcom-filmil-tap2junit.svg?style=shield)](https://circleci.com/gh/circleci/cci-demo-docker) [![MIT Licensed](https://img.shields.io/badge/license-MIT-blue.svg)](https://raw.githubusercontent.com/circleci/cci-demo-react/master/LICENSE)

This is a go implementation of a converter from the TAP test format (from
www.testanything.org) to the jUnit format.

As both formats are somewhat loosely specified, the conversion is somewhat
of an interpretative dance.  We use a comprehensive test suite to guard known
good functionality.

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



