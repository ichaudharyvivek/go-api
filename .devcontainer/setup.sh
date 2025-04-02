# install curl, git, ...
apt update
apt install -y curl git jq

useradd -m user
su user

# install go
VERSION='1.24.2'
OS='linux'
ARCH='amd64'

curl -OL https://dl.google.com/go/go${VERSION}.${OS}-${ARCH}.tar.gz
tar -C /usr/local -xzf go$VERSION.$OS-$ARCH.tar.gz
rm go$VERSION.$OS-$ARCH.tar.gz

INSTALLED_GO_VERSION=$(go version)
echo "Go version ${INSTALLED_GO_VERSION} is installed"

# install direnv
curl -sfL https://direnv.net/install.sh | bash
echo "Installed direnv for configurations"

# install gopls, dlv, hey, air
echo "Getting development tools"
go install golang.org/x/tools/gopls@latest
go install github.com/go-delve/delve/cmd/dlv@latest
go install github.com/rakyll/hey@latest
go install github.com/air-verse/air@latest

# vscode-go dependencies 
echo "Getting dependencies for the vscode-go plugin"
go install github.com/acroca/go-symbols@latest
go install github.com/cweill/gotests/...@latest
go install github.com/davidrjenni/reftools/cmd/fillstruct@latest
go install github.com/haya14busa/goplay/cmd/goplay@latest
go install github.com/mdempsky/gocode@latest
go install github.com/ramya-rao-a/go-outline@latest
go install github.com/rogpeppe/godef@latest
go install github.com/sqs/goreturns@latest
go install github.com/uudashr/gopkgs/cmd/gopkgs@latest
go install github.com/zmb3/gogetdoc@latest
go install golang.org/x/lint/golint@latest
go install golang.org/x/tools/cmd/gorename@latest
