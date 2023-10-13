package internal

import (
	"strings"
	"time"
)

// Initialises an empty Tests struct
func NewTests() Tests {
	return Tests{
		Packages: make(map[string]TestPackage),
	}
}

func newTestPackage(f TestOutputLine) TestPackage {
	return TestPackage{
		Name:    f.Pack,
		Tests:   make(map[string]*Test),
		Elapsed: f.Elapsed,
		Time:    f.Time,
		Output:  f.Output,
		Status:  Unknown,
	}
}

func newNode(f TestOutputLine) *Test {
	return &Test{
		Name:        f.Test,
		Elapsed:     f.Elapsed,
		Time:        f.Time,
		Output:      f.Output,
		NestedTests: make(map[string]*Test),
		Status:      Unknown,
	}
}

/*
Converts the collection of tests to a format that can be used by the template
*/
func (t *Tests) ToTemplateData() (*TemplateData, error) {
	passed := 0
	failed := 0
	skipped := 0

	tests, err := t.ToFlatStructure()
	if err != nil {
		return nil, err
	}

	filteredTests := make([]*SingleTest, 0)
	for _, test := range tests {
		if test.Name == "" {
			continue
		}

		filteredTests = append(filteredTests, test)
	}

	for _, test := range filteredTests {
		switch test.Status {
		case "pass":
			passed++
		case "fail":
			failed++
		case "skip":
			skipped++
		}
	}

	packages := make([]string, 0)
	for _, p := range t.Packages {
		packages = append(packages, p.Name)
	}

	return &TemplateData{
		Packages: packages,
		Tests:    filteredTests,
		Passed:   passed,
		Failed:   failed,
		Skipped:  skipped,
	}, nil
}

func (t *Tests) ToFlatStructure() ([]*SingleTest, error) {
	flatTests := make([]*SingleTest, 0)
	for _, pack := range t.Packages {
		for _, test := range pack.Tests {
			tests, err := test.getFlatTests(pack.Name)
			if err != nil {
				return nil, err
			}

			flatTests = append(flatTests, tests...)
		}
	}

	return flatTests, nil
}

func (t *Tests) AddField(f TestOutputLine) error {
	switch f.Action {
	case "start":
		return t.addStart(f)
	case "run":
		return t.addRun(f)
	case "pass":
		return t.addPass(f)
	case "fail":
		return t.addFail(f)
	case "skip":
		return t.addSkip(f)
	case "output":
		return t.addOutput(f)
	case "bench":
		return t.addBench(f)
	case "pause":
		return t.addPause(f)
	case "cont":
		return t.addCont(f)
	}

	return nil
}

func (t *Tests) addBench(f TestOutputLine) error {
	return nil
}

func (t *Tests) addPause(f TestOutputLine) error {
	return t.addOutput(f)
}

func (t *Tests) addCont(f TestOutputLine) error {
	return t.addOutput(f)
}
func (t *Tests) addStart(f TestOutputLine) error {
	t.Packages[f.Pack] = newTestPackage(f)
	return nil
}

func (t *Tests) addRun(f TestOutputLine) error {
	// We assume that the package exists
	pack := t.Packages[f.Pack]
	tokens := strings.Split(f.Test, "/")
	if len(tokens) == 1 {
		// We have a top level test
		pack.Tests[f.Test] = newNode(f)
	} else {
		// We have a nested test
		test := pack.Tests[tokens[0]]
		for i := 1; i < len(tokens); i++ {
			if _, ok := test.NestedTests[tokens[i]]; !ok {
				test.NestedTests[tokens[i]] = newNode(f)
			}

			test = test.NestedTests[tokens[i]]
		}
	}

	return nil
}

func (t *Tests) addPass(f TestOutputLine) error {
	// We assume that the package exists
	// We only set passed on the leaf
	pack := t.Packages[f.Pack]
	tokens := strings.Split(f.Test, "/")
	if len(tokens) == 1 {

		if f.Test == "" {
			// The package itself passed
			pack.Elapsed = f.Elapsed
			pack.Status = Pass
			return nil
		}

		if pack.Tests[f.Test].Elapsed < f.Elapsed {
			pack.Tests[f.Test].Elapsed = f.Elapsed
		}

		// We have a top level test
		pack.Tests[f.Test].Status = Pass
	} else {
		// We have a nested test
		test := pack.Tests[tokens[0]]
		for i := 1; i < len(tokens); i++ {
			test = test.NestedTests[tokens[i]]
			if test == nil {
				test = newNode(f)
				test.NestedTests[tokens[i]] = test
			}
		}

		test.Status = Pass
		if test.Elapsed < f.Elapsed {
			test.Elapsed = f.Elapsed
		}
	}

	return nil
}

func (t *Tests) addFail(f TestOutputLine) error {
	// We assume that the package exists, we flag all nodes until the leaf as failed
	pack := t.Packages[f.Pack]
	tokens := strings.Split(f.Test, "/")
	if len(tokens) == 1 {
		// We have a top level test

		if f.Test == "" {
			// The package itself failed
			pack.Status = Fail
			pack.Elapsed = f.Elapsed
			return nil
		}

		if pack.Tests[f.Test] == nil {
			pack.Tests[f.Test] = newNode(f)
		}

		if pack.Tests[f.Test].Elapsed < f.Elapsed {
			pack.Tests[f.Test].Elapsed = f.Elapsed
		}

		pack.Tests[f.Test].Status = Fail
	} else {
		// We have a nested test
		test := pack.Tests[tokens[0]]
		for i := 1; i < len(tokens); i++ {
			test = test.NestedTests[tokens[i]]
			test.Status = Fail
		}

		if test.Elapsed < f.Elapsed {
			test.Elapsed = f.Elapsed
		}
		test.Status = Fail
	}

	return nil
}

func (t *Tests) addSkip(f TestOutputLine) error {
	// We assume that the package exists, we flag only the given test as skipped
	pack := t.Packages[f.Pack]
	tokens := strings.Split(f.Test, "/")
	if len(tokens) == 1 {
		// We have a top level test
		if pack.Tests[f.Test] == nil {
			pack.Tests[f.Test] = newNode(f)
		}

		pack.Tests[f.Test].Status = Skip
	} else {
		// We have a nested test
		test := pack.Tests[tokens[0]]
		for i := 1; i < len(tokens); i++ {
			test = test.NestedTests[tokens[i]]
		}

		test.Status = Skip
	}

	return nil
}

func (t *Tests) addOutput(f TestOutputLine) error {
	// We assume that the package exists. We append output to the leaf
	pack := t.Packages[f.Pack]

	if f.Test == "" {
		// Output for a whole package, append to package itself
		pack.Output += f.Output
		return nil
	}

	tokens := strings.Split(f.Test, "/")
	if len(tokens) == 1 {
		// We have a top level test
		if pack.Tests[f.Test] == nil {
			pack.Tests[f.Test] = newNode(f)
		}
		pack.Tests[f.Test].Output += f.Output
	} else {
		// We have a nested test
		test := pack.Tests[tokens[0]]
		if test == nil {
			test = newNode(f)
			pack.Tests[tokens[0]] = test
		}

		for i := 1; i < len(tokens); i++ {
			test = test.NestedTests[tokens[i]]
			if test == nil {
				test = newNode(f)
				test.NestedTests[tokens[i]] = test
			}
		}

		test.Output += f.Output
	}

	return nil
}

func (t *Test) ToSingleTest(p string) *SingleTest {
	return &SingleTest{
		Name:        t.Name,
		Elapsed:     time.Duration(t.Elapsed * float64(time.Second)).String(),
		Time:        t.Time,
		Output:      strings.Split(t.Output, "\n"),
		Status:      t.Status.String(),
		FancyStatus: t.Status.FancyString(),
		Package:     p,
	}
}

func (t *Test) getFlatTests(p string) ([]*SingleTest, error) {
	tests := make([]*SingleTest, 0)
	tests = append(tests, t.ToSingleTest(p))

	for _, t := range t.NestedTests {
		ts, err := t.getFlatTests(p)
		if err != nil {
			return nil, err
		}
		tests = append(tests, ts...)
	}

	return tests, nil
}
