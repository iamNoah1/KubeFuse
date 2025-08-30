# KubeFuse

A CLI tool that allows safe, temporary live-patching of Kubernetes resources with features like TTL-based rollback, dry-run, and audit annotations.

## Key Features (MVP):
- Patch Kubernetes resources, ie `spec.replicas=2`
- Optional `-ttl` flag to auto-revert the change
- Add patch reason annotations
- Dry-run mode for previewing changes
- Optional: TTL manager as in-cluster controller or background job

## Intended Example Usage (not implemented yet)

kubefuse deploy/web -n prod \\
  set spec.template.spec.containers[0].image=nginx:1.21 \\
  --ttl 10m \\
  --reason "hotfix for prod"

## Usage 
* `go run main.go <command> <subcommand>`
* or `go build -o kubefuse` and then `./kubefuse <command> <subcommand>`