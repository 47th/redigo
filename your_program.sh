#!/bin/sh
set -e


cd "$(dirname "$0")" # Ensure compile steps are run within the repository directory

if [ "$#" -eq 1 ]
then
  if [ "$1" = "server" ]
  then
      go run app/main.go
  else
      go run client/client.go
  fi

else
  #normal compilation and run

  #compilation
    go build -o /tmp/redigo app/*.go

  #run
  exec /tmp/redigo "$@"
fi
