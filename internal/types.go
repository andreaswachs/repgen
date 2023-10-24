package internal

import "time"

/*
Represents the total collection of tests having been tested.
*/
type Tests struct {
	Packages map[string]TestPackage
}

/*
Represents the tests that a single package contains, being aware of the package itself
*/
type TestPackage struct {
	Name    string
	Tests   map[string]*Test
	Elapsed float64
	Time    time.Time
	Output  string
	Status  Status
}

/*
representation of a go test, where nested tests are able to exists
*/
type Test struct {
	Name        string
	Elapsed     float64
	Time        time.Time
	Output      string
	NestedTests map[string]*Test
	Status      Status
}

/*
Representation of a test as a data transfer object.
This exists for inputting content into the template, where we have a hard time
traversing these nested recursive structures
*/
type TestDTO struct {
	Name        string
	Elapsed     string
	ElapsedSec  float64
	Time        time.Time
	Output      []string
	Status      string
	FancyStatus string
	Package     string
}

/*
A data transfer object that moves test into the report template
*/
type TemplateData struct {
	Dependencies map[string]string
	Packages     []string
	Tests        []*TestDTO
	Passed       int
	Failed       int
	Skipped      int
}
