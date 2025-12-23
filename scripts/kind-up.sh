#!/usr/bin/env bash
set -euo pipefail

CLUSTER_NAME="${KIND_CLUSTER_NAME:-kubefuse}"

if ! command -v kind >/dev/null 2>&1; then
  echo "kind not found. Install from https://kind.sigs.k8s.io/"
  exit 1
fi

kind create cluster --name "${CLUSTER_NAME}"
echo "Cluster '${CLUSTER_NAME}' is ready."
