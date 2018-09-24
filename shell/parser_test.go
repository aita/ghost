package shell

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseCommandNode(t *testing.T) {
	input := "echo hello world;"
	prog, err := Parse(strings.NewReader(input))
	assert.Nil(t, err)
	assert.Len(t, prog.Body, 1)

	stmt, ok := prog.Body[0].(*CommandNode)
	if !ok {
		t.Fatalf("expected *CommandNode, got=%T", prog.Body[0])
	}
	assert.Len(t, stmt.List, 3)
	assert.Equal(t, "echo", stmt.List[0].Value)
	assert.Equal(t, "hello", stmt.List[1].Value)
	assert.Equal(t, "world", stmt.List[2].Value)
}

func TestParseIfNode(t *testing.T) {
	input := `if test 1;
		echo line1;
		echo line2
	end
	`
	prog, err := Parse(strings.NewReader(input))
	assert.Nil(t, err, "got err=%s", err)
	assert.Len(t, prog.Body, 1)

	stmt, ok := prog.Body[0].(*IfNode)
	if !ok {
		t.Fatalf("expected *IfNode, got=%T", prog.Body[0])
	}
	cond, ok := stmt.Cond.(*CommandNode)
	if !ok {
		t.Fatalf("expected *CommandNode, got=%T", cond)
	}
	assert.Len(t, cond.List, 2)
	assert.Equal(t, "test", cond.List[0].Value)
	assert.Equal(t, "1", cond.List[1].Value)

	assert.Len(t, stmt.Body.List, 2)
	stmt1 := stmt.Body.List[0].(*CommandNode)
	if !ok {
		t.Fatalf("expected *CommandNode, got=%T", stmt1)
	}
	assert.Len(t, stmt1.List, 2)
	assert.Equal(t, "echo", stmt1.List[0].Value)
	assert.Equal(t, "line1", stmt1.List[1].Value)
	stmt2 := stmt.Body.List[1].(*CommandNode)
	if !ok {
		t.Fatalf("expected *CommandNode, got=%T", stmt2)
	}
	assert.Len(t, stmt1.List, 2)
	assert.Equal(t, "echo", stmt2.List[0].Value)
	assert.Equal(t, "line2", stmt2.List[1].Value)

	assert.Nil(t, stmt.Else)
}

func TestParseIfNodeWithElse(t *testing.T) {
	input := `
	if test 1;
		echo if
	else
		echo else
	end
	`
	prog, err := Parse(strings.NewReader(input))
	assert.Nil(t, err, "got err=%s", err)
	assert.Len(t, prog.Body, 1)

	stmt, ok := prog.Body[0].(*IfNode)
	if !ok {
		t.Fatalf("expected *IfNode, got=%T", prog.Body[0])
	}
	cond, ok := stmt.Cond.(*CommandNode)
	if !ok {
		t.Fatalf("expected *CommandNode, got=%T", cond)
	}
	assert.Len(t, cond.List, 2)
	assert.Equal(t, "test", cond.List[0].Value)
	assert.Equal(t, "1", cond.List[1].Value)

	assert.Len(t, stmt.Body.List, 1)
	thenNode := stmt.Body.List[0].(*CommandNode)
	if !ok {
		t.Fatalf("expected *CommandNode, got=%T", thenNode)
	}
	assert.Len(t, thenNode.List, 2)
	assert.Equal(t, "echo", thenNode.List[0].Value)
	assert.Equal(t, "if", thenNode.List[1].Value)

	elseBlock := stmt.Else.(*BlockNode)
	if !ok {
		t.Fatalf("expected *BlockNode, got=%T", elseBlock)
	}
	assert.Len(t, elseBlock.List, 1)
	elseNode := elseBlock.List[0].(*CommandNode)
	assert.Len(t, elseNode.List, 2)
	assert.Equal(t, "echo", elseNode.List[0].Value)
	assert.Equal(t, "else", elseNode.List[1].Value)
}
