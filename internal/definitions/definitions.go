package definitions

import "math/big"

type ExpandMessageImplementationType string
type HashImplementationType string

const ExpandMessageXofImplementation ExpandMessageImplementationType = "expand_message_xof"
const HashShake256Implementation HashImplementationType = "SHAKE-256"

type CipherSuite struct {
	Id                string
	OctetScalarLength uint
	OctetPointLength  uint
	ExpandLen         uint
	HashToCurveSuite  *HashToCurveSuite
}

type HashToCurveSuite struct {
	Id                          string
	Prime                       *big.Int
	Order                       *big.Int
	SecurityLevel               uint
	ExpandMessageImplementation ExpandMessageImplementationType
	HashImplementationType      HashImplementationType
}

func CreateBls12_381_Shake_256() *CipherSuite {
	return &CipherSuite{
		Id:                "BBS_BLS12381G1_XOF:SHAKE-256_SSWU_RO_",
		OctetScalarLength: 32,
		OctetPointLength:  48,
		ExpandLen:         48,
		HashToCurveSuite: &HashToCurveSuite{
			Id:                          "BLS12-381_SHAKE-256",
			Prime:                       BigIntFromHexString("1a0111ea397fe69a4b1ba7b6434bacd764774b84f38512bf6730d2a0f6b0f6241eabfffeb153ffffb9feffffffffaaab"),
			Order:                       BigIntFromHexString("73eda753299d7d483339d80809a1d80553bda402fffe5bfeffffffff00000001"),
			SecurityLevel:               128,
			ExpandMessageImplementation: ExpandMessageXofImplementation,
			HashImplementationType:      HashShake256Implementation,
		},
	}
}

func BigIntFromHexString(s string) *big.Int {
	newInt := big.NewInt(0)

	newInt.SetString(s, 16)

	return newInt
}
