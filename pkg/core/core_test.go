package core

import (
	"bbsplus/internal/definitions"
	"bytes"
	"math/big"
	"testing"
)

// TestHashToScalar uses test vectors from the RFC
//
// For now only for SHAKE-256
// https://www.ietf.org/archive/id/draft-irtf-cfrg-bbs-signatures-05.html#name-hash-to-scalar-test-vectors
func TestHashToScalar(t *testing.T) {
	bbsCore := NewBbsCore(definitions.CreateBls12_381_Shake_256())
	msg, _ := big.NewInt(0).SetString("9872ad089e452c7b6e283dfac2a80d58e8d0ff71cc4d5e310a1debdda4a45f02", 16)
	dst, _ := big.NewInt(0).SetString("4242535f424c53313233383147315f584f463a5348414b452d3235365f535357555f524f5f4832475f484d32535f4832535f", 16)

	actual, err := bbsCore.HashToScalar(msg.Bytes(), dst.Bytes())
	if err != nil {
		t.Errorf("HashToScalar failed: %v", err)
	}

	expected, _ := big.NewInt(0).SetString("0500031f786fde5326aa9370dd7ffe9535ec7a52cf2b8f432cad5d9acfb73cd3", 16)

	if !bytes.Equal(expected.Bytes(), actual) {
		t.Errorf("HashToScalar failed: expected %s, actual %s", expected.Bytes(), actual)
	}
}
