package internal

import (
	"fmt"
	"time"
)

type TestOutputLine struct {
	Time    time.Time `json:"time"`
	Action  string    `json:"action"`
	Test    string    `json:"test,omitempty"`
	Pack    string    `json:"package"`
	Output  string    `json:"output,omitempty"`
	Elapsed float64   `json:"elapsed,omitempty"`
}

type TestOutputLines []TestOutputLine

func (f TestOutputLine) ShortString() string {
	return fmt.Sprintf("%s %s", f.Action, f.Test)
}
func (f TestOutputLine) String() string {
	return fmt.Sprintf("<time: %s, action: %s, test: %s, package: %s, output: %s, elapsed: %f>", f.Time, f.Action, f.Test, f.Pack, f.Output, f.Elapsed)
}
