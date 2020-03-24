package Terminal

var Bashrc = `#!/usr/bin/env bash

# SSSH defined variables
# $SSSH The path for the sssh_server executable
# $SSSH_USER user id string
# $HIST_FILE_NAME session history file

createIfNotExist(){
    if [ ! -e "$1" ] ; then
        touch "$1"
    fi
}


#history -anrw $HIST_FILE_NAME

createIfNotExist ~/.bashrc
[[ -r ~/.bashrc ]] && . ~/.bashrc
[[ -r ~/.profile ]] && . ~/.profile

# Colors in the terminal

export PS1="\[\033[36m\]\u\[\033[m\]@\[\033[32m\]\h:\[\033[33;1m\]\w\[\033[m\]\$ "
export CLICOLOR=1
export LSCOLORS=ExFxBxDxCxegedabagacad
alias ls='ls -GFh'
export TERM=xterm-color

# End Colors

# Set up of history File
createIfNotExist $HIST_FILE_NAME

export HISTFILE=$HIST_FILE_NAME
export HISTCONTROL=ignorespace ; history -d 1

# End Set up of history File


# This is for production
$SSSH -mode=prompt -userid=$SSSH_USER -history="$(history)" -pwd="$HOME"
#curl --data "$(history)" http://localhost:2000/newcommand?SSSH_USER=$SSSH_USER &> /dev/null

export PROMPT_COMMAND='& $SSSH -mode=prompt -userid=$SSSH_USER -history="$(history 1)" -pwd="$(pwd)"'

 cd
`
