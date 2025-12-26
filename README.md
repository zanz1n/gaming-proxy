# gaming-proxy

A simple proxy application that exposes a local game server instance to the internet.
It also sets SRV DNS records to facilitate public access.

## Use case

You have a game server (tipically minecraft) that you want to expose publically through a secure tunnel and a domain name in cloudflare that you want to automatically point to this server.

## Configuration

```sh
PROXIED_HOST=127.0.0.1 # Domain names or ip addresses
PROXIED_PORT=12345

CLOUDFLARE_TOKEN="" # API token

# Zone = the root domain registred in cloudflare
CLOUDFLARE_ZONE_ID=""

# In this example you can access the minecraft server in `play.domain.com`
CLOUDFLARE_SUB_DOMAIN="play"

# The game protocol
CLOUDFLARE_SERVICE="minecraft"

# If set to true, an existent DNS record may be replaced
CLOUDFLARE_OVERWRITE=true

# Currently only ngrok is supported
NGROK_TOKEN="" # API token
```

With this configuration, people may connect to your minecraft server accessing `play.domain.com`.

To get the cloudflare Zone ID for your `domain.com` and API token, you can search for online tutorials.

## Running

### Linux

```sh
curl -L -o ./gaming-proxy https://github.com/zanz1n/gaming-proxy/releases/latest/download/proxy-$(uname -s | tr '[:upper:]' '[:lower:]')-$(uname -m)
```

**Create the .env file in the `.` directory based in the [configuration example](./.env.example)**

```sh
./gaming-proxy
```

### MacOS or Windows

Download the latest executable [here](https://github.com/zanz1n/gaming-proxy/releases/latest).
