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