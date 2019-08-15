package SSH

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"io"
	"sssh_server/CustomUtils"
)

type User struct {
	// Later I have to add here the keys
	ID string
}

func NewUser() *User {
	user := new(User)
	seed, err := GenerateRandomBytes(32)
	CustomUtils.CheckPanic(err, "Could not generate random string wtf")
	user.ID = hex.EncodeToString(seed)
	return user
}

func GenerateRandomBytes(n int) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	// Note that err == nil only if we read len(b) bytes.
	if err != nil {
		return nil, err
	}

	return b, nil
}

func ProofClaims(key []byte) error {
	return nil
}

//func SSHDecode(msg []byte, user *User) ([]byte, error) {
//	// TODO: encryption not implemented
//	// If the decoding is not successful return an error
//	var encryptedJSON EncryptedJSON
//	key, _ := hex.DecodeString("6368616e679520746869731070617373")
//	err := json.Unmarshal(msg, &encryptedJSON)
//	CustomUtils.CheckPanic(err, "Could not unmarshal json")
//	iv, e := hex.DecodeString(encryptedJSON.Iv)
//	CustomUtils.CheckPanic(e, "Could not decode iv")
//	decodedMessage, e := hex.DecodeString(encryptedJSON.EncryptedMessage)
//
//	return decrypt(decodedMessage, key, iv)[:encryptedJSON.MessageSize], nil
//}

//func decrypt(msg []byte, key []byte, iv []byte) []byte {
//	block, err := aes.NewCipher(key)
//	if err != nil {
//		panic(err)
//	}
//	//ciphertext := make([]byte, len(msg))
//	mode := cipher.NewCBCDecrypter(block, iv)
//	mode.CryptBlocks(msg, msg)
//	return msg
//}
//
//func getIV(plaintext string) []byte {
//	ciphertext := make([]byte, aes.BlockSize+len(plaintext))
//	iv := ciphertext[:aes.BlockSize]
//	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
//		panic(err)
//	}
//	return iv
//}

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
