package mockgossie

import (
	"bytes"
	"testing"

	"github.com/kr/pretty"
	"github.com/stretchr/testify/assert"
	"github.com/wadey/gossie/src/gossie"
)

type BetweenStruct struct {
	Key       string `cf:"cf" mapping:"compact" key:"Key" cols:"Timestamp,UserID"`
	Timestamp int64
	UserID    string
}

var mappingBetweenStruct = gossie.MustNewMapping(&BetweenStruct{})

func makeColumn(t int64, u string) string {
	b1, _ := gossie.Marshal(t, gossie.LongType)
	b2, _ := gossie.Marshal(u, gossie.UTF8Type)
	var bb bytes.Buffer
	bb.Write(packComposite(b1, eocEquals))
	bb.Write(packComposite(b2, eocEquals))
	return bb.String()
}

func getAll(t *testing.T, q gossie.Query) []*BetweenStruct {
	result, err := q.Get("key1")
	if err != nil {
		t.Fatal(err)
	}
	rs := make([]*BetweenStruct, 0)
	i := &BetweenStruct{}
	for err = result.Next(i); err == nil; err = result.Next(i) {
		ix := *i
		rs = append(rs, &ix)
	}
	if err != gossie.Done {
		t.Fatal(err)
	}
	pretty.Log(rs)
	return rs
}

func int64Ptr(i int64) *int64 {
	return &i
}

func TestQueryBetween(t *testing.T) {
	m := NewMockConnectionPool()

	m.Load(Dump{
		"cf": {
			"key1": {
				makeColumn(100, "u1"):  []byte{},
				makeColumn(200, "u2"):  []byte{},
				makeColumn(300, "u3a"): []byte{},
				makeColumn(300, "u3b"): []byte{},
				makeColumn(400, "u4"):  []byte{},
			},
		},
	})

	q := m.Query(mappingBetweenStruct)

	assert.Equal(t, getAll(t, q), []*BetweenStruct{
		&BetweenStruct{"key1", 100, "u1"},
		&BetweenStruct{"key1", 200, "u2"},
		&BetweenStruct{"key1", 300, "u3a"},
		&BetweenStruct{"key1", 300, "u3b"},
		&BetweenStruct{"key1", 400, "u4"},
	})

	q.Between(int64Ptr(100), int64Ptr(300))
	assert.Equal(t, getAll(t, q), []*BetweenStruct{
		&BetweenStruct{"key1", 100, "u1"},
		&BetweenStruct{"key1", 200, "u2"},
	})

	q.Between(nil, int64Ptr(300))
	assert.Equal(t, getAll(t, q), []*BetweenStruct{
		&BetweenStruct{"key1", 100, "u1"},
		&BetweenStruct{"key1", 200, "u2"},
	})

	q.Between(int64Ptr(300), nil)
	assert.Equal(t, getAll(t, q), []*BetweenStruct{
		&BetweenStruct{"key1", 300, "u3a"},
		&BetweenStruct{"key1", 300, "u3b"},
		&BetweenStruct{"key1", 400, "u4"},
	})

	// TODO is it correct that this is start exclusive?
	q.Reversed(true)

	q.Between(nil, nil)
	assert.Equal(t, getAll(t, q), []*BetweenStruct{
		&BetweenStruct{"key1", 400, "u4"},
		&BetweenStruct{"key1", 300, "u3b"},
		&BetweenStruct{"key1", 300, "u3a"},
		&BetweenStruct{"key1", 200, "u2"},
		&BetweenStruct{"key1", 100, "u1"},
	})

	q.Between(int64Ptr(300), nil)
	assert.Equal(t, getAll(t, q), []*BetweenStruct{
		&BetweenStruct{"key1", 200, "u2"},
		&BetweenStruct{"key1", 100, "u1"},
	})

	q.Between(int64Ptr(300), int64Ptr(200))
	assert.Equal(t, getAll(t, q), []*BetweenStruct{
		&BetweenStruct{"key1", 200, "u2"},
	})

	q.Between(nil, int64Ptr(200))
	assert.Equal(t, getAll(t, q), []*BetweenStruct{
		&BetweenStruct{"key1", 400, "u4"},
		&BetweenStruct{"key1", 300, "u3b"},
		&BetweenStruct{"key1", 300, "u3a"},
		&BetweenStruct{"key1", 200, "u2"},
	})

	q.Limit(2, 1) // limit with reverse

	q.Between(int64Ptr(400), nil)
	assert.Equal(t, getAll(t, q), []*BetweenStruct{
		&BetweenStruct{"key1", 300, "u3b"},
		&BetweenStruct{"key1", 300, "u3a"},
	})

	q.Reversed(false) // limit no reverse

	q.Between(int64Ptr(300), nil)
	assert.Equal(t, getAll(t, q), []*BetweenStruct{
		&BetweenStruct{"key1", 300, "u3a"},
		&BetweenStruct{"key1", 300, "u3b"},
	})
}
