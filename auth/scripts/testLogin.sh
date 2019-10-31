#!/bin/bash

EMAIL="test@test.com"
PASSWORD="asdf"

echo "Email:    ${EMAIL}"
echo "Password: ${PASSWORD}"

BODY='{"email": "'"${EMAIL}"'", "password": "'"${PASSWORD}"'" }' 

RESPONSE=$(curl -s -H "Host:events-demo.localhost" -d "${BODY}" localhost/api/auth/login)
TOKEN=$(echo "${RESPONSE}" | jq -r '.token')

echo "Token:    ${TOKEN}"

# echo ''
# echo ${RESPONSE}

export TOKEN

