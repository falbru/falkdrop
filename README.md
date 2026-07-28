# FalkDrop

Send and receive files across devices

## Features

- Upload files or clipboard to a drop which can be shared with a link

- Download a drop using a code, link or QR code

- Drops expire after a certain time period

- Authentication and authorization to ensure only authorized users can create drops

## Quick Start

1. In `./deployments/local/`, copy the `.env.example` to `.env` and fill out the values

```sh
cd ./deployments/local/
cp .env.example .env
$EDITOR .env
```

2. Compose up all the docker services

```sh
cd ./deployments/local/
docker compose up -d
```

3. Configure Keycloak as described in [Keycloak Setup](./doc/keycloak.md).

4. Set up the object store as described in [Object Store Setup](./doc/objectstore.md).

5. Set up the server as described in [Server Setup](./server/README.md).

6. Set up the webapp as described in [Webapp Setup](./webapp/README.md).
