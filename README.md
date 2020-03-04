sensu-operator
==============

Sensu Operator is a open-source project based on [Sensu Go][1] using [operator-sdk][2] framework. It exposes Sensu API in Kubernetes as K8S objects, like sensuBackend or sensuCheck. 

In these implementation it uses 2 credentials:  
* default admin user: only used with sensuBackend object;  
* operator-user: used for all objects to access Sensu Backend API.   

**NOTE:** This sensu-operator only works with sensu backend with ssl enabled.

# Operational Modes

There are 2 operational modes weather you want to run the sensu backend on the same kubernetes cluster or just use the remote backend which is already running. If you choose to run the dedicated sensu backend on the cluster order sensu-operator to create your own.  

## Owner's Sensu Backend API

By default, sensu-operator running inside Kubernetes creates it own Sensu Backend Deployment and keeps polling it. It also checks if sensu backend api always running and responsive, if not, it will kill and recreate.   

## Remote Sensu Backend API

Using sensu-operator user, it access a remote Sensu Backend API to create all objects from Kubernetes.  

If you want to deploy Sensu Backend separetely, look into these 2 repositories as starting point:  
* [Sensu Backend using Statefulset][4]  
* [Sensu Backend using Deployment and ETCD outside][5]  

## Environment variables

SENSU_BACKEND_CLUSTER_ADMIN_USERNAME: default value "admin". Only used in sensuBackend object.  
SENSU_BACKEND_CLUSTER_ADMIN_PASSWORD: default value "P@ssw0rd!2GO". Only used in sensuBackend object.  
OPERATOR_SENSU_USER: default value "sensu-operator". All Kubernetes objects.  
OPERATOR_SENSU_PASSWORD: default value "P@ssw0rd!2GO". All Kubernetes objects.  

## Custom Resource Definition

|NAME              |SHORTNAMES |  APIGROUP           | NAMESPACED |  KIND| Example |
|------------------|-----------|---------------------|------------|------|---------|
|sensuagents       |           |  sensu.k8s.sensu.io  | true       |  SensuAgent| [agent](deploy/crds/sensu.k8s.sensu.io_v1alpha1_sensuagent_cr.yaml) |
|sensuassets       |           |  sensu.k8s.sensu.io  | true       |  SensuAsset| [asset](deploy/crds/sensu.k8s.sensu.io_v1alpha1_sensuasset_cr.yaml) |
|sensubackends     |           |  sensu.k8s.sensu.io  | true       |  SensuBackend| [backend](deploy/crds/sensu.k8s.sensu.io_v1alpha1_sensubackend_cr.yaml) |
|sensuchecks       |           |  sensu.k8s.sensu.io  | true       |  SensuCheck| [check](deploy/crds/sensu.k8s.sensu.io_v1alpha1_sensucheck_cr.yaml) |
|sensufilters      |           |  sensu.k8s.sensu.io  | true       |  SensuFilter| [filter](deploy/crds/sensu.k8s.sensu.io_v1alpha1_sensufilter_cr.yaml) |
|sensuhandlers     |           |  sensu.k8s.sensu.io  | true       |  SensuHandler| [handler](deploy/crds/sensu.k8s.sensu.io_v1alpha1_sensuhandler_cr.yaml) |
|sensumutators     |           |  sensu.k8s.sensu.io  | true       |  SensuMutator| [mutator](deploy/crds/sensu.k8s.sensu.io_v1alpha1_sensumutator_cr.yaml) |
|sensunamespaces   |           |  sensu.k8s.sensu.io  | true       |  SensuNamespace| [namespace](deploy/crds/sensu.k8s.sensu.io_v1alpha1_sensunamespace_cr.yaml) |

# Development

More information in [operator guide][7].  

Scripts: run.sh (regenerate crds and k8s apis, deploy all CRDs and run operator locally using kubectl configuration) and remove.sh (removes everything).

Order:  
1. Install operator-sdk and golang.
2. Generate all secrets using bellow instructions.
3. You must have kubectl and kubectx (or commented out kubectx lines)
4. Execute: `bash run.sh`  

After any code changes execute again `bash run.sh`.


# Build and run

Install [operator-sdk][3] and run:
```
operator-sdk build repository/sensu-operator:version
```

Push it to your docker repository:
```
docker push repository/sensu-operator:version
```

Modify [operator.yaml](deploy/operator.yaml) to use your own image.


# Deployment


## Create Sensu Certificates using cfssl

More information in [Sensu Secure][6].  

```sh
cd sensu-certs/
cfssl gencert -initca sensu-ca.json | cfssljson -bare ca
cfssl gencert -ca=ca.pem -ca-key=ca-key.pem -config=sensu-ca.json -profile=server sensu-backend.json | cfssljson -bare sensu-backend
```

### Create secrets from Certificates

```sh
kubectl create secret generic sensu-backend-pem --from-file=sensu-backend.pem=sensu-backend.pem \
    --from-file=sensu-backend-key.pem=sensu-backend-key.pem -n sensu --dry-run -o yaml > ../sensu-backend-secrets.yaml
kubectl create secret generic sensu-ca-pem --from-file=sensu-ca.pem=ca.pem -n sensu \
    --dry-run -o yaml > ../sensu-ca-secrets.yaml
```

### Create secret to keep admin and operator passwords

```
kubectl create secret generic sensu-operator --from-literal=adminpassword='P@ssw0rd!2GO' \
  --from-literal=operatorpassword='P@ssw0rd!2GO' \
  -n sensu --dry-run -o yaml > sensuoperator-secret.yaml
```

## To Deploy it in Kubernetes


### Deploy all CRD's

```
kubectl apply -f deploy/crds/sensu.k8s.sensu.io_sensuagents_crd.yaml 
kubectl apply -f deploy/crds/sensu.k8s.sensu.io_sensuassets_crd.yaml 
kubectl apply -f deploy/crds/sensu.k8s.sensu.io_sensubackends_crd.yaml 
kubectl apply -f deploy/crds/sensu.k8s.sensu.io_sensuchecks_crd.yaml 
kubectl apply -f deploy/crds/sensu.k8s.sensu.io_sensufilters_crd.yaml 
kubectl apply -f deploy/crds/sensu.k8s.sensu.io_sensuhandlers_crd.yaml 
kubectl apply -f deploy/crds/sensu.k8s.sensu.io_sensumutators_crd.yaml 
kubectl apply -f deploy/crds/sensu.k8s.sensu.io_sensunamespaces_crd.yaml 
```


### Deploy Operator

```
kubectl create -f k8s-namespace.yaml
kubectl create -f deploy/service_account.yaml -n sensu
kubectl create -f deploy/role.yaml -n sensu
kubectl create -f deploy/role_binding.yaml -n sensu
kubectl create -f deploy/operator.yaml -n sensu
```

## Contributing

Any help are welcome!


[1]: https://github.com/sensu/sensu-go  
[2]: https://github.com/operator-framework/operator-sdk  
[3]: https://github.com/operator-framework/operator-sdk/blob/master/doc/user/install-operator-sdk.md
[4]: https://github.com/betorvs/sensu-go-statefulset
[5]: https://github.com/betorvs/sensu-go-deployment  
[6]: https://docs.sensu.io/sensu-go/latest/guides/securing-sensu/#secure-the-api-and-dashboard  
[7]: https://github.com/operator-framework/operator-sdk/blob/master/doc/user-guide.md