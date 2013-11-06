package gossie

import (
	"sync"
)

//Stack implementation for connection pool
type lifo struct {
	l []*connection
	m sync.Mutex
}

func (q *lifo) Push(value *connection) {
	q.m.Lock()
	q.l = append(q.l, value)
	q.m.Unlock()
}

func (q *lifo) Pop() (*connection, bool) {
	q.m.Lock()
	if len(q.l) == 0 {
		q.m.Unlock()
		return nil, false
	}
	value := q.l[len(q.l)-1]
	q.l = q.l[0 : len(q.l)-1]
	q.m.Unlock()
	return value, true
}

//This function return the item from the bottom of the stack
//if size is more than n
//Should not be used very extensively because it creates garbage in memory
//We are using it for bleeder
func (q *lifo) PopBottom(n int) (*connection, bool) {
	q.m.Lock()
	if len(q.l) < n {
		q.m.Unlock()
		return nil, false
	}
	value := q.l[0]
	q.l = q.l[1:]
	q.m.Unlock()
	return value, true
}
