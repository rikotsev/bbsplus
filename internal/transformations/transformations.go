package transformations

import (
	"errors"
	"math/big"
)

var ErrIntegerTooLarge = errors.New("integer too large. x < 256^xLen")
var ErrNegativeNumber = errors.New("integer is negative. x >= 0")
var ErrIllegalSize = errors.New("size should be at least 1. xLen > 0")

// Itosp Converts a non-negative integer to an octet string of a specified length
//
// https://www.rfc-editor.org/rfc/rfc8017.html#section-4.1
func Itosp(x *big.Int, xLen int) ([]byte, error) {
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

var ErrLenInBytesTooBig = errors.New("length of bytes is too big. lenInBytes <= 65535")
var ErrDstTooBig = errors.New("dst is too big. len(dst) <= 255")

// ExpandMessageXof produces a uniformly random byte string using an
// extendable-output function(XOF) H.
//
// https://www.rfc-editor.org/rfc/rfc9380.html#section-5.3.2
func ExpandMessageXof(msg []byte, dst []byte, lenInBytes int) ([]byte, error) {
	if lenInBytes > 65535 {
		return nil, ErrLenInBytesTooBig
	}

	if len(dst) > 255 {
		return nil, ErrDstTooBig
	}
	return nil, nil
}
