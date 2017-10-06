#!/bin/bash
GOPATH=/home/pers/go go run /home/pers/go/src/jira-test-counter/*.go >> /var/log/jira-count.log 2>&1
