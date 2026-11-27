#!/usr/bin/env bash

if (( $# == 0 )); then
    printf 'Usage: %s COMMAND [ARG ...]\n' "$0" >&2
    exit 2
fi

# RESP bulk string lengths count bytes.
export LC_ALL=C

{
    printf '*%d\r\n' "$#"
    for arg in "$@"; do
        printf '$%d\r\n%s\r\n' "${#arg}" "$arg"
    done
} | nc -w 1 "${REDIS_HOST:-127.0.0.1}" "${REDIS_PORT:-6379}"
