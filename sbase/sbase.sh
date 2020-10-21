#!/usr/bin/env bash
source .env
# https://splunkbase.splunk.com/app/2890/ 
BASE_URL= 
# we should check this is installed, but really, this is a POC and I will mgrate this into better code later.
docker-compose up -d
AUTH=`curl -sS -d "username=${SBASE_U}&password=${SBASE_P}" -X POST https://splunkbase.splunk.com/api/account:login/ | grep -o -P '(?<=<id>).*(?=</id>)'`
VER=`curl -sS -k https://splunkbase.splunk.com/app/${1}/| grep -oP '(?<=<sb-release-select u-for="download-modal" sb-selector="release-version" sb-target=").*(?=" )' | head -1`
curl -vvv -k -u admin:stuff123 -d name=https://splunkbase.splunk.com/app/${1}/release/${VER}/download/ -d update=true -d filename=true -d auth=${AUTH} https://localhost:8089/services/apps/local/


# lets grab some info from the API so we know what the folder name should be.
#curl -sS -k https://splunkbase.splunk.com/api/v1/app/${1}/
# ToDo
# Mount etc/apps to local file system via a splunk container.
