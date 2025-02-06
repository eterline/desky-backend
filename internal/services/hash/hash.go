package hash

import (
	"crypto/rand"
	"crypto/sha512"

	"golang.org/x/crypto/bcrypt"
)

const (
	_ Algorithm = iota
	SHA512
	Bcrypt
)

func NewHashService(algo Algorithm, salt []byte) *HashService {
	return &HashService{
		algo: algo,
		salt: salt,
	}
}

func (h *HashService) StringHash(value string) HashStream {

	v := []byte(value)

	switch h.algo {

	case Bcrypt:
		return hBcrypt(v, h.salt)

	case SHA512:
		return hSHA512(v, h.salt)

	default:
		return hSHA512(v, h.salt)
	}
}

func (h *HashService) EqStrings(source, stringFromRequest string) bool {
	return source == h.StringHash(stringFromRequest).String()
}

// hash funcs =====================

func hBcrypt(stream, salt []byte) []byte {
	stream = append(stream, salt...)

	stream, err := bcrypt.GenerateFromPassword(stream, bcrypt.DefaultCost)
	if err != nil {
		panic(err)
	}

	return stream
}

func hSHA512(stream, salt []byte) []byte {
	var sha512Hasher = sha512.New()
	stream = append(stream, salt...)
	sha512Hasher.Write(stream)
	return sha512Hasher.Sum(nil)
}

func GenerateCryptoSalt(size uint) ([]byte, error) {
	var salt = make([]byte, size)
	_, err := rand.Read(salt[:])
	if err != nil {
		return nil, err
	}
	return salt, nil
}
