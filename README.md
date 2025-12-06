# KubeFuse

A CLI tool that allows safe, temporary live-patching of Kubernetes resources with features like TTL-based rollback, dry-run, and audit annotations.

## Key Features (MVP):
- Patch Kubernetes resources, ie `spec.replicas=2`
- Optional `-ttl` flag to auto-revert the change
- Add patch reason annotations
- Dry-run mode for previewing changes
- Optional: TTL manager as in-cluster controller or background job

## Usage

kubefuse set <kind/name> <path=value>... [--ttl 10m] [--reason "..."] [--dry-run]

Example: 

```
kubefuse deploy/web -n prod \\
  set spec.template.spec.containers[0].image=nginx:1.21 spec.template.spec.containers[1].image=nginx:1.19\\
  --ttl 10m \\
  --reason "hotfix for prod"
  --namespace kubefuse
```

### Value Mapping 

| CLI Value    | Internal represenation | 
|--------------|------------------------|
| true / false | bool                   |
| 1,2,3,4 ..   | int                    |
| null         | nil                    |


## Build/Run
* `go run main.go <command> <subcommand>`
* or `go build -o kubefuse` and then `./kubefuse <command> <subcommand>`

## Run Tests 
* `go test ./...`

## Additional Resources
* https://github.com/spf13/cobra
* 