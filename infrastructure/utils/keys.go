package utils

import (
	"math/rand"
	"time"
)

const (
	DIGITS                   = "0123456789"
	SPECIALS_CHARACTER       = "~=+%^*/()[]{}/!@#$?|<>"
	ALL_CHARACTER_SPECIAL    = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz" + DIGITS + SPECIALS_CHARACTER
	ALL_CHARACTER_UNSPECIALS = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz" + DIGITS
)

//GenerateSymmetricKey creates symmetric and secure key
func GenerateSymmetricKey(length int) []byte {
	rand.Seed(time.Now().UnixNano())
	digits := DIGITS
	specials := SPECIALS_CHARACTER
	all := ALL_CHARACTER_SPECIAL
	buf := make([]byte, length)
	buf[0] = digits[rand.Intn(len(digits))]
	buf[1] = specials[rand.Intn(len(specials))]
	for i := 2; i < length; i++ {
		buf[i] = all[rand.Intn(len(all))]
	}
	rand.Shuffle(len(buf), func(i, j int) {
		buf[i], buf[j] = buf[j], buf[i]
	})
	str := string(buf)

	return []byte(str)
}

//GenerateSymmetricKey create symmetric and secure key without spicials
func GenerateSymmetricKeyU(length int) []byte {
	rand.Seed(time.Now().UnixNano())
	digits := DIGITS
	all := ALL_CHARACTER_UNSPECIALS
	buf := make([]byte, length)
	buf[0] = digits[rand.Intn(len(digits))]
	for i := 1; i < length; i++ {
		buf[i] = all[rand.Intn(len(all))]
	}
	rand.Shuffle(len(buf), func(i, j int) {
		buf[i], buf[j] = buf[j], buf[i]
	})
	str := string(buf)

	return []byte(str)
}
