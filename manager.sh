#!/bin/bash
manager_file='slender-manager.zip'

echo clear manager directory
find ./web/manager -type f ! -name '.gitkeep' -delete

echo download slender-manager...
response=$(curl -s 'https://api.github.com/repos/dragonish/slender-manager/releases/latest')
download_url=$(echo "$response" | jq -r ".assets[0].browser_download_url")
curl -o "$manager_file" -JLO "$download_url" --header 'Host: github.com'
# or
# wget -O "$manager_file" "$download_url"
echo download done.

echo unzip files...
unzip "$manager_file" -d ./web/manager/
rm -f "$manager_file"
echo All done.
