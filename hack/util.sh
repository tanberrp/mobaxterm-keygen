# Copyright 2024 The mobaxterm-keygen Authors.
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

function util::get_oss() {
  echo linux darwin
}

function util::get_archs() {
  echo amd64 arm64
}

function util::build_binary_for_platform() {
  local -r target=$1
  local -r platform=$2
  local -r os=${target}
  local -r arch=${platform}

  if [[ -n "${3+x}" ]]; then
    local -r output=$3
  else
    local -r output=_output/bin/${target}/${platform}
  fi

  set -x
  mkdir -p "${REPO_ROOT}/${output}"
  GO111MODULE=on CGO_ENABLED=0 GOOS=${os} GOARCH=${arch} go build \
    -mod vendor \
    -o "${REPO_ROOT}/${output}/" \
    "${REPO_ROOT}/cmd/mobaxterm-keygen/mobaxterm-keygen.go"
  set +x
}

function util::package_for_platform() {
  local -r os=$1
  local -r arch=$2

  release_dir="${REPO_ROOT}/_output/release"
  mkdir -p "${release_dir}"
  tar_file="mobaxterm-keygen-${os}-${arch}.tgz"
  echo "Packaging ${tar_file}"
  if [[ ${os} == "windows" ]]; then
    tar czf "${release_dir}/${tar_file}" -C "${REPO_ROOT}" LICENSE -C "${REPO_ROOT}/_output/bin/${os}/${arch}" "mobaxterm-keygen.exe"
  else
    tar czf "${release_dir}/${tar_file}" -C "${REPO_ROOT}" LICENSE -C "${REPO_ROOT}/_output/bin/${os}/${arch}" "mobaxterm-keygen"
  fi
  cd "${release_dir}" || exit
  sha256sum "${tar_file}" >"${tar_file}.sha256"
  cd - >/dev/null || exit
}
