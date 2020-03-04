package sensuagent

import (
	"context"
	"fmt"
	"strings"

	sensuv1alpha1 "github.com/betorvs/sensu-operator/pkg/apis/sensu/v1alpha1"
	"github.com/betorvs/sensu-operator/pkg/config"
	appsv1 "k8s.io/api/apps/v1"
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

var log = logf.Log.WithName("controller_sensuagent")

/**
* USER ACTION REQUIRED: This is a scaffold file intended for the user to modify with their own Controller
* business logic.  Delete these comments after modifying this file.*
 */

// Add creates a new SensuAgent Controller and adds it to the Manager. The Manager will set fields on the Controller
// and Start it when the Manager is Started.
func Add(mgr manager.Manager) error {
	return add(mgr, newReconciler(mgr))
}

// newReconciler returns a new reconcile.Reconciler
func newReconciler(mgr manager.Manager) reconcile.Reconciler {
	return &ReconcileSensuAgent{client: mgr.GetClient(), scheme: mgr.GetScheme()}
}

// add adds a new Controller to mgr with r as the reconcile.Reconciler
func add(mgr manager.Manager, r reconcile.Reconciler) error {
	// Create a new controller
	c, err := controller.New("sensuagent-controller", mgr, controller.Options{Reconciler: r})
	if err != nil {
		return err
	}

	// Watch for changes to primary resource SensuAgent
	err = c.Watch(&source.Kind{Type: &sensuv1alpha1.SensuAgent{}}, &handler.EnqueueRequestForObject{})
	if err != nil {
		return err
	}

	return nil
}

// blank assignment to verify that ReconcileSensuAgent implements reconcile.Reconciler
var _ reconcile.Reconciler = &ReconcileSensuAgent{}

// ReconcileSensuAgent reconciles a SensuAgent object
type ReconcileSensuAgent struct {
	// This client, initialized using mgr.Client() above, is a split client
	// that reads objects from the cache and writes to the apiserver
	client client.Client
	scheme *runtime.Scheme
}

// Reconcile reads that state of the cluster for a SensuAgent object and makes changes based on the state read
// and what is in the SensuAgent.Spec
// The Controller will requeue the Request to be processed again if the returned error is non-nil or
// Result.Requeue is true, otherwise upon completion it will remove the work from the queue.
func (r *ReconcileSensuAgent) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	reqLogger := log.WithValues("Request.Namespace", request.Namespace, "Request.Name", request.Name)
	reqLogger.Info("Reconciling SensuAgent")

	// Fetch the SensuAgent instance
	instance := &sensuv1alpha1.SensuAgent{}
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

	// Define a new Deployment object
	deploy := r.newDeploymentForCR(instance)

	// Check if this Deployment already exists
	found := &appsv1.Deployment{}
	err = r.client.Get(context.TODO(), types.NamespacedName{Name: deploy.Name, Namespace: deploy.Namespace}, found)
	if err != nil && errors.IsNotFound(err) {
		reqLogger.Info("Creating a new Deployment", "Deployment.Namespace", deploy.Namespace, "Deployment.Name", deploy.Name)
		err = r.client.Create(context.TODO(), deploy)
		if err != nil {
			return reconcile.Result{}, err
		}

		// Pod created successfully - don't requeue
		return reconcile.Result{}, nil
	} else if err != nil {
		return reconcile.Result{}, err
	}

	// Pod already exists - don't requeue
	reqLogger.Info("Skip reconcile: Deployment already exists", "Deployment.Namespace", found.Namespace, "Deployment.Name", found.Name)
	return reconcile.Result{}, nil
}

// newDeploymentForCR returns a Deployment with the same name/namespace as the cr
func (r *ReconcileSensuAgent) newDeploymentForCR(cr *sensuv1alpha1.SensuAgent) *appsv1.Deployment {
	labels := map[string]string{
		"app": cr.Name,
	}
	var replicas int32
	if cr.Spec.Replicas == 0 {
		replicas = 1
	} else {
		replicas = cr.Spec.Replicas
	}
	var setImage string
	if cr.Spec.Image == "" {
		setImage = fmt.Sprintf("sensu/sensu:%s", config.SensuImageTag)
	} else {
		setImage = cr.Spec.Image
	}
	var setLogLevel string
	if cr.Spec.LogLevel == "" {
		setLogLevel = "info"
	} else {
		setLogLevel = cr.Spec.LogLevel
	}
	var sensuBackendWebsocket string
	if cr.Spec.SensuBackendWebsocket == "" {
		sensuBackendWebsocket = "wss://sensu-api.sensu.svc.cluster.local:8081"
	} else {
		sensuBackendWebsocket = fmt.Sprintf("wss://%s", cr.Spec.SensuBackendWebsocket)
	}
	var secretCertificate string
	if cr.Spec.CACertificate == "" {
		secretCertificate = "sensu-ca-pem"
	} else {
		secretCertificate = cr.Spec.CACertificate
	}
	var fileNameCA string
	if cr.Spec.CAFileName == "" {
		fileNameCA = "--trusted-ca-file /certs-ca/sensu-ca.pem"
	} else {
		fileNameCA = fmt.Sprintf("--trusted-ca-file /certs-ca/%s", cr.Spec.CAFileName)
	}
	var subscriptions string
	if len(cr.Spec.Subscriptions) == 0 {
		subscriptions = "--subscriptions agent"
	} else {
		subscriptions = fmt.Sprintf("--subscriptions %s", strings.Join(cr.Spec.Subscriptions, ","))
	}
	deploy := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      cr.Name + "-agent",
			Namespace: cr.Namespace,
			Labels:    labels,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: labels,
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: labels,
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:    "agent",
							Image:   setImage,
							Command: []string{"sensu-agent", "start", "--deregister", "--statsd-disable", subscriptions, "--log-level", setLogLevel, fileNameCA},
							Ports: []corev1.ContainerPort{
								{
									ContainerPort: 3031,
									Name:          "agent",
								},
							},
							Env: []corev1.EnvVar{
								{
									Name:  "SENSU_BACKEND_URL",
									Value: sensuBackendWebsocket,
								},
							},
							VolumeMounts: []corev1.VolumeMount{
								{
									MountPath: "/certs-ca",
									Name:      secretCertificate,
									ReadOnly:  true,
								},
							},
						},
					},
					Volumes: []corev1.Volume{
						{
							Name: secretCertificate,
							VolumeSource: corev1.VolumeSource{
								Secret: &corev1.SecretVolumeSource{
									SecretName: secretCertificate,
								},
							},
						},
					},
				},
			},
		},
	}
	// Set SensuBackend instance as the owner and controller
	_ = controllerutil.SetControllerReference(cr, deploy, r.scheme)
	return deploy
}
