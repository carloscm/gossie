package gossie

import (
	"testing"
)

func BenchmarkReasonableTwoMapping(b *testing.B) {
	m2, err := NewMapping(&ReasonableTwo{})
	if err != nil {
		b.Fatal("Error building mapping:", err)
	}
	r100 := &ReasonableTwo{
		Username: "batchuser1",
		TweetID:  int64(1),
		Version:  int64(2),
		Lat:      1.00002,
		Lon:      -38.11,
		Body:     "hey this thing appears to work, nice!",
	}

	var fakecnt int

	for i := 0; i < b.N; i++ {
		row, err := m2.Map(r100)
		if err != nil {
			b.Fatal("Error marshalling:", err)
		}
		fakecnt += len(row.Columns)
	}
}
