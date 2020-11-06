# ArgoCD codefresh event listener
> Tool for listen argocd CRDs and store it on codefresh side

[![Codefresh build status]( https://g.codefresh.io/api/badges/pipeline/codefresh-inc/argo%2Fagent?type=cf-1&key=eyJhbGciOiJIUzI1NiJ9.NTY3MmQ4ZGViNjcyNGI2ZTM1OWFkZjYy.AN2wExsAsq7FseTbVxxWls8muNx_bBUnQWQVS8IgDTI)]( https://g.codefresh.io/pipelines/edit/new/builds?id=5f21305719d46c880abeeeb5&pipeline=agent&projects=argo&projectId=5f16b786f25d80a086a56bcb)


## Installation     
 

```sh
codefresh install gitops argocd-agent 
```

<img src="/art/installation.gif?raw=true" width="1200px">

## Usage
This tool require for use new codefresh argocd integration and environment view  

## Run tests
`go test -cover ./...`
