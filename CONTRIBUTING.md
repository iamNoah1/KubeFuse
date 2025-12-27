# Contributing

Thanks for your interest in KubeFuse.

## Development

Build/run locally:

```sh
go run main.go <command> <subcommand>
```

Or build the binary:

```sh
go build -o kubefuse
./kubefuse <command> <subcommand>
```

## Tests

```sh
go test ./...
```

### Integration tests (envtest)

The integration tests run against controller-runtime's envtest binaries (local API server + etcd)
instead of a real cluster. This is faster than kind and good for API-level behavior that doesn't
need real nodes, CRI, or controllers.

Use kind when you want a full Kubernetes control plane and nodes for end-to-end testing.

More info:
- https://pkg.go.dev/sigs.k8s.io/controller-runtime/pkg/envtest
- https://book.kubebuilder.io/reference/envtest.html

Use the helper script to print the `KUBEBUILDER_ASSETS` export line:

```sh
./scripts/setup-envtest.sh
export KUBEBUILDER_ASSETS=... # copy from script output
go test -tags=integration ./internal/integration
```

### End-to-end tests (kind)

End-to-end tests should run against a real cluster (kind) and verify KubeFuse behavior
across create → patch → TTL rollback. These are heavier and can be reserved for release
validation.

```sh
./scripts/e2e-kind.sh
```

## Local Cluster (quick start)

### kind (Kubernetes-in-Docker)
Install kind:

```sh
go install sigs.k8s.io/kind@latest
```

Create a cluster:

```sh
./scripts/kind-up.sh
kubectl config use-context kind-kubefuse
```

Delete the cluster:

```sh
./scripts/kind-down.sh
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

## Releases

Releases are built and published automatically from git tags (e.g., `v0.2.0`) using GitHub Actions + GoReleaser.

Create a release:

```sh
git tag v0.1.0
git push origin v0.1.0
```
