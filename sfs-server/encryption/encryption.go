package encryption

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha256"
	"errors"
	"fmt"
	"io/ioutil"
	"net/url"
	"os"
	"strings"
)

const Key = "sCzMZFTPXVwXvQLitBLC4qFl3J2eii3c"

func CheckSum(filepath string) ([]byte, error) {
	f, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}
	content, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, err
	}
	sum := sha256.Sum256(content)
	return sum[:], nil
}

func EncryptMany(values ...*string) error {
	for _, val := range values {
		if *val == "." || *val == ".." || *val == "~" {
			continue
		}
		if strings.Contains(*val, "/") {
			err := EncryptPath(val)
			if err != nil {
				return err
			}
		} else {
			if err := encrypt(val); err != nil {
				return err
			}
		}
	}
	return nil
}

func DecryptMany(values ...*string) error {
	for _, val := range values {
		if strings.Contains(*val, "/") {
			err := DecryptPath(val)
			if err != nil {
				return err
			}
		}
		if err := decrypt(val); err != nil {
			return err
		}
	}
	return nil
}

func DecryptPath(value *string) error {
	tokens := strings.Split(*value, "/")
	for _, token := range tokens {
		if token == "." || token == ".." || token == "~" {
			continue
		}
		err := decrypt(&token)
		if err != nil {
			return err
		}
	}
	*value = strings.Join(tokens, "/")
	return nil
}

func EncryptPath(value *string) error {
	tokens := strings.Split(*value, "/")
	for i, token := range tokens {
		if token == "." || token == ".." || token == "~" {
			continue
		}
		err := encrypt(&tokens[i])
		if err != nil {
			return err
		}
	}
	*value = strings.Join(tokens, "/")
	return nil
}

func encrypt(value *string) error {
	text := []byte(*value)
	key := []byte(Key)
	cipherBlock, err := aes.NewCipher(key)
	if err != nil {
		return err
	}
	gcm, err := cipher.NewGCM(cipherBlock)
	if err != nil {
		return err
	}
	nounce := []byte("3dPWjxlMI7sQ")
	*value = url.PathEscape(string(gcm.Seal(nounce, nounce, text, nil)))
	return nil
}

func decrypt(value *string) error {
	unescaped, err := url.PathUnescape(*value)
	cipherText := []byte(unescaped)
	key := []byte(Key)
	cipherBlock, err := aes.NewCipher(key)
	if err != nil {
		return err
	}
	gcm, err := cipher.NewGCM(cipherBlock)
	if err != nil {
		return err
	}
	nonceSize := gcm.NonceSize()
	if len(cipherText) < nonceSize {
		return errors.New(fmt.Sprintf("Invalid cipher text %s", cipherText))
	}
	nonce, cipherText := cipherText[:nonceSize], cipherText[nonceSize:]
	plainText, err := gcm.Open(nil, nonce, cipherText, nil)
	*value = string(plainText)
	return nil
}
