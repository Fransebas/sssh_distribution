package Terminal

var Bashrc = `#!/usr/bin/env bash
# SSSH defined variables
# $SSSH The path for the sssh_server executable
# $SSSH_USER user id string
# $HIST_FILE_NAME session history file

export TERM="xterm"

# The terminal doesn't start on the user directory so we use cd'
cd

createIfNotExist() {
    if [ ! -e "$1" ] ; then
        touch "$1"
    fi
}

[[ -r ~/.bash_profile ]] && . ~/.bash_profile
[[ -r ~/.bashrc ]] && . ~/.bashrc
[[ -r ~/.profile ]] && . ~/.profile


############ OH MY ZSH

ZSH_AUTOSUGGEST_HIGHLIGHT_STYLE="fg=#ff00ff,bg=cyan,bold,underline"

ZSH_AUTOSUGGEST_HIGHLIGHT_STYLE="fg=240"

export ZSH=/etc/sssh_zsh
ZSH_THEME="custom"

plugins=(zsh-autosuggestions colorize colored-man-pages cp)

source $ZSH/oh-my-zsh.sh

# ! ########### END OH MY ZSH

############ Set up of history File

HISTSIZE=1000
SAVEHIST=1000

createIfNotExist $HIST_FILE_NAME

export HISTFILE=$HIST_FILE_NAME

# ! ########### End Set up of history File

############ Hooks

chpwdcmd() { 
	setopt LOCAL_OPTIONS NO_NOTIFY NO_MONITOR
	$SSSH -mode=prompt -userid=$SSSH_USER -pwd="$(pwd)" &> /dev/null & disown
	setopt LOCAL_OPTIONS NOTIFY MONITOR
}

prmptcmd() {
	setopt LOCAL_OPTIONS NO_NOTIFY NO_MONITOR
	$SSSH -mode=prompt -userid=$SSSH_USER -history="$(history 1)" -pwd="$(pwd)" &> /dev/null & disown
	setopt LOCAL_OPTIONS NOTIFY MONITOR
}
precmd_functions=(prmptcmd)

# ! ########### End of Hooks


history &> /dev/null

if [ $? -eq 0 ] ; then
	setopt LOCAL_OPTIONS NO_NOTIFY NO_MONITOR
	$SSSH -mode=prompt -userid=$SSSH_USER -history="$(history)" &> /dev/null & disown
	setopt LOCAL_OPTIONS NOTIFY MONITOR
fi



## THEME



`
