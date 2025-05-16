package api

import (
	"math/big"
	"testing"
)

func TestKeyGen(t *testing.T) {

	sk, err := KeyGen(
		[]byte(""),
		KeyGenOpts{},
	)

	if err != nil {
		t.Errorf("KeyGen() error = %v", err)
		return
	}

	var expectedSk = new(big.Int)
	expectedSk.SetString("2eee0f60a8a3a8bec0ee942bfd46cbdae9a0738ee68f5a64e7238311cf09a079", 16)

	var actualSk = new(big.Int)
	actualSk.SetBytes(sk)

	if expectedSk.Cmp(actualSk) != 0 {
		t.Errorf("KeyGen() sk = %v, want %v", sk, expectedSk)
	}
}
