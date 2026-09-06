// Package fixtures provides go:generate directives to build test fixture binaries.
package fixtures

//go:generate go build -o ../bin/fixture-good ./good/
//go:generate go build -o ../bin/fixture-no-nocolor ./no-nocolor/
//go:generate go build -o ../bin/fixture-no-help ./no-help/
