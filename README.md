# goapp-demo

Small Go service + the manifests/infra to run it on Kubernetes.

## Contents

- `app/` – Go app, exposes `/health`, `/version`, `/env`, and CRUD on `/config`
- `Dockerfile` – multi-stage build, `golang:alpine` → `alpine` runtime
- `helm/` – chart to deploy the app (`goapp-demo`)
- `terraform/` – creates the namespace and installs the Helm release
- `.gitlab-ci.yml` – test (go test / hadolint / helm lint, run in parallel) → build (kaniko) → deploy (opentofu, manual)
- `build.sh` – local multi-arch (amd64/arm64) image build

## Running it locally

A local cluster (kind/minikube/k3d) and a kubeconfig is required for deploy.

```sh
# test
cd app
go test

# compilte
cd app
go build .

# build the image
./build.sh

# deploy (namespace + helm release)
cd terraform
tofu init
tofu apply

# reach the app
kubectl port-forward svc/goapp-demo 8080:80
curl localhost:8080/health
```

## CI

DOCKER_AUTH variable is required to configure docker login for ci to work properly.

## TODO / Improvements

### Top

- Vulnerability scan - trivy?
- Deploy stage in gitlab-ci (tfstate), or a better deploy solution altogether (gitops)
- Config param management: add db support, and/or persistence with an optional helm hook to fill the json with
  predefined key/value pairs
- Better ci rules - required lint pass on merge-request
- Observability - metrics endpoint

### Additional

- Make more stuff configurable via envvar, eg. http portlistener
- Better logging
- Release stage in gitlab-ci, manage tags, bumps automatically
- Better versioning overall
- Helm chart tidy
- Liveness/readyness probe separation, depends if bigger complexity is added later
