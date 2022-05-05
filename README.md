# assets-manager

## Overview

The repository contains a several subprojects:

- `API` - RESTful API for managing/validating assets.
- `Consumer` - Github Bot for Pull Requests to Trust Wallet Assets repository. It receives Github events from MQ and does appropriate logic.
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

### Consumer Service

Developing and testing Consumer service is a bit tricky. Here is a guide on how you can do this.

**Webhook deliveries**

Github Bot receives Github events via webhook. But Github can't call your local webhook while you're developing.

So we need some tool for forwarding HTTP requests to our local service.

Install [SMEE](https://smee.io/).
Use CLI

```sh
npm install --global smee-client
```

Once you installed SMEE, create a channel on [website](https://smee.io/) and use a channel's link to launch this in a separate terminal window.

For example:

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

After all, you will need to copy your `App ID` and generate/download a private key of your app. You should rename your private key file to this name `github-private-key.pem` (this name is set in .gitignore).

Run

``` sh
GITHUB_APP_PRIVATE_KEY=`cat github-private-key.pem` GITHUB_APP_ID=167859 make go-build start-consumer
```

**The most common cases from moderators**

A lot of of things for assets management you can control via [config](https://github.com/trustwallet/assets-manager/blob/main/config.yml).
For example:

- adding moderator to the list if moderators to turn off fee requirement for them
- updating fee amount
- updating holders amount
- validation requirements and etc

***Supporting new token type***

go-primitives

1. Define this token type in [constants](https://github.com/trustwallet/go-primitives/blob/master/types/token.go)

2. Add this token type to the func [GetTokenType](https://github.com/trustwallet/go-primitives/blob/master/types/token.go) + case in unit test

3. Add this token type to the func [GetTokenTypes](https://github.com/trustwallet/go-primitives/blob/master/types/token.go)

4. Add explorer link for this token to the func [GetCoinExploreURL](https://github.com/trustwallet/go-primitives/blob/da50809e2a612d4a32cb0824f21653da3661801e/coin/models.go) + case in unit test

5. Add this token type to the func [GetChainFromAssetType](https://github.com/trustwallet/go-primitives/blob/master/types/chain.go)

6. After merging to the master, create a tag of package version and make `go get github.com/trustwallet/go-primitives@YOUR_VERSION` in [assets-go-libs](https://github.com/trustwallet/assets-go-libs), [assets-manager](https://github.com/trustwallet/assets-manager) and [assets](https://github.com/trustwallet/assets)

***Increasing or decreasing the fee for PRs***

Sometimes, due to activity of users and huge amount of open PR-s, moderators suggest to increase the amount of required fee.
Example of PR: https://github.com/trustwallet/assets/pull/20021#issuecomment-1107364308

In order to change this amount you either need to update this value on Heroku in ENVs or change it in [config](https://github.com/trustwallet/assets-manager/blob/main/config.yml#L50)

***Adding new tags***

All tags are placed in [config](https://github.com/trustwallet/assets-manager/blob/61ce41b86f925397aecfae046c0ea0c3f6150648/config.yml#L117). Just update it.

***Removing tokens***

All tokens can not be removed from assets repository.
Moderators have to change `status` field from `active` to `abandoned`. The bot warns about it but if they merge a PR with removed `info.json` files - then backend worker can be broken. To fix it - you need to revert changes and then mark these tokens as `abandoned`.



