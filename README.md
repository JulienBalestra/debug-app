# Getting started locally

```bash
curl -Lfs https://github.com/DataDog/pupernetes/releases/download/v0.8.0/pupernetes -o ./pupernetes
chmod +x pupernetes 
sudo apt-get update && sudo apt-get install -qqy tar unzip libseccomp2
sudo ./pupernetes daemon run /opt/sandbox --container-runtime containerd --skip-probes --kubectl-link /usr/local/bin/kubectl

kubectl apply -f probe-failures/daemonset.yaml
```


## dlv

```bash
GOCACHE=off CGO_ENABLED=0 go build -gcflags='all=-N -l' -o bin/containerd-shim -ldflags '-X github.com/containerd/containerd/version.Version=v1.2.0-beta.2-43-g6ca8355a.m -X github.com/containerd/containerd/version.Revision=6ca8355a4e8ae1857c0577fcd648976916c7ce33.m -X github.com/containerd/containerd/version.Package=github.com/containerd/containerd -extldflags "-static" -linkmode internal' -tags "seccomp apparmor" ./cmd/containerd-shim
sudo cp -v bin/containerd-shim ~/go/src/github.com/DataDog/pupernetes/sandbox/bin/
```

```bash
sudo ~/go/bin/dlv --headless -l=:2345 --api-version=2 attach 13026
```
