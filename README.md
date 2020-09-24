# KLoops

KLoops is a modern CI/CD platform built for Kubernetes. It is strongly inspired from [Prow](https://github.com/kubernetes/test-infra/tree/master/prow) and [Lighthouse](https://github.com/jenkins-x/lighthouse).

It is based on a ChatOps based webhook handler which can trigger build pipelines and report status to the scm provider.

A dashboard is available to configure KLoops and visualize running or past jobs.

Supported providers:
- Github
- Github entreprise
- Gitlab (WIP)
- BitBucket cloud (WIP)
- Gitea

Behind the scene, we leverage [Tekton pipelines](https://github.com/tektoncd/pipeline) to run build pipelines execution. The [Tekton dashboard](https://github.com/tektoncd/dashboard) provides pipelines execution and logs visualization.

## Why KLoops ?

KLoops was created after working with Prow and Lighthouse for several years. Both have limitations that KLoops adresses:
- Prow was built specifically for Kubernetes and does not apply for general purpose
- Prow and Lighthouse manage configuration centrally, that makes it difficult to use as a self service
- Prow supports GitHub only
- Lighthouse started as a Prow fork and therefore inherited a lot of the Prow limitations, KLoops on the other hand started from scratch
- Lighthouse was created for Jenkins-X first, KLoops completely abstracts the underlying engine
- Tekton is reeally powerfull at running pipelines but is unopiniionated, a CICD platform can probably be built with Tekton only but it will require a lot of work
- Lighthouse lacks user interface

KLoops was strongly inspired from Prow, Lighthouse and to some extents Tekton. It is self service oriented, supports all major scm providers, stores its configuration in a non centralized way and can receive webhooks from all scm providers at the same time (no scm specific deployment).

All the configuration is stored in Custom Resources and different configuration per team/repository is by nature possible.

We are planning to support some very advanced network security features to enable full multi tenancy soon.

## KLoops philosophy

The KLoops platform strongly believes in ownership and simplicity of use. This mantra is the most important rule that we apply when building KLoops.

For these reasons, every aspect of KLoops is simple, deployment, usage, access...

Deployment is as simple as deploying a helm chart.

Once deployed, you can start building your repository by enlisting it in your KLoops dashboard, it will take care of setting up everything for you and will start watching for events in your repository.

## Getting started

You can deploy a fully working solution on a local cluster with `kind`, using `Gitea` as a scm provider.

Clone the repository and run `./scripts/up.sh`, it will create a `kind` cluster, deploy everything in it and initialize `Gitea` with sample datas.

You can look at our [ROADMAP](./ROADMAP.md), feel free ro create issues or open pull requests to contribute.

## Docker images

KLoops docker images are available on [dockerhub](https://hub.docker.com/).

You will find the following components images:
- [CHATBOT](https://hub.docker.com/r/eddycharly/kloops-chatbot)
- [DASHBOARD](https://hub.docker.com/r/eddycharly/kloops-dashboard)

## Helm charts

KLoops helm charts are available at https://eddycharly.github.io/kloops.

You can browse the chart [docs](./charts/kloops/README.md) to find install instructions and supported configuration.

## Docs

Docs are located in the [docs](./docs/README.md) folder, you will find user docs, install docs, dev docs and more.
