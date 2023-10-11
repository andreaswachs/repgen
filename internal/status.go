package internal

import (
	"encoding/json"
)

type Status int

const (
	Unknown Status = iota
	Pass
	Fail
	Bench
	Skip
)

/* extendLock makes sure that any implementing type cannot be extended outside of this package */
type extendLock interface{ xx() }

// Make sure that Status cannot be extended outside of this package, only valid values are the ones defined above
func (s Status) xx() {}

func (s Status) String() string {
	switch s {
	case Pass:
		return "pass"
	case Fail:
		return "fail"
	case Bench:
		return "bench"
	case Skip:
		return "skip"
	}

	return ""
}

func (s Status) FancyString() string {
	switch s {
	case Pass:
		return "✅ pass"
	case Fail:
		return "❌ fail"
	case Bench:
		return "⏱ bench"
	case Skip:
		return "⚠️ skip"
	}

	// Unknown
	return "⚠️ unknown"
}

func (s Status) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.String())
}
