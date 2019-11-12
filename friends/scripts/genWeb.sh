#!/bin/bash

NUM_USERS=50
NUM_INVITES=$(( NUM_USERS * 4 ))

# Parallel data hype
EMAIL_LIST=()
PASSWORD_LIST=()
TOKEN_LIST=()

register() {
  local EMAIL=$1
  local PASSWORD=$2

  local BODY='{"email": "'"${EMAIL}"'", "password": "'"${PASSWORD}"'" }' 
  echo Doing
  local RESPONSE=$(curl -s -H "Host:events-demo.localhost" -d "${BODY}" localhost/api/auth/register)

  echo "${RESPONSE}" | jq -r '.token'
}

login() {
  local EMAIL=$1
  local PASSWORD=$2

  local BODY='{"email": "'"${EMAIL}"'", "password": "'"${PASSWORD}"'" }' 
  local RESPONSE=$(curl -s -H "Host:events-demo.localhost" -d "${BODY}" localhost/api/auth/login)

  echo "${RESPONSE}" | jq -r '.token'
}

invite() {
  local TOKEN=$1
  local TARGET=$2

  local BODY='{"email":"'"${TARGET}"'"}'
  local RESPONSE=$(curl -s -H "Host:events-demo.localhost" -H "X-Auth-Token: ${TOKEN}" -d "${BODY}" localhost/api/friends/invite)
}

for (( i=0; i < NUM_USERS; i++ ))
do
  EMAIL=$(uuidgen)@genweb.com
  PASSWORD=$(uuidgen)

  EMAIL_LIST+=( ${EMAIL} )
  PASSWORD_LIST+=( ${PASSWORD} )
  TOKEN_LIST+=( $(register ${EMAIL} ${PASSWORD}) )

  echo "Registered ${EMAIL}"
done

for (( i=0; i < NUM_INVITES; i++ ))
do
  FROM_INDEX=$((RANDOM % ${#EMAIL_LIST[@]}))
  TO_INDEX=$((RANDOM % ${#EMAIL_LIST[@]}))

  echo "${EMAIL_LIST[$FROM_INDEX]} is inviting ${EMAIL_LIST[${TO_INDEX}]}"

  invite ${TOKEN_LIST[$FROM_INDEX]} ${EMAIL_LIST[$TO_INDEX]}
done

