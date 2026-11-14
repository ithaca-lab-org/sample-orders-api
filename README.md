# sample-orders-api

The first tenant app of the [ithaca IDP](https://github.com/ithaca-lab-org/idp-infra): a tiny Go service that proves the whole delivery path (image, hand-made Spinnaker pipeline, portal secret in the environment, route over HTTPS). Plan and runbook: `idp-infra/docs/order-04-first-app.md` and `order-04-first-run.md`; design: ADR 0011.

- `GET /` returns `{"service":"sample-orders-api","secret":{"ORDERS_API_KEY":"set|unset","ORDERS_GREETING":"set|unset"}}`: it reports whether the portal secret `orders-config` reached the pod, never the values. `GET /healthz` is the probe.
- `deploy/base/` is the Kustomize base Spinnaker bakes (Rosco): a Deployment and a Service named `orders`, running as the ServiceAccount `orders` of the WorkloadIdentity with the label `idp.ithaca-idp.com/workload: orders`, `envFrom` the Secret `orders-config`. Only a Deployment and a Service: the platform owns the route, ServiceAccount, policies and Secrets.
- The workflow `image.yml` builds, scans (Trivy) and pushes `ghcr.io/ithaca-lab-org/sample-orders-api:<sha>` on every push to `main` and prints the digest in the job summary. Paste the digest into `deploy/base/deployment.yaml` (it names the image by digest) and commit, then run the pipeline with that commit SHA.

Local: `go test ./...`, `kubectl kustomize deploy/base`.
