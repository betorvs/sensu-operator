package sensuhandler

import (
	"context"
	"fmt"

	sensuv1alpha1 "github.com/betorvs/sensu-operator/pkg/apis/sensu/v1alpha1"
	"github.com/betorvs/sensu-operator/pkg/config"
	"github.com/betorvs/sensu-operator/pkg/usecase"
	"github.com/betorvs/sensu-operator/pkg/utils"
	"github.com/go-logr/logr"
	v2 "github.com/sensu/sensu-go/api/core/v2"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

var log = logf.Log.WithName("controller_sensuhandler")

// Add creates a new SensuHandler Controller and adds it to the Manager. The Manager will set fields on the Controller
// and Start it when the Manager is Started.
func Add(mgr manager.Manager) error {
	return add(mgr, newReconciler(mgr))
}

// newReconciler returns a new reconcile.Reconciler
func newReconciler(mgr manager.Manager) reconcile.Reconciler {
	return &ReconcileSensuHandler{client: mgr.GetClient(), scheme: mgr.GetScheme()}
}

// add adds a new Controller to mgr with r as the reconcile.Reconciler
func add(mgr manager.Manager, r reconcile.Reconciler) error {
	// Create a new controller
	c, err := controller.New("sensuhandler-controller", mgr, controller.Options{Reconciler: r})
	if err != nil {
		return err
	}

	// Watch for changes to primary resource SensuHandler
	err = c.Watch(&source.Kind{Type: &sensuv1alpha1.SensuHandler{}}, &handler.EnqueueRequestForObject{})
	if err != nil {
		return err
	}

	return nil
}

// blank assignment to verify that ReconcileSensuHandler implements reconcile.Reconciler
var _ reconcile.Reconciler = &ReconcileSensuHandler{}

const sensuHandlerFinalizer = "finalizer.sensu.k8s.sensu.io"

// ReconcileSensuHandler reconciles a SensuHandler object
type ReconcileSensuHandler struct {
	// This client, initialized using mgr.Client() above, is a split client
	// that reads objects from the cache and writes to the apiserver
	client client.Client
	scheme *runtime.Scheme
}

// Reconcile reads that state of the cluster for a SensuHandler object and makes changes based on the state read
// and what is in the SensuHandler.Spec
// The Controller will requeue the Request to be processed again if the returned error is non-nil or
// Result.Requeue is true, otherwise upon completion it will remove the work from the queue.
func (r *ReconcileSensuHandler) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	reqLogger := log.WithValues("Request.Namespace", request.Namespace, "Request.Name", request.Name)
	reqLogger.Info("Reconciling SensuHandler")

	// Fetch the SensuHandler instance
	instance := &sensuv1alpha1.SensuHandler{}
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

	sensuURL := fmt.Sprintf("https://%s", instance.Spec.SensuBackendAPI)
	var sensuBackendToken string
	if usecase.SensuHealth(sensuURL) {
		sensuBackendToken = usecase.GetOperatorUserSensuAPIToken(sensuURL)
		sensubackendClusterID := usecase.GetClusterID(sensuURL, sensuBackendToken)

		if !usecase.CheckHandlerExist(sensuURL, sensuBackendToken, instance.Spec.Namespace, instance.Spec.Name) || instance.Status.Status == "" {
			reqLogger.Info("Creating a new Handler in Sensu", "Sensu.Namespace", instance.Spec.Namespace)
			sensuHandler := &v2.Handler{
				Type:     instance.Spec.Type,
				Command:  instance.Spec.Command,
				Handlers: instance.Spec.Handlers,
				ObjectMeta: v2.ObjectMeta{
					Name:      instance.Spec.Name,
					Namespace: instance.Spec.Namespace,
					Annotations: map[string]string{
						"CreatedBy": "sensu-operator",
					},
				},
			}
			if instance.Spec.Timeout != 0 {
				sensuHandler.Timeout = uint32(instance.Spec.Timeout)
			}
			if instance.Spec.Handlers != nil {
				sensuHandler.Handlers = instance.Spec.Handlers
			}
			if instance.Spec.RuntimeAssets != nil {
				sensuHandler.RuntimeAssets = instance.Spec.RuntimeAssets
			}
			if instance.Spec.Filters != nil {
				sensuHandler.Filters = instance.Spec.RuntimeAssets
			}
			if instance.Spec.Annotations != nil {
				for k, v := range instance.Spec.Annotations {
					sensuHandler.Annotations[k] = v
				}
			}
			err := usecase.AddHandler(sensuURL, sensuBackendToken, instance.Spec.Namespace, sensuHandler)
			if err != nil {
				if err := r.statusPending(reqLogger, instance); err != nil {
					return reconcile.Result{}, err
				}
				return reconcile.Result{}, err
			}
			if err := r.statusCreated(reqLogger, instance, sensubackendClusterID); err != nil {
				return reconcile.Result{}, err
			}
			return reconcile.Result{}, nil
		}
		if instance.Status.OwnerID != sensubackendClusterID {
			if err := r.statusPending(reqLogger, instance); err != nil {
				return reconcile.Result{}, err
			}
			return reconcile.Result{}, nil
		}

	}

	// Handler if the Sensu Handler instance is marked to be deleted, which is
	// indicated by the deletion timestamp being set.
	isSensuHandlerMarkedToBeDeleted := instance.GetDeletionTimestamp() != nil
	if isSensuHandlerMarkedToBeDeleted {
		if utils.Contains(instance.GetFinalizers(), sensuHandlerFinalizer) {
			// Run finalization logic for sensuHandlerFinalizer. If the
			// finalization logic fails, don't remove the finalizer so
			// that we can retry during the next reconciliation.
			if err := r.finalizeSensuHandler(reqLogger, instance); err != nil {
				return reconcile.Result{}, err
			}
			// fmt.Printf("INFO: Starting looking into sensu api for Deletion %s, %s", instance.Spec.SensuBackend, sensuBackend.Spec.SensuBackendURL)
			if usecase.CheckHandlerExist(sensuURL, sensuBackendToken, instance.Spec.Namespace, instance.Spec.Name) {
				reqLogger.Info("Deleting a Handler in Sensu", "Sensu.Handler", instance.Spec.Name)
				err := usecase.DeleteHandler(sensuURL, sensuBackendToken, instance.Spec.Namespace, instance.Spec.Name)
				if err != nil {
					return reconcile.Result{}, err
				}
			}
			// Remove sensuHandlerFinalizer. Once all finalizers have been
			// removed, the object will be deleted.
			instance.SetFinalizers(utils.Remove(instance.GetFinalizers(), sensuHandlerFinalizer))
			err := r.client.Update(context.TODO(), instance)
			if err != nil {
				return reconcile.Result{}, err
			}
		}
		return reconcile.Result{}, nil
	}

	// Add finalizer for this CR
	if !utils.Contains(instance.GetFinalizers(), sensuHandlerFinalizer) {
		if err := r.addFinalizer(reqLogger, instance); err != nil {
			return reconcile.Result{}, err
		}
	}

	// reqLogger.Info("Skip reconcile: Pod already exists", "Pod.Namespace", found.Namespace, "Pod.Name", found.Name)
	reqLogger.Info("Skip reconcile: Sensu Handler already exists", "Sensu.Handler", instance.Spec.Name)

	return reconcile.Result{Requeue: true, RequeueAfter: config.RequeueTime}, nil
}

func (r *ReconcileSensuHandler) finalizeSensuHandler(reqLogger logr.Logger, cr *sensuv1alpha1.SensuHandler) error {
	// Log finalizer
	reqLogger.Info("Successfully finalized Sensu Handler")
	return nil
}

func (r *ReconcileSensuHandler) addFinalizer(reqLogger logr.Logger, cr *sensuv1alpha1.SensuHandler) error {
	reqLogger.Info("Adding Finalizer for the Sensu Handler")
	cr.SetFinalizers(append(cr.GetFinalizers(), sensuHandlerFinalizer))

	// Update CR
	err := r.client.Update(context.TODO(), cr)
	if err != nil {
		reqLogger.Error(err, "Failed to update Sensu Handler with finalizer")
		return err
	}
	return nil
}

func (r *ReconcileSensuHandler) statusPending(reqLogger logr.Logger, cr *sensuv1alpha1.SensuHandler) error {
	reqLogger.Info("Adding Status pending for the Sensu Handler")
	cr.Status.Status = ""
	cr.Status.OwnerID = ""

	// Update CR
	err := r.client.Status().Update(context.TODO(), cr)
	if err != nil {
		reqLogger.Error(err, "Failed to update Status pending Sensu Handler")
		return err
	}
	return nil
}

func (r *ReconcileSensuHandler) statusCreated(reqLogger logr.Logger, cr *sensuv1alpha1.SensuHandler, sensubackendClusterID string) error {
	reqLogger.Info("Adding Status created for the Sensu Handler")
	cr.Status.Status = "created"
	cr.Status.OwnerID = sensubackendClusterID

	// Update CR
	err := r.client.Status().Update(context.TODO(), cr)
	if err != nil {
		reqLogger.Error(err, "Failed to update Status created for Sensu Handler")
		return err
	}
	return nil
}
