#!/bin/bash

# Get the directory of this script
# https://stackoverflow.com/questions/59895/getting-the-source-directory-of-a-bash-script-from-within
DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )"

# Get the name of the repository
# https://stackoverflow.com/questions/23162299/how-to-get-the-last-part-of-dirname-in-bash/23162553
REPO="$(basename "$DIR")"

# For non-interactive shells (e.g. the server running this script to build itself),
# the "HOME" environment variable must be specified or there will be a cache error when compiling
# the Go code (but don't do this in Travis, since doing this will cause it to break)
if [[ -z $HOME ]] && [[ -z $CI ]]; then
  export HOME=/root
fi

"$DIR/build.sh"
if [[ $? -ne 0 ]]; then
  exit 1
fi
cd "$DIR"
"$DIR/$REPO"
