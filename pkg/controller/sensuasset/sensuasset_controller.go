package sensuasset

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

var log = logf.Log.WithName("controller_sensuasset")

/**
* USER ACTION REQUIRED: This is a scaffold file intended for the user to modify with their own Controller
* business logic.  Delete these comments after modifying this file.*
 */

// Add creates a new SensuAsset Controller and adds it to the Manager. The Manager will set fields on the Controller
// and Start it when the Manager is Started.
func Add(mgr manager.Manager) error {
	return add(mgr, newReconciler(mgr))
}

// newReconciler returns a new reconcile.Reconciler
func newReconciler(mgr manager.Manager) reconcile.Reconciler {
	return &ReconcileSensuAsset{client: mgr.GetClient(), scheme: mgr.GetScheme()}
}

// add adds a new Controller to mgr with r as the reconcile.Reconciler
func add(mgr manager.Manager, r reconcile.Reconciler) error {
	// Create a new controller
	c, err := controller.New("sensuasset-controller", mgr, controller.Options{Reconciler: r})
	if err != nil {
		return err
	}

	// Watch for changes to primary resource SensuAsset
	err = c.Watch(&source.Kind{Type: &sensuv1alpha1.SensuAsset{}}, &handler.EnqueueRequestForObject{})
	if err != nil {
		return err
	}

	// TODO(user): Modify this to be the types you create that are owned by the primary resource
	// Watch for changes to secondary resource Pods and requeue the owner SensuAsset
	// err = c.Watch(&source.Kind{Type: &corev1.Pod{}}, &handler.EnqueueRequestForOwner{
	// 	IsController: true,
	// 	OwnerType:    &sensuv1alpha1.SensuAsset{},
	// })
	// if err != nil {
	// 	return err
	// }

	return nil
}

// blank assignment to verify that ReconcileSensuAsset implements reconcile.Reconciler
var _ reconcile.Reconciler = &ReconcileSensuAsset{}

const sensuAssetFinalizer = "finalizer.sensu.k8s.sensu.io"

// ReconcileSensuAsset reconciles a SensuAsset object
type ReconcileSensuAsset struct {
	// This client, initialized using mgr.Client() above, is a split client
	// that reads objects from the cache and writes to the apiserver
	client client.Client
	scheme *runtime.Scheme
}

// Reconcile reads that state of the cluster for a SensuAsset object and makes changes based on the state read
// and what is in the SensuAsset.Spec
// TODO(user): Modify this Reconcile function to implement your Controller logic.  This example creates
// a Pod as an example
// Note:
// The Controller will requeue the Request to be processed again if the returned error is non-nil or
// Result.Requeue is true, otherwise upon completion it will remove the work from the queue.
func (r *ReconcileSensuAsset) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	reqLogger := log.WithValues("Request.Namespace", request.Namespace, "Request.Name", request.Name)
	reqLogger.Info("Reconciling SensuAsset")

	// Fetch the SensuAsset instance
	instance := &sensuv1alpha1.SensuAsset{}
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

	// check if sensu Backend api is health
	sensuURL := fmt.Sprintf("https://%s", instance.Spec.SensuBackendAPI)
	var sensuBackendToken string
	if usecase.SensuHealth(sensuURL) {
		sensuBackendToken = usecase.GetOperatorUserSensuAPIToken(sensuURL)
		sensubackendClusterID := usecase.GetClusterID(sensuURL, sensuBackendToken)
		if !usecase.CheckAssetExist(sensuURL, sensuBackendToken, instance.Spec.Namespace, instance.Spec.Name) || instance.Status.Status == "" {
			reqLogger.Info("Creating a new Asset in Sensu", "Sensu.Namespace", instance.Spec.Namespace)
			sensuAsset := &v2.Asset{
				URL:    instance.Spec.AssetURL,
				Sha512: instance.Spec.Sha512,
				ObjectMeta: v2.ObjectMeta{
					Name:      instance.Spec.Name,
					Namespace: instance.Spec.Namespace,
					Annotations: map[string]string{
						"CreatedBy": "sensu-operator",
					},
				},
			}
			err := usecase.AddAsset(sensuURL, sensuBackendToken, instance.Spec.Namespace, sensuAsset)
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

	// Check if the Sensu Asset instance is marked to be deleted, which is
	// indicated by the deletion timestamp being set.
	isSensuAssetMarkedToBeDeleted := instance.GetDeletionTimestamp() != nil
	if isSensuAssetMarkedToBeDeleted {
		if utils.Contains(instance.GetFinalizers(), sensuAssetFinalizer) {
			// Run finalization logic for sensuAssetFinalizer. If the
			// finalization logic fails, don't remove the finalizer so
			// that we can retry during the next reconciliation.
			if err := r.finalizeSensuAsset(reqLogger, instance); err != nil {
				return reconcile.Result{}, err
			}
			if usecase.CheckAssetExist(sensuURL, sensuBackendToken, instance.Spec.Namespace, instance.Spec.Name) {
				reqLogger.Info("Deleting a Asset in Sensu", "Sensu.Asset", instance.Spec.Name)
				err := usecase.DeleteAsset(sensuURL, sensuBackendToken, instance.Spec.Namespace, instance.Spec.Name)
				if err != nil {
					return reconcile.Result{}, err
				}
			}
			// Remove sensuAssetFinalizer. Once all finalizers have been
			// removed, the object will be deleted.
			instance.SetFinalizers(utils.Remove(instance.GetFinalizers(), sensuAssetFinalizer))
			err := r.client.Update(context.TODO(), instance)
			if err != nil {
				return reconcile.Result{}, err
			}
		}
		return reconcile.Result{}, nil
	}

	// Add finalizer for this CR
	if !utils.Contains(instance.GetFinalizers(), sensuAssetFinalizer) {
		if err := r.addFinalizer(reqLogger, instance); err != nil {
			// fmt.Printf("Add finalizers %s", instance.GetFinalizers())
			return reconcile.Result{}, err
		}
	}

	// reqLogger.Info("Skip reconcile: Pod already exists", "Pod.Namespace", found.Namespace, "Pod.Name", found.Name)
	reqLogger.Info("Skip reconcile: Sensu Asset already exists", "Sensu.Asset", instance.Spec.Name)
	return reconcile.Result{Requeue: true, RequeueAfter: config.RequeueTime}, nil
}

func (r *ReconcileSensuAsset) finalizeSensuAsset(reqLogger logr.Logger, cr *sensuv1alpha1.SensuAsset) error {
	// Log finalizer
	reqLogger.Info("Successfully finalized Sensu Asset")
	return nil
}

func (r *ReconcileSensuAsset) addFinalizer(reqLogger logr.Logger, cr *sensuv1alpha1.SensuAsset) error {
	reqLogger.Info("Adding Finalizer for the Sensu Asset")
	cr.SetFinalizers(append(cr.GetFinalizers(), sensuAssetFinalizer))
	// fmt.Printf("%s", cr.GetFinalizers())
	// Update CR
	err := r.client.Update(context.TODO(), cr)
	//r.client.Update(context.TODO(), cr)
	if err != nil {
		reqLogger.Error(err, "Failed to update Sensu Asset with finalizer")
		return err
	}
	return nil
}

func (r *ReconcileSensuAsset) statusPending(reqLogger logr.Logger, cr *sensuv1alpha1.SensuAsset) error {
	reqLogger.Info("Adding Status pending for the Sensu Asset")
	cr.Status.Status = ""
	cr.Status.OwnerID = ""

	// Update CR
	err := r.client.Status().Update(context.TODO(), cr)
	if err != nil {
		reqLogger.Error(err, "Failed to update Status pending for Sensu Asset")
		return err
	}
	return nil
}

func (r *ReconcileSensuAsset) statusCreated(reqLogger logr.Logger, cr *sensuv1alpha1.SensuAsset, sensubackendClusterID string) error {
	reqLogger.Info("Adding Status created for the Sensu Asset")
	cr.Status.Status = "created"
	cr.Status.OwnerID = sensubackendClusterID
	// fmt.Printf("INFO: %s, %s", cr.Status.Status, cr.Status.OwnerID)
	// Update CR
	err := r.client.Status().Update(context.TODO(), cr)
	if err != nil {
		reqLogger.Error(err, "Failed to update Status created for Sensu Asset")
		return err
	}
	return nil
}
