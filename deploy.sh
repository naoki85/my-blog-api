#!/bin/bash
HOST=""

while getopts h: OPT; do
  case $OPT in
    "h") HOST="$OPTARG" ;;
  esac
done

if [[ "$HOST" == "" ]]; then
  echo "Parameter must be given."
  exit 1
fi

echo "Start build"
GOOS=linux GOARCH=amd64 go build main.go
echo "Finish build and start stopping service"
ssh $HOST "sudo systemctl stop my_blog_api"
echo "Start uploading file"
scp ./main "$HOST:/home/naoki_yoneyama/my_blog_api/main"
echo "Restart service"
ssh $HOST "sudo systemctl start my_blog_api"
echo "Finish!!"
