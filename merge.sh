#!/bin/bash

echo merge translated files...
goi18n merge -outdir ./internal/local/i18n ./internal/local/i18n/active.*.toml ./internal/local/i18n/translate.*.toml
rm -f ./internal/local/i18n/translate.*.toml
echo merged and deleted translate.*.toml
