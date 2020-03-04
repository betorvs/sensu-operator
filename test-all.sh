#!/usr/bin/env bash
# Simple script to check all objects created by sensu-operator
# use these commands: kubectl, sensuctl, kubectx

create() {
    kubectx docker-desktop
    kubectl apply -f deploy/crds/sensu.k8s.sensu.io_v1alpha1_sensuagent_cr.yaml -n sensu 
    kubectl apply -f deploy/crds/sensu.k8s.sensu.io_v1alpha1_sensuasset_cr.yaml -n sensu
    kubectl apply -f deploy/crds/sensu.k8s.sensu.io_v1alpha1_sensucheck_cr.yaml -n sensu
    kubectl apply -f deploy/crds/sensu.k8s.sensu.io_v1alpha1_sensufilter_cr.yaml -n sensu
    kubectl apply -f deploy/crds/sensu.k8s.sensu.io_v1alpha1_sensuhandler_cr.yaml -n sensu
    kubectl apply -f deploy/crds/sensu.k8s.sensu.io_v1alpha1_sensumutator_cr.yaml -n sensu
    kubectl apply -f deploy/crds/sensu.k8s.sensu.io_v1alpha1_sensunamespace_cr.yaml -n sensu
}

check() {
    kubectx docker-desktop
    kubectl get sensuBackend -n sensu -o yaml
    kubectl get sensuAgent -n sensu -o yaml
    kubectl get sensuAsset -n sensu -o yaml
    kubectl get sensuCheck -n sensu -o yaml
    kubectl get sensuFilter -n sensu -o yaml
    kubectl get sensuHandler -n sensu -o yaml
    kubectl get sensuMutator -n sensu -o yaml
    kubectl get sensuNamespace -n sensu -o yaml
}

delete_objects() {
    kubectx docker-desktop
    kubectl delete -f deploy/crds/sensu.k8s.sensu.io_v1alpha1_sensuasset_cr.yaml -n sensu
    kubectl delete -f deploy/crds/sensu.k8s.sensu.io_v1alpha1_sensucheck_cr.yaml -n sensu
    kubectl delete -f deploy/crds/sensu.k8s.sensu.io_v1alpha1_sensufilter_cr.yaml -n sensu
    kubectl delete -f deploy/crds/sensu.k8s.sensu.io_v1alpha1_sensuhandler_cr.yaml -n sensu
    kubectl delete -f deploy/crds/sensu.k8s.sensu.io_v1alpha1_sensumutator_cr.yaml -n sensu
    kubectl delete -f deploy/crds/sensu.k8s.sensu.io_v1alpha1_sensunamespace_cr.yaml -n sensu
}

sensuctl_init() {
    sensuctl configure -n --username 'admin' --password 'P@ssw0rd!2GO' --namespace default --trusted-ca-file $PWD/ca.pem --url 'https://sensu-api.sensu.svc.cluster.local:8080'
}

sensuctl_list() (
    for value in asset check filter handler mutator namespace 
      do 
        sensuctl $value list 
      done
)

kubectl_get_all() {
    kubectx docker-desktop
    kubectl get all,secret -n sensu
}

case $1 in
    "create")
        create
    ;;
    "check")
        check 
    ;;
    "delete")
        delete_objects
    ;;
    "ctl")
        sensuctl_init
    ;;
    "verify")
        sensuctl_list
    ;;
    "get-all")
        kubectl_get_all
    ;;
    *)
        echo "$0 create|check|delete|ctl|verify|get-all"
        echo "create: all CR objects"
        echo "check: kubectl get sensuObject* -n sensu -o yaml"
        echo "delete: delete all CRs"
        echo "ctl: sensuctl configure"
        echo "verify: sensuctl Objects* list"
        echo "get-all: kubectl get all,secrets -n sensu"
esac