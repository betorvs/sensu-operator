package sensuasset

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

func TestSensuAssetController(t *testing.T) {
	logf.SetLogger(logf.ZapLogger(true))
	var (
		name      = "test"
		namespace = "default"
		assetname = "test-asset"
		asseturl  = "https://bonsai.sensu.io/"
		sha512    = "asadadasdsads"
	)
	// Create a MockRepository to test
	repo := utiltests.SensuRepositoryMock{}
	appcontext.Current.Add(appcontext.SensuRepository, repo)
	// Create a SensuAsset CR
	asset := &sensuv1alpha1.SensuAsset{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		Spec: sensuv1alpha1.SensuAssetSpec{
			Name:     assetname,
			AssetURL: asseturl,
			Sha512:   sha512,
		},
	}
	// Objects to track in the fake client.
	objs := []runtime.Object{
		asset,
	}

	// Register operator types with the runtime scheme.
	s := scheme.Scheme
	s.AddKnownTypes(sensuv1alpha1.SchemeGroupVersion, asset)
	// Create a fake client to mock API calls.
	cl := fake.NewFakeClient(objs...)
	// Create a ReconcileSensuAsset object with the scheme and fake client.
	r := &ReconcileSensuAsset{client: cl, scheme: s}

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

	// Check if SensuAsset has been created
	dep := &sensuv1alpha1.SensuAsset{}
	err = cl.Get(context.TODO(), req.NamespacedName, dep)
	if err != nil {
		t.Fatalf("Get Sensu Asset: (%v)", err)
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
