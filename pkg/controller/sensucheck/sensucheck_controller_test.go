package sensucheck

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
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	logf "sigs.k8s.io/controller-runtime/pkg/runtime/log"
)

func TestSensuCheckController(t *testing.T) {
	logf.SetLogger(logf.ZapLogger(true))
	var (
		name          = "test"
		namespace     = "default"
		checkname     = "test-check"
		command       = "ruby -v"
		interval      = 60
		publish       = true
		handlers      = []string{"default"}
		subscriptions = []string{"linux"}
		runtimeAssets = []string{"sensu-runtime-ruby"}
	)
	// Create a MockRepository to test
	repo := utiltests.SensuRepositoryMock{}
	appcontext.Current.Add(appcontext.SensuRepository, repo)
	// Create a SensuCheck CR
	check := &sensuv1alpha1.SensuCheck{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		Spec: sensuv1alpha1.SensuCheckSpec{
			Name:          checkname,
			Command:       command,
			Interval:      interval,
			Publish:       publish,
			Handlers:      handlers,
			Subscriptions: subscriptions,
			RuntimeAssets: runtimeAssets,
		},
	}
	// Objects to track in the fake client.
	objs := []runtime.Object{
		check,
	}

	// Register operator types with the runtime scheme.
	s := scheme.Scheme
	s.AddKnownTypes(sensuv1alpha1.SchemeGroupVersion, check)
	// Create a fake client to mock API calls.
	cl := fake.NewFakeClient(objs...)
	// Create a ReconcileSensuCheck object with the scheme and fake client.
	r := &ReconcileSensuCheck{client: cl, scheme: s}

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

	// Check if SensuCheck has been created
	dep := &sensuv1alpha1.SensuCheck{}
	err = cl.Get(context.TODO(), req.NamespacedName, dep)
	if err != nil {
		t.Fatalf("Get Sensu Check: (%v)", err)
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
