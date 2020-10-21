# sbase shell script

Right now this is just a testing ground, a collection of information I have pieced together to try and be able to give an APPID, and get the downloaded package.  

Ideally I would like to move this to a go binary/tool/library to be used in a larger web ui based tool.

It contains random postings of info I have found, discoveries I have made simply plugging in values to API's and things I have guessed, or found in other scripts etc.

SBASE_U is your Splunkbase Username
SBASE_P is your Splunkbase password

Stick them in a .env file (that does not get commited to git) and sun `./sbase.sh`
Pray it's not broken.

## Requirments

- Needs docker and docker-compose
- Needs a username for splunkbase
- Might need your username to have the API turned on via <https://dev.splunk.com/enterprise/docs/releaseapps/splunkbase/submitcontentrestapi>

## To Do

- [ ] What I really want to know is how to recreate the request the splunk binary does to download an app.  It's using the auth token we get from the API, but if I could remove needing the splunk binary/splunk docker container, that would be great.
