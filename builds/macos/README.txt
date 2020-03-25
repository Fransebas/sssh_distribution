#Instalation

Installation is pretty simple just run

`sudo ./install.sh`

This will move the executable `sssh_server` to  `/usr/local/bin/sssh_server`

Move the sssh server configuration file `sssh.conf` to `/etc/sssh.conf`

Move the sssh_server.service to `/lib/systemd/system/sssh_server.service` or in mac `com.ssshserver.app.plist` to `/Library/LaunchDaemons/com.ssshserver.app.plist`

Then it will crete a rsa key to be use to recognize your host machine, it will be located in `/etc/sssh/rsa_host` .
You could also wish to use the host keys already in you computer, they are located in `/etc/ssh/...` to do so, modify the conf file.

# Alternative use

If you do not wish to use the installer, you can simply run

`sudo ./sssh_server`

But you need to have a rsa key, which you can create with the OpenSSH suite or using the `sssh_server` like this:

`./sssh_server -mode=keygen`

You can specify the location of the keys:

`sudo ./sssh_server -keyfile=./id_rsa`