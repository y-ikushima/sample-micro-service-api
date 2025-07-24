#!/bin/bash
#
# repo: https://github.com/sparanoid/env.sh
# related: https://github.com/vercel/next.js/discussions/17641
#
# env.sh - Dead-simple .env file reader and generator
# Tunghsiao Liu (t@sparanoid.com)
#
# Inspired by:
# - https://github.com/andrewmclagan/react-env
# - https://www.freecodecamp.org/news/7f9d42a91d70/
# - https://github.com/kunokdev/cra-runtime-environment-variables
#
# Features:
# - Designed to be used for Next.js app inside a Docker container (General
#   React app should also work)
# - Read env file and generate __env.js file for runtime client use
# - Merge current environment variables passing to it (Useful for Docker images)
# - No dependencies (More like a lite version of react-env). This is important
#   to keep our container image as small as possible.
#
# Usage:
# - General usage:
#   $ ./env.sh
#
# - Replacing varaible:
#   $ NEXT_PUBLIC_API_BASE=xxx ./env.sh
#
# - Enviroment variable not in whitelist will be discarded:
#   $ BAD_ENV=zzz ./env.sh
#
# - Change script options:
#   $ ENVSH_ENV_FILE="./.env.staging" ENVSH_OUTPUT="./public/config.js" ./env.sh
#
# - Use it inside Dockerfile:
#   RUN chmod +x ./env.sh
#   ENTRYPOINT ["./env.sh"]
#
# Debug:
# NEXT_PUBLIC_OB_ENV=123_from_fish NEXT_BAD_ENV=zzz NEXT_PUBLIC_OB_TESTNEW=testenv NEXT_PUBLIC_CODE_UPLOAD_SIZE_LIMIT=6666 ./env.sh

echo -e "env.sh loaded"

# Config
ENVSH_READ_ENV_VARS="${ENVSH_READ_ENV_VARS:-"false"}"
ENVSH_READ_ENV_FILE="${ENVSH_READ_ENV_FILE:-"false"}"
ENVSH_ENV_FILE="${ENVSH_ENV_FILE:-"./.env"}"
ENVSH_PREFIX="${ENVSH_PREFIX:-"NEXT_PUBLIC_"}"

# Can be `window.__env = {` or `const ENV = {` or whatever you want
ENVSH_PREPEND="${ENVSH_PREPEND:-"window.__env = {"}"
ENVSH_APPEND="${ENVSH_APPEND:-"}"}"
ENVSH_OUTPUT="${ENVSH_OUTPUT:-"./public/__env.js"}"

# Utils
__green() {
  printf '\033[1;31;32m%b\033[0m' "$1"
}

__yellow() {
  printf '\033[1;31;33m%b\033[0m' "$1"
}

__red() {
  printf '\033[1;31;40m%b\033[0m' "$1"
}

__info() {
  printf "%s\n" "$1"
}

__debug() {
  ENVSH_VERBOSE="${ENVSH_VERBOSE:-"false"}"
  if [ "$ENVSH_VERBOSE" == "true" ]; then
    printf "ENVSH_VERBOSE: %s\n" "$1"
  fi
}

# Recreate config file
rm -f "$ENVSH_OUTPUT"
touch "$ENVSH_OUTPUT"

# Add assignment
echo "$ENVSH_PREPEND" >> "$ENVSH_OUTPUT"

if [ "$ENVSH_READ_ENV_VARS" == "true" ]; then
  # Create an array from inline variables
  matched_envs=$(env | grep "^${ENVSH_PREFIX}")
  IFS=$'\n' read -r -d '' -a matched_envs_arr <<< "$matched_envs"
  __info "Matched inline env:"
  for matched_env in "${matched_envs_arr[@]}"; do
    echo "$matched_env"
    awk -F '=' '{print $1 ": \"" (ENVIRON[$1] ? ENVIRON[$1] : $2) "\","}' \
        <<< "$matched_env" >> "$ENVSH_OUTPUT"
  done
fi


if [ "$ENVSH_READ_ENV_FILE" == "true" ]; then
  # Check if file exists
  [[ -f "$ENVSH_ENV_FILE" ]] || { echo "$ENVSH_ENV_FILE does not exist" ; exit 1 ;}

  # Process .env for runtime client use
  __info "$(__green "Reading ${ENVSH_ENV_FILE}...")"
  while IFS= read -r line
  do
    # Check if this line is a valid environment variable and matches our prefix
    if printf '%s' "$line" | grep -e "=" | grep -e "^$ENVSH_PREFIX"; then

      # Read and apply environment variable if exists
      # NOTE: <<< here operator not working with `sh`
      awk -F '=' '{print $1 ": \"" (ENVIRON[$1] ? ENVIRON[$1] : $2) "\","}' \
        <<< "$line" >> "$ENVSH_OUTPUT"
    fi
  done < "$ENVSH_ENV_FILE"
fi

echo "$ENVSH_APPEND" >> "$ENVSH_OUTPUT"

# Print result
__debug "$(__green "Done! Final result in ${ENVSH_OUTPUT}:")"
__debug "$(cat "$ENVSH_OUTPUT")"

__info "$(__green "env.sh done\n")"

# Accepting commands (for Docker)
exec "$@"
