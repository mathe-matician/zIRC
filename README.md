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
4. Enter any valid commands for the chat server.

For example:

```bash
NICK test-user
USER test-user
```

The server response should tell you that you are now registered.