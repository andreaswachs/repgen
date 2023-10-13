package test_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestPass(t *testing.T) {}

func TestFailure(t *testing.T) {
	t.Fail()
}

func TestPassNested(t *testing.T) {
	t.Run("nested 1", func(t *testing.T) {})
	t.Run("nested 2", func(t *testing.T) {})
	t.Run("nested 3", func(t *testing.T) {})
	t.Run("nested 4", func(t *testing.T) {})
	t.Run("nested 5", func(t *testing.T) {})
}

func TestTakesTime(t *testing.T) {
	time.Sleep(2 * time.Second)
}

func TestTakesTimeFails(t *testing.T) {
	time.Sleep(2 * time.Second)
	t.Fail()
}

func TestFailNested(t *testing.T) {
	t.Run("nested 1", func(t *testing.T) { t.Fail() })
	t.Run("nested 2", func(t *testing.T) { t.Fail() })
	t.Run("nested 3", func(t *testing.T) { t.Fail() })
	t.Run("nested 4", func(t *testing.T) { t.Fail() })
	t.Run("nested 5", func(t *testing.T) { t.Fail() })
}

func TestMixedNested(t *testing.T) {
	t.Run("nested 1", func(t *testing.T) {})
	t.Run("nested 2", func(t *testing.T) { t.Fail() })
	t.Run("nested 3", func(t *testing.T) {})
	t.Run("nested 4", func(t *testing.T) { t.Fail() })
	t.Run("nested 5", func(t *testing.T) {})
}

func TestFailWithTestify(t *testing.T) {
	assert.NotNil(t, nil)
}

func TestNestedFailWithTestify(t *testing.T) {
	t.Run("nested 1", func(t *testing.T) { assert.NotNil(t, nil) })
	t.Run("nested 2", func(t *testing.T) { assert.NotNil(t, nil) })
	t.Run("nested 3", func(t *testing.T) { assert.NotNil(t, nil) })
}
