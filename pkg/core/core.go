package core

import (
	"bbsplus/internal/definitions"
	"bbsplus/internal/transformations"
	"fmt"
)

type BbsCore struct {
	cipherSuite *definitions.CipherSuite
}

func NewBbsCore(cipherSuite *definitions.CipherSuite) *BbsCore {
	return &BbsCore{
		cipherSuite: cipherSuite,
	}
}

// HashToScalar hash an arbitrary octet string to a scalar values in the multiplicative group of integers mod r
//
// https://www.ietf.org/archive/id/draft-irtf-cfrg-bbs-signatures-05.html#name-hash-to-scalar
func (bbsCore *BbsCore) HashToScalar(msg []byte, dst []byte) ([]byte, error) {
	uniformBytes, err := transformations.ExpandMessage(msg, dst, bbsCore.cipherSuite.ExpandLen, transformations.ExpandMessageOpts{
		ExtendableOutputFunction: bbsCore.cipherSuite.HashToCurveSuite.HashImplementationType,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to expand message: %w", err)
	}
	positiveInteger := transformations.OspToIp(uniformBytes)
	positiveInteger.Mod(positiveInteger, bbsCore.cipherSuite.HashToCurveSuite.Order)

	return positiveInteger.Bytes(), nil
}
