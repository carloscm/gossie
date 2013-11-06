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

func Testpushpop2(t *testing.T) {
	var l lifo
	conn := &connection{}
	conn2 := &connection{}
	l.Push(conn)
	l.Push(&connection{})
	l.Push(conn2)

	cb, okb := l.PopBottom()
	assert.NotNil(t, cb)
	assert.True(t, okb)
	assert.Equal(t, cb, conn)

	c, ok := l.Pop()
	assert.NotNil(t, c)
	assert.True(t, ok)
	assert.Equal(t, c, conn2)
	c2, ok2 := l.Pop()
	assert.Nil(t, c2)
	assert.False(t, ok2)

	assert.Equal(t, 1, len(l.l))
}
