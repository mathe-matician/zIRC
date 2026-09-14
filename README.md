# zIRC

A Chat Application based on IRC v3 protocol spec.

## Purpose

This project was a learning exercise for myself to practice implementing RFCs to code and understanding network protocols better. Part of the objectives were to hand write all code, but use AI to understand the IRC spec and get direction where needed.

### What I'd do differently

If I were starting this project again, I'd change the following:

1. Global state. I have a number of areas where I basically designed myself into relying on global state (`g_` or `G_` prefixed state). Starting from scratch, I would focus on designing this cleaner; passing that required state for each owning component instead of accessing it globally.

2. Reliance on `map[string]interface{}` for one-place-get-all for state passed to commands. This ultimately made for really clunky interactions for every single command. I essentially originally did this by trying to be too generic in the main handler to pass arbitrary state to commands. It then also led to missed type assertions which could lead to issues down the road - not to mention it isn't really readable. Starting from scratch, I wouldn't have focused so much on making one generic callable entrypoint for all commands.

3. Structs holding `sync.Mutex` are being copied by value, so they don't actually guard the shared data. Starting from scratch, things like `command_map`, `server_manager_commands`, and `ctcp_map` should be `map[string]*Command` instead of `map[string]Command`.

4. Significantly more TDD / test focus. It would have been better to understand the protocol more and write some tests first instead of doing large implementation chunks and then only then deciding I wanted to write tests.

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

Prereqs:
- Docker installed and running

1. Navigate to the repo root

```bash
cd $CHAT_ROOT
make dev
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

Prereqs:
- Node must be installed

1. Navigate to the test client dir

```bash
cd $CHAT_ROOT
make testclient
```

3. Enter `connect` to connect to the chat server running in docker
4. Enter any valid commands for the chat server.

For example, enter each of these commands in order and press enter after them:

```bash
NICK test-user
USER test-user
JOIN #general
PRIVMSG #general :whats up y'all
```

The server response should tell you that you are now registered.
