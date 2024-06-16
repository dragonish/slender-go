#!/bin/bash

echo extract untranslated messages...
goi18n extract -outdir ./internal/local/i18n
goi18n merge -outdir ./internal/local/i18n ./internal/local/i18n/active.*.toml
echo translate.*.toml with the messages to be translated.
