#!/bin/bash
TAILWIND_VERSION="v3.3.5"

install_aws_cdk() {
    cecho "Installing aws-cdk"
    npm install -g aws-cdk
}

install_tailwindcss() {
    cecho "Installing tailwindcss $TAILWIND_VERSION"
    curl -sLO "https://github.com/tailwindlabs/tailwindcss/releases/download/$TAILWIND_VERSION/tailwindcss-linux-x64"
    chmod +x tailwindcss-linux-x64

    BIN_DIR="/usr/local/bin"
    if [ ! -d "$BIN_DIR" ];  then
        mkdir $BIN_DIR
    fi

    mv tailwindcss-linux-x64 /usr/local/bin/tailwindcss
}

install_gotools() {
    
    cecho "Installing go tool github.com/cosmtrek/air@latest"
    go install github.com/cosmtrek/air@latest
    
    cecho "Installing go tool github.com/a-h/templ/cmd/templ@latestt"
    go install github.com/a-h/templ/cmd/templ@latest
}

GREEN="\033[1;32m"
NOCOLOR="\033[0m"
cecho() {
    echo -e "$GREEN$1$NOCOLOR"
}

install_tailwindcss
install_gotools
install_aws_cdk