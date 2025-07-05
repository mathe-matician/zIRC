# zIRC

A Chat Application based on IRC v3 protocol spec.

## Install

### Prereqs

1. Ensure you have docker desktop installed
2. To run the cli test client, you will also have to have Node installed. (currently Node v22)
3. (optional) For convenience, create an env variable which contains the path to where you have cloned this repo. The command below assums you are using bash. Replace `<path to repo>` with the actual path.

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

## Testing with clients

There are two ways to test against the zirc server:
1. With an existing IRC client
2. With the test client

Both are useful for various reasons. Using any existing client should work as the server is compliant with the protocol - but it may be worth mentioning I've only tested it with [Halloy](https://halloy.squidowl.org/index.html) and [Emacs IRC client](https://www.gnu.org/software/emacs/erc.html)

Another useful way to test is to use both clients. For example:
1. Start the chat server
2. Start the "Test cli client". Register and join a channel, e.g. `JOIN #general`
3. Start the "Halloy" client and see that `#general` exists. Join it and start chatting with the existing user
4. Send `PRIVMSG` between both clients and see that it works

### Halloy

To [configure](https://halloy.squidowl.org/configuration/index.html) Halloy to work locally with zirc, set this in your Halloy configuration `config.toml`: 

```toml
[servers.zirc]
nickname = "halloy9478"
server = "localhost"
port = 6667
use_tls = false
```

### Test cli client

Note the test cli is just that, a testing tool. It doesn't accept high level user commands like `/msg` for `PRIVMSG` - it only accepts raw protocol commands directly.

#### Docker compose use

There is a test client that is spun up as part of the docker compose build. The advantage of using the docker compose client is that it makes it easier to test a multi-server environment where you can use the direct docker compose service name or docker network static ip.

1. Exec into the test client container
2. Run `node test_client.js` to start it
3. Set the host name of the irc server you want to connect to (see options)
4. Run `connect` to connect to it.

#### Local use

1. Navigate to the test client dir

```bash
cd $CHAT_ROOT/Chat_Server/TestClient
```

2. Run the test client

```bash
node test_client.js
```

3. Enter `connect` to connect to the chat server running in docker
4. Enter any valid commands for the chat server.

For example:

```bash
NICK test-user
USER test-user
```

The server response should tell you that you are now registered.
