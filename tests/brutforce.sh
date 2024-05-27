#!/bin/bash

login=user1
password=regjkfd
ip="3.45.67.89"

echo $login $password $ip
curl -X 'GET' \
  'http://0.0.0.0:8000/hello' \
  -H 'accept: application/json' \
  -H 'Content-Type: application/json' \
  -d '{ "ip":"'$login'",  "login":"'$login'",  "password":"'$password'" }'