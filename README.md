# GKE Authentication Plugin

This plugin provides a standalone way to generate an ExecCredential for use by k8s.io/client-go applications.

Google already provides a [gke-gcloud-auth-plugin](https://cloud.google.com/blog/products/containers-kubernetes/kubectl-auth-changes-in-gke); however, that plugin depends on the gcloud CLI, which is written in Python. This dependency graph is Large if you want to authenticate and interact with a GKE cluster from a go application.

The plugin is for use outside of a cluster; when running in the cluster, mount a service account and use that token to interact with the Kubernetes API.

## Changes from upstream

- Add `--credential/-c` argument to override default json credentials.
- Put the cache within the same KUBECONFIG file directory.
- Named after original Google provided plugin.

## Build

```shell
make
```

Or with Docker:
```shell
docker build -f Dockerfile.dev -t gke-gcloud-auth-plugin-dev .

docker run -it --rm --name gke-gcloud-auth-plugin-dev-container -v ${PWD}:/home/nonroot gke-gcloud-auth-plugin-dev

make
```

## Run

```shell
# generate ExecCredential
bin/gke-gcloud-auth-plugin

# version
bin/gke-gcloud-auth-plugin version
```

### Example Exec Section of Kubeconfig

```yaml
users:
- name: user_id
  user:
    exec:
      apiVersion: client.authentication.k8s.io/v1beta1
      command: gke-gcloud-auth-plugin
      args:
      - -c
      - /path/to/adc.json
      provideClusterInfo: true
      interactiveMode: Never
```
## TODO

- Add unit tests
