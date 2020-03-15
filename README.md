# SSSH Alpha version

This repository is to distribute the code for SSSH

* [Server](#server)
  * [Running the server](#running-the-server)
  * [Creating an RSA key](#creating-an-rsa-key)
  * [Getting the key fingerprint](#getting-the-key-fingerprint)

# Server

The server currently only supports Unix-like systems, it has been tested in Macos, Ubuntu, Debian Jessy (raspberry pi). Any compatibility issue with a different Unix system, please open an issue.

## Running the server

Note: I'm working on the right way to run the server as a service on boot but in the meantime, you need to run it manually.

Because the server should keep open it's recommended to run it using a session manager program like `tmux` or `screen`, or using `&` to run it as a separate process like this (for multiuser access the server should run as sudo):

`& sudo sssh_server`

The current flags for the `server` mode are:
  - `keyfile` path for the RSA key to use to authenticate the server, the default value is `./id_rsa`
  - `port` Port for the sssh server, default 2222
  
Example:

`& sudo sssh_server -keyfile=/etc/ssh/id_rsa -port=22` With this flags, the commad will search for the auth key in the directory `/etc/ssh/id_rsa` and use the port 22 for comunication.

## Creating an RSA key

Note: don't confuse this key with the user authentication key, this key is to trust the server.

Currently, the only available key type is RSA, you can use the existing key under `/etc/ssh` or create a new key using:

`sssh_server -mode=keygen` 

This will generate two files in the current directory `id_rsa` which is the private key and `id_rsa.pub` the public key that you can put in your `know_host`.

The current flags for the `keygen` mode are:

- `filename` The filename of the generated key i.e. `Filename` and `Filename.pub`
- `type` For the key type, currently only `rsa` type is supported


## Getting the key fingerprint

In order to verify the identity of your server, you'll need to have physical access to it and generate the fingerprint of your key with the following command:

`sssh_server -mode=fingerprint -file=/keylocation`

The current flags for the `fingerprint` mode are:

- `file` The location of your key, the default value is the current directory i.e. `./id_rsa` and `./id_rsa.pub`
- `server` The URL for the server (use this only for debugging, this is not secure for remote servers), default `localhost`
