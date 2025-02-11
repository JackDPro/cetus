package test

import (
	"crypto/rand"
	"github.com/pjebs/optimus-go"
	"math/big"
	"testing"
)

//http://primes.utm.edu/lists/small/millions/

func RandN(N int64) uint64 {
	b49 := *big.NewInt(N)
	n, _ := rand.Int(rand.Reader, &b49)
	in := n.Uint64() + 1
	return in
}

func GenerateSeed() (optimus.Optimus, error) {
	selectedPrime := uint64(104393867)
	modInverse := optimus.ModInverse(selectedPrime)
	random := RandN(int64(optimus.MAX_INT - 1))
	o := optimus.New(selectedPrime, modInverse, random)
	return o, nil
}

func TestGenerateId(t *testing.T) {
	o, _ := GenerateSeed()
	t.Logf("\nOPTIMUS_PRIME=%v\nOPTIMUS_INVERSE=%v\nOPTIMUS_RANDOM=%v\n", o.Prime(), o.ModInverse(), o.Random())
	encodeNum := o.Encode(1)
	t.Logf("encodeId=%v\n", encodeNum)
	decodeNum := o.Decode(encodeNum)
	t.Logf("decodeId=%v\n", decodeNum)
}
