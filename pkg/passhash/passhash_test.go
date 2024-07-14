package passhash

import (
	"log"
	"testing"
)

func TestStringString(t *testing.T) {
	plainText := "This is a test."

	hash, err := HashString(plainText)

	if err != nil {
		t.Error(err)
	}

	if !MatchString(hash, plainText) {
		t.Error("Password does not match")
	}
}

func TestByteByte(t *testing.T) {
	plainText := []byte("This is a test.")

	hash, err := HashBytes(plainText)

	if err != nil {
		t.Error(err)
	}

	if !MatchBytes(hash, plainText) {
		t.Error("Password does not match")
	}
}

// go test -run TestStringByte -v
func TestStringByte(t *testing.T) {
	plainText := "123"

	hash, err := HashString(plainText)

	if err != nil {
		t.Error(err)
	}

	log.Println("plainText: ", plainText)
	log.Println("hash: ", hash)

	if !MatchString("$2a$10$HWTZXG6kpZTEGPAzxbcuKewnGZiPh6VbGD1GCRUojgOQVCSO.8O.i", plainText) {
		t.Error("Password does not match")
	}
}

func TestByteString(t *testing.T) {
	plainText := []byte("This is a test.")

	hash, err := HashBytes(plainText)

	if err != nil {
		t.Error(err)
	}

	if !MatchString(string(hash), string(plainText)) {
		t.Error("Password does not match")
	}
}
