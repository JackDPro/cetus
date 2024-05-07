package config

import (
	"os"
	"strconv"
	"sync"
)

type HashId struct {
	Prime   uint64
	Inverse uint64
	Random  uint64
}

var hashIdConfigInstance *HashId
var hashIdConfigOnce sync.Once

func GetHashIdConfig() *HashId {
	hashIdConfigOnce.Do(func() {
		prime, err := strconv.ParseUint(os.Getenv("OPTIMUS_PRIME"), 10, 64)
		if err != nil {
			panic("HashId prime is invalid")
		}
		inverse, err := strconv.ParseUint(os.Getenv("OPTIMUS_INVERSE"), 10, 64)
		if err != nil {
			panic("HashId inverse is invalid")
		}
		random, err := strconv.ParseUint(os.Getenv("OPTIMUS_RANDOM"), 10, 64)
		if err != nil {
			panic("HashId random is invalid")
		}

		hashIdConfigInstance = &HashId{
			Prime:   prime,
			Inverse: inverse,
			Random:  random,
		}
	})
	return hashIdConfigInstance
}
