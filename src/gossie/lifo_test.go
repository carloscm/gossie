package gossie

import (
	"github.com/stretchrcom/testify/assert"
	"testing"
)

func TestInitEmpty(t *testing.T) {
	var l lifo
	c, ok := l.Pop()
	assert.Nil(t, c)
	assert.False(t, ok)
}

func Testpushpop(t *testing.T) {
	var l lifo
	conn := &connection{}
	l.Push(conn)
	c, ok := l.Pop()
	assert.NotNil(t, c)
	assert.True(t, ok)
	c2, ok2 := l.Pop()
	assert.Nil(t, c2)
	assert.False(t, ok2)
}
