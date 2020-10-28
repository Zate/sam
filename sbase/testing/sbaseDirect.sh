#!/usr/bin/env bash
source .env

AUTH=`curl -sS -d "username=${SBASE_U}&password=${SBASE_P}" -X POST https://splunkbase.splunk.com/api/account:login/ | grep -o -P '(?<=<id>).*(?=</id>)'`
VER=`curl -sS -k https://splunkbase.splunk.com/app/${1}/| grep -oP '(?<=<sb-release-select u-for="download-modal" sb-selector="release-version" sb-target=").*(?=" )' | head -1`
curl -vvv -L -k  -H "X-Auth-Token: ${AUTH}" -X GET https://splunkbase.splunk.com/app/${1}/release/${VER}/download/ --output ${1}.tar.gz
# --referer https://192.168.0.16:8089/server/apps/local/ -A "Splunkd/8.1.0 (Linux; version=3.14; arch=x86_64; build=1.6; 3.7)"
