package api

import (
	"bbsplus/internal/definitions"
	"bbsplus/internal/transformations"
	"bbsplus/pkg/core"
	"fmt"
	"math/big"
)

type KeyGenOpts struct {
	KeyInfo []byte
	KeyDst  []byte
}

type Bbs struct {
	CipherSuite *definitions.CipherSuite
	bbsCore     *core.BbsCore
}

func NewBbs(cipherSuite *definitions.CipherSuite) *Bbs {
	return &Bbs{
		CipherSuite: cipherSuite,
		bbsCore:     core.NewBbsCore(cipherSuite),
	}
}

// KeyGen generates an SK (secret key)
//
// https://www.ietf.org/archive/id/draft-irtf-cfrg-bbs-signatures-05.html#name-secret-key
func (bbs *Bbs) KeyGen(keyMaterial []byte, opts KeyGenOpts) ([]byte, error) {
	//TODO maybe error out
	if opts.KeyDst == nil {
		opts.KeyDst = []byte{}
	}

	if opts.KeyInfo == nil {
		opts.KeyInfo = []byte{}
	}

	sizedOct, _ := transformations.ItoOsp(big.NewInt(int64(len(opts.KeyInfo))), 2)
	deriveInput := append(keyMaterial, sizedOct...)
	deriveInput = append(deriveInput, opts.KeyInfo...)

	secretKey, err := bbs.bbsCore.HashToScalar(deriveInput, opts.KeyDst)
	if err != nil {
		return nil, fmt.Errorf("failed to hash to scalar: %w", err)
	}

	return secretKey, nil
}
