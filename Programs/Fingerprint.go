package Programs

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"sssh_server/CustomUtils"
	"sssh_server/Modules/Configuration"
	"sssh_server/Modules/SSH"
	"sssh_server/SessionModules/SessionLayer"
)

func print(pubKey SessionLayer.PubKeyShare) {
	fmt.Printf("Fingerprint SHA256: %v \nMnemonic : %v \n", pubKey.Hash, pubKey.Mnemonic)
}

func file(config Configuration.FingerprintConfig) {
	b, e := ioutil.ReadFile(config.Filename)
	CustomUtils.CheckPanic(e, "Couldn't read pub key")
	hash, e := SSH.GetKeyHash(b)
	mnemonic, e := SSH.MakeMnemonic(b)

	pubKey := SessionLayer.PubKeyShare{
		Hash:     base64.RawStdEncoding.EncodeToString(hash),
		Mnemonic: mnemonic,
		Key:      string(b),
	}

	print(pubKey)
}

func server(config Configuration.FingerprintConfig) {
	r, e := http.Get(fmt.Sprintf("http://%v:%v/pubKey", config.Url, config.Port))
	if e != nil {
		fmt.Println("Couln't connect to the running instance")
		return
	}
	b, e := ioutil.ReadAll(r.Body)
	CustomUtils.CheckPanic(e, "Couldn't read response")
	var pubKey SessionLayer.PubKeyShare
	e = json.Unmarshal(b, &pubKey)
	CustomUtils.CheckPanic(e, "Couldn't read response")

	print(pubKey)

}

func Fingerprint(config Configuration.Configuration) {
	fingerprintConfig := config.FingerprintConfig
	if fingerprintConfig.FingerprintMode == "file" {
		file(fingerprintConfig)
	} else if fingerprintConfig.FingerprintMode == "server" {
		server(fingerprintConfig)
	} else {
		panic(fmt.Sprintf("Fingerprint mode %v doesn't exist", fingerprintConfig.FingerprintMode))
	}
}
