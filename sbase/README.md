# sbase - Splunkbase Downloader

---

## About

`sbase` is the beginnings of a download manager for Splunkbase apps and eventually apps from other locations also (such as S3, Gitlab and GutHub). It uses a local Splunk Server spun up via docker-compose.

---

## Requirements

- docker
- docker-compose
- jq
- .env file laid out similar to env.example with your Splunkbase username, password and the IP of the local Splunk Server (will default to default)
  - `SBASE_P`
  - `SBASE_U`
  - `SPLUNK_SERVER`

---

## Usage

Build and Show Command Line switches

```console
sam/sbase $ go build

sam/sbase $ ./sbase
        -a string   AppID to Download
        -d          Turn on Debug
```

---

Start up the Splunk Server

```console
sam/sbase $ cd testing
sam/sbase/testing $ docker-compose up -d
Creating network "testing_splunknet" with driver "bridge"
Creating testing_so1_1 ... done
```

---

Tear down the container and delete it's files.

```console
sam/sbase/testing $ ./killit.sh
Stopping testing_so1_1 ... done
Removing testing_so1_1 ... done
Removing network testing_splunknet
```

---

Running sbase to download an App (<https://splunkbase.splunk.com/app/2890/>)

```console
sam/sbase $ ./sbase -a 2890
```
