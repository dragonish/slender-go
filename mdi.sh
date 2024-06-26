#!/bin/bash

index_json="./assets/material-design-icons/index.json"
index_js="./assets/material-design-icons/index.js"

echo download material-design-icons...
curl -o $index_json https://registry.npmjs.org/@mdi/js/latest
curl -o $index_js https://gcore.jsdelivr.net/npm/@mdi/js@latest/mdi.js
echo download done.

if command -v jq &>/dev/null; then
  echo "material-design-icons version: $(cat $index_json | jq '.version')"
fi
