# ArgoCD agent for codefresh gitops dashboard 

[![Codefresh build status]( https://g.codefresh.io/api/badges/pipeline/codefresh-inc/argo%2Fagent?type=cf-1&key=eyJhbGciOiJIUzI1NiJ9.NTY3MmQ4ZGViNjcyNGI2ZTM1OWFkZjYy.AN2wExsAsq7FseTbVxxWls8muNx_bBUnQWQVS8IgDTI)]( https://g.codefresh.io/pipelines/edit/new/builds?id=5f21305719d46c880abeeeb5&pipeline=agent&projects=argo&projectId=5f16b786f25d80a086a56bcb)


Codefresh providing [dashboard](https://codefresh.io/docs/docs/ci-cd-guides/gitops-deployments/) for watching on all activities that happens on argocd side. Codefresh argocd agent important part for check all argocd CRD use watch api and notify codefresh about all changes. 

Like: 
* Application created/removed/updated
* Project created/removed/updated
* Your manifest repo information base on context that you provide to us during installation

In addition this agent do automatic application sync between argocd and codefresh 



## Prerequisites

Make sure that you have

* a [Codefresh account](https://codefresh.io/docs/docs/getting-started/create-a-codefresh-account/) with enabled gitops feature
* a [Codefresh API token](https://codefresh.io/docs/docs/integrations/codefresh-api/#authentication-instructions) that will be used as a secret in the agent
* a [Codefresh CLI](https://codefresh-io.github.io/cli/) that will be used for install agent
* a [ArgoCD Server](https://argoproj.github.io/argo-cd/cli_installation/)

ArgoCD agent has following resource requirements 
```
 requests:
   memory: "128Mi"
   cpu: "0.2"
 limits:
   memory: "256Mi"
   cpu: "0.4"
```

## Installation     
 

```sh
codefresh install gitops argocd-agent 
```

<img src="/art/installation.gif?raw=true" width="1200px">

## Uninstall     
 

```sh
codefresh uninstall gitops argocd-agent 
```

## Upgrade     

Codefresh will show you indicator inside your [gitops integration](https://g.codefresh.io/account-admin/account-conf/integration/gitops) when you need upgrade your agent

<img src="/art/upgrade.png?raw=true" width="800px"> 

```sh
codefresh upgrade gitops argocd-agent 
```

## How to use the ArgoCD agent

<img src="/art/dashboard.png" width="1200px">

## How it works ( Diagram )

<img src="/art/high_level.png" width="1200px">

<img src="/art/detailed.png" width="1200px">

## Local execution

### Environment variables 

* ARGO_HOST - Argocd host (like https://34.71.103.174/)
* ARGO_USERNAME - Argocd username ( Need provide if ARGO_TOKEN empty )
* ARGO_PASSWORD - Argocd password ( Need provide if ARGO_TOKEN empty )
* ARGO_TOKEN - Argocd user token
* CODEFRESH_TOKEN - [Codefresh user token](https://codefresh.io/docs/docs/integrations/codefresh-api/#authentication-instructions)
* CODEFRESH_INTEGRATION - Codefresh gitops integration name
* CODEFRESH_HOST - Codefresh host ( prodution https://g.codefresh.io)
* GIT_PASSWORD - Git token

## Run tests
`go test -cover ./...`
