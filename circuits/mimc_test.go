package circuits

import (
	"fmt"
	"github.com/consensys/gnark-crypto/ecc"
	"github.com/consensys/gnark/backend"
	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gnark/frontend/cs/r1cs"
	"github.com/consensys/gnark/test"
	"github.com/iden3/go-iden3-crypto/mimc7"
	"math/big"
	"testing"
)

type mimcTest struct {
	A    frontend.Variable `gnark:",public"`
	Hash frontend.Variable
}

func (t *mimcTest) Define(api frontend.API) error {
	hash := MiMC7(api, 90, t.A, 5)
	api.AssertIsEqual(hash, t.Hash)
	return nil
}

func Test_MiMC(t *testing.T) {
	assert := test.NewAssert(t)
	var circuit mimcTest

	a := big.NewInt(1)
	hash := mimc7.MIMC7HashGeneric(a, big.NewInt(5), 90)

	t.Logf("hash: %v", hash)

	assert.ProverSucceeded(&circuit, &mimcTest{
		A:    a,
		Hash: hash,
	}, test.WithCurves(ecc.BN254), test.WithBackends(backend.GROTH16))

	_r1cs, _ := frontend.Compile(ecc.BN254, r1cs.NewBuilder, &circuit)
	internal, secret, public := _r1cs.GetNbVariables()
	fmt.Printf("public, secret, internal %v, %v, %v\n", public, secret, internal)
}
