package gossie

func NewInt32(i int) *int32 {
	i32 := int32(i)
	return &i32
}

func NewInt64(i int64) *int64 {
	i64 := i
	return &i64
}

func NewBytes(in []byte) *[]byte {
	return &in
}
