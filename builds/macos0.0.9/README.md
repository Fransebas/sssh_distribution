## Instalation

Installation is pretty simple just run

`sudo ./install.sh`

After that run the following:
#### Linux

`sudo systemctl start sssh_server`

`sudo systemctl statys sssh_server`

The status should show something like this:

```
● sssh_server.service - Service for sssh server
   Loaded: loaded (/lib/systemd/system/sssh_server.service; static; vendor preset: enabled)
   Active: active (running) since Mon 2020-03-30 22:39:57 CST; 5s ago
 Main PID: 19346 (sssh_server)
    Tasks: 7 (limit: 4496)
   Memory: 3.3M
   CGroup: /system.slice/sssh_server.service
           └─19346 /usr/local/bin/sssh_server

mar 30 22:39:57 fransebasUbuntu systemd[1]: Started Service for sssh server.
```

Any problem running the server, please open an Issue.

#### Macos

`sudo launchctl load /Library/LaunchDaemons/com.ssshserver.app.plist`

`sudo launchctl start com.ssshserver.app`

`sudo launchctl list | grep com.ssshserver.app` 

The status should show something like this:

`-	2	com.ssshserver.app`

If you don't see any output that means that the service is not running.

Any problem running the server, please open an Issue.

## What does the installation script does?


This will move the executable `sssh_server` to  `/usr/local/bin/sssh_server`

Move the sssh server configuration file `sssh.conf` to `/etc/sssh.conf`

Move the `sssh_server.service` to `/lib/systemd/system/sssh_server.service` or in mac `com.ssshserver.app.plist` to `/Library/LaunchDaemons/com.ssshserver.app.plist`

Then it will crete a rsa key to be use to recognize your host machine, it will be located in `/etc/sssh/rsa_host` .
You could also wish to use the host keys already in you computer, they are located in `/etc/ssh/...` to do so, modify the conf file.


## Alternative use

If you do not wish to use the installer, you can simply run

`sudo ./sssh_server`

But you need to have a rsa key, which you can create with the OpenSSH suite or using the `sssh_server` like this:

`./sssh_server -mode=keygen`

You can specify the location of the keys:

`sudo ./sssh_server -keyfile=./id_rsa`