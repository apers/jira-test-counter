#!/bin/bash
GOPATH=/home/pers/go go run /home/pers/go/src/jira-test-counter/*.go >> /var/log/jirs-count.log 2>&1
