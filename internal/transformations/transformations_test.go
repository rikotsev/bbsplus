package transformations

import (
	"bbsplus/internal/definitions"
	"bytes"
	"errors"
	"math/big"
	"testing"
)

func TestItoOspErrors(t *testing.T) {

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
		_, err := ItoOsp(testCase.x, testCase.xLen)

		if !errors.Is(err, testCase.expected) {
			t.Errorf("ItoOsp(%d, %d) = %v, want %v", testCase.x, testCase.xLen, err, testCase.expected)
		}
	}
}

func pad(input []byte, size int) []byte {
	tmp := make([]byte, size)
	copy(tmp[size-len(input):], input)

	return tmp
}

func TestIToOsp(t *testing.T) {
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
		actual, err := ItoOsp(testCase.x, testCase.xLen)
		if err != nil {
			t.Errorf("ItoOsp(%d, %d) produced unexpected error %v", testCase.x, testCase.xLen, err)
			continue
		}

		if !bytes.Equal(actual, testCase.expected) {
			t.Errorf("ItoOsp(%d, %d) = %v, want %v", testCase.x, testCase.xLen, actual, testCase.expected)
		}
	}
}

func TestExpandMessageXofErrors(t *testing.T) {
	testCases := []struct {
		msg        []byte
		dst        []byte
		lenInBytes uint
		opts       ExpandMessageOpts
		expected   error
	}{
		{
			msg:        make([]byte, 0),
			dst:        make([]byte, 0),
			lenInBytes: 65536,
			opts: ExpandMessageOpts{
				ExtendableOutputFunction: definitions.HashShake256Implementation,
			},
			expected: ErrLenInBytesTooBig,
		},
		{
			msg:        make([]byte, 0),
			dst:        make([]byte, 256),
			lenInBytes: 100,
			opts: ExpandMessageOpts{
				ExtendableOutputFunction: definitions.HashShake256Implementation,
			},
			expected: ErrDstTooBig,
		},
		{
			msg:        make([]byte, 0),
			dst:        make([]byte, 255),
			lenInBytes: 100,
			opts: ExpandMessageOpts{
				ExtendableOutputFunction: definitions.HashImplementationType("TEST"),
			},
			expected: ErrUnsupportedHashType,
		},
	}

	for _, testCase := range testCases {
		_, err := expandMessageXof(testCase.msg, testCase.dst, testCase.lenInBytes, testCase.opts)
		if !errors.Is(err, testCase.expected) {
			t.Errorf("expandMessageXof(%v, %d, %d) produced unexpected error %v", testCase.msg, testCase.dst, testCase.lenInBytes, err)
		}
	}
}

// TestExpandMessageXof uses the test vectors in the RFC.
// For the time being only the SHAKE-256 hash implementation
//
// https://www.rfc-editor.org/rfc/rfc9380.html#name-expand_message_xofshake256
func TestExpandMessageXof(t *testing.T) {
	dstConst := []byte("QUUX-V01-CS02-with-expander-SHAKE256")
	testCases := []struct {
		msg        []byte
		dst        []byte
		lenInBytes uint
		opts       ExpandMessageOpts
		expected   []byte
	}{
		{
			msg:        []byte(""),
			dst:        dstConst,
			lenInBytes: 32,
			opts: ExpandMessageOpts{
				ExtendableOutputFunction: definitions.HashShake256Implementation,
			},
			expected: definitions.BigIntFromHexString("2ffc05c48ed32b95d72e807f6eab9f7530dd1c2f013914c8fed38c5ccc15ad76").Bytes(),
		},
		{
			msg:        []byte("abc"),
			dst:        dstConst,
			lenInBytes: 32,
			opts: ExpandMessageOpts{
				ExtendableOutputFunction: definitions.HashShake256Implementation,
			},
			expected: definitions.BigIntFromHexString("b39e493867e2767216792abce1f2676c197c0692aed061560ead251821808e07").Bytes(),
		},
		{
			msg:        []byte("abcdef0123456789"),
			dst:        dstConst,
			lenInBytes: 32,
			opts: ExpandMessageOpts{
				ExtendableOutputFunction: definitions.HashShake256Implementation,
			},
			expected: definitions.BigIntFromHexString("245389cf44a13f0e70af8665fe5337ec2dcd138890bb7901c4ad9cfceb054b65").Bytes(),
		},
		{
			msg:        []byte("q128_qqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqq"),
			dst:        dstConst,
			lenInBytes: 32,
			opts: ExpandMessageOpts{
				ExtendableOutputFunction: definitions.HashShake256Implementation,
			},
			expected: definitions.BigIntFromHexString("719b3911821e6428a5ed9b8e600f2866bcf23c8f0515e52d6c6c019a03f16f0e").Bytes(),
		},
		{
			msg:        []byte("a512_aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"),
			dst:        dstConst,
			lenInBytes: 32,
			opts: ExpandMessageOpts{
				ExtendableOutputFunction: definitions.HashShake256Implementation,
			},
			expected: definitions.BigIntFromHexString("9181ead5220b1963f1b5951f35547a5ea86a820562287d6ca4723633d17ccbbc").Bytes(),
		},
		{
			msg:        []byte(""),
			dst:        dstConst,
			lenInBytes: 128,
			opts: ExpandMessageOpts{
				ExtendableOutputFunction: definitions.HashShake256Implementation,
			},
			expected: definitions.BigIntFromHexString("7a1361d2d7d82d79e035b8880c5a3c86c5afa719478c007d96e6c88737a3f631dd74a2c88df79a4cb5e5d9f7504957c70d669ec6bfedc31e01e2bacc4ff3fdf9b6a00b17cc18d9d72ace7d6b81c2e481b4f73f34f9a7505dccbe8f5485f3d20c5409b0310093d5d6492dea4e18aa6979c23c8ea5de01582e9689612afbb353df").Bytes(),
		},
		{
			msg:        []byte("abc"),
			dst:        dstConst,
			lenInBytes: 128,
			opts: ExpandMessageOpts{
				ExtendableOutputFunction: definitions.HashShake256Implementation,
			},
			expected: definitions.BigIntFromHexString("a54303e6b172909783353ab05ef08dd435a558c3197db0c132134649708e0b9b4e34fb99b92a9e9e28fc1f1d8860d85897a8e021e6382f3eea10577f968ff6df6c45fe624ce65ca25932f679a42a404bc3681efe03fcd45ef73bb3a8f79ba784f80f55ea8a3c367408f30381299617f50c8cf8fbb21d0f1e1d70b0131a7b6fbe").Bytes(),
		},
		{
			msg:        []byte("abcdef0123456789"),
			dst:        dstConst,
			lenInBytes: 128,
			opts: ExpandMessageOpts{
				ExtendableOutputFunction: definitions.HashShake256Implementation,
			},
			expected: definitions.BigIntFromHexString("e42e4d9538a189316e3154b821c1bafb390f78b2f010ea404e6ac063deb8c0852fcd412e098e231e43427bd2be1330bb47b4039ad57b30ae1fc94e34993b162ff4d695e42d59d9777ea18d3848d9d336c25d2acb93adcad009bcfb9cde12286df267ada283063de0bb1505565b2eb6c90e31c48798ecdc71a71756a9110ff373").Bytes(),
		},
		{
			msg:        []byte("q128_qqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqq"),
			dst:        dstConst,
			lenInBytes: 128,
			opts: ExpandMessageOpts{
				ExtendableOutputFunction: definitions.HashShake256Implementation,
			},
			expected: definitions.BigIntFromHexString("4ac054dda0a38a65d0ecf7afd3c2812300027c8789655e47aecf1ecc1a2426b17444c7482c99e5907afd9c25b991990490bb9c686f43e79b4471a23a703d4b02f23c669737a886a7ec28bddb92c3a98de63ebf878aa363a501a60055c048bea11840c4717beae7eee28c3cfa42857b3d130188571943a7bd747de831bd6444e0").Bytes(),
		},
		{
			msg:        []byte("a512_aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"),
			dst:        dstConst,
			lenInBytes: 128,
			opts: ExpandMessageOpts{
				ExtendableOutputFunction: definitions.HashShake256Implementation,
			},
			expected: definitions.BigIntFromHexString("09afc76d51c2cccbc129c2315df66c2be7295a231203b8ab2dd7f95c2772c68e500bc72e20c602abc9964663b7a03a389be128c56971ce81001a0b875e7fd17822db9d69792ddf6a23a151bf470079c518279aef3e75611f8f828994a9988f4a8a256ddb8bae161e658d5a2a09bcfe839c6396dc06ee5c8ff3c22d3b1f9deb7e").Bytes(),
		},
	}

	for _, testCase := range testCases {
		actual, err := expandMessageXof(testCase.msg, testCase.dst, testCase.lenInBytes, testCase.opts)
		if err != nil {
			t.Errorf("expandMessageXof(%v, %d, %d, %v) produces unexpected error %v", testCase.msg, testCase.dst, testCase.lenInBytes, testCase.opts, err)
		}

		if !bytes.Equal(actual, testCase.expected) {
			t.Errorf("expandMessageXof(%v, %d, %d, %v) = %v; want %v", testCase.msg, testCase.dst, testCase.lenInBytes, testCase.opts, actual, testCase.expected)
		}
	}

}
