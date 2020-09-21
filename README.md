# KLoops

KLoops is a modern CI/CD platform built for Kubernetes. It is strongly inspired from [Prow](https://github.com/kubernetes/test-infra/tree/master/prow) and [Lighthouse](https://github.com/jenkins-x/lighthouse).

It is based on a ChatOps based webhook handler which can trigger build pipelines and report status to the scm provider.

A dashboard is available to configure KLoops and visualize running or past jobs.

Supported providers:
- Github
- Github entreprise
- Gitlab
- BitBucket cloud
- Gitea

Behind the scene, we leverage [Tekton pipelines](https://github.com/tektoncd/pipeline) to run build pipelines execution. The [Tekton dashboard](https://github.com/tektoncd/dashboard) provides pipelines execution and logs visualization.

## KLoops philosophy

The KLoops platform strongly believes in ownership and simplicity of use. This mantra is the most important rule that we apply when building KLoops.

For these reasons, every aspect of KLoops is simple, deployment, usage, access...

Deployment is as simple as deploying a helm chart.

Once deployed, you can start building your repository by enlisting it in your KLoops dashboard, it will take care of setting up everything for you and will start watching for events in your repository.

## Getting started

You can deploy a fully working solution on a local cluster with `kind`, using `Gitea` as a scm provider.

Clone the repository and run `./scripts/up.sh`, it will create a `kind` cluster, deploy everything in it and initialize `Gitea` with sample datas.
