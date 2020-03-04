#!/usr/bin/env bash

kubectx docker-desktop

operator-sdk generate k8s

operator-sdk generate crds 
# go mod vendor

kubectl apply -f deploy/crds/sensu.k8s.sensu.io_sensuagents_crd.yaml 
kubectl apply -f deploy/crds/sensu.k8s.sensu.io_sensuassets_crd.yaml 
kubectl apply -f deploy/crds/sensu.k8s.sensu.io_sensubackends_crd.yaml 
kubectl apply -f deploy/crds/sensu.k8s.sensu.io_sensuchecks_crd.yaml 
kubectl apply -f deploy/crds/sensu.k8s.sensu.io_sensufilters_crd.yaml 
kubectl apply -f deploy/crds/sensu.k8s.sensu.io_sensuhandlers_crd.yaml 
kubectl apply -f deploy/crds/sensu.k8s.sensu.io_sensumutators_crd.yaml 
kubectl apply -f deploy/crds/sensu.k8s.sensu.io_sensunamespaces_crd.yaml 

export OPERATOR_NAME=sensu-operator
export SENSU_CA_CERTIFICATE=./ca.pem

kubectl apply -f k8s-namespace.yaml
kubectl apply -f sensu-backend-secrets.yaml
kubectl apply -f sensu-ca-secrets.yaml

operator-sdk up local --namespace=sensu