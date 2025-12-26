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

## Local Cluster (quick start)

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

## Releases

Releases are built and published automatically from git tags (e.g., `v0.2.0`) using GitHub Actions + GoReleaser.

Create a release:

```sh
git tag v0.1.0
git push origin v0.1.0
```
