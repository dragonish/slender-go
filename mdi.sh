#!/bin/bash

echo download material-design-icons...
curl -o ./assets/material-design-icons/index.json https://registry.npmjs.org/@mdi/js/latest
curl -o ./assets/material-design-icons/index.js https://gcore.jsdelivr.net/npm/@mdi/js@latest/mdi.js
echo download done.
