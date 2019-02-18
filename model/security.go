package model

import (
	"crypto/hmac"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/hex"
	"io"

	log "github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

func hashAndSalt(pwd []byte) string {
	// Use GenerateFromPassword to hash & salt pwd
	// MinCost is just an integer constant provided by the bcrypt
	// package along with DefaultCost & MaxCost.
	// The cost can be any value you want provided it isn't lower
	// than the MinCost (4)
	hash, err := bcrypt.GenerateFromPassword(pwd, bcrypt.MinCost)
	if err != nil {
		log.WithError(err).Error("bcrypt.GenerateFromPassword")
	}
	// GenerateFromPassword returns a byte slice so we need to
	// convert the bytes to a string and return it
	return string(hash)
}

func comparePasswords(hashedPwd string, plainPwd []byte) bool {
	// Since we'll be getting the hashed password from the DB it
	// will be a string so we'll need to convert it to a byte slice
	byteHash := []byte(hashedPwd)
	err := bcrypt.CompareHashAndPassword(byteHash, plainPwd)
	if err != nil {
		log.WithError(err).Error("bcrypt.CompareHashAndPassword")
		return false
	}
	return true
}

func encryptPassword(pwd string) string {
	key := []byte("zouweicheng@gmail.com")
	mac := hmac.New(sha256.New, key)
	mac.Write([]byte(pwd))
	h := sha1.New()
	io.WriteString(h, hex.EncodeToString(mac.Sum(nil)))
	return hex.EncodeToString(h.Sum(nil))
}
