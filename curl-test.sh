#!/usr/bin/env bash

SENSU_ACCESS_TOKEN="$1"

SENSUAPI="sensu-api.sensu.svc.cluster.local:8080"

asset() {
curl -X GET -k \
https://${SENSUAPI}/api/core/v2/namespaces/default/assets \
-H "Authorization: Key $SENSU_ACCESS_TOKEN" | jq .
}

check() {
curl -v -X GET -k \
https://${SENSUAPI}/api/core/v2/namespaces/default/checks \
-H "Authorization: Key $SENSU_ACCESS_TOKEN" | jq .
}

handler() {
curl -X GET -k \
https://${SENSUAPI}/api/core/v2/namespaces/default/handlers \
-H "Authorization: Key $SENSU_ACCESS_TOKEN" | jq .
}
hook() {
curl -X GET -k \
https://${SENSUAPI}/api/core/v2/namespaces/default/hooks \
-H "Authorization: Key $SENSU_ACCESS_TOKEN" | jq .
}

case $2 in
    "check")
        check 
    ;;
    "asset")
        asset
    ;;
    "handler")
        handler
    ;;
    "hook")
        hook
    ;;
    *)
        echo "$0 token-long-string check|asset|handler|hook"
esac
