#!/bin/bash

# Function to check if a command exists
command_exists() {
    command -v "$1" >/dev/null 2>&1
}

# Function to check Go version compatibility
check_go_version() {
    REQUIRED_GO_VERSION="1.20"
    INSTALLED_GO_VERSION=$(go version | awk '{print $3}' | sed 's/go//')

    if [ "$(printf '%s\n' "$REQUIRED_GO_VERSION" "$INSTALLED_GO_VERSION" | sort -V | head -n1)" = "$REQUIRED_GO_VERSION" ]; then 
        echo "Go version is compatible: $INSTALLED_GO_VERSION"
    else
        echo "Go version $REQUIRED_GO_VERSION or higher is required. Installed version is $INSTALLED_GO_VERSION."
        exit 1
    fi
}

# Function to check OS and architecture compatibility
check_os_arch() {
    OS=$(uname -s)
    ARCH=$(uname -m)

    case "$OS" in
        Linux*|Darwin*)
            echo "Operating system is compatible: $OS"
            ;;
        *)
            echo "This script only supports Linux and macOS."
            exit 1
            ;;
    esac

    case "$ARCH" in
        x86_64|arm64)
            echo "Architecture is compatible: $ARCH"
            ;;
        *)
            echo "This script only supports x86_64 and arm64 architectures."
            exit 1
            ;;
    esac
}

# Function to add alias to bash or zsh
add_alias() {
    SHELL_NAME=$(basename "$SHELL")
    ALIAS_CMD="alias log-watcher='$(pwd)/log-watcher'"

    case "$SHELL_NAME" in
        bash)
            echo "$ALIAS_CMD" >> ~/.bashrc
            source ~/.bashrc
            ;;
        zsh)
            echo "$ALIAS_CMD" >> ~/.zshrc
            source ~/.zshrc
            ;;
        *)
            echo "Shell $SHELL_NAME is not supported for automatic alias addition."
            ;;
    esac

    echo "Alias 'log-watcher' added to $SHELL_NAME."
}

# Check if the OS and architecture are compatible
check_os_arch

# Check if Go is installed
if command_exists go; then
    echo "Go is installed."
    check_go_version
else
    echo "Go is not installed. Installing Go..."

    # Install Go (example assumes Ubuntu/Debian or macOS with Homebrew)
    if [ "$OS" = "Linux" ]; then
        if [ "$ARCH" = "x86_64" ]; then
            wget https://golang.org/dl/go1.20.linux-amd64.tar.gz
            sudo tar -C /usr/local -xzf go1.20.linux-amd64.tar.gz
        elif [ "$ARCH" = "arm64" ]; then
            wget https://golang.org/dl/go1.20.linux-arm64.tar.gz
            sudo tar -C /usr/local -xzf go1.20.linux-arm64.tar.gz
        fi
        export PATH=$PATH:/usr/local/go/bin
    elif [ "$OS" = "Darwin" ]; then
        brew install go
    fi

    # Verify installation
    if ! command_exists go; then
        echo "Go installation failed."
        exit 1
    fi
    check_go_version
fi

# Build the project based on the detected architecture
echo "Building the project..."
if [ "$ARCH" = "x86_64" ]; then
    GOARCH=amd64
elif [ "$ARCH" = "arm64" ]; then
    GOARCH=arm64
fi

GOOS=$(uname | tr '[:upper:]' '[:lower:]') GOARCH=$GOARCH go build -o log-watcher ./cmd

echo "Installation completed successfully."

# Prompt to add alias
read -p "Do you want to add 'log-watcher' as an alias to your $SHELL_NAME? (y/n): " -n 1 -r
echo
if [[ $REPLY =~ ^[Yy]$ ]]; then
    add_alias
else
    echo "Alias not added. You can manually add it later if needed."
fi