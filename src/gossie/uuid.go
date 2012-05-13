package gossie

import (
	"crypto/rand"
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"
)

/*
	to do:
		implement concurrent get with atomic increment of prev timestamp, like in Java
		check correctness of NewTimeUUID
		comparison
			comp func
			NewTimeUUIDMax and NewTimeUUIDMin for building slice comps
*/

type UUID [16]byte

var ZeroUUID UUID = [16]byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}
var LowestTimeUUID UUID = [16]byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x10, 0x00, 0x80, 0x80, 0x80, 0x80, 0x80, 0x80, 0x80, 0x80}
var HighestTimeUUID UUID = [16]byte{0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0x1f, 0xff, 0xbf, 0x7f, 0x7f, 0x7f, 0x7f, 0x7f, 0x7f, 0x7f}

func (value UUID) String() string {
	var r []string
	var s int
	for _, size := range [5]int{4, 2, 2, 2, 6} {
		var v int64
		for i := 0; i < size; i++ {
			v = v << 8
			v = v | int64(value[s+i])
		}
		r = append(r, fmt.Sprintf("%0*x", size*2, v))
		s += size
	}
	return strings.Join(r, "-")
}

func ParseUUID(value string) (UUID, error) {
	var r []byte
	var ru UUID

	if len(value) != 36 {
		return ZeroUUID, ErrorUnsupportedMarshaling
	}
	ints := strings.Split(value, "-")
	if len(ints) != 5 {
		return ZeroUUID, ErrorUnsupportedMarshaling
	}

	for i, size := range [5]int{4, 2, 2, 2, 6} {
		t, err := strconv.ParseInt(ints[i], 16, 64)
		if err != nil {
			return ZeroUUID, ErrorUnsupportedMarshaling
		}
		b, _ := marshalInt(t, size, BytesType)
		r = append(r, b...)
	}

	unmarshalUUID(r, BytesType, &ru)
	return ru, nil
}

func randomBase() ([]byte, error) {
	r := make([]byte, 16)
	_, err := io.ReadFull(rand.Reader, r)
	if err != nil {
		return nil, err
	}
	return r, nil
}

func NewRandomUUID() (UUID, error) {
	var ru UUID

	r, err := randomBase()
	if err != nil {
		return ZeroUUID, err
	}
	r[6] = (r[6] & 0x0F) | 0x40
	r[8] = (r[8] &^ 0x40) | 0x80

	unmarshalUUID(r, BytesType, &ru)
	return ru, nil
}

// http://johannburkard.de/software/uuid/

func NewTimeUUID(t time.Time) (UUID, error) {
	var ru UUID

	r, err := randomBase()
	if err != nil {
		return ZeroUUID, err
	}

	var hundredNSBlock int64 = (t.UnixNano() / 100) + 0x01B21DD213814000

	stamp := hundredNSBlock << 32
	stamp = stamp | ((hundredNSBlock & 0xFFFF00000000) >> 16)
	stamp = stamp | (0x1000 | ((hundredNSBlock >> 48) & 0x0FFF))

	r[0] = byte(stamp >> 56)
	r[1] = byte((stamp >> 48) & 0xff)
	r[2] = byte((stamp >> 40) & 0xff)
	r[3] = byte((stamp >> 32) & 0xff)
	r[4] = byte((stamp >> 24) & 0xff)
	r[5] = byte((stamp >> 16) & 0xff)
	r[6] = byte((stamp >> 8) & 0xff)
	r[7] = byte(stamp & 0xff)

	unmarshalUUID(r, BytesType, &ru)
	return ru, nil
}
