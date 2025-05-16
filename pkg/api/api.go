package api

import (
	"bbsplus/internal/transformations"
	"math/big"
)

type KeyGenOpts struct {
	KeyInfo []byte
	KeyDst  []byte
}

func KeyGen(keyMaterial []byte, opts KeyGenOpts) ([]byte, error) {
	if opts.KeyDst == nil {
		opts.KeyDst = []byte("BBS_BLS12381G1_XOF:SHAKE-256_SSWU_RO_H2G_HM2S_KEYGEN_DST_")
	}

	if opts.KeyInfo == nil {
		opts.KeyInfo = []byte{}
	}

	sizedOct, _ := transformations.ItoOsp(big.NewInt(int64(len(opts.KeyInfo))), 2)
	deriveInput := append(keyMaterial, sizedOct...)
	deriveInput = append(deriveInput, opts.KeyInfo...)

	return nil, nil
}
