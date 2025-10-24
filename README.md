Server for deploying my services by triggering from API.

# Installation

## 1. Build or download the binary

Ready build releases can be found at [github](https://github.com/Ben-PP/publy/releases).

## 2. Create directories

Publy reads configs from `/etc/publy/config.yaml` so you have to create this
directory and file.

```bash
sudo mkdir -p /etc/publy
```

Then either copy the `example-config.yaml` to the directory as `config.yaml` or
create the `config.yaml` from scratch.

#### config.yaml

```yaml
host: localhost # Address/domain which the server listens. Use 0.0.0.0 for all.
port: 8000 # Port on which the server listens.
script-dir: scripts # Directory where the scripts to be executed are stored.
proxies: # List of addresses or CIDR notations for trusted proxies.
  - localhost
ssl:
  enabled: true # Use SSL
  certificate: /etc/ssl/certs/server.crt # Server certificate file
  key: /etc/ssl/private/server.key # Certificate key file

pubs: # Any number of pubs, which define what is executed when API is called.
  example-service: # Name of the pub. This is given as query parameter when calling the API.
    script: example.sh # Name of the script to be executed when publishing the pub.
    token-hash: # bcrypt hash of token chosen

  another-service:
    script: another.sh
    token-hash:
```

# Usage

Publy is used for automating application deployment on selfhosted environment.
It runs bash scripts based on API calls made to `GET /api/v1/publish?pub=example-service`
endpoint. Calls are authenticated via user generated tokens which are stored in
the config file hashed with bcrypt. For hashing you can use
`POST /api/v1/generate-hash` endpoint and giving the string to hash in the body
`{"string": "example-token"}`. This endpoint can only be called from the private
IPv4 ranges and localhost.

## Pub

Pub is a single application to be published. These are defined in the config.yaml
under the `pubs` key. See above `config.yaml` example.

## Publishing a pub

To run a publish, make a `GET` request to `/api/v1/publish`. You have to specify
query parameter `pub` for which the value must be one of the pubs declared in the
config. You will also have to provide correct token in Authorization header as a
Bearer token and it must match the hash given for the pub in the config.
