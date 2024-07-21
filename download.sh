#!/bin/bash

mdi_json='./assets/material-design-icons/index.json'
mdi_js='./assets/material-design-icons/index.js'
mdi_json_url='https://registry.npmjs.org/@mdi/js/latest'
mdi_js_url='https://gcore.jsdelivr.net/npm/@mdi/js@latest/mdi.js'

si_json='./assets/simple-icons/index.json'
si_js='./assets/simple-icons/index.js'
si_json_url='https://registry.npmjs.org/simple-icons/latest'
si_js_url='https://gcore.jsdelivr.net/npm/simple-icons@latest/index.js'

GetVersion() {
  if command -v jq &>/dev/null; then
    cat $1 | jq -r '.version'
  else
    sed -n 's/.*"version":"\([^"]\+\)".*/\1/p' $1
  fi
}

GetTagName() {
  if command -v jq &>/dev/null; then
    cat $1 | jq -r '.tag_name'
  else
    awk '/"tag_name"/ {print $2}' $1 | grep -o '[^",]*'
  fi
}


GetDownloadUrl() {
  if command -v jq &>/dev/null; then
    cat $1 | jq -r ".assets[0].browser_download_url"
  else
    awk '/"browser_download_url"/ {print $2}' $1 | grep -o '[^"]*'
  fi
}


DownloadIcon() {
  old_versoin="0.0.0"
  if [ -f $1 ]; then
    old_versoin=$(GetVersion $1)
  fi

  curl -o $1 $3
  new_version=$(GetVersion $1)

  echo "latest version: $new_version"

  if [ "$old_versoin" != "$new_version" ]; then
    curl -o $2 $4
  fi
}

DownloadManager() {
  manager_json='./latest.json'
  manager_file='slender-manager.zip'
  manager_latest_url='https://api.github.com/repos/dragonish/slender-manager/releases/latest'

  old_tag="v0.0.0"
  if [ -f $manager_json ]; then
    old_tag=$(GetTagName $manager_json)
  fi

  curl -o $manager_json $manager_latest_url
  new_tag=$(GetTagName $manager_json)

  echo "latest version: $new_tag"

  if [ "$old_tag" != "$new_tag" ]; then
    download_url=$(GetDownloadUrl $manager_json)
    curl -o "$manager_file" -JLO "$download_url" --header 'Host: github.com'
    # or
    # wget -O "$manager_file" "$download_url"
    echo download done.

    echo clear manager directory
    find ./web/manager -type f ! -name '.gitkeep' -delete

    echo unzip files...
    unzip "$manager_file" -d ./web/manager/
    rm -f "$manager_file"
    echo unzip done.
  fi
}

echo download material-design-icons...
DownloadIcon $mdi_json $mdi_js $mdi_json_url $mdi_js_url
echo download simple-icons...
DownloadIcon $si_json $si_js $si_json_url $si_js_url
echo download manager...
DownloadManager
echo all done.
