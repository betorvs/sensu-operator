package sensuhandler

import (
	"context"
	"testing"

	sensuv1alpha1 "github.com/betorvs/sensu-operator/pkg/apis/sensu/v1alpha1"
	"github.com/betorvs/sensu-operator/pkg/appcontext"
	"github.com/betorvs/sensu-operator/utiltests"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes/scheme"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

func TestSensuHandlerController(t *testing.T) {
	logf.SetLogger(zap.Logger(true))
	var (
		name        = "test"
		namespace   = "default"
		handlerName = "test-handler"
		typeHandler = "pipe"
		command     = "echo"
		handlers    = []string{"test"}
	)
	// Create a MockRepository to test
	repo := utiltests.SensuRepositoryMock{}
	appcontext.Current.Add(appcontext.SensuRepository, repo)
	// Create a SensuHandler CR
	handler := &sensuv1alpha1.SensuHandler{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		Spec: sensuv1alpha1.SensuHandlerSpec{
			Name:     handlerName,
			Type:     typeHandler,
			Command:  command,
			Handlers: handlers,
		},
	}
	// Objects to track in the fake client.
	objs := []runtime.Object{
		handler,
	}

	// Register operator types with the runtime scheme.
	s := scheme.Scheme
	s.AddKnownTypes(sensuv1alpha1.SchemeGroupVersion, handler)
	// Create a fake client to mock API calls.
	cl := fake.NewFakeClient(objs...)
	// Create a ReconcileSensuHandler object with the scheme and fake client.
	r := &ReconcileSensuHandler{client: cl, scheme: s}

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
		t.Fatalf("First Reconcile Check: (%v)", err)
	}

	// Check if SensuHandler has been created
	dep := &sensuv1alpha1.SensuHandler{}
	err = cl.Get(context.TODO(), req.NamespacedName, dep)
	if err != nil {
		t.Fatalf("Get Sensu Handler: (%v)", err)
	}

	// last reconcile check
	res, err = r.Reconcile(req)
	if err != nil {
		t.Fatalf("Second Reconcile Check: (%v)", err)
	}
	// Check the result of reconciliation was requeue.
	if !res.Requeue {
		t.Error("Reconcile Requeue which is expected")
	}
}
