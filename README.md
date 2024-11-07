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

Both are useful for various reasons. The only existing IRC client zirc has been tested with is [Halloy](https://halloy.squidowl.org/index.html).

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

Note the test cli is just that, a testing tool. It doesn't accept high level user commands like `/msg` for `PRIVMSG` - it only accepts protocol commands directly.

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