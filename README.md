# assets-manager

## Overview

The repository contains a several subprojects:

- `API` - RESTful API for managing/validating assets.
- `Worker` - Github Bot for Pull Requests to Trust Wallet Assets repository. It receives Github events from MQ and does appropriate logic.
- `Web` - Web application to help adding new tokens and managing the Trust Wallet assets repository.

## Local development

Use docker compose to have queues and databases locally.

`docker-compose up -d`

**Rabbit MQ:** Go to [rabbitmq admin dashboard](http://localhost:15672) using default credentials (username: guest, password: guest).

### API Service

Run

``` sh
make go-build start-api
```

### Worker Service

Developing and testing of Worker service is a bit tricky. But here is a guide how you can do this.

**Webhook deliveries**
Github Bot receives Github events via webhook. But Github can't call you local webhook while you're developing.
So we need some tool for forwarding HTTP requests to our local service.

Install [SMEE](https://smee.io/).
Use CLI

```sh
npm install --global smee-client
```

Once you installed SMEE, create a channel on [website](https://smee.io/) and use a chanel's link to launch this in a separate terminal window. For example:

```sh
# Use flag -P to forward the request to appropriate API's endpoint
smee -P "/v1/github/events/webhook" -u https://smee.io/1yMiWzuRE3hPMeFP
```

**Github App**
First of all, you need to create your own [Github App](https://github.com/settings/apps).

- Use your SMEE link to set up a `Webhook URL`.
- In [Permissions & events](https://github.com/settings/apps/merge-fee-bot-test/permissions) give `Read & Write` access to your app.
- Subscribe to the same events as [Merge-Fee-Bot](https://github.com/organizations/trustwallet/settings/apps/merge-fee-bot).
- [Install](https://github.com/settings/apps/merge-fee-bot-test/installations) the app to your test repository.

After all, you will need to copy you App ID and generate/download a private key of your app. You should rename you private key file to this name github-private-key.pem (this name is set in .gitignore).

Run

``` sh
GITHUB_APP_PRIVATE_KEY=`cat github-private-key.pem` GITHUB_APP_ID=167859 make go-build start-worker
```
