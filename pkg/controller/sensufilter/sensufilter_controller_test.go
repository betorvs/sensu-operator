package sensufilter

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

func TestSensuFilterController(t *testing.T) {
	logf.SetLogger(zap.Logger(true))
	var (
		name        = "test"
		namespace   = "default"
		filtername  = "test-asset"
		action      = "allow"
		expressions = []string{"event.entity.labels['environment'] == 'production'"}
	)
	// Create a MockRepository to test
	repo := utiltests.SensuRepositoryMock{}
	appcontext.Current.Add(appcontext.SensuRepository, repo)
	// Create a SensuFilter CR
	filter := &sensuv1alpha1.SensuFilter{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		Spec: sensuv1alpha1.SensuFilterSpec{
			Name:        filtername,
			Action:      action,
			Expressions: expressions,
		},
	}
	// Objects to track in the fake client.
	objs := []runtime.Object{
		filter,
	}

	// Register operator types with the runtime scheme.
	s := scheme.Scheme
	s.AddKnownTypes(sensuv1alpha1.SchemeGroupVersion, filter)
	// Create a fake client to mock API calls.
	cl := fake.NewFakeClient(objs...)
	// Create a ReconcileSensuFilter object with the scheme and fake client.
	r := &ReconcileSensuFilter{client: cl, scheme: s}

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

	// Check if SensuFilter has been created
	dep := &sensuv1alpha1.SensuFilter{}
	err = cl.Get(context.TODO(), req.NamespacedName, dep)
	if err != nil {
		t.Fatalf("Get Sensu Filter: (%v)", err)
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
