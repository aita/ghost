package shell

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestShellExec(t *testing.T) {
	for _, tt := range []struct {
		script   string
		expected string
	}{
		{
			`echo hello world`,
			"hello world\n",
		},
		{
			`set x hello; echo $x`,
			"hello\n",
		},
	} {
		buf := bytes.NewBuffer(nil)
		sh := &Shell{
			Out: buf,
		}
		sh.Init()

		sh.Exec(tt.script)
		assert.Equal(t, tt.expected, buf.String())
	}
}
