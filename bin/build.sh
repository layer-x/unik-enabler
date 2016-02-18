#!/bin/bash

if [[ "$(which gox)X" == "X" ]]; then
  echo "Please install gox. https://github.com/mitchellh/gox#readme"
  exit 1
fi


rm -f unik-enabler*

gox -os linux -os windows -arch 386 --output="unik-enabler_{{.OS}}_{{.Arch}}"
gox -os darwin -os linux -os windows -arch amd64 --output="unik-enabler_{{.OS}}_{{.Arch}}"

rm -rf out
mkdir -p out
mv unik-enabler* out/
