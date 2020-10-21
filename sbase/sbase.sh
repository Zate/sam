#!/usr/bin/env bash
source .env
# we should check this is installed, but really, this is a POC and I will mgrate this into better code later.
docker-compose up -d
AUTH=`curl -sS -d "username=${SBASE_U}&password=${SBASE_P}" -X POST https://splunkbase.splunk.com/api/account:login/ | grep -o -P '(?<=<id>).*(?=</id>)'`
VER=`curl -sS -k https://splunkbase.splunk.com/app/2890/ | grep -oP '(?<=<sb-release-select u-for="download-modal" sb-selector="release-version" sb-target=").*(?=" )' | head -1`
APPURL="${1}/release/${VER}/download/"
curl -vvv -k -u admin:stuff123 -d name=${APPURL} -d update=true -d filename=true -d auth=${AUTH} https://localhost:8089/services/apps/local/

# ToDo
# Mount etc/apps to local file system via a splunk container.

#https://splunkbase.splunk.com/api/v1/app/2890/release/

# {
# "id": 16431, <-- package_id ?
# "app": 2890, <--- can start with this.
# "name": "5.2.0", <-- get this
# "release_notes": "For the latest release notes, see http://docs.splunk.com/Documentation/MLApp/latest/User/Whatsnew",
# "CIM_versions": [],
# "splunk_versions": [
# 31,
# 25,
# 30
# ],
# "public": true,
# "public_ever_true": true,
# "created_datetime": "2020-05-09T00:07:29.055882Z",
# "published_datetime": "2020-05-09T00:07:29.055918Z",
# "size": 16524687,
# "filename": "splunk-machine-learning-toolkit_520.tgz", <-- useful.
# "platform": "independent",
# "is_bundle": true,
# "has_ui": false,
# "approved": false,
# "appinspect_status": false,
# "install_method_single": "simple",
# "install_method_distributed": "appmgmt_phase",
# "requires_cloud_vetting": true,
# "appinspect_request_id": null,
# "cloud_vetting_request_id": "9841b444-9a82-4254-9f77-7ab73107ef36",
# "python3_acceptance": true,
# "python3_acceptance_datetime": "2020-05-15T04:00:45.552000Z",
# "python3_acceptance_user": 235284,
# "fedramp_validation": "no",
# "cloud_compatible": true
# },


#https://splunkbase.splunk.com/api/v1/app/2890/
# {
# "uid": 2890,
# "appid": "Splunk_ML_Toolkit",
# "title": "Splunk Machine Learning Toolkit",
# "created_time": "2015-09-15T15:40:42+00:00",
# "published_time": "2015-09-15T15:40:42+00:00",
# "updated_time": "2020-05-09T00:07:29+00:00",
# "license_name": "Splunk Software License Agreement",
# "type": "app",
# "license_url": "http://www.splunk.com/en_us/legal/splunk-software-license-agreement.html",
# "description": "Splunk Machine Learning Toolkit\n\nThe Splunk Machine Learning Toolkit App delivers new SPL commands, custom visualizations, assistants, and examples to explore a variety of ml concepts. \n\nEach assistant includes end-to-end examples with datasets, plus the ability to apply the visualizations and SPL commands to your own data. You can inspect the assistant panels and underlying code to see how it all works.\n\nML Youtube Playlist http://tiny.cc/splunkmlvideos\nML Cheat Sheet http://tiny.cc/mltkcheatsheet\n\nAssistants:\n* Predict Numeric Fields (Linear Regression): e.g. predict median house values.\n* Predict Categorical Fields (Logistic Regression): e.g. predict customer churn.\n* Detect Numeric Outliers (distribution statistics): e.g. detect outliers in IT Ops data.\n* Detect Categorical Outliers (probabilistic measures): e.g. detect outliers in diabetes patient records.\n* Forecast Time Series: e.g. forecast data center growth and capacity planning.\n* Cluster Numeric Events: e.g. Cluster Hard Drives by SMART Metrics\n\nSmart Assistants (new assistants with revamped UI and better ml pipeline/experiment management):\n*Smart Forecasting Assistant (provides enhanced time-series analysis for users with little to no SPL knowledge and leverages the StateSpaceForecasting algorithm): e.g. forecasting app logons with special days\n\nAvailable on both on-premise and cloud.\n\nDeep Learning Toolkit for Splunk\nIntegrate with advanced custom machine learning systems using the Deep Learning Toolkit for Splunk (https://splunkbase.splunk.com/app/4607/). It extends Splunk’s Machine Learning Toolkit with prebuilt Docker containers for TensorFlow 2.0, PyTorch and a collection of NLP libraries. Python expertise is required to create your own neural networks.\nAvailable only for on-premise customers.\n\nSplunk Community for MLTK Algorithms on GitHub\nCheck out our Open Source community on Github that lets you share your algorithms with the community of Splunk MLTK users or import one of the algorithms that have been shared by the community: https://github.com/splunk/mltk-algo-contrib\n\nThe GitHub repo algorithms are also available as an app which provides access to custom algorithms. Cloud customers can use GitHub algorithms via this app and need to create a support ticket to have this installed:https://splunkbase.splunk.com/app/4403/\nAvailable on cloud and on-premise",
# "access": "unrestricted",
# "appinspect_passed": false,
# "path": "https://splunkbase.splunk.com/app/2890/",
# "install_method_distributed": "appmgmt_phase",
# "install_method_single": "simple",
# "download_count": 84203,
# "install_count": 12333,
# "archive_status": "live",
# "is_archived": false,
# "fedramp_validation": "no"
# }

# id is package id?
# https://splunkbase.splunk.com/api/v1/app/2890/release/16431/
# {
# "app": 2890,
# "name": "5.2.0",
# "release_notes": "For the latest release notes, see http://docs.splunk.com/Documentation/MLApp/latest/User/Whatsnew",
# "splunk_versions": [
# "8.1",
# "8.0",
# "8.1.2008"
# ],
# "CIM_versions": [],
# "public": true
# }
