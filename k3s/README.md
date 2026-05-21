# Example Deployment

this is for example deployment

## prerequisites

makesure you have a k3s cluster running that already have:
* traefik default ingress

## how to deploy (devbox)

if you are using devbox, you can simply do:

1. save your kubeconfig to `{this folder}/kubeconfig.yaml`
2. run `devbox run apply` or `kubectl apply -k ./manifests`

you can also deploy them separately:
* backing services: `kubectl apply -k ./manifests/backing-services`
* app services: `kubectl apply -k ./manifests/app-services`

all your stuff will be deployed in namespace=gaman.

## todo

* make secret not exposed
* pvc removed when we delete using kubectl remove -k