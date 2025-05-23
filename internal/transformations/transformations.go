package transformations

import (
	"bbsplus/internal/definitions"
	"crypto/sha3"
	"errors"
	"fmt"
	"math/big"
)

var ErrIntegerTooLarge = errors.New("integer too large. x < 256^xLen")
var ErrNegativeNumber = errors.New("integer is negative. x >= 0")
var ErrIllegalSize = errors.New("size should be at least 1. xLen > 0")

// ItoOsp Converts a non-negative integer to an octet string of a specified length
//
// https://www.rfc-editor.org/rfc/rfc8017.html#section-4.1
func ItoOsp(x *big.Int, xLen int) ([]byte, error) {
	if xLen <= 0 {
		return nil, ErrIllegalSize
	}

	if x.Sign() < 0 {
		return nil, ErrNegativeNumber
	}

	if len(x.Bytes()) > xLen {
		return nil, ErrIntegerTooLarge
	}

	if len(x.Bytes()) == xLen {
		return x.Bytes(), nil
	}

	if len(x.Bytes()) > xLen {
		return x.Bytes()[:xLen], nil
	}

	result := make([]byte, xLen)
	copy(result[xLen-len(x.Bytes()):], x.Bytes())
	return result, nil
}

// OspToIp Converts an octet-string into a non-negative integer
//
// https://www.rfc-editor.org/rfc/rfc8017.html#section-4.2
func OspToIp(input []byte) *big.Int {
	return big.NewInt(0).SetBytes(input)
}

var ErrLenInBytesTooBig = errors.New("length of bytes is too big. lenInBytes <= 65535")
var ErrDstTooBig = errors.New("dst is too big. len(dst) <= 255")
var ErrUnsupportedHashType = errors.New("unsupported hash type")

type ExpandMessageOpts struct {
	ExtendableOutputFunction definitions.HashImplementationType
}

// ExpandMessage generates a uniformly random byte string
//
// https://www.rfc-editor.org/rfc/rfc9380.html#name-expand_message
func ExpandMessage(msg []byte, dst []byte, lenInBytes uint, opts ExpandMessageOpts) ([]byte, error) {
	if opts.ExtendableOutputFunction == definitions.HashShake256Implementation {
		return expandMessageXof(msg, dst, lenInBytes, opts)
	}

	return nil, ErrUnsupportedHashType
}

// expandMessageXof produces a uniformly random byte string using an
// extendable-output function(XOF) H.
//
// https://www.rfc-editor.org/rfc/rfc9380.html#section-5.3.2
func expandMessageXof(msg []byte, dst []byte, lenInBytes uint, opts ExpandMessageOpts) ([]byte, error) {
	if lenInBytes > 65535 {
		return nil, ErrLenInBytesTooBig
	}

	if len(dst) > 255 {
		return nil, ErrDstTooBig
	}

	dstOsp, err := ItoOsp(big.NewInt(int64(len(dst))), 1)
	if err != nil {
		return nil, fmt.Errorf("failed to create dstOsp: %w", err)
	}

	lenInBytesOsp, err := ItoOsp(big.NewInt(int64(lenInBytes)), 2)
	if err != nil {
		return nil, fmt.Errorf("failed to create lenInBytesOsp: %w", err)
	}

	dstPrime := append(dst, dstOsp...)
	msgPrime := append(msg, lenInBytesOsp...)
	msgPrime = append(msgPrime, dstPrime...)

	//TODO add support for other implementations
	if opts.ExtendableOutputFunction == definitions.HashShake256Implementation {
		shake256 := sha3.NewSHAKE256()
		_, err = shake256.Write(msgPrime)
		if err != nil {
			return nil, fmt.Errorf("failed to shake (write) msgPrime: %w", err)
		}

		result := make([]byte, lenInBytes)
		_, err = shake256.Read(result)
		if err != nil {
			return nil, fmt.Errorf("failed to shake (read) msgPrime: %w", err)
		}

		return result, nil
	}

	return nil, ErrUnsupportedHashType
}
