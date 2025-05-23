package definitions

import (
	"testing"
)

func TestBigIntFromHexString(t *testing.T) {
	actual := BigIntFromHexString("f")

	if !actual.IsInt64() || actual.Int64() != 15 {
		t.Errorf("BigIntFromHexString(\"f\") = %v; want 15", actual)
	}
}
