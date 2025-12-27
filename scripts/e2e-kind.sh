#!/usr/bin/env bash
set -euo pipefail

CLUSTER_NAME="${KIND_CLUSTER_NAME:-kubefuse}"
KUBEFUSE_BIN="${KUBEFUSE_BIN:-./kubefuse}"
NAMESPACE="${E2E_NAMESPACE:-kubefuse-e2e}"
TTL="${E2E_TTL:-5s}"
APPLY_REPLICAS="${E2E_APPLY_REPLICAS:-2}"
ORIGINAL_REPLICAS="${E2E_ORIGINAL_REPLICAS:-1}"
TIMEOUT_SECONDS="${E2E_TIMEOUT_SECONDS:-60}"

ensure_command() {
  if ! command -v "$1" >/dev/null 2>&1; then
    echo "$1 not found in PATH"
    exit 1
  fi
}

ensure_kind() {
  if command -v kind >/dev/null 2>&1; then
    return 0
  fi
  echo "kind not found in PATH, installing via go install..."
  ensure_command go
  GOBIN="$(go env GOPATH)/bin"
  go install sigs.k8s.io/kind@latest
  export PATH="${GOBIN}:${PATH}"
  ensure_command kind
}

wait_for_replicas() {
  local expected="$1"
  local deadline="$((SECONDS + TIMEOUT_SECONDS))"
  while [ "$SECONDS" -lt "$deadline" ]; do
    current="$(kubectl -n "${NAMESPACE}" get deployment web -o jsonpath='{.spec.replicas}' 2>/dev/null || true)"
    if [ "${current}" = "${expected}" ]; then
      return 0
    fi
    sleep 1
  done
  return 1
}

cleanup() {
  kubectl -n "${NAMESPACE}" delete deployment web --ignore-not-found >/dev/null 2>&1 || true
  kubectl delete namespace "${NAMESPACE}" --ignore-not-found >/dev/null 2>&1 || true
  if [ "${CREATED_CLUSTER}" = "true" ]; then
    kind delete cluster --name "${CLUSTER_NAME}" >/dev/null 2>&1 || true
  fi
  if [ -n "${PREV_CONTEXT}" ]; then
    kubectl config use-context "${PREV_CONTEXT}" >/dev/null 2>&1 || true
  fi
}

ensure_command kubectl
ensure_kind

CREATED_CLUSTER="false"
PREV_CONTEXT="$(kubectl config current-context 2>/dev/null || true)"

if ! kind get clusters | grep -qx "${CLUSTER_NAME}"; then
  kind create cluster --name "${CLUSTER_NAME}"
  CREATED_CLUSTER="true"
fi

trap cleanup EXIT

kubectl config use-context "kind-${CLUSTER_NAME}"

if [ ! -x "${KUBEFUSE_BIN}" ]; then
  ensure_command go
  go build -o "${KUBEFUSE_BIN}"
fi

kubectl create namespace "${NAMESPACE}" >/dev/null 2>&1 || true
kubectl -n "${NAMESPACE}" create deployment web --image=nginx --replicas="${ORIGINAL_REPLICAS}" >/dev/null

"${KUBEFUSE_BIN}" set deployment/web -n "${NAMESPACE}" "spec.replicas=${APPLY_REPLICAS}" --ttl "${TTL}" --reason "e2e" &
KUBEFUSE_PID=$!

if ! wait_for_replicas "${APPLY_REPLICAS}"; then
  echo "timed out waiting for apply replicas=${APPLY_REPLICAS}"
  kill "${KUBEFUSE_PID}" >/dev/null 2>&1 || true
  exit 1
fi

if ! wait_for_replicas "${ORIGINAL_REPLICAS}"; then
  echo "timed out waiting for rollback replicas=${ORIGINAL_REPLICAS}"
  kill "${KUBEFUSE_PID}" >/dev/null 2>&1 || true
  exit 1
fi

wait "${KUBEFUSE_PID}"

echo "e2e kind test passed"
