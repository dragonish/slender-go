#!/bin/bash

index_json="./assets/simple-icons/index.json"
index_js="./assets/simple-icons/index.js"

echo download simple-icons file...
curl -o $index_json https://registry.npmjs.org/simple-icons/latest
curl -o $index_js https://gcore.jsdelivr.net/npm/simple-icons@latest/index.js
echo download done.

if command -v jq &>/dev/null; then
  echo "simple-icons version: $(cat $index_json | jq '.version')"
fi
