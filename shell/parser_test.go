package shell

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseCommandStmt(t *testing.T) {
	input := "echo hello world;"
	prog, err := Parse(strings.NewReader(input))
	assert.Nil(t, err)
	assert.Len(t, prog.Body, 1)

	stmt, ok := prog.Body[0].(*CommandStmt)
	if !ok {
		t.Fatalf("expected *CommandStmt, got=%T", prog.Body[0])
	}
	assert.Len(t, stmt.List, 3)
	assert.Equal(t, "echo", stmt.List[0].Value)
	assert.Equal(t, "hello", stmt.List[1].Value)
	assert.Equal(t, "world", stmt.List[2].Value)
}

func TestParseIfStmt(t *testing.T) {
	input := `if test 1;
		echo line1;
		echo line2
	end
	`
	prog, err := Parse(strings.NewReader(input))
	assert.Nil(t, err, "got err=%s", err)
	assert.Len(t, prog.Body, 1)

	stmt, ok := prog.Body[0].(*IfStmt)
	if !ok {
		t.Fatalf("expected *IfStmt, got=%T", prog.Body[0])
	}
	cond, ok := stmt.Cond.(*CommandStmt)
	if !ok {
		t.Fatalf("expected *CommandStmt, got=%T", cond)
	}
	assert.Len(t, cond.List, 2)
	assert.Equal(t, "test", cond.List[0].Value)
	assert.Equal(t, "1", cond.List[1].Value)

	assert.Len(t, stmt.Body.List, 2)
	stmt1 := stmt.Body.List[0].(*CommandStmt)
	if !ok {
		t.Fatalf("expected *CommandStmt, got=%T", stmt1)
	}
	assert.Len(t, stmt1.List, 2)
	assert.Equal(t, "echo", stmt1.List[0].Value)
	assert.Equal(t, "line1", stmt1.List[1].Value)
	stmt2 := stmt.Body.List[1].(*CommandStmt)
	if !ok {
		t.Fatalf("expected *CommandStmt, got=%T", stmt2)
	}
	assert.Len(t, stmt1.List, 2)
	assert.Equal(t, "echo", stmt2.List[0].Value)
	assert.Equal(t, "line2", stmt2.List[1].Value)

	assert.Nil(t, stmt.Else)
}

func TestParseIfStmtWithElse(t *testing.T) {
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

	stmt, ok := prog.Body[0].(*IfStmt)
	if !ok {
		t.Fatalf("expected *IfStmt, got=%T", prog.Body[0])
	}
	cond, ok := stmt.Cond.(*CommandStmt)
	if !ok {
		t.Fatalf("expected *CommandStmt, got=%T", cond)
	}
	assert.Len(t, cond.List, 2)
	assert.Equal(t, "test", cond.List[0].Value)
	assert.Equal(t, "1", cond.List[1].Value)

	assert.Len(t, stmt.Body.List, 1)
	thenStmt := stmt.Body.List[0].(*CommandStmt)
	if !ok {
		t.Fatalf("expected *CommandStmt, got=%T", thenStmt)
	}
	assert.Len(t, thenStmt.List, 2)
	assert.Equal(t, "echo", thenStmt.List[0].Value)
	assert.Equal(t, "if", thenStmt.List[1].Value)

	elseBlock := stmt.Else.(*BlockStmt)
	if !ok {
		t.Fatalf("expected *BlockStmt, got=%T", elseBlock)
	}
	assert.Len(t, elseBlock.List, 1)
	elseStmt := elseBlock.List[0].(*CommandStmt)
	assert.Len(t, elseStmt.List, 2)
	assert.Equal(t, "echo", elseStmt.List[0].Value)
	assert.Equal(t, "else", elseStmt.List[1].Value)
}
