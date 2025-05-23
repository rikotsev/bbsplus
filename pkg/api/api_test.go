package api

import (
	"bbsplus/internal/definitions"
	"math/big"
	"testing"
)

func TestKeyGen(t *testing.T) {
	bbs := NewBbs(definitions.CreateBls12_381_Shake_256())

	keyMaterial := big.NewInt(0)
	keyInfo := big.NewInt(0)
	keyDst := big.NewInt(0)

	keyMaterial.SetString("746869732d49532d6a7573742d616e2d546573742d494b4d2d746f2d67656e65726174652d246528724074232d6b6579", 16)
	keyInfo.SetString("746869732d49532d736f6d652d6b65792d6d657461646174612d746f2d62652d757365642d696e2d746573742d6b65792d67656e", 16)
	keyDst.SetString("4242535f424c53313233383147315f584f463a5348414b452d3235365f535357555f524f5f4832475f484d32535f4b455947454e5f4453545f", 16)

	sk, err := bbs.KeyGen(
		keyMaterial.Bytes(),
		KeyGenOpts{
			KeyDst:  keyDst.Bytes(),
			KeyInfo: keyInfo.Bytes(),
		},
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
		t.Errorf("KeyGen() sk = %v, want %v", actualSk, expectedSk)
	}
}
