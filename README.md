# zIRC

A Chat Application based on IRC v3 protocol spec.

## Install

### Prereqs

1. Ensure you have docker desktop installed
2. (optional) For convenience, create an env variable which contains the path to where you have cloned this repo. The command below assums you are using bash. Replace `<path to repo>` with the actual path.

```bash
echo "export CHAT_ROOT=<path to repo>" >> $HOME/.bashrc
source $HOME/.bashrc
```

### Steps

1. Navigate to the repo root

```bash
cd $CHAT_ROOT
```

2. Start via docker compose

```bash
docker compose up --build -d
```

## Using the cli test client

1. Navigate to the test client dir

```bash
cd $CHAT_ROOT/Chat_Server/TestClient
```

2. Run the test client

```bash
node test_client.js
```

3. Enter `connect` to connect to the chat server running in docker
4. Enter any valid commands for the chat server. For example:

```bash
HELP USERCMDS
```

will output:

> :Z.IRC 704 * :** Help System **
>
>:Z.IRC 705 * :
>:Z.IRC 705 * :Try HELP <commmand> for specific help
>:Z.IRC 705 * :HELP USERCMDS to list available commands,
>:Z.IRC 706 * :or join the #help channel 
>
>HELP USERCMDS
>:Z.IRC 704 * :** Help User Commands **
>
>:Z.IRC 705 * :
>:Z.IRC 705 * :CAP
>:Z.IRC 705 * :JOIN
>
>:Z.IRC 705 * :NICK
>
>:Z.IRC 705 * :SEND
>:Z.IRC 705 * :PASS
>:Z.IRC 705 * :AUTHENTICATE
>:Z.IRC 705 * :USER
>:Z.IRC 705 * :LIST
>:Z.IRC 705 * :PRIVMSG
>:Z.IRC 705 * :HELP
>:Z.IRC 706 * :REGISTER 

You can then type in `HELP <cmd>` to output the help description for that particular command.