package seqx

// integer constrains the integer types Range can count over.
type integer interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 |
		~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr
}

// Number constrains the types Sum can total: any integer or float.
type Number interface {
	integer | ~float32 | ~float64
}

// Ordered constrains the types Min and Max can compare.
type Ordered interface {
	Number | ~string
}
