#!/bin/bash

Extract() {
  echo extract untranslated messages...
  goi18n extract -outdir ./internal/local/i18n
  goi18n merge -outdir ./internal/local/i18n ./internal/local/i18n/active.*.toml
  echo translate.*.toml with the messages to be translated.
}

Merge() {
  echo merge translated files...
  goi18n merge -outdir ./internal/local/i18n ./internal/local/i18n/active.*.toml ./internal/local/i18n/translate.*.toml
  rm -f ./internal/local/i18n/translate.*.toml
  echo merged and deleted translate.*.toml
}


if ! command -v goi18n &>/dev/null; then
  echo "No goi18n tool found!"
  echo "Please run the command: go install -v github.com/nicksnyder/go-i18n/v2/goi18n@latest"
  exit 3
fi


echo 'Select the operation: '
PS3='Pick an option: '
options=("Extract untranslated messages" "Merge translated files")

select opt in "${options[@]}" "Quit"; do
  case "$REPLY" in
    1)
      Extract
      break;;
    2)
      Merge
      break;;
    
    $(( ${#options[@]}+1 )) ) echo "$opt"; break;;
    *) echo "Invalid option. Try another one."; continue;;
  esac
done
