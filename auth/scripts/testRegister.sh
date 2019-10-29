#!/bin/bash

EMAIL=$(uuidgen)@tst.com
PASSWORD=asdf9314

BODY='{"email": "'"${EMAIL}"'", "password": "'"${PASSWORD}"'" }' 

curl -s -H "Host:events-demo.localhost" -d "${BODY}" localhost/api/auth/register

