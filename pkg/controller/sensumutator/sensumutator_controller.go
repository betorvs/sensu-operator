package sensumutator

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

var log = logf.Log.WithName("controller_sensumutator")

/**
* USER ACTION REQUIRED: This is a scaffold file intended for the user to modify with their own Controller
* business logic.  Delete these comments after modifying this file.*
 */

// Add creates a new SensuMutator Controller and adds it to the Manager. The Manager will set fields on the Controller
// and Start it when the Manager is Started.
func Add(mgr manager.Manager) error {
	return add(mgr, newReconciler(mgr))
}

// newReconciler returns a new reconcile.Reconciler
func newReconciler(mgr manager.Manager) reconcile.Reconciler {
	return &ReconcileSensuMutator{client: mgr.GetClient(), scheme: mgr.GetScheme()}
}

// add adds a new Controller to mgr with r as the reconcile.Reconciler
func add(mgr manager.Manager, r reconcile.Reconciler) error {
	// Create a new controller
	c, err := controller.New("sensumutator-controller", mgr, controller.Options{Reconciler: r})
	if err != nil {
		return err
	}

	// Watch for changes to primary resource SensuMutator
	err = c.Watch(&source.Kind{Type: &sensuv1alpha1.SensuMutator{}}, &handler.EnqueueRequestForObject{})
	if err != nil {
		return err
	}

	return nil
}

// blank assignment to verify that ReconcileSensuMutator implements reconcile.Reconciler
var _ reconcile.Reconciler = &ReconcileSensuMutator{}

const sensuMutatorFinalizer = "finalizer.sensu.k8s.sensu.io"

// ReconcileSensuMutator reconciles a SensuMutator object
type ReconcileSensuMutator struct {
	// This client, initialized using mgr.Client() above, is a split client
	// that reads objects from the cache and writes to the apiserver
	client client.Client
	scheme *runtime.Scheme
}

// Reconcile reads that state of the cluster for a SensuMutator object and makes changes based on the state read
// and what is in the SensuMutator.Spec
// The Controller will requeue the Request to be processed again if the returned error is non-nil or
// Result.Requeue is true, otherwise upon completion it will remove the work from the queue.
func (r *ReconcileSensuMutator) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	reqLogger := log.WithValues("Request.Namespace", request.Namespace, "Request.Name", request.Name)
	reqLogger.Info("Reconciling SensuMutator")

	// Fetch the SensuMutator instance
	instance := &sensuv1alpha1.SensuMutator{}
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

		if !usecase.CheckMutatorExist(sensuURL, sensuBackendToken, instance.Spec.Namespace, instance.Spec.Name) || instance.Status.Status == "" {
			reqLogger.Info("Creating a new Mutator in Sensu", "Sensu.Namespace", instance.Spec.Namespace)
			sensuMutator := &v2.Mutator{
				Command: instance.Spec.Command,
				ObjectMeta: v2.ObjectMeta{
					Name:      instance.Spec.Name,
					Namespace: instance.Spec.Namespace,
					Annotations: map[string]string{
						"CreatedBy": "sensu-operator",
					},
				},
			}
			if instance.Spec.Timeout != 0 {
				sensuMutator.Timeout = uint32(instance.Spec.Timeout)
			}
			if instance.Spec.EnvVars != nil {
				sensuMutator.EnvVars = instance.Spec.EnvVars
			}
			if instance.Spec.RuntimeAssets != nil {
				sensuMutator.RuntimeAssets = instance.Spec.RuntimeAssets
			}
			err := usecase.AddMutator(sensuURL, sensuBackendToken, instance.Spec.Namespace, sensuMutator)
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

	// Check if the Sensu Mutator instance is marked to be deleted, which is
	// indicated by the deletion timestamp being set.
	isSensuMutatorMarkedToBeDeleted := instance.GetDeletionTimestamp() != nil
	if isSensuMutatorMarkedToBeDeleted {
		if utils.Contains(instance.GetFinalizers(), sensuMutatorFinalizer) {
			// Run finalization logic for sensuMutatorFinalizer. If the
			// finalization logic fails, don't remove the finalizer so
			// that we can retry during the next reconciliation.
			if err := r.finalizeSensuMutator(reqLogger, instance); err != nil {
				return reconcile.Result{}, err
			}
			// fmt.Printf("INFO: Starting looking into sensu api for Deletion %s, %s", instance.Spec.SensuBackend, sensuBackend.Spec.SensuBackendURL)
			if usecase.CheckMutatorExist(sensuURL, sensuBackendToken, instance.Spec.Namespace, instance.Spec.Name) {
				reqLogger.Info("Deleting a Mutator in Sensu", "Sensu.Mutator", instance.Spec.Name)
				err := usecase.DeleteMutator(sensuURL, sensuBackendToken, instance.Spec.Namespace, instance.Spec.Name)
				if err != nil {
					return reconcile.Result{}, err
				}
			}
			// Remove sensuMutatorFinalizer. Once all finalizers have been
			// removed, the object will be deleted.
			instance.SetFinalizers(utils.Remove(instance.GetFinalizers(), sensuMutatorFinalizer))
			err := r.client.Update(context.TODO(), instance)
			if err != nil {
				return reconcile.Result{}, err
			}
		}
		return reconcile.Result{}, nil
	}

	// Add finalizer for this CR
	if !utils.Contains(instance.GetFinalizers(), sensuMutatorFinalizer) {
		if err := r.addFinalizer(reqLogger, instance); err != nil {
			return reconcile.Result{}, err
		}
	}

	// reqLogger.Info("Skip reconcile: Pod already exists", "Pod.Namespace", found.Namespace, "Pod.Name", found.Name)
	reqLogger.Info("Skip reconcile: Sensu Mutator already exists", "Sensu.Mutator", instance.Spec.Name)

	return reconcile.Result{Requeue: true, RequeueAfter: config.RequeueTime}, nil
}

func (r *ReconcileSensuMutator) finalizeSensuMutator(reqLogger logr.Logger, cr *sensuv1alpha1.SensuMutator) error {
	// Log finalizer
	reqLogger.Info("Successfully finalized Sensu Mutator")
	return nil
}

func (r *ReconcileSensuMutator) addFinalizer(reqLogger logr.Logger, cr *sensuv1alpha1.SensuMutator) error {
	reqLogger.Info("Adding Finalizer for the Sensu Mutator")
	cr.SetFinalizers(append(cr.GetFinalizers(), sensuMutatorFinalizer))

	// Update CR
	err := r.client.Update(context.TODO(), cr)
	if err != nil {
		reqLogger.Error(err, "Failed to update Sensu Mutator with finalizer")
		return err
	}
	return nil
}

func (r *ReconcileSensuMutator) statusPending(reqLogger logr.Logger, cr *sensuv1alpha1.SensuMutator) error {
	reqLogger.Info("Adding Status pending for the Sensu Mutator")
	cr.Status.Status = ""
	cr.Status.OwnerID = ""

	// Update CR
	err := r.client.Status().Update(context.TODO(), cr)
	if err != nil {
		reqLogger.Error(err, "Failed to update Status pending Sensu Mutator")
		return err
	}
	return nil
}

func (r *ReconcileSensuMutator) statusCreated(reqLogger logr.Logger, cr *sensuv1alpha1.SensuMutator, sensubackendClusterID string) error {
	reqLogger.Info("Adding Status created for the Sensu Mutator")
	cr.Status.Status = "created"
	cr.Status.OwnerID = sensubackendClusterID

	// Update CR
	err := r.client.Status().Update(context.TODO(), cr)
	if err != nil {
		reqLogger.Error(err, "Failed to update Status created Sensu Mutator")
		return err
	}
	return nil
}
