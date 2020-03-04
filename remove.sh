#!/usr/bin/env bash

kubectx docker-desktop

objects() {
    kubectl delete -f deploy/crds/sensu.k8s.sensu.io_v1alpha1_sensuagent_cr.yaml -n sensu 
    kubectl delete -f deploy/crds/sensu.k8s.sensu.io_v1alpha1_sensuasset_cr.yaml -n sensu
    kubectl delete -f deploy/crds/sensu.k8s.sensu.io_v1alpha1_sensucheck_cr.yaml -n sensu
    kubectl delete -f deploy/crds/sensu.k8s.sensu.io_v1alpha1_sensufilter_cr.yaml -n sensu
    kubectl delete -f deploy/crds/sensu.k8s.sensu.io_v1alpha1_sensuhandler_cr.yaml -n sensu
    kubectl delete -f deploy/crds/sensu.k8s.sensu.io_v1alpha1_sensumutator_cr.yaml -n sensu
    kubectl delete -f deploy/crds/sensu.k8s.sensu.io_v1alpha1_sensunamespace_cr.yaml -n sensu
    kubectl delete -f deploy/crds/sensu.k8s.sensu.io_v1alpha1_sensubackend_cr.yaml -n sensu
}

crds() {
    kubectl delete -f deploy/crds/sensu.k8s.sensu.io_sensuagents_crd.yaml 
    kubectl delete -f deploy/crds/sensu.k8s.sensu.io_sensuassets_crd.yaml 
    kubectl delete -f deploy/crds/sensu.k8s.sensu.io_sensuchecks_crd.yaml 
    kubectl delete -f deploy/crds/sensu.k8s.sensu.io_sensufilters_crd.yaml 
    kubectl delete -f deploy/crds/sensu.k8s.sensu.io_sensuhandlers_crd.yaml 
    kubectl delete -f deploy/crds/sensu.k8s.sensu.io_sensumutators_crd.yaml 
    kubectl delete -f deploy/crds/sensu.k8s.sensu.io_sensunamespaces_crd.yaml 
    kubectl delete -f deploy/crds/sensu.k8s.sensu.io_sensubackends_crd.yaml 
}


# unset OPERATOR_NAME
# unset SENSU_CA_CERTIFICATE

objects
crds

