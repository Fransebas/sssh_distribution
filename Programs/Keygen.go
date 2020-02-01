package Programs

import (
	"fmt"
	"sssh_server/CustomUtils"
	"sssh_server/Modules/Configuration"
	"sssh_server/SessionModules/SSH"
)

func Keygen(config Configuration.Configuration) {
	keygenConfig := config.KeygenConfig
	if keygenConfig.Type == SSH.RSA {
		privatePath := keygenConfig.Filename
		publicPath := privatePath + ".pub"
		e := SSH.MakeSSHKeyPair(privatePath, publicPath)
		CustomUtils.CheckPanic(e, "could not generate key in path "+privatePath)
	} else if keygenConfig.Type == SSH.ECDSA {
		SSH.GenerateNewECSDAKey()
	} else {
		panic(fmt.Sprintf("Key type %v still not supported =(", keygenConfig.Type))
	}

}
