package transformations

import (
	"bytes"
	"errors"
	"math/big"
	"testing"
)

func TestItospErrors(t *testing.T) {

	testCases := []struct {
		x        *big.Int
		xLen     int
		expected error
	}{
		{
			x:        big.NewInt(-1),
			xLen:     100,
			expected: ErrNegativeNumber,
		},
		{
			x:        big.NewInt(257),
			xLen:     1,
			expected: ErrIntegerTooLarge,
		},
		{
			x:        big.NewInt(1),
			xLen:     -1,
			expected: ErrIllegalSize,
		},
	}

	for _, testCase := range testCases {
		_, err := Itosp(testCase.x, testCase.xLen)

		if !errors.Is(err, testCase.expected) {
			t.Errorf("Itosp(%d, %d) = %v, want %v", testCase.x, testCase.xLen, err, testCase.expected)
		}
	}
}

func pad(input []byte, size int) []byte {
	tmp := make([]byte, size)
	copy(tmp[size-len(input):], input)

	return tmp
}

func TestItosp(t *testing.T) {
	testCases := []struct {
		x        *big.Int
		xLen     int
		expected []byte
	}{
		{
			x:        big.NewInt(52000),
			xLen:     2,
			expected: big.NewInt(52000).Bytes()[:2],
		},
		{
			x:        big.NewInt(10),
			xLen:     2,
			expected: pad(big.NewInt(10).Bytes(), 2),
		},
		{
			x:        big.NewInt(65535),
			xLen:     2,
			expected: pad(big.NewInt(65535).Bytes(), 2),
		},
	}

	for _, testCase := range testCases {
		actual, err := Itosp(testCase.x, testCase.xLen)
		if err != nil {
			t.Errorf("Itosp(%d, %d) produced unexpected error %v", testCase.x, testCase.xLen, err)
			continue
		}

		if !bytes.Equal(actual, testCase.expected) {
			t.Errorf("Itosp(%d, %d) = %v, want %v", testCase.x, testCase.xLen, actual, testCase.expected)
		}
	}
}

func TestExpandMessageXofErrors(t *testing.T) {
	testCases := []struct {
		msg        []byte
		dst        []byte
		lenInBytes int
		expected   error
	}{
		{
			msg:        make([]byte, 0),
			dst:        make([]byte, 0),
			lenInBytes: 65536,
			expected:   ErrLenInBytesTooBig,
		},
		{
			msg:        make([]byte, 0),
			dst:        make([]byte, 256),
			lenInBytes: 100,
			expected:   ErrDstTooBig,
		},
	}

	for _, testCase := range testCases {
		_, err := ExpandMessageXof(testCase.msg, testCase.dst, testCase.lenInBytes)
		if !errors.Is(err, testCase.expected) {
			t.Errorf("ExpandMessageXof(%v, %d, %d) produced unexpected error %v", testCase.msg, testCase.dst, testCase.lenInBytes, err)
		}
	}
}
