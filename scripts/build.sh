#!/usr/bin/env bash

# This script builds the application from source for multiple platforms.

# Get the parent directory of where this script is.
SOURCE="${BASH_SOURCE[0]}"
while [[ -h "$SOURCE" ]] ; do SOURCE="$(readlink "$SOURCE")"; done
ROOT_DIR="$( cd -P "$( dirname "$SOURCE" )/.." && pwd )"

readonly LOCAL_OUTPUT_ROOT="${ROOT_DIR}/${OUT_DIR:-_output}"
readonly LOCAL_OUTPUT_IMAGE_STAGING="${LOCAL_OUTPUT_ROOT}/images"
readonly LOCAL_OUTPUT_BINPATH="${LOCAL_OUTPUT_ROOT}/bin"

# Change into that directory
cd "$ROOT_DIR"

# Get the git commit
GIT_COMMIT=$(git rev-parse HEAD)
GIT_DIRTY=$(test -n "`git status --porcelain`" && echo "+CHANGES" || true)

# Determine the arch/os combos we're building for
XC_ARCH=${XC_ARCH:-"amd64"}
XC_OS=${XC_OS:-linux}
XC_EXCLUDE_OSARCH="!darwin/arm !darwin/386"

# Delete the old dir
echo "==> Removing old directory..."
rm -rf "$LOCAL_OUTPUT_BINPATH"
rm -rf "$LOCAL_OUTPUT_IMAGE_STAGING"
mkdir -p "$LOCAL_OUTPUT_BINPATH"

if ! which gox > /dev/null; then
    echo "==> Installing gox..."
    go get -u github.com/mitchellh/gox
fi

# instruct gox to build statically linked binaries
export CGO_ENABLED=0

# Allow LD_FLAGS to be appended during development compilations
LD_FLAGS="-X main.GitCommit=${GIT_COMMIT}${GIT_DIRTY} $LD_FLAGS"

# Ensure all remote modules are downloaded and cached before build so that
# the concurrent builds launched by gox won't race to redundantly download them.
echo "==> Start go mod download"
go mod download

function teddy::build() {
    pushd $1
    echo "==> Enter $1"
    echo "==> Building..."
    gox \
        -os="${XC_OS}" \
        -arch="${XC_ARCH}" \
        -osarch="${XC_EXCLUDE_OSARCH}" \
        -ldflags "${LD_FLAGS}" \
        -output "$LOCAL_OUTPUT_BINPATH/{{.OS}}_{{.Arch}}/${PWD##*/}" \
        .
    popd >/dev/null 2>&1
}

function teddy::build_docker_for_api() {
    CMD_NAME=$(basename $1)
    REAL_BUILD=false
    echo "==> Building docker image for $CMD_NAME..."
    for PLATFORM in $(find ${LOCAL_OUTPUT_BINPATH} -mindepth 1 -maxdepth 1 -type d); do
        OSARCH=$(basename ${PLATFORM})
        if [[ ${OSARCH} == "linux_amd64" ]]; then
            docker build --no-cache "--build-arg=TEDDY_CMD=${CMD_NAME}" \
                -t teddy/${CMD_NAME} -f "${ROOT_DIR}/scripts/docker-release/Dockerfile.api" ${PLATFORM}
            break
        fi
    done
    if [[ "${REAL_BUILD}" == false ]]; then
        echo "==> Please build bin for linux_amd64"
    fi
}

function teddy::build_docker_for_srv() {
    CMD_NAME=$(basename $1)
    REAL_BUILD=false
    echo "==> Building docker image for $CMD_NAME..."
    for PLATFORM in $(find ${LOCAL_OUTPUT_BINPATH} -mindepth 1 -maxdepth 1 -type d); do
        OSARCH=$(basename ${PLATFORM})
        if [[ ${OSARCH} == "linux_amd64" ]]; then
            docker build --no-cache "--build-arg=TEDDY_CMD=${CMD_NAME}" \
                -t teddy/${CMD_NAME} -f "${ROOT_DIR}/scripts/docker-release/Dockerfile.srv" ${PLATFORM}
            break
        fi
    done
    if [[ "${REAL_BUILD}" == false ]]; then
        echo "==> Please build bin for linux_amd64"
    fi
}

teddy::build "cmd/uaa"
teddy::build_docker_for_srv "cmd/uaa"

teddy::build "cmd/message"
teddy::build_docker_for_srv "cmd/message"

teddy::build "cmd/content"
teddy::build_docker_for_srv "cmd/content"

teddy::build "cmd/captcha"
teddy::build_docker_for_srv "cmd/captcha"

teddy::build "cmd/api-base"
teddy::build_docker_for_api "cmd/api-base"

teddy::build "cmd/api-content"
teddy::build_docker_for_api "cmd/api-content"

teddy::build "cmd/api-uaa"
teddy::build_docker_for_api "cmd/api-uaa"
