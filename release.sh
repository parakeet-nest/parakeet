#!/bin/bash
set -o allexport; source release.env; set +o allexport

: <<'COMMENT'
Todo:
- update of release.env:
  - TAG
  - ABOUT
COMMENT

read -r -d '' PARAKEET_MAIN <<- EOM
package parakeet

var Version = "${TAG}"
var About = "${ABOUT}"
EOM

echo "${PARAKEET_MAIN}" > parakeet.go

echo "Generating release: ${TAG} ${ABOUT}"

find . -name '.DS_Store' -type f -delete

git add .
git commit -m "📦 ${ABOUT}"
git tag -a ${TAG} -m "${ABOUT}"
git push origin ${TAG}

