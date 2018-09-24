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
			`\ `,
			" ",
		},
		{
			`hello\ world`,
			"hello world",
		},
		{
			"first\\\nsecond\\\n",
			"first\nsecond\n",
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
			"${var}",
			map[string]string{
				"var": "hello",
			},
			"hello",
		},
		{
			"$a${b} $c $d",
			map[string]string{
				"a": "1",
				"b": "2",
				"c": "3",
				"d": "4",
			},
			"12 3 4",
		},
	} {
		env := &Environment{
			store: tt.store,
		}
		s, err := expandDollar(env, tt.input)
		assert.Equal(t, tt.expected, s)
		assert.Nil(t, err)
	}
}
