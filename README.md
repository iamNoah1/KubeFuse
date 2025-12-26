# KubeFuse

A CLI tool for safe, temporary live-patching of Kubernetes resources. KubeFuse applies a patch, waits for the TTL, and rolls the resource back to its original values.

## Key Features (MVP):
- Patch Kubernetes resources with dot-paths, e.g. `spec.replicas=2`
- TTL-based rollback (the CLI waits and reverts the change)
- Audit annotations for reason and TTL
- Dry-run mode for previewing changes
- Resource-aware shell completion

## Usage

kubefuse set <kind/name> <path=value>... [--ttl 10m] [--reason "..."] [--dry-run]

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

### Value Mapping 

| CLI Value    | Internal represenation | 
|--------------|------------------------|
| true / false | bool                   |
| 1,2,3,4 ..   | int                    |
| null         | nil                    |


## Shell Completion

Generate completion scripts:

```sh
kubefuse completion bash
kubefuse completion zsh
kubefuse completion fish
kubefuse completion powershell
```

The `set` command completes resource kinds and names using your current kubeconfig and `--namespace` flag.

## Build/Run
* `go run main.go <command> <subcommand>`
* or `go build -o kubefuse` and then `./kubefuse <command> <subcommand>`

## Releases
Releases are built and published automatically from git tags (e.g., `v0.2.0`) using GitHub Actions + GoReleaser.

Create a release:

```sh
git tag v0.1.0
git push origin v0.1.0
```

## Download a Release
Prebuilt binaries are published on GitHub Releases.

1) Open the latest release: `https://github.com/iamNoah1/KubeFuse/releases/latest`
2) Download the asset for your OS/arch (e.g., `kubefuse_v0.1.0_darwin_arm64.tar.gz`).

Example (macOS arm64):

```sh
curl -LO https://github.com/iamNoah1/KubeFuse/releases/download/v0.1.0/kubefuse_v0.1.0_darwin_arm64.tar.gz
tar -xzf kubefuse_v0.1.0_darwin_arm64.tar.gz
./kubefuse --help
```

## Install
Choose one of the standard install options below.

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

## Try It Locally (quick clusters)

### kind (Kubernetes-in-Docker)
Install kind:

```sh
go install sigs.k8s.io/kind@latest
```

Create a cluster:

```sh
kind create cluster --name kubefuse
kubectl config use-context kind-kubefuse
```

### Docker Desktop Kubernetes
Enable Kubernetes in Docker Desktop, then select the `docker-desktop` context:

```sh
kubectl config use-context docker-desktop
```

## Smoke Test
```sh
kubectl create deployment web --image=nginx
kubefuse set deployment/web spec.replicas=2 --ttl 30s --reason "smoke test"
kubectl get deployment web -w
```

## Run Tests 
* `go test ./...`

## Limitations (current)
* Path segments are dot-separated (e.g., `spec.replicas`). Array indexes like `containers[0]` are not supported yet.
* Rollback is handled by the running CLI process; there is no in-cluster controller yet.

## Additional Resources
* https://github.com/spf13/cobra
* 
