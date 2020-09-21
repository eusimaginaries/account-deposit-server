#!/bin/bash
for cmd in "$@"
do
  case $cmd in
    run)
      echo "* performing run *"
      ./scripts/build.sh install test
      $GOPATH/bin/account-deposit-server
  esac
done