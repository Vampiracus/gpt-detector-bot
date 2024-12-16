# Bot for https://github.com/BratsunKD/chat-gpt-detector-app

## Setup

```Bash
export BOT_TOKEN  = <your_bot_token>
export SERVICE_IP = <service_ip>
go run .
```

## Installing GO

```Bash
sudo apt update
sudo apt install golang-go
wget https://go.dev/dl/go1.21.3.linux-amd64.tar.gz
sudo tar -C /usr/local -xzf go1.21.3.linux-amd64.tar.gz
export PATH=$PATH:/usr/local/go/bin
export GOPATH=$HOME/goproject
export PATH=$PATH:$GOPATH/bin
```
