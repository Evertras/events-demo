#!/bin/bash

EMAIL=$(uuidgen)@tst.com
PASSWORD=$(uuidgen)

echo "Email:    ${EMAIL}"
echo "Password: ${PASSWORD}"

BODY='{"email": "'"${EMAIL}"'", "password": "'"${PASSWORD}"'" }' 

RESPONSE=$(curl -s -H "Host:events-demo.localhost" -d "${BODY}" localhost/api/auth/register)
TOKEN=$(echo "${RESPONSE}" | jq -r '.token')

echo "Token:    ${TOKEN}"

export TOKEN

