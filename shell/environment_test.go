package shell

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEnvironment(t *testing.T) {
	topLevel := &Environment{
		store: map[string]string{
			"x": "1",
			"z": "2",
		},
	}
	localEnv := &Environment{
		outer: topLevel,
		store: map[string]string{
			"z": "3",
		},
	}

	for _, tt := range []struct {
		name   string
		val    string
		exists bool
	}{
		{"x", "1", true},
		{"y", "", false},
		{"z", "3", true},
	} {
		val, exists := localEnv.Get(tt.name)
		assert.Equal(t, tt.val, val)
		assert.Equal(t, tt.exists, exists)
	}
}
