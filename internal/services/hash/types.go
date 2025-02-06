package hash

type Algorithm int

type HashStream []byte

func (s HashStream) String() string {
	return string(s)
}

func (s HashStream) Bytes() []byte {
	return []byte(s)
}

type HashService struct {
	algo Algorithm
	salt []byte
}
