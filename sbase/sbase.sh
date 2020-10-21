#!/usr/bin/env bash
source .env
# https://splunkbase.splunk.com/app/2890/ 
BASE_URL= 
# we should check this is installed, but really, this is a POC and I will mgrate this into better code later.
docker-compose up -d
AUTH=`curl -sS -d "username=${SBASE_U}&password=${SBASE_P}" -X POST https://splunkbase.splunk.com/api/account:login/ | grep -o -P '(?<=<id>).*(?=</id>)'`
VER=`curl -sS -k https://splunkbase.splunk.com/app/${1}/| grep -oP '(?<=<sb-release-select u-for="download-modal" sb-selector="release-version" sb-target=").*(?=" )' | head -1`
curl --write-out '%{http_code}' --silent --output /dev/null -k -u admin:stuff123 -d name=https://splunkbase.splunk.com/app/${1}/release/${VER}/download/ -d update=true -d filename=true -d auth=${AUTH} https://localhost:8089/services/apps/local/?output_mode=json | jq -r .


# lets grab some info from the API so we know what the folder name should be.
APPID=`curl -sS -k https://splunkbase.splunk.com/api/v1/app/${1}/ | jq -r .appid`
# Going to put in a couple of quick things here to export it as a spl and copy it out, unsure if that format will be useful later or not.
#curl -sS -k -u admin:stuff123 https://localhost:8089/services/apps/local/${APPID}/package?output_mode=json | jq -r .
docker cp sbase_so1_1:/opt/splunk/etc/apps/${APPID} ${APPID}
tar -zcf ${APPID}.tar.gz ${APPID}
sudo rm -rf ${APPID}
# cp -rp apps/${APPID} /tmp/
# ls -al /tmp/${APPID}
# ToDo
# Mount etc/apps to local file system via a splunk container.
