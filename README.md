# SSSH Alpha version

SSSH is an augmentation for the SSH protocol which adds GUI elements to improve the usability of a terminal. The current functionalities are (all the communication are transmitted over an ssh connection):

* Manage servers configurations:
  * Username
  * Password
  * Authentication key
  * Host IP
* Bash terminal.
* File explorer.
* History GUI:
  * List of recently used commands.
  * Search.
  * Locally save commands.
* List of all commands installed:
  * Search.
  * Locally save commands.
* Text editor for remote files (with command support for those vim lovers).
* List of current variables in the bash terminal (working on more improvements).
* Manual Visualizer for commands.
* *The ability to create Plug-ins.

*The Plug-in support is on development but the current implementation takes this into consideration and every functionality above is created like a plug-in that uses an API.

The idea for the SSSH protocol is to create a GUI for every Unix utility that would improve the use of a terminal, for example, the `History` functionality is a GUI to manage the `history` command and can help terminal users to locally save frequently used commands for later use or help search the command history in a more visual way than using `ctrl + r` or hitting the up arrow multiple times.

As a dummy example to illustrate this idea, let's image a Plug-in that has a simple GUI to control `chmod` which is not a hard command to use but this GUI could help people that don't know by memory all the modifier values for `chmod` also it could help people that use that command in a regular basis by making it faster to use (maybe remembering the last used modifiers) and also having a file selector to quickly select the file to apply the modifiers.

# Index

* [Server](#server)
  * [Instalation](#instalation)
     * [Linux](#linux)
     * [Macos](#macos)
  * [What does the installation script does?](#what-does-the-installation-script-does)
  * [Alternative use](#alternative-use)
* [Client](#client)

# Server

## Instalation

Installation is pretty simple just run

`sudo ./install.sh`

After that run the following:
#### Linux

`systemctl start sssh_server`

`systemctl statys sssh_server`

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

`launchctl load /Library/LaunchDaemons/com.ssshserver.app.plist`

`launchctl start com.ssshserver.app`

`launchctl list | grep com.ssshserver.app` 

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


# Client

The client currently only supports Unix-like systems and Windows 10 (<= 1809), it has been tested in macos (10.14.6), Ubuntu (16), Windows (latest version on 3/14/2020). Any compatibility problem please reported it in the issue section.

If you find a Unix system wher the client doesn't work, please open an Issue to fix it.
