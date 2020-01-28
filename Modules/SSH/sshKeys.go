package SSH

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"github.com/tyler-smith/go-bip39"
	"golang.org/x/crypto/ssh"
	"io/ioutil"
	"log"
	"math/big"
	"net"
	"os"
	"sssh_server/CustomUtils"
	"strings"
	"time"
)

// Valid KEY types
const (
	RSA   = "rsa"
	ECDSA = "ecdsa"
)

func publicKey(priv interface{}) interface{} {
	switch k := priv.(type) {
	case *rsa.PrivateKey:
		return &k.PublicKey
	case *ecdsa.PrivateKey:
		return &k.PublicKey
	default:
		return nil
	}
}

func pemBlockForKey(priv interface{}) *pem.Block {
	switch k := priv.(type) {
	case *rsa.PrivateKey:
		return &pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(k)}
	case *ecdsa.PrivateKey:
		b, err := x509.MarshalECPrivateKey(k)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Unable to marshal ECDSA private key: %v", err)
			os.Exit(2)
		}
		return &pem.Block{Type: "EC PRIVATE KEY", Bytes: b}
	default:
		return nil
	}
}

// "TODO: figure out the right directory to store everything"
var KEYS_DIR = ""

func SaveFiles(priv interface{}, derBytes []byte) {
	certOut, err := os.Create(KEYS_DIR + "cert.pem")
	if err != nil {
		log.Fatalf("failed to open cert.pem for writing: %s", err)
	}
	if err := pem.Encode(certOut, &pem.Block{Type: "CERTIFICATE", Bytes: derBytes}); err != nil {
		log.Fatalf("failed to write data to cert.pem: %s", err)
	}
	if err := certOut.Close(); err != nil {
		log.Fatalf("error closing cert.pem: %s", err)
	}
	log.Print("wrote cert.pem\n")

	keyOut, err := os.OpenFile(KEYS_DIR+"key.pem", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		log.Print("failed to open key.pem for writing:", err)
		return
	}
	if err := pem.Encode(keyOut, pemBlockForKey(priv)); err != nil {
		log.Fatalf("failed to write data to key.pem: %s", err)
	}
	if err := keyOut.Close(); err != nil {
		log.Fatalf("error closing key.pem: %s", err)
	}
	log.Print("wrote key.pem\n")
}

func GenerateNewECSDAKey() (certPEMBlock, keyPEMBlock []byte) {
	priv, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	notBefore := time.Now()

	serialNumberLimit := new(big.Int).Lsh(big.NewInt(1), 128)
	serialNumber, err := rand.Int(rand.Reader, serialNumberLimit)
	if err != nil {
		log.Fatalf("failed to generate serial number: %s", err)
	}
	validFor := 10 * 365 * 24 * time.Hour
	notAfter := notBefore.Add(validFor)

	template := x509.Certificate{
		SerialNumber: serialNumber,
		Subject: pkix.Name{
			CommonName: "127.0.0.1:2000",
		},
		NotBefore: notBefore,
		NotAfter:  notAfter,

		//KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign | x509.KeyUsageCRLSign,
		//ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		//BasicConstraintsValid: true,
		IsCA:        true,
		IPAddresses: []net.IP{net.ParseIP("127.0.0.1")},
		//DNSNames : []string{"localhost:2000"},
	}

	certPEMBlock, err = x509.CreateCertificate(rand.Reader, &template, &template, publicKey(priv), priv)
	if err != nil {
		log.Fatalf("Failed to create certificate: %s", err)
	}
	keyPEMBlock = pemBlockForKey(priv).Bytes
	SaveFiles(priv, certPEMBlock)
	return
}

func MakeMnemonic(pubKey []byte) (string, error) {
	hash, e := GetKeyHash(pubKey)
	if e != nil {
		return "", e
	}
	return bip39.NewMnemonic(hash)
}

func GetKeyHash(pubKey []byte) ([]byte, error) {
	var res = CustomUtils.ExecuteCommand(fmt.Sprintf(`ssh-keygen -q -l -F localhost -f /dev/stdin <<<"localhost %v"`, string(pubKey)))
	hash := strings.Split(strings.Split(res, " ")[2], ":")[1]
	return base64.RawStdEncoding.DecodeString(hash)
}

// MakeSSHKeyPair make a pair of public and private keys for SSH access.
// Public key is encoded in the format for inclusion in an OpenSSH authorized_keys file.
// Private Key generated is PEM encoded
func MakeSSHKeyPair(privateKeyPath, pubKeyPath string) error {
	privateKey, err := rsa.GenerateKey(rand.Reader, 1024)
	if err != nil {
		return err
	}

	// generate and write private key as PEM
	privateKeyFile, err := os.Create(privateKeyPath)
	defer privateKeyFile.Close()
	if err != nil {
		return err
	}
	privateKeyPEM := &pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(privateKey)}
	if err := pem.Encode(privateKeyFile, privateKeyPEM); err != nil {
		return err
	}

	// generate and write public key
	pub, err := ssh.NewPublicKey(&privateKey.PublicKey)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(pubKeyPath, ssh.MarshalAuthorizedKey(pub), 0655)
}
