#!/usr/bin/env bash

# Copyright The zerok Authors.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

: "${GITHUB_REPO:="zk-cli"}"
: "${GITHUB_OWNER:="zerok-ai"}"
: "${BINARY_NAME:="zkctl"}"
: "${INSTALL_DIR:="${HOME}/.zerok/bin"}"


BOLD="$(tput bold 2>/dev/null || printf '')"
UNDERLINE="$(tput smul 2>/dev/null || printf '')"
GREY="$(tput setaf 0 2>/dev/null || printf '')"
REV_BG="$(tput rev 2>/dev/null || printf '')"
RED="$(tput setaf 1 2>/dev/null || printf '')"
GREEN="$(tput setaf 2 2>/dev/null || printf '')"
BLUE="$(tput setaf 4 2>/dev/null || printf '')"
YELLOW="$(tput setaf 3 2>/dev/null || printf '')"
NO_COLOR="$(tput sgr0 2>/dev/null || printf '')"

newline() {
  printf "\n"
}

info() {
  printf '%s\n' "${BOLD}${GREY}>${NO_COLOR} $*"
}

warn() {
  printf '%s\n' "${YELLOW}! $*${NO_COLOR}"
}

error() {
  printf '%s\n' "${RED}✕ $*${NO_COLOR}" >&2
}

completed() {
  printf '%s\n' "${GREEN}✔${NO_COLOR} $*"
}

printBanner() {
cat << 'BANNER'
❄❄❄❄❄❄❄❄❄❄❄❄❄❄❄❄❄❄❄❄❄❄❄❄❄❄❄❄❄❄❄❄❄❄❄❄❄❄❄❄❄❄❄❄❄❄❄❄❄❄❄❄❄
	__  ___  __   __           ___                     
	 / |__  |__) /  \    |__/ |__  |    \  / | |\ |    
	/_ |___ |  \ \__/    |  \ |___ |___  \/  | | \| 

	                                ❄ As Kool as it G8s!
❄❄❄❄❄❄❄❄❄❄❄❄❄❄❄❄❄❄❄❄❄❄❄❄❄❄❄❄❄❄❄❄❄❄❄❄❄❄❄❄❄❄❄❄❄❄❄❄❄❄❄❄❄

BANNER
}

parseArguments() {
  while [ "$#" -gt 0 ]; do
    case "$1" in
    -t | --token)
      token="$2"
      shift 2
      ;;
    -t=* | --token=*)
      token="${1#*=}"
      shift 1
      ;;
    *)
      error "Unknown option: $1"
      exit 1
      ;;
    esac
  done
}

# initArch discovers the architecture for this system.
initArch() {
  ARCH=$(uname -m)
  case $ARCH in
    armv5*) ARCH="armv5";;
    armv6*) ARCH="armv6";;
    armv7*) ARCH="arm";;
    aarch64) ARCH="arm64";;
    x86) ARCH="386";;
    x86_64) ARCH="amd64";;
    i686) ARCH="386";;
    i386) ARCH="386";;
  esac
}

# initOS discovers the operating system for this system.
initOS() {
  OS=$(uname |tr '[:upper:]' '[:lower:]')

  case "$OS" in
    # Minimalist GNU for Windows
    mingw*|cygwin*) OS='windows';;
  esac
}

latestReleaseMetaData() {
  
  if [[ -n "${GITHUB_TOKEN}" ]]; then 
    git_token="${GITHUB_TOKEN}@"
  fi

  LATEST_TAG=$(curl -Ls \
              -H "Accept: application/vnd.github+json" \
              -H "Authorization: Bearer ${GITHUB_TOKEN}" \
              -H "X-GitHub-Api-Version: 2022-11-28" \
              "https://api.github.com/repos/${GITHUB_OWNER}/${GITHUB_REPO}/releases/latest")
  
  echo "${LATEST_TAG}"
}

# initLatestTag discovers latest version on GitHub releases.
initLatestTag() {

  latestReleaseData="${1}"
  if [[ -z "${latestReleaseData}" ]]; then 
    latestReleaseData=$(latestReleaseMetaData)
  fi

  LATEST_TAG=$( echo "${latestReleaseData}" | awk -F\" '/tag_name/{print $(NF-1)}')

  if [ -z "${LATEST_TAG}" ]; then
    error "Failed to fetch latest version from ${latest_release_url}"
    exit 1
  fi
  echo "LATEST_TAG = $LATEST_TAG"
}

# appendShellPath append our install bin directory to PATH on bash, zsh and fish shells
appendShellPath() {
  local bashrc_file="${HOME}/.bashrc"
  if [ -f "${bashrc_file}" ]; then
    local export_path_expression="export PATH=${INSTALL_DIR}:\${PATH}"
    if ! grep -q "${export_path_expression}" "${bashrc_file}"; then
      printf "\n%s\n" "${export_path_expression}" >> "${bashrc_file}"
      completed "Added ${INSTALL_DIR} to \$PATH in ${bashrc_file}"
    fi    
  fi

  local zshrc_file="${HOME}/.zshrc"
  if [ -f "${zshrc_file}" ] || [ "${OS}" = "darwin" ]; then
    local export_path_expression="export PATH=${INSTALL_DIR}:\${PATH}"
    if ! grep -q "${export_path_expression}" "${zshrc_file}"; then
      printf "\n%s\n" "${export_path_expression}" >> "${zshrc_file}"
      completed "Added ${INSTALL_DIR} to \$PATH in ${zshrc_file}"
    fi
  fi

  local fish_config_file="${HOME}/.config/fish/config.fish"
  if [ -f "${fish_config_file}" ]; then
    local export_path_expression="set -U fish_user_paths ${INSTALL_DIR} \$fish_user_paths"
    if ! grep -q "${export_path_expression}" "${fish_config_file}"; then
      printf "\n%s\n" "${export_path_expression}" >> "${fish_config_file}"
      completed "Added ${INSTALL_DIR} to \$PATH in ${fish_config_file}"
    fi
  fi
}

# verifySupported checks that the os/arch combination is supported for
# binary builds, as well whether or not necessary tools are present.
verifySupported() {
  local supported="darwin-amd64\ndarwin-arm64\nlinux-amd64\nlinux-arm64"
  if ! echo "${supported}" | grep -q "${OS}-${ARCH}"; then
    error "No prebuilt binary for ${OS}-${ARCH}."
    exit 1
  fi
}

# checkInstalledVersion checks which version of cli is installed and
# if it needs to be changed.
checkInstalledVersion() {
  if [ -f "${INSTALL_DIR}/${BINARY_NAME}" ]; then
    local version
    version=$("${INSTALL_DIR}/${BINARY_NAME}" 
    # --skip-cli-update version
    )
    if [ "${version}" = "${LATEST_TAG#v}" ]; then
      completed "zerok ${version} is already latest"
      return 0
    else
      info "zerok ${LATEST_TAG} is available. Updating from version ${version}."
      return 1
    fi
  else
    return 1
  fi
}

# downloadFile downloads the latest binary package.
initAssetUrl() {

  latestReleaseData="${1}"
  if [[ -z "${latestReleaseData}" ]]; then 
    latestReleaseData=$(latestReleaseMetaData)
  fi

  ARCHIVE_NAME="${BINARY_NAME}_${LATEST_TAG#v}_${OS}_${ARCH}.tar.gz"
  eval $(echo "$latestReleaseData" | grep -C3 "name.:.\+$ARCHIVE_NAME" | grep -w id | tr : = | tr -cd '[[:alnum:]]=')

  DOWNLOAD_URL="https://api.github.com/repos/${GITHUB_OWNER}/${GITHUB_REPO}/releases/assets/${id}"
}

downloadFile() {
  TMP_ROOT="$(mktemp -dt groundcover-installer-XXXXXX)"
  ARCHIVE_TMP_PATH="${TMP_ROOT}/${ARCHIVE_NAME}"

  if [[ "${2}"=="private" ]]; then
    curl -SsL \
      -H "Accept: Accept:application/octet-stream" \
      -H "Authorization: Bearer $GITHUB_TOKEN" \
      -H "X-GitHub-Api-Version: 2022-11-28" \
      "${DOWNLOAD_URL}" \
      -o "${ARCHIVE_TMP_PATH}"
  else
    curl -SsL "${DOWNLOAD_URL}" -o "${ARCHIVE_TMP_PATH}"
  fi
  
}

# installFile installs the cli binary.
installFile() {
  tar xf "${ARCHIVE_TMP_PATH}" -C "${TMP_ROOT}"
  BIN_PATH="${INSTALL_DIR}/${BINARY_NAME}"
  BIN_TMP_PATH="${TMP_ROOT}/${BINARY_NAME}"
  info "Preparing to install ${BINARY_NAME} into ${INSTALL_DIR}"
  mkdir -p "${INSTALL_DIR}"
  cp "${BIN_TMP_PATH}" "${BIN_PATH}"
  chmod +x "${BIN_PATH}"
  completed "${BINARY_NAME} installed into ${BIN_PATH}"
}

# cleanup temporary files
cleanup() {
  if [ -d "${TMP_ROOT:-}" ]; then
    rm -rf "${TMP_ROOT}"
  fi
}

printWhatNow() {
  printf "\n%s\
what now?\n\
* run ${GREEN}zerok auth login${NO_COLOR}\n\
* then ${GREEN}zerok deploy${NO_COLOR}\n\
* ${REV_BG}let the magic begin.${NO_COLOR}\n\n\
run ${BLUE}zerok help${NO_COLOR}, or dive deeper with ${BLUE}${UNDERLINE}https://docs.zerok.ai/docs${NO_COLOR}.\n"
}

deployWithToken() {
  "${INSTALL_DIR}/${BINARY_NAME}" deploy --token "${token}"
}

# fail_trap is executed if an error occurs.
fail_trap() {
  result=$?
  if [ "$result" != "0" ]; then
    error "Failed to install ${BINARY_NAME}"
    info "For support, go to ${BLUE}${UNDERLINE}https://github.com/zerok-com/cli${NO_COLOR}"
  fi
  cleanup
  exit $result
}

# Execution

#Stop execution on any error
trap "fail_trap" EXIT
set -e

latestRelease=$(latestReleaseMetaData)

printBanner
parseArguments "$@"
initArch
initOS

initLatestTag "${latestRelease}"

if ! checkInstalledVersion; then
  # downloadFile
  initAssetUrl "${latestRelease}"
  downloadFile "private"
  installFile
fi
appendShellPath
completed "zerok cli was successfully installed!"
if [ -z "${token}" ]
then
  printWhatNow
  cleanup
  exec "${SHELL}" # Reload shell
else
  newline
  deployWithToken
fi