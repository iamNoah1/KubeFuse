#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
SETUP_ENVTEST_BIN="${SETUP_ENVTEST_BIN:-setup-envtest}"
K8S_VERSION="${K8S_VERSION:-1.29.x}"

if ! command -v "${SETUP_ENVTEST_BIN}" >/dev/null 2>&1; then
  echo "setup-envtest not found. Install with:"
  echo "  go install sigs.k8s.io/controller-runtime/tools/setup-envtest@latest"
  exit 1
fi

ASSETS="$("${SETUP_ENVTEST_BIN}" use -p path "${K8S_VERSION}")"
echo "export KUBEBUILDER_ASSETS=${ASSETS}"
echo "Example:"
echo "  export KUBEBUILDER_ASSETS=${ASSETS}"
echo "  go test -tags=integration ./internal/integration"
echo "Root: ${ROOT_DIR}"
