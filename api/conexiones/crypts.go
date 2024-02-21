package conexiones

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
)

var key string = "hereismycipherke"

// var key string = "thisis32bitlongp"

func TestAES() {

	// cipher key
	//	key := "thisis32bitlongpassphraseimusing"

	// plaintext
	pt := "This is a secret"

	// c := EncryptAES([]byte(key), pt)
	c := Encriptar(pt)

	// plaintext
	fmt.Println(pt)

	// ciphertext
	fmt.Println(c)

	// decrypt
	pt = Desencriptar(c)
	// plaintext
	fmt.Println(pt)
}

func Encriptar(ct string) string {
	return hex.EncodeToString(encrypt([]byte(key), []byte(ct)))
}

func Desencriptar(ct string) string {
	ctbyte, _ := hex.DecodeString(ct)
	return decrypt([]byte(key), ctbyte)
}

/*
package main

import (

	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"

)

	func main() {
		text := []byte("TEXT TO ENCRYPT AND DECRYPT")

		key := []byte("5v8y/B?E(G+KbPeShVmYq3t6w9z$C&12")

		secret := encrypt(key, text)
		//str1 := string(secret[:])
		myString := hex.EncodeToString(secret)
		fmt.Println(myString)
		secret, _ = hex.DecodeString(myString)
		plainString := decrypt(key, secret)

		fmt.Println(plainString)
	}
*/
func encrypt(key []byte, text []byte) []byte {
	c, _ := aes.NewCipher(key)
	gcm, _ := cipher.NewGCM(c)

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		panic(err.Error())
	}

	result := gcm.Seal(nonce, nonce, text, nil)

	return result
}

func decrypt(key []byte, ciphertext []byte) string {
	c, _ := aes.NewCipher(key)
	gcm, _ := cipher.NewGCM(c)

	nonceSize := gcm.NonceSize()
	if len(ciphertext) < nonceSize {
		panic("ciphertext size is less than nonceSize")
	}

	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]
	plaintext, _ := gcm.Open(nil, nonce, ciphertext, nil)

	return string(plaintext)
}
