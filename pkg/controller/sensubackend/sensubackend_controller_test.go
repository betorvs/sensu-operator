package sensubackend

import (
	"context"
	"testing"

	sensuv1alpha1 "github.com/betorvs/sensu-operator/pkg/apis/sensu/v1alpha1"
	"github.com/betorvs/sensu-operator/pkg/appcontext"
	"github.com/betorvs/sensu-operator/utiltests"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes/scheme"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

func TestSensuBackendController(t *testing.T) {
	logf.SetLogger(zap.Logger(true))
	var (
		name            = "test"
		namespace       = "default"
		replicas  int32 = 1
	)
	// Create a MockRepository to test
	repo := utiltests.SensuRepositoryMock{}
	appcontext.Current.Add(appcontext.SensuRepository, repo)
	// Create backend Object
	backend := &sensuv1alpha1.SensuBackend{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		Spec: sensuv1alpha1.SensuBackendSpec{
			Replicas:        replicas,
			SensuBackendURL: utiltests.FakeSensuURL,
		},
	}
	// Objects to track in the fake client.
	objs := []runtime.Object{
		backend,
	}
	// Register operator types with the runtime scheme.
	s := scheme.Scheme
	s.AddKnownTypes(sensuv1alpha1.SchemeGroupVersion, backend)
	// Create a fake client to mock API calls.
	cl := fake.NewFakeClient(objs...)
	// Create a ReconcileSensuBackend object with the scheme and fake client.
	r := &ReconcileSensuBackend{client: cl, scheme: s}

	// Mock request to simulate Reconcile() being called on an event for a
	// watched resource .
	req := reconcile.Request{
		NamespacedName: types.NamespacedName{
			Name:      name,
			Namespace: namespace,
		},
	}

	res, err := r.Reconcile(req)
	if err != nil {
		t.Fatalf("reconcile: (%v)", err)
	}

	// Check if Deployment has been created and has the correct size.
	dep := &sensuv1alpha1.SensuBackend{}
	err = cl.Get(context.TODO(), req.NamespacedName, dep)
	if err != nil {
		t.Fatalf("Get Sensu Backend: (%v)", err)
	}
	dsize := dep.Spec.Replicas
	if dsize != replicas {
		t.Errorf("dep size (%d) is not the expected size (%d)", dsize, replicas)
	}

	// last reconcile check
	res, err = r.Reconcile(req)
	if err != nil {
		t.Fatalf("reconcile: (%v)", err)
	}
	// Check the result of reconciliation to make sure it has the desired state.
	if res.Requeue {
		t.Error("reconcile requeue which is not expected")
	}

	// Check if Headless Service has been created.
	reqSvc := reconcile.Request{
		NamespacedName: types.NamespacedName{
			Name:      "sensu",
			Namespace: namespace,
		},
	}

	ser := &corev1.Service{}
	err = cl.Get(context.TODO(), reqSvc.NamespacedName, ser)
	if err != nil {
		t.Fatalf("Get Headless Service: (%v)", err)
	}

	// Check if API Service has been created.
	reqSvcAPI := reconcile.Request{
		NamespacedName: types.NamespacedName{
			Name:      "sensu-api",
			Namespace: namespace,
		},
	}

	svcAPI := &corev1.Service{}
	err = cl.Get(context.TODO(), reqSvcAPI.NamespacedName, svcAPI)
	if err != nil {
		t.Fatalf("Get API Service: (%v)", err)
	}

	// Check if Sensu Secret has been created.
	reqSecret := reconcile.Request{
		NamespacedName: types.NamespacedName{
			Name:      "sensu-secret",
			Namespace: namespace,
		},
	}

	secret := &corev1.Secret{}
	err = cl.Get(context.TODO(), reqSecret.NamespacedName, secret)
	if err != nil {
		t.Fatalf("Get Sensu Secret: (%v)", err)
	}
}
