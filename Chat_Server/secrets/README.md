# Secrets

All secrets are bcrypt hashes and generally can be read through the files in this directory.

## Fallbacks

File secrets are the default, but the config code will fall back to trying to find the secrets via env vars if the file doesn't exist via:
- `IRC_S2S_PASSWORD`
- `IRC_SERVER_PASSWORD`
- `IRC_DB_PASSWORD`

## Types

### S2S

It is **highly recommended** that your server has a password for S2S communication. Otherwise any client would be able to send commands to the S2S port. There are additional layers of security which would probably prevent anything from happening (IP whitelisting, code that checks if `isServer` flag is set), but you should still add one in production.

#### Steps

1. Ensure the server you plan to connect to has its password set accordingly. This password can be set via:

   - The `s2s.password_file` key in `config.yaml`. This is the path to the secret file. This is the default method.
   - `IRC_S2S_PASSWORD` env variable.

2. Add file here with secret value

3. Open the server's S2S config file. E.g. the main server uses `s2s_config.yaml`, while the `animalhouse` test server uses `animalhouse_s2s_config.yaml`. This similar type of configuration can be replicated any number of times to be used with any arbitrary number of servers. The server finds this config file via the path set in the env var `IRC_S2S_CONFIG_FILE` or will default to the path `/chat_server/s2s_config.yaml`.

4. Add servers to the config file. The file consists of yaml array items found at the root indentation level.

```yaml
- host: "zirc_animalhouse"
  port: "7000"
  ip: "172.20.0.6"
  auto_connect: false
  # password_file: "/run/secrets/s2s_password"
  # password: "$2a$10$7w82iJDU7941avC65WYLJuJL0p7hIaGVye8pq95farNt7DD5dDGga"
```

| config | type | description |
| -- | -- | -- |
| host | string | The DNS host name of the server. |
| port | string | The port which the server is listening on |
| ip | string | The IP address of the server. In Docker use the `networks.<subnet>.ipv4_address` configuration to choose an IP. This configuration is only used because of built in IP whitelisting that each server has. |
| auto_connect | bool | Whether to automatically connect to the listed server. For example, if you just want to whitelist a server for future connections, use this. |
| password_file | string | The path in the Docker container where to read a mounted password file from. If both `password_file` and `password` are set, `password_file` will be used. |
| password | string | A password to use when sending `PASS` to the server. |