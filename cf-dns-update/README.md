# cf-dns-update

update Cloudflare DNS A & CNAME records

[![License](https://img.shields.io/github/license/seankhliao/cf-dns-update.svg?style=for-the-badge&maxAge=31536000)](LICENSE)
[![Build](https://badger.seankhliao.com/i/github_seankhliao_cf-dns-update)](https://badger.seankhliao.com/l/github_seankhliao_cf-dns-update)

## About

initContainer for deployments in my k8s [kluster](https://github.com/seankhliao/kluster) to update DNS records as deployments happen

- `A`: update records to point to your current IP
  - sort existing records by age, deletes more than `REPLICAS`, updates oldest
- `CNAME`: update records to point to `CONTENT`

## Usage

#### Prerequisites

- cloudflare
- ipify.org
  - external ip address

Env vars to set:

```
RECORD=seankhliao.com::badger:CNAME:kluster.seankhliao.com
RECORD=seankhliao.com:noproxy:wg:CNAME:kluster.seankhliao.com
RECORD=seankhliao.com::kluster:A:1
RECORD = [zone] : [proxy] : [name] : [type] : [content / replicas]

X_AUTH_EMAIL = Email for cloudflare api
X_AUTH_KEY = Api key for cloudflare api
````

#### Run

after setting all the env vars

```
cf-dns-update
```
