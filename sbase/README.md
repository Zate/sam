# sbase - Splunkbase Downloader

---

## About

`sbase` is the beginnings of a download manager for Splunkbase apps and eventually apps from other locations also (such as S3, Gitlab and GutHub).

---

## Requirements

- .env file laid out similar to env.example with your Splunkbase username and password
  - `SBASE_P`
  - `SBASE_U`

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

Running sbase to download an App (<https://splunkbase.splunk.com/app/2890/>)

```console
sam/sbase $ ./sbase -a 2890
Download for 2890 - Splunk Machine Learning Toolkit Version: 5.2.0 is complete!
Files located at apps/2890/Splunk_ML_Toolkit/5.2.0/
```
