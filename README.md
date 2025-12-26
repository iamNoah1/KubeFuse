# KubeFuse

A CLI tool for safe, temporary live-patching of Kubernetes resources. KubeFuse applies a patch, waits for the TTL, and rolls the resource back to its original values.

## Overview
### Key Features (MVP)
- Patch Kubernetes resources with dot-paths, e.g. `spec.replicas=2`
- TTL-based rollback (the CLI waits and reverts the change)
- Audit annotations for reason and TTL
- Dry-run mode for previewing changes
- Resource-aware shell completion

## Install

### Go install (recommended for Go users)
```sh
go install github.com/iamNoah1/KubeFuse@latest
```

### Release binary (no Go toolchain required)
```sh
curl -LO https://github.com/iamNoah1/KubeFuse/releases/download/v0.1.0/kubefuse_v0.1.0_darwin_arm64.tar.gz
tar -xzf kubefuse_v0.1.0_darwin_arm64.tar.gz
./kubefuse --help
```

## Prerequisites
You need access to a Kubernetes cluster and a working `kubectl` context (KubeFuse uses your `KUBECONFIG`).

## Quickstart
### Command syntax

kubefuse set <kind/name> <path=value>... [--ttl 10m] [--reason "..."] [--dry-run]

### Shell completion
Generate completion scripts:

```sh
kubefuse completion bash
kubefuse completion zsh
kubefuse completion fish
kubefuse completion powershell
```

The `set` command completes resource kinds and names using your current kubeconfig and `--namespace` flag.

### Examples
Example:

```sh
kubefuse set deployment/web spec.replicas=3 --ttl 5m --reason "scale for peak"
```

Example (namespace + labels):

```sh
kubefuse set deploy/api -n prod metadata.labels.tier=backend --ttl 30m --reason "temporary label"
```

## How It Works

1. KubeFuse reads the current values at each patch path (and the existing audit annotations).
2. It applies your merge patch and adds `kubefuse.dev/reason` + `kubefuse.dev/ttl`.
3. If `--ttl` is non-zero and not `--dry-run`, KubeFuse waits for the TTL and applies a rollback patch.

Note: The CLI process stays running until the rollback is complete.

## Contributing
See `CONTRIBUTING.md`.
