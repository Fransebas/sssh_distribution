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

![preview](https://i.imgur.com/EoHIJJv.png)

# Index

* [Download](#download)
* [Server](#server)
  * [Installation](#installation)
  * [Alternative Installation](#alternative-installation)
* [Client](#client)
* [Contact](#contact)


# Downloads

## Client

| OS | Link |
|---|---|
| Linux (Red Hat) | [sssh-client-0.0.16.x86_64.rpm](https://github.com/Fransebas/sssh_distribution/releases/download/v0.0.16-alpha/sssh-client-0.0.16.x86_64.rpm) |
| Linux (Debian) |  [sssh-client-0.0.16.x86_64.deb](https://github.com/Fransebas/sssh_distribution/releases/download/v0.0.16-alpha/sssh-client_0.0.16_amd64.deb) |
| Linux |  [sssh-client_0.0.16_amd64.snap](https://github.com/Fransebas/sssh_distribution/releases/download/v0.0.16-alpha/sssh-client_0.0.16_amd64.snap) |
| macOS |  [sssh-client-0.0.16.dmg](https://github.com/Fransebas/sssh_distribution/releases/download/v0.0.16-alpha/SSSH.Client-0.0.16.dmg) |
| windows |  [sssh-client-0.0.16.msi](https://github.com/Fransebas/sssh_distribution/releases/download/v0.0.16-alpha/SSSH.Client.0.0.16.msi) |




## Server

| OS | Link |
|---|---|
| Linux (Debian) | [sssh-server_x86-64_0.0.16.deb](https://github.com/Fransebas/sssh_distribution/releases/download/v0.0.16-alpha-server/sssh-server_x86-64_0.0.16.deb) |
| macOS |  [sssh_server0.0.16.pkg](https://github.com/Fransebas/sssh_distribution/releases/download/v0.0.16-alpha-server/sssh_server0.0.16.pkg) |
| Raspberrypi (Debian) |  [sssh-server_raspberrypi_0.0.16.deb](https://github.com/Fransebas/sssh_distribution/releases/download/v0.0.16-alpha-server/sssh-server_raspberrypi_0.0.16.deb) |
| Linux Manual Install |  [manual.sssh_server-linux_x86-64_0.0.16.tar.gz](https://github.com/Fransebas/sssh_distribution/releases/download/v0.0.16-alpha-server/manual.sssh_server-linux_x86-64_0.0.16.tar.gz) |


# Server

## Installation

The installation has been updated and now you can just double click the installer or use dpkg (linux).


## Alternative installation

If you do not wish to use the installer, you can simply run

`sudo ./sssh_server`

But you need to have a rsa key, which you can create with the OpenSSH suite or using the `sssh_server` like this:

`./sssh_server -mode=keygen`

You can specify the location of the keys:

`sudo ./sssh_server -keyfile=./id_rsa`


# Client

The client currently only supports Unix-like systems and Windows 10 (<= 1809), it has been tested in macos (10.14.6), Ubuntu (16), Windows (latest version on 3/14/2020). Any compatibility problem please reported it in the issue section.

If you find a Unix system wher the client doesn't work, please open an Issue to fix it.

# Contact

ffransebas@gmail.com
