# Getting started locally

```bash
curl -Lfs https://github.com/DataDog/pupernetes/releases/download/v0.8.0/pupernetes -o ./pupernetes
chmod +x pupernetes 
sudo apt-get update && sudo apt-get install -qqy tar unzip libseccomp2
sudo ./pupernetes daemon run /opt/sandbox --container-runtime containerd --skip-probes --kubectl-link /usr/local/bin/kubectl

kubectl apply -f probe-failures/daemonset.yaml
```