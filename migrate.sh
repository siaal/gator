#!/bin/bash
connstring="$(jq .db_url ~/.gatorconfig.json)"
connstring="${connstring%?sslmode=disable\"}"
connstring="${connstring#\"}"
rootdir="$(realpath "$(dirname "$0")")"
schemadir="$rootdir/sql/schema"
echo "$schemadir $connstring"
goose-up() {
    cd "$schemadir" &&
        goose postgres "$connstring" up
}
goose-down() {
    cd "$schemadir" &&
        goose postgres "$connstring" down
}

case $1 in
    up)
        goose-up
        ;;
    down)
        goose-down
        ;;
    reset)
        goose-down && goose-up
        ;;
    *)
        echo "Unknown: $1"
        ;;
esac
