#!/bin/bash
#
current_dir=$(dirname "$0")
"$SPLUNK_HOME/bin/splunk" cmd python3.7 "$current_dir/command.py" $@
