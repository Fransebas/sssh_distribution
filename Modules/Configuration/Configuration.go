package Configuration

import "flag"

var modePtr = flag.String("mode", "server", `Select a mode for the program, available modes are: 
	server : running the sssh server
	prompt : system only function (the user shouldn't use it), it send a request to the server indicating the user typed a command, it should be use it conjunction with userid 
	keygen : generate a new key
	stop : stop running instance
	fingerprint : use to get the associated fingerprint with a running server or a given public key file
`)

// Internal use flags (prompt)
var userIdPtr = flag.String("userid", "error", "Send the id of the user should be used with the mode flag set to prompt")
var historyPtr = flag.String("history", "error", "The history of the bash, should be used with the model flag set to prompt")
var pwdPtr = flag.String("pwd", "-", "The current working directory, should be used with the model flag set to prompt")

// END: Internal use flags

// SSSH server
var httPortPtr = flag.Int("httpport", 2000, "HTTPPort for the http server")
var portPtr = flag.Int("port", 2222, "Port for the sssh server, default 2222")
var rpcPortPtr = flag.Int("rpcport", 2001, "Select a port for the rpc (internal process communication)")
var keyFile = flag.String("keyfile", "id_rsa", "If this flag is set, this key will be used to authenticate the host")

// END: SSSH server

// Keygen
var kryTypePtr = flag.String("type", "rsa", `Type of keys to generate, valid values include:
	rsa : generates a rsa key specified in PKCS#1 with the format used by open-ssh
	ecdsa : (experimental please do not use) generates a ECDSA256 key with the format used by open-ssl`)
var filenamePtr = flag.String("filename", "id_rsa", `Filename to store or read the key`)

// END: KeygenFlags

// Fingerprint
var FingerprintModeStrPtr = flag.String("fingerPrintMode", "file", `Whatever you want to get the fingerprint of the server or of a given file. Values are:
	file : get the finger of a file (only rsa supported right now )
	server : get the finger print of a given server (localhost and port, for port use port flag)
`)
var fingerprintURL = flag.String("fingerPrintUrl", "localhost", `url to send the request to get the fingerprint (url other than localhost not recommended for security reasons)`)

// END: Fingerprint

type KeygenConfig struct {
	Filename string
	Type     string
}

type FingerprintConfig struct {
	FingerprintMode string
	Port            int
	Filename        string
	Url             string
}

type PromptConfig struct {
	Pwd     string
	History string
	UserId  string
}

type Configuration struct {
	Mode              string
	UserId            string
	HTTPPort          int
	Port              int
	RPCPort           int
	KeyFile           string
	KeygenConfig      KeygenConfig
	FingerprintConfig FingerprintConfig
	PromptConfig      PromptConfig
}

func (c *Configuration) initKeygen() {
	c.KeygenConfig.Type = *kryTypePtr
	c.KeygenConfig.Filename = *filenamePtr
}

func (c *Configuration) initFingerprint() {
	c.FingerprintConfig.FingerprintMode = *FingerprintModeStrPtr
	c.FingerprintConfig.Port = *portPtr
	c.FingerprintConfig.Filename = *filenamePtr
	c.FingerprintConfig.Url = *fingerprintURL
}

func (c *Configuration) initPrompt() {
	c.PromptConfig.Pwd = *pwdPtr
	c.PromptConfig.History = *historyPtr
	c.PromptConfig.UserId = *userIdPtr
}

func (c *Configuration) Init() {
	flag.Parse()
	c.Mode = *modePtr
	c.UserId = *userIdPtr
	c.HTTPPort = *httPortPtr
	c.RPCPort = *rpcPortPtr
	c.KeyFile = *keyFile
	c.Port = *portPtr
	c.initKeygen()
	c.initFingerprint()
	c.initPrompt()
}