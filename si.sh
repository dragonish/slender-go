#!/bin/bash

echo download simple-icons file...
curl -o ./assets/simple-icons/index.json https://registry.npmjs.org/simple-icons/latest
curl -o ./assets/simple-icons/index.js https://gcore.jsdelivr.net/npm/simple-icons@latest/index.js
echo download done.
