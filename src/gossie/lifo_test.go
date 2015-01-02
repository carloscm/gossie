package gossie

import (
	"testing"

	"github.com/stretchr/testify/assert"
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

func TestPushPop2(t *testing.T) {
	var l lifo
	conn := &connection{}
	conn2 := &connection{}
	l.Push(conn)
	l.Push(&connection{})
	l.Push(conn2)

	cb, okb := l.PopBottom(0)
	assert.NotNil(t, cb)
	assert.True(t, okb)
	assert.Equal(t, cb, conn)

	c, ok := l.Pop()
	assert.NotNil(t, c)
	assert.True(t, ok)
	assert.Equal(t, c, conn2)

	c2, ok2 := l.Pop()
	assert.NotNil(t, c2)
	assert.True(t, ok2)

	assert.Equal(t, 0, len(l.l))
}

func TestKeepN(t *testing.T) {
	var l lifo
	conn := &connection{}
	conn2 := &connection{}
	l.Push(conn)
	l.Push(&connection{})
	l.Push(conn2)

	cb, okb := l.PopBottom(5)
	assert.Nil(t, cb)
	assert.False(t, okb)

	c, ok := l.Pop()
	assert.NotNil(t, c)
	assert.True(t, ok)
	assert.Equal(t, c, conn2)

	assert.Equal(t, 2, len(l.l))
}
