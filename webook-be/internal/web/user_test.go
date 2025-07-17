package web_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
)

func TestEncrypt(t *testing.T) {
	password := "testingPwd@123"
	// 密码加密
	encrypt, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		t.Fatal(err)
	}
	// 密码校验
	err = bcrypt.CompareHashAndPassword(encrypt, []byte(password))
	assert.NoError(t, err)
}
