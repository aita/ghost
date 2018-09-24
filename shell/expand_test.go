package shell

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExpandEscape(t *testing.T) {
	for _, tt := range []struct {
		input    string
		expected string
	}{
		{
			`\n`,
			"\n",
		},
		{
			`hello world\n`,
			"hello world\n",
		},
		{
			`first\n second\n`,
			"first\n second\n",
		},
	} {
		assert.Equal(t, tt.expected, expandEscape(tt.input))
	}
}

func TestExpandDollar(t *testing.T) {
	for _, tt := range []struct {
		input    string
		store    map[string]string
		expected string
	}{
		{
			"hello",
			nil,
			"hello",
		},
		{
			"$var",
			map[string]string{
				"var": "hello",
			},
			"hello",
		},
		{
			"$x$y $z",
			map[string]string{
				"x": "1",
				"y": "2",
				"z": "3",
			},
			"12 3",
		},
	} {
		env := &Environment{
			store: tt.store,
		}
		assert.Equal(t, tt.expected, expandDollar(env, tt.input))
	}
}
