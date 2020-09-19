#!/bin/bash
for cmd in "$@"
do
  case $cmd in
    install)
      echo "* performing install all *"
      go install
    ;;
    test)
      echo "* performing tests and coverage *"
      go test ./...  -timeout 30s -coverprofile=$GOPATH/bin/go-code-cover
    ;;
    tidy)
      echo "* tidying up project *"
      go mod tidy
    ;;
    lint)
      echo "* linting project *"
      $GOPATH/bin/golint ./...
  esac
done