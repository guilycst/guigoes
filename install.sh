#!/bin/bash
TAILWIND_VERSION="v3.3.5"

install_tailwindcss() {
    wget https://github.com/tailwindlabs/tailwindcss/releases/download/$TAILWIND_VERSION/tailwindcss-macos-arm64
    chmod +x tailwindcss-macos-arm64
    mv tailwindcss-macos-arm64 /usr/local/bin/tailwindcss
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