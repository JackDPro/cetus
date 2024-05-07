package provider

import "golang.org/x/crypto/bcrypt"

func HashMake(str string) (hashedStr string, err error) {
	var hashed []byte
	hashed, err = bcrypt.GenerateFromPassword([]byte(str), bcrypt.DefaultCost)
	if err != nil {
		return
	}
	hashedStr = string(hashed)
	return
}

func HashCheck(plainStr, hashedStr string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hashedStr), []byte(plainStr))
	if err != nil {
		return err
	}
	return nil
}

func HashNeedRefresh(hashedStr string) bool {
	hashCost, err := bcrypt.Cost([]byte(hashedStr))
	return err != nil || hashCost != bcrypt.DefaultCost
}
