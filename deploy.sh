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

ssh $HOST "sudo systemctl stop my_blog_api"
scp ./main "$HOST:/home/naoki_yoneyama/my_blog_api/main"
ssh $HOST "sudo systemctl start my_blog_api"
