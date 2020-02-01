package SSH

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"io"
)

type User struct {
	// Later I have to add here the keys
	ID string
}

func ProofClaims(key []byte) error {
	return nil
}

func encrypt(msg []byte, key []byte, iv []byte) []byte {
	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}
	//ciphertext := make([]byte, len(msg))
	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(msg, msg)
	return msg
}

/*
 * The message should be a multiple of 16 in size,
 * if its not, then fill with random chars
 */
func completeMessage(msg []byte) (message []byte, size int) {
	rndmBytes := make([]byte, (aes.BlockSize-len(msg)%aes.BlockSize)%aes.BlockSize)
	if _, err := io.ReadFull(rand.Reader, rndmBytes); err != nil {
		panic(err)
	}
	message = append(msg, rndmBytes...)
	return message, len(message)
}

func SSHEncode(msg []byte, user *User) ([]byte, error) {
	// TODO: encryption not implemented
	// If the encoding is not successful return an error
	// TODO: User proper key
	var encryptedJSON EncryptedJSON
	encryptedJSON.MessageSize = len(msg)
	msg, _ = completeMessage(msg)

	key, _ := hex.DecodeString("6368616e679520746869731070617373")

	iv := make([]byte, aes.BlockSize)
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		panic(err)
	}
	encryptedJSON.EncryptedMessage = hex.EncodeToString(encrypt(msg, key, iv))
	encryptedJSON.Iv = hex.EncodeToString(iv)

	return json.Marshal(encryptedJSON)
}

type EncryptedJSON struct {
	Iv               string `json:"iv"`               // initialization vector in hexadecimal format
	EncryptedMessage string `json:"encryptedMessage"` // encrypted message in hex format
	MessageSize      int    `json:"messageSize"`      // This states the size of the message because the program
	// needs to add random data to complete a multiple of 16 and this size tell us how much was added
}
