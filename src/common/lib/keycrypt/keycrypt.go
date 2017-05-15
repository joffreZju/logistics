package keycrypt

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"strings"
)

func Encode(keyStr, plaintextStr string) (ciphertextStr string, err error) {
	key, plaintext := []byte(keyStr), []byte(plaintextStr)
	for len(plaintext)%aes.BlockSize != 0 {
		plaintext = append(plaintext, ' ')
	}
	// CBC mode works on blocks so plaintexts may need to be padded to the
	// next whole block. For an example of such padding, see
	// https://tools.ietf.org/html/rfc5246#section-6.2.3.2. Here we'll
	// assume that the plaintext is already of the correct length.
	if len(plaintext)%aes.BlockSize != 0 {
		err = errors.New("plaintext is not a multiple of the block size")

		fmt.Println(len(plaintext), string(plaintext))
		return
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return
	}

	// The IV needs to be unique, but not secure. Therefore it's common to
	// include it at the beginning of the ciphertext.
	ciphertext := make([]byte, aes.BlockSize+len(plaintext))
	iv := ciphertext[:aes.BlockSize]
	if _, err = io.ReadFull(rand.Reader, iv); err != nil {
		return
	}

	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(ciphertext[aes.BlockSize:], plaintext)

	// It's important to remember that ciphertexts must be authenticated
	// (i.e. by using crypto/hmac) as well as being encrypted in order to
	// be secure.

	ciphertextStr = fmt.Sprintf("%x", ciphertext)

	return
}

func Decode(keyStr, ciphertextStr string) (plaintextStr string, err error) {
	key := []byte(keyStr)
	ciphertext, _ := hex.DecodeString(ciphertextStr)

	block, err := aes.NewCipher(key)
	if err != nil {
		return
	}

	// The IV needs to be unique, but not secure. Therefore it's common to
	// include it at the beginning of the ciphertext.
	if len(ciphertext) < aes.BlockSize {
		err = errors.New("ciphertext too short")
		return
	}
	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]

	// CBC mode always works in whole blocks.
	if len(ciphertext)%aes.BlockSize != 0 {
		err = errors.New("ciphertext is not a multiple of the block size")
		return
	}

	mode := cipher.NewCBCDecrypter(block, iv)

	// CryptBlocks can work in-place if the two arguments are the same.
	mode.CryptBlocks(ciphertext, ciphertext)

	// If the original plaintext lengths are not a multiple of the block
	// size, padding would have to be added when encrypting, which would be
	// removed at this point. For an example, see
	// https://tools.ietf.org/html/rfc5246#section-6.2.3.2. However, it's
	// critical to note that ciphertexts must be authenticated (i.e. by
	// using crypto/hmac) before being decrypted in order to avoid creating
	// a padding oracle.

	plaintextStr = strings.TrimSpace(string(ciphertext))
	return
}

const salt string = "suanpeizai"

func Sha256Cal(src string) string {
	h := sha256.New()
	h.Write([]byte(src + salt))
	return hex.EncodeToString(h.Sum(nil))
}
func CheckSha256(src, dst string) bool {
	h := sha256.New()
	h.Write([]byte(src + salt))
	v := hex.EncodeToString(h.Sum(nil))
	return v == dst
}
