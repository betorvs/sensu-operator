package sensubackend

import (
	"context"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"

	sensuv1alpha1 "github.com/betorvs/sensu-operator/pkg/apis/sensu/v1alpha1"
	"github.com/betorvs/sensu-operator/pkg/config"
	"github.com/betorvs/sensu-operator/pkg/usecase"
	"github.com/betorvs/sensu-operator/pkg/utils"
	"github.com/go-logr/logr"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

const (
	versionString = "latest"
)

var log = logf.Log.WithName("controller_sensubackend")

// Add creates a new SensuBackend Controller and adds it to the Manager. The Manager will set fields on the Controller
// and Start it when the Manager is Started.
func Add(mgr manager.Manager) error {
	return add(mgr, newReconciler(mgr))
}

// newReconciler returns a new reconcile.Reconciler
func newReconciler(mgr manager.Manager) reconcile.Reconciler {
	return &ReconcileSensuBackend{client: mgr.GetClient(), scheme: mgr.GetScheme()}
}

// add adds a new Controller to mgr with r as the reconcile.Reconciler
func add(mgr manager.Manager, r reconcile.Reconciler) error {
	// Create a new controller
	c, err := controller.New("sensubackend-controller", mgr, controller.Options{Reconciler: r})
	if err != nil {
		return err
	}

	// Watch for changes to primary resource SensuBackend
	err = c.Watch(&source.Kind{Type: &sensuv1alpha1.SensuBackend{}}, &handler.EnqueueRequestForObject{})
	if err != nil {
		return err
	}

	// TODO(user): Modify this to be the types you create that are owned by the primary resource
	// Watch for changes to secondary resource Pods and requeue the owner SensuBackend
	err = c.Watch(&source.Kind{Type: &corev1.Pod{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &sensuv1alpha1.SensuBackend{},
	})
	if err != nil {
		return err
	}

	// watch headless service
	err = c.Watch(&source.Kind{Type: &corev1.Service{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &sensuv1alpha1.SensuBackend{},
	})
	if err != nil {
		return err
	}

	// watch secret
	err = c.Watch(&source.Kind{Type: &corev1.Secret{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &sensuv1alpha1.SensuBackend{},
	})
	if err != nil {
		return err
	}

	return nil
}

// blank assignment to verify that ReconcileSensuBackend implements reconcile.Reconciler
var _ reconcile.Reconciler = &ReconcileSensuBackend{}

const sensuBackendFinalizer = "finalizer.sensu.k8s.sensu.io"

// ReconcileSensuBackend reconciles a SensuBackend object
type ReconcileSensuBackend struct {
	// This client, initialized using mgr.Client() above, is a split client
	// that reads objects from the cache and writes to the apiserver
	client client.Client
	scheme *runtime.Scheme
}

// Reconcile reads that state of the cluster for a SensuBackend object and makes changes based on the state read
// and what is in the SensuBackend.Spec
// The Controller will requeue the Request to be processed again if the returned error is non-nil or
// Result.Requeue is true, otherwise upon completion it will remove the work from the queue.
func (r *ReconcileSensuBackend) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	reqLogger := log.WithValues("Request.Namespace", request.Namespace, "Request.Name", request.Name)
	reqLogger.Info("Reconciling SensuBackend")

	// Fetch the SensuBackend instance
	instance := &sensuv1alpha1.SensuBackend{}
	err := r.client.Get(context.TODO(), request.NamespacedName, instance)
	if err != nil {
		if errors.IsNotFound(err) {
			// Request object not found, could have been deleted after reconcile request.
			// Owned objects are automatically garbage collected. For additional cleanup logic use finalizers.
			// Return and don't requeue
			return reconcile.Result{}, nil
		}
		// Error reading the object - requeue the request.
		return reconcile.Result{}, err
	}

	// Check for headless service
	headless := r.newHeadlessService(instance)

	err = r.client.Get(context.TODO(), types.NamespacedName{Name: "sensu", Namespace: headless.Namespace}, headless)
	if err != nil && errors.IsNotFound(err) {
		reqLogger.Info("Creating a new Headless Service", "Service.Namespace", headless.Namespace, "service.Name", headless.Name)
		err = r.client.Create(context.TODO(), headless)
		if err != nil {
			return reconcile.Result{}, err
		}
	} else if err != nil {
		return reconcile.Result{}, err
	}

	// Check for sensuSecret
	sensuSecret := r.newSensuSecret(instance)

	err = r.client.Get(context.TODO(), types.NamespacedName{Name: "sensu-secret", Namespace: sensuSecret.Namespace}, sensuSecret)
	if err != nil && errors.IsNotFound(err) {
		reqLogger.Info("Creating a new Sensu Secret", "Secret.Namespace", sensuSecret.Namespace, "secret.Name", sensuSecret.Name)
		err = r.client.Create(context.TODO(), sensuSecret)
		if err != nil {
			return reconcile.Result{}, err
		}
	} else if err != nil {
		return reconcile.Result{}, err
	}

	// Define a new Pod object
	newpod := r.newPodForCR(instance, "new", 1)

	// Check if this Pod already exists
	found := &corev1.Pod{}
	err = r.client.Get(context.TODO(), types.NamespacedName{Name: newpod.Name, Namespace: newpod.Namespace}, found)
	if err != nil && errors.IsNotFound(err) {
		reqLogger.Info("Creating a new Pod", "Pod.Namespace", newpod.Namespace, "Pod.Name", newpod.Name)
		// newpod := r.newPodForCR(instance, "new", 1)
		err = r.client.Create(context.TODO(), newpod)
		if err != nil {
			return reconcile.Result{}, err
		}

		// Pod created successfully - don't requeue
		return reconcile.Result{}, nil
	} else if err != nil {
		return reconcile.Result{}, err
	}

	// Check for sensuService service
	sensuService := r.newSensuAPIService(instance)

	err = r.client.Get(context.TODO(), types.NamespacedName{Name: "sensu-api", Namespace: sensuService.Namespace}, sensuService)
	if err != nil && errors.IsNotFound(err) {
		reqLogger.Info("Creating a new Sensu API Service", "Service.Namespace", sensuService.Namespace, "service.Name", sensuService.Name)
		err = r.client.Create(context.TODO(), sensuService)
		if err != nil {
			return reconcile.Result{}, err
		}
	} else if err != nil {
		return reconcile.Result{}, err
	}

	// Update the pod status with the pod names
	// List the pods for this pod's deployment
	podList := &corev1.PodList{}
	listOpts := []client.ListOption{
		client.InNamespace(instance.Namespace),
		client.MatchingLabels(map[string]string{
			"app": instance.Name,
		}),
	}
	err = r.client.List(context.TODO(), podList, listOpts...)
	if err != nil {
		reqLogger.Error(err, "Failed to list pods.", "pod.Namespace", instance.Namespace, "pod.Name", instance.Name)
		return reconcile.Result{}, err
	}
	podNames := getPodNames(podList.Items)
	// podNames := []string{}
	// Count the pods that are pending or running as available
	// for _, pod := range podList.Items {
	// 	if pod.GetObjectMeta().GetDeletionTimestamp() != nil {
	// 		continue
	// 	}
	// 	if pod.Status.Phase == corev1.PodPending || pod.Status.Phase == corev1.PodRunning {
	// 		podNames = append(podNames, pod.GetObjectMeta().GetName())
	// 	}
	// }

	// Update status.Nodes if needed
	if !reflect.DeepEqual(podNames, instance.Status.Nodes) {
		instance.Status.Nodes = podNames
		err := r.client.Status().Update(context.TODO(), instance)
		if err != nil {
			reqLogger.Error(err, "Failed to update SensuBackend Nodes Status.")
			return reconcile.Result{}, err
		}
	}

	// Update the svc status with the pod names
	// List the svcs for this svc's deployment
	svcList := &corev1.ServiceList{}
	listsvcOpts := []client.ListOption{
		client.InNamespace(instance.Namespace),
		client.MatchingLabels(map[string]string{
			"app": instance.Name,
		}),
	}
	err = r.client.List(context.TODO(), svcList, listsvcOpts...)
	if err != nil {
		reqLogger.Error(err, "Failed to list pods.", "svc.Namespace", instance.Namespace, "svc.Name", instance.Name)
		return reconcile.Result{}, err
	}
	svcNames := getSvcNames(svcList.Items)

	// Update status.Nodes if needed
	if !reflect.DeepEqual(svcNames, instance.Status.Nodes) {
		instance.Status.Services = svcNames
		err := r.client.Status().Update(context.TODO(), instance)
		if err != nil {
			reqLogger.Error(err, "Failed to update SensuBackend Services Status.")
			return reconcile.Result{}, err
		}
	}

	// access sensu api and get a new api token
	if instance.Status.AdminToken == "" && instance.Status.ClusterID == "" {
		sensuURL := fmt.Sprintf("https://%s", instance.Spec.SensuBackendURL)
		for {
			time.Sleep(5 * time.Second)
			if usecase.SensuHealth(sensuURL) {
				break
			}
		}
		token := usecase.GetSensuAPIToken(sensuURL)
		instance.Status.AdminToken = token
		err := r.client.Status().Update(context.TODO(), instance)
		if err != nil {
			reqLogger.Error(err, "Failed to update SensuBackend Status AdminToken.")
			return reconcile.Result{}, err
		}

	} else {
		sensuURL := fmt.Sprintf("https://%s", instance.Spec.SensuBackendURL)
		if !usecase.SensuTestToken(sensuURL, instance.Status.AdminToken) {
			instance.Status.AdminToken = ""
			err := r.client.Status().Update(context.TODO(), instance)
			if err != nil {
				reqLogger.Error(err, "Failed to update SensuBackend Status AdminToken with empty status.")
				return reconcile.Result{}, err
			}
		}
	}

	if instance.Status.OperatorToken == "" && instance.Status.ClusterID == "" {
		sensuURL := fmt.Sprintf("https://%s", instance.Spec.SensuBackendURL)
		for {
			time.Sleep(5 * time.Second)
			if usecase.SensuHealth(sensuURL) {
				break
			}
		}
		token := usecase.CreateOperatorUserGetToken(sensuURL)
		instance.Status.OperatorToken = token
		err := r.client.Status().Update(context.TODO(), instance)
		if err != nil {
			reqLogger.Error(err, "Failed to update SensuBackend Status OperatorToken.")
			return reconcile.Result{}, err
		}

	} else {
		sensuURL := fmt.Sprintf("https://%s", instance.Spec.SensuBackendURL)
		if !usecase.SensuTestToken(sensuURL, instance.Status.OperatorToken) {
			instance.Status.OperatorToken = ""
			err := r.client.Status().Update(context.TODO(), instance)
			if err != nil {
				reqLogger.Error(err, "Failed to update SensuBackend Status OperatorToken with empty status.")
				return reconcile.Result{}, err
			}
		}
	}

	if instance.Status.Token == "" && instance.Status.ClusterID == "" {
		sensuURL := fmt.Sprintf("https://%s", instance.Spec.SensuBackendURL)
		for {
			time.Sleep(5 * time.Second)
			if usecase.SensuHealth(sensuURL) {
				break
			}
		}
		operatorToken := instance.Status.OperatorToken
		adminToken := instance.Status.AdminToken
		if usecase.SensuTestToken(sensuURL, operatorToken) {
			instance.Status.Token = operatorToken
			err := r.client.Status().Update(context.TODO(), instance)
			if err != nil {
				reqLogger.Error(err, "Failed to update SensuBackend Status Token with Operator Token.")
				return reconcile.Result{}, err
			}
		} else if usecase.SensuTestToken(sensuURL, adminToken) {
			instance.Status.Token = adminToken
			err := r.client.Status().Update(context.TODO(), instance)
			if err != nil {
				reqLogger.Error(err, "Failed to update SensuBackend Status Token with Admin Token.")
				return reconcile.Result{}, err
			}
		} else {
			reqLogger.Info("Both Tokens are invalids", "pod.Namespace", instance.Namespace, "pod.Name", instance.Name)
			instance.Status.Token = ""
			err := r.client.Status().Update(context.TODO(), instance)
			if err != nil {
				reqLogger.Error(err, "Failed to update SensuBackend Status Token with empty status.")
				return reconcile.Result{}, err
			}
		}
	} else {
		sensuURL := fmt.Sprintf("https://%s", instance.Spec.SensuBackendURL)
		if !usecase.SensuTestToken(sensuURL, instance.Status.Token) {
			instance.Status.Token = ""
			err := r.client.Status().Update(context.TODO(), instance)
			if err != nil {
				reqLogger.Error(err, "Failed to update SensuBackend Status Token with empty status token.")
				return reconcile.Result{}, err
			}
		}
	}

	if instance.Status.ClusterID == "" || instance.Status.ClusterID == "pending" || instance.Status.ClusterID == "replacing" {
		sensuURL := fmt.Sprintf("https://%s", instance.Spec.SensuBackendURL)
		clusterID := usecase.GetClusterID(sensuURL, instance.Status.Token)
		instance.Status.ClusterID = clusterID
		err := r.client.Status().Update(context.TODO(), instance)
		if err != nil {
			reqLogger.Error(err, "Failed to update SensuBackend ClusterID Status.")
			return reconcile.Result{}, err
		}
	}

	// Scale Up Pods
	numberOfPods := int32(len(podNames))
	if instance.Spec.Replicas > numberOfPods {
		reqLogger.Info("Waiting for 15 seconds to proceed Sensu Backend Scale UP", "Cluster ID", instance.Status.ClusterID)
		time.Sleep(15 * time.Second)
		sensuURL := fmt.Sprintf("https://%s", instance.Spec.SensuBackendURL)
		newMemberNumber := numberOfPods + 1
		newMember := fmt.Sprintf("%s-%d.sensu.sensu.svc.cluster.local", instance.Name, newMemberNumber)
		if !usecase.CheckMemberExist(sensuURL, instance.Status.Token, newMember) {
			reqLogger.Info("Adding a new pod to SensuBackend", "Expected Replicas", instance.Spec.Replicas, "Pod.Names", podNames)

			err = usecase.AddNewMember(sensuURL, instance.Status.Token, newMember)
			if err != nil {
				return reconcile.Result{}, err
			}
			newpod := r.newPodForCR(instance, "existing", newMemberNumber)
			err = r.client.Create(context.TODO(), newpod)
			if err != nil {
				return reconcile.Result{}, err
			}
			reqLogger.Info("Waiting for 15 seconds to proceed", "Cluster ID", instance.Status.ClusterID)
			time.Sleep(15 * time.Second)

		}
		return reconcile.Result{}, nil
	} else if instance.Spec.Replicas < numberOfPods {
		reqLogger.Info("Waiting for 15 seconds to proceed Sensu Backend Downsize", "Cluster ID", instance.Status.ClusterID)
		time.Sleep(15 * time.Second)
		sensuURL := fmt.Sprintf("https://%s", instance.Spec.SensuBackendURL)
		newMemberNumber := numberOfPods
		member := fmt.Sprintf("%s-%d.sensu.sensu.svc.cluster.local", instance.Name, newMemberNumber)
		if usecase.CheckMemberExist(sensuURL, instance.Status.Token, member) {
			reqLogger.Info("Removing a pod to SensuBackend", "Expected Replicas", instance.Spec.Replicas, "Pod.Names", podNames)
			err = usecase.RemoveMember(sensuURL, instance.Status.Token, member)
			if err != nil {
				return reconcile.Result{}, err
			}
			deletedPod := r.newPodForCR(instance, "existing", newMemberNumber)
			err = r.client.Delete(context.TODO(), deletedPod)
			if err != nil {
				reqLogger.Error(err, "Failed to delete pod", "pod.name", deletedPod.Name)
				return reconcile.Result{}, err
			}
		}

		return reconcile.Result{}, nil
	}

	// Check if sensu api still running
	if numberOfPods > 1 && instance.Status.ClusterID != "" {
		sensuURL := fmt.Sprintf("https://%s", instance.Spec.SensuBackendURL)
		reqLogger.Info("Waiting for 15 seconds to proceed with Sensu API Health Check", "Cluster ID", instance.Status.ClusterID)
		time.Sleep(15 * time.Second)
		reqLogger.Info("Checking if Sensu API is Health First Check", "Cluster ID", instance.Status.ClusterID)
		if !usecase.SensuHealth(sensuURL) {
			time.Sleep(15 * time.Second)
			reqLogger.Info("Checking if Sensu API is Health Second Check after 15 seconds", "Cluster ID", instance.Status.ClusterID)
			if !usecase.SensuHealth(sensuURL) {
				reqLogger.Info("Sensu API Not Health, killing", "Cluster ID", instance.Status.ClusterID)
				reqLogger.Info("Deleting all pods to SensuBackend", "Expected Replicas", instance.Spec.Replicas, "Pod.Names", podNames)
				for _, v := range podNames {
					newMemberNumber := strings.Split(v, "-")
					targetNumber, _ := strconv.Atoi(newMemberNumber[1])
					parsedNumber := int32(targetNumber)
					deletedPod := r.newPodForCR(instance, "existing", parsedNumber)
					reqLogger.Info("Deleting", "Sensu Node", v)
					err := r.client.Delete(context.TODO(), deletedPod)
					if err != nil {
						return reconcile.Result{}, err
					}
				}
				// instance.Status.Nodes = []string{}
				instance.Status.ClusterID = "pending"
				err := r.client.Status().Update(context.TODO(), instance)
				if err != nil {
					reqLogger.Error(err, "Failed to update SensuBackend Status ClusterID with pendind status.")
					return reconcile.Result{}, err
				}
				return reconcile.Result{Requeue: true, RequeueAfter: 60}, nil

			}
		}
		// check if sensu version match
		testVersion := strings.Split(instance.Spec.Image, ":")
		version := testVersion[1]
		if version == versionString {
			version = config.SensuImageTag
		}
		reqLogger.Info("Checking if Sensu API version match", "Version", version)
		if !usecase.SensuVersion(sensuURL, version) {
			reqLogger.Info("Versions does not match", "Version Requested", version)
			reqLogger.Info("Deleting all pods to SensuBackend", "Expected Replicas", instance.Spec.Replicas, "Pod.Names", podNames)
			for _, v := range podNames {
				newMemberNumber := strings.Split(v, "-")
				targetNumber, _ := strconv.Atoi(newMemberNumber[1])
				parsedNumber := int32(targetNumber)
				deletedPod := r.newPodForCR(instance, "existing", parsedNumber)
				reqLogger.Info("Deleting", "Sensu Node", v)
				err := r.client.Delete(context.TODO(), deletedPod)
				if err != nil {
					return reconcile.Result{}, err
				}
			}
			// instance.Status.Nodes = []string{""}
			instance.Status.ClusterID = "replacing"
			err := r.client.Status().Update(context.TODO(), instance)
			if err != nil {
				reqLogger.Error(err, "Failed to update SensuBackend Status ClusterID with pending status.")
				return reconcile.Result{}, err
			}
			return reconcile.Result{Requeue: true}, nil
		}

	}

	// Check if the Sensu Backend instance is marked to be deleted, which is
	// indicated by the deletion timestamp being set.
	isSensuBackendMarkedToBeDeleted := instance.GetDeletionTimestamp() != nil
	if isSensuBackendMarkedToBeDeleted {
		if utils.Contains(instance.GetFinalizers(), sensuBackendFinalizer) {
			// Run finalization logic for sensuBackendFinalizer. If the
			// finalization logic fails, don't remove the finalizer so
			// that we can retry during the next reconciliation.
			if err := r.finalizeSensuBackend(reqLogger, instance); err != nil {
				return reconcile.Result{}, err
			}
			reqLogger.Info("Starting finalizer", "Pod.Namespace", found.Namespace)

			// Remove sensuBackendFinalizer. Once all finalizers have been
			// removed, the object will be deleted.
			instance.SetFinalizers(utils.Remove(instance.GetFinalizers(), sensuBackendFinalizer))
			err := r.client.Update(context.TODO(), instance)
			if err != nil {
				return reconcile.Result{}, err
			}
		}
		return reconcile.Result{}, nil
	}

	// Add finalizer for this CR
	if !utils.Contains(instance.GetFinalizers(), sensuBackendFinalizer) {
		if err := r.addFinalizer(reqLogger, instance); err != nil {
			return reconcile.Result{}, err
		}
	}
	// Pod already exists - don't requeue
	reqLogger.Info("Skip reconcile: Pod already exists", "Pod.Namespace", found.Namespace, "Pod.Name", podNames)
	return reconcile.Result{}, nil
}

// newPodForCR returns a busybox pod with the same name/namespace as the cr
func (r *ReconcileSensuBackend) newPodForCR(cr *sensuv1alpha1.SensuBackend, clusterInitialState string, ordinal int32) *corev1.Pod {
	labels := map[string]string{
		"app": cr.Name,
	}
	var setImage string
	if cr.Spec.Image == "" {
		setImage = fmt.Sprintf("sensu/sensu:%s", config.SensuImageTag)
	} else {
		setImage = cr.Spec.Image
	}
	var setDebug string
	if cr.Spec.DebugSensu {
		setDebug = "--log-level debug --debug"
	} else {
		setDebug = "--log-level info"
	}
	addEnvFrom := []corev1.EnvFromSource{}
	if cr.Spec.SecretEnvFrom != "" {
		addEnvFrom = []corev1.EnvFromSource{
			{
				SecretRef: &corev1.SecretEnvSource{
					LocalObjectReference: corev1.LocalObjectReference{
						Name: cr.Spec.SecretEnvFrom,
					},
				},
			},
		}
	}
	volumeMounts := []corev1.VolumeMount{
		{
			MountPath: "/certs",
			Name:      "sensu-backend-pem",
			ReadOnly:  true,
		},
		{
			MountPath: "/certs-ca",
			Name:      "sensu-ca-pem",
			ReadOnly:  true,
		},
	}
	volume := []corev1.Volume{
		{
			Name: "sensu-backend-pem",
			VolumeSource: corev1.VolumeSource{
				Secret: &corev1.SecretVolumeSource{
					SecretName: "sensu-backend-pem",
				},
			},
		},
		{
			Name: "sensu-ca-pem",
			VolumeSource: corev1.VolumeSource{
				Secret: &corev1.SecretVolumeSource{
					SecretName: "sensu-ca-pem",
				},
			},
		},
	}
	if cr.Spec.SecretVolume != "" {
		volumeMounts = []corev1.VolumeMount{
			{
				MountPath: "/certs",
				Name:      "sensu-backend-pem",
				ReadOnly:  true,
			},
			{
				MountPath: "/certs-ca",
				Name:      "sensu-ca-pem",
				ReadOnly:  true,
			},
			{
				MountPath: "/etc/secrets",
				Name:      cr.Spec.SecretVolume,
				ReadOnly:  true,
			},
		}
		volume = []corev1.Volume{
			{
				Name: "sensu-backend-pem",
				VolumeSource: corev1.VolumeSource{
					Secret: &corev1.SecretVolumeSource{
						SecretName: "sensu-backend-pem",
					},
				},
			},
			{
				Name: "sensu-ca-pem",
				VolumeSource: corev1.VolumeSource{
					Secret: &corev1.SecretVolumeSource{
						SecretName: "sensu-ca-pem",
					},
				},
			},
			{
				Name: cr.Spec.SecretVolume,
				VolumeSource: corev1.VolumeSource{
					Secret: &corev1.SecretVolumeSource{
						SecretName: cr.Spec.SecretVolume,
					},
				},
			},
		}
	}
	pod := &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      cr.Name + fmt.Sprintf("-%d", ordinal),
			Namespace: cr.Namespace,
			Labels:    labels,
		},
		Spec: corev1.PodSpec{
			Hostname:  cr.Name + fmt.Sprintf("-%d", ordinal),
			Subdomain: "sensu",
			InitContainers: []corev1.Container{{
				Name:    "initdns",
				Image:   "busybox:1.28",
				Command: []string{"sh", "-c", "until nslookup sensu." + cr.Namespace + ".svc.cluster.local; do echo waiting for myservice; sleep 2; done;"},
			}},
			Containers: []corev1.Container{
				{
					Name:    "sensu-backend",
					Image:   setImage,
					Command: []string{"sh", "-ec", "sensu-backend start --etcd-name ${HOSTNAME}.sensu." + cr.Namespace + ".svc.cluster.local --etcd-discovery-srv sensu." + cr.Namespace + ".svc.cluster.local --etcd-initial-advertise-peer-urls http://${HOSTNAME}.sensu." + cr.Namespace + ".svc.cluster.local:2380 --etcd-initial-cluster-token sensu --etcd-initial-cluster-state " + clusterInitialState + " --etcd-advertise-client-urls http://${HOSTNAME}.sensu." + cr.Namespace + ".svc.cluster.local:2379 --etcd-listen-client-urls http://0.0.0.0:2379 --etcd-listen-peer-urls http://0.0.0.0:2380 --state-dir /var/lib/sensu/sensu-backend/${HOSTNAME} " + setDebug + " --trusted-ca-file /certs-ca/sensu-ca.pem --cert-file /certs/sensu-backend.pem --key-file /certs/sensu-backend-key.pem --api-url https://" + cr.Spec.SensuBackendURL},
					Ports: []corev1.ContainerPort{
						{
							ContainerPort: 2379,
							Name:          "client",
						},
						{
							ContainerPort: 2380,
							Name:          "server",
						},
						{
							ContainerPort: 8080,
							Name:          "api",
						},
						{
							ContainerPort: 8081,
							Name:          "websocket",
						},
						{
							ContainerPort: 3000,
							Name:          "dashboard",
						},
					},
					VolumeMounts: volumeMounts,
					EnvFrom:      addEnvFrom,
					Env: []corev1.EnvVar{
						{
							Name:  "SENSU_BACKEND_CLUSTER_ADMIN_USERNAME",
							Value: config.DefaultUser,
						},
						{
							Name: "SENSU_BACKEND_CLUSTER_ADMIN_PASSWORD",
							ValueFrom: &corev1.EnvVarSource{
								SecretKeyRef: &corev1.SecretKeySelector{
									LocalObjectReference: corev1.LocalObjectReference{Name: "sensu-secret"},
									Key:                  "SENSU_BACKEND_CLUSTER_ADMIN_PASSWORD",
								},
							},
						},
					},
				},
			},
			Volumes: volume,
		},
	}
	// Set SensuBackend instance as the owner and controller
	_ = controllerutil.SetControllerReference(cr, pod, r.scheme)
	return pod
}

// newHeadlessService func
func (r *ReconcileSensuBackend) newHeadlessService(cr *sensuv1alpha1.SensuBackend) *corev1.Service {
	labels := map[string]string{
		"app": cr.Name,
	}
	headless := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "sensu",
			Namespace: cr.Namespace,
			Labels:    labels,
		},
		Spec: corev1.ServiceSpec{
			PublishNotReadyAddresses: true,
			ClusterIP:                "None",
			Selector:                 labels,
			Ports: []corev1.ServicePort{
				{
					Port: 2379,
					Name: "etcd-client",
				},
				{
					Port: 2380,
					Name: "etcd-server",
				},
			},
		},
	}
	// Set SensuBackend instance as the owner and controller
	_ = controllerutil.SetControllerReference(cr, headless, r.scheme)
	return headless
}

// newSensuAPIService func
func (r *ReconcileSensuBackend) newSensuAPIService(cr *sensuv1alpha1.SensuBackend) *corev1.Service {
	labels := map[string]string{
		"app": cr.Name,
	}
	sensuAPI := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "sensu-api",
			Namespace: cr.Namespace,
			Labels:    labels,
		},
		Spec: corev1.ServiceSpec{
			Selector: labels,
			Type:     "LoadBalancer",
			Ports: []corev1.ServicePort{
				{
					Port: 8080,
					Name: "api",
				},
				{
					Port: 8081,
					Name: "websocket",
				},
				{
					Port: 3000,
					Name: "dashboard",
				},
			},
		},
	}
	// Set SensuBackend instance as the owner and controller
	_ = controllerutil.SetControllerReference(cr, sensuAPI, r.scheme)
	return sensuAPI
}

// newSensuSecret func
func (r *ReconcileSensuBackend) newSensuSecret(cr *sensuv1alpha1.SensuBackend) *corev1.Secret {
	labels := map[string]string{
		"app": cr.Name,
	}
	sensuSecret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "sensu-secret",
			Namespace: cr.Namespace,
			Labels:    labels,
		},
		StringData: map[string]string{
			"SENSU_BACKEND_CLUSTER_ADMIN_PASSWORD": config.DefaultPassword,
		},
	}
	_ = controllerutil.SetControllerReference(cr, sensuSecret, r.scheme)
	return sensuSecret
}

// getPodNames returns the pod names of the array of pods passed in
func getPodNames(pods []corev1.Pod) []string {
	var podNames []string
	for _, pod := range pods {
		if pod.GetObjectMeta().GetDeletionTimestamp() != nil {
			continue
		}
		if pod.Status.Phase == corev1.PodPending || pod.Status.Phase == corev1.PodRunning {
			podNames = append(podNames, pod.GetObjectMeta().GetName())
		}
	}
	return podNames
}

// getSvcNames returns the svc names of the array of pods passed in
func getSvcNames(svcs []corev1.Service) []string {
	var svcNames []string
	for _, svc := range svcs {
		svcNames = append(svcNames, svc.Name)
	}
	return svcNames
}

func (r *ReconcileSensuBackend) finalizeSensuBackend(reqLogger logr.Logger, cr *sensuv1alpha1.SensuBackend) error {
	// Log finalizer
	reqLogger.Info("Successfully finalized Sensu Backend")
	return nil
}

func (r *ReconcileSensuBackend) addFinalizer(reqLogger logr.Logger, cr *sensuv1alpha1.SensuBackend) error {
	reqLogger.Info("Adding Finalizer for the Sensu Backend")
	cr.SetFinalizers(append(cr.GetFinalizers(), sensuBackendFinalizer))

	// Update CR
	err := r.client.Update(context.TODO(), cr)
	if err != nil {
		reqLogger.Error(err, "Failed to update Sensu Backend with finalizer")
		return err
	}
	return nil
}
