package usecase

import (
	"testing"

	"github.com/betorvs/sensu-operator/pkg/appcontext"
	"github.com/betorvs/sensu-operator/utiltests"
	v2 "github.com/sensu/sensu-go/api/core/v2"
)

func TestSensuHealth(t *testing.T) {
	repo := utiltests.SensuRepositoryMock{}
	appcontext.Current.Add(appcontext.SensuRepository, repo)
	_ = SensuHealth(utiltests.FakeSensuURL)
	expected := 1
	if utiltests.SensuBackendHealthCalls != expected {
		t.Fatalf("Invalid 1.1 TestSensuHealth %d", utiltests.SensuBackendHealthCalls)
	}
}
func TestSensuVersion(t *testing.T) {
	repo := utiltests.SensuRepositoryMock{}
	appcontext.Current.Add(appcontext.SensuRepository, repo)
	_ = SensuVersion(utiltests.FakeSensuURL, "5.18.0")
	expected := 1
	if utiltests.SensuVersionCalls != expected {
		t.Fatalf("Invalid 2.1 TestSensuVersion %d", utiltests.SensuVersionCalls)
	}
}
func TestGetSensuAPIToken(t *testing.T) {
	repo := utiltests.SensuRepositoryMock{}
	appcontext.Current.Add(appcontext.SensuRepository, repo)
	_ = GetSensuAPIToken(utiltests.FakeSensuURL)
	expected := 1
	if utiltests.GetSensuAPITokenCalls != expected {
		t.Fatalf("Invalid 3.1 TestGetSensuAPIToken %d", utiltests.GetSensuAPITokenCalls)
	}
}
func TestCreateOperatorUserGetToken(t *testing.T) {
	repo := utiltests.SensuRepositoryMock{}
	appcontext.Current.Add(appcontext.SensuRepository, repo)
	_ = CreateOperatorUserGetToken(utiltests.FakeSensuURL)
	expected := 1
	if utiltests.CreateOperatorUserGetTokenCalls != expected {
		t.Fatalf("Invalid 4.1 TestCreateOperatorUserGetToke %d", utiltests.CreateOperatorUserGetTokenCalls)
	}
}
func TestGetOperatorUserSensuAPIToken(t *testing.T) {
	repo := utiltests.SensuRepositoryMock{}
	appcontext.Current.Add(appcontext.SensuRepository, repo)
	_ = GetOperatorUserSensuAPIToken(utiltests.FakeSensuURL)
	expected := 1
	if utiltests.GetOperatorUserSensuAPITokenCalls != expected {
		t.Fatalf("Invalid 5.1 TestGetOperatorUserSensuAPIToken %d", utiltests.GetOperatorUserSensuAPITokenCalls)
	}
}
func TestSensuTestToken(t *testing.T) {
	repo := utiltests.SensuRepositoryMock{}
	appcontext.Current.Add(appcontext.SensuRepository, repo)
	_ = SensuTestToken(utiltests.FakeSensuURL, "aaaaaqqqqqwwww")
	expected := 1
	if utiltests.SensuTestTokenCalls != expected {
		t.Fatalf("Invalid 6.1 TestSensuTestToken %d", utiltests.SensuTestTokenCalls)
	}
}
func TestGetClusterID(t *testing.T) {
	repo := utiltests.SensuRepositoryMock{}
	appcontext.Current.Add(appcontext.SensuRepository, repo)
	_ = GetClusterID(utiltests.FakeSensuURL, "aaaaaqqqqqwwww")
	expected := 1
	if utiltests.GetClusterIDCalls != expected {
		t.Fatalf("Invalid 7.1 TestGetClusterID %d", utiltests.GetClusterIDCalls)
	}
}
func TestCheckMemberExist(t *testing.T) {
	repo := utiltests.SensuRepositoryMock{}
	appcontext.Current.Add(appcontext.SensuRepository, repo)
	_ = CheckMemberExist(utiltests.FakeSensuURL, "aaaaaqqqqqwwww", "sensubackend-3")
	expected := 1
	if utiltests.CheckMemberExistCalls != expected {
		t.Fatalf("Invalid 8.1 TestCheckMemberExist %d", utiltests.CheckMemberExistCalls)
	}
}

func TestAddNewMember(t *testing.T) {
	repo := utiltests.SensuRepositoryMock{}
	appcontext.Current.Add(appcontext.SensuRepository, repo)
	_ = AddNewMember(utiltests.FakeSensuURL, "aaaaaqqqqqwwww", "sensubackend-3")
	expected := 1
	if utiltests.AddNewMemberCalls != expected {
		t.Fatalf("Invalid 9.1 TestAddNewMember %d", utiltests.AddNewMemberCalls)
	}
}

func TestRemoveMember(t *testing.T) {
	repo := utiltests.SensuRepositoryMock{}
	appcontext.Current.Add(appcontext.SensuRepository, repo)
	_ = RemoveMember(utiltests.FakeSensuURL, "aaaaaqqqqqwwww", "sensubackend-3")
	expected := 1
	if utiltests.RemoveMemberCalls != expected {
		t.Fatalf("Invalid 10.1 TestRemoveMember %d", utiltests.RemoveMemberCalls)
	}
}

func TestCheckAssetExist(t *testing.T) {
	repo := utiltests.SensuRepositoryMock{}
	appcontext.Current.Add(appcontext.SensuRepository, repo)
	_ = CheckAssetExist(utiltests.FakeSensuURL, "aaaaaqqqqqwwww", "default", "test")
	expected := 1
	if utiltests.CheckAssetsExistCalls != expected {
		t.Fatalf("Invalid 11.1 TestCheckAssetExist %d", utiltests.CheckAssetsExistCalls)
	}
}

func TestCheckExist(t *testing.T) {
	repo := utiltests.SensuRepositoryMock{}
	appcontext.Current.Add(appcontext.SensuRepository, repo)
	_ = CheckExist(utiltests.FakeSensuURL, "aaaaaqqqqqwwww", "default", "test")
	expected := 1
	if utiltests.ChecksExistCalls != expected {
		t.Fatalf("Invalid 11.2 TestCheckExist %d", utiltests.ChecksExistCalls)
	}
}

func TestCheckNamespaceExist(t *testing.T) {
	repo := utiltests.SensuRepositoryMock{}
	appcontext.Current.Add(appcontext.SensuRepository, repo)
	_ = CheckNamespaceExist(utiltests.FakeSensuURL, "aaaaaqqqqqwwww", "default")
	expected := 1
	if utiltests.CheckNamespacesExistCalls != expected {
		t.Fatalf("Invalid 11.3 TestCheckNamespaceExist %d", utiltests.CheckNamespacesExistCalls)
	}
}

func TestCheckHandlerExist(t *testing.T) {
	repo := utiltests.SensuRepositoryMock{}
	appcontext.Current.Add(appcontext.SensuRepository, repo)
	_ = CheckHandlerExist(utiltests.FakeSensuURL, "aaaaaqqqqqwwww", "default", "test")
	expected := 1
	if utiltests.CheckHandlersExistCalls != expected {
		t.Fatalf("Invalid 11.4 TestHandlerAssetExist %d", utiltests.CheckHandlersExistCalls)
	}
}

func TestCheckFilterExist(t *testing.T) {
	repo := utiltests.SensuRepositoryMock{}
	appcontext.Current.Add(appcontext.SensuRepository, repo)
	_ = CheckFilterExist(utiltests.FakeSensuURL, "aaaaaqqqqqwwww", "default", "test")
	expected := 1
	if utiltests.CheckFiltersExistCalls != expected {
		t.Fatalf("Invalid 11.5 TestCheckFilterExist %d", utiltests.CheckFiltersExistCalls)
	}
}

func TestCheckMutatorExist(t *testing.T) {
	repo := utiltests.SensuRepositoryMock{}
	appcontext.Current.Add(appcontext.SensuRepository, repo)
	_ = CheckMutatorExist(utiltests.FakeSensuURL, "aaaaaqqqqqwwww", "default", "test")
	expected := 1
	if utiltests.CheckMutatorsExistCalls != expected {
		t.Fatalf("Invalid 11.6 TestCheckMutatorExist %d", utiltests.CheckMutatorsExistCalls)
	}
}

func TestAddAsset(t *testing.T) {
	repo := utiltests.SensuRepositoryMock{}
	appcontext.Current.Add(appcontext.SensuRepository, repo)
	sensuAsset := &v2.Asset{
		URL:    "https://bonsai.sensu.io",
		Sha512: "qqqaaazzzqqqaazzz",
		ObjectMeta: v2.ObjectMeta{
			Name:      "test",
			Namespace: "default",
		},
	}
	_ = AddAsset(utiltests.FakeSensuURL, "aaaaaqqqqqwwww", "default", sensuAsset)
	expected := 1
	if utiltests.AddAssetsCalls != expected {
		t.Fatalf("Invalid 12.1 TestAddAsset %d", utiltests.AddAssetsCalls)
	}
}

func TestAddNamespace(t *testing.T) {
	repo := utiltests.SensuRepositoryMock{}
	appcontext.Current.Add(appcontext.SensuRepository, repo)
	_ = AddNamespace(utiltests.FakeSensuURL, "aaaaaqqqqqwwww", "test")
	expected := 1
	if utiltests.AddNamespacesCalls != expected {
		t.Fatalf("Invalid 12.2 TestAddNamespace %d", utiltests.AddNamespacesCalls)
	}
}

func TestAddCheck(t *testing.T) {
	repo := utiltests.SensuRepositoryMock{}
	appcontext.Current.Add(appcontext.SensuRepository, repo)
	sensuCheck := &v2.Check{
		Subscriptions: []string{"test"},
		Command:       "echo",
		Interval:      uint32(60),
		Publish:       true,
		Handlers:      []string{"test"},
		ObjectMeta: v2.ObjectMeta{
			Name:      "test",
			Namespace: "default",
		},
	}
	_ = AddCheck(utiltests.FakeSensuURL, "aaaaaqqqqqwwww", "default", sensuCheck)
	expected := 1
	if utiltests.AddChecksCalls != expected {
		t.Fatalf("Invalid 12.3 TestAddChecks %d", utiltests.AddChecksCalls)
	}
}

func TestAddHandler(t *testing.T) {
	repo := utiltests.SensuRepositoryMock{}
	appcontext.Current.Add(appcontext.SensuRepository, repo)
	sensuHandler := &v2.Handler{
		Type:     "pipe",
		Command:  "echo",
		Handlers: []string{"test"},
		ObjectMeta: v2.ObjectMeta{
			Name:      "test",
			Namespace: "default",
		},
	}
	_ = AddHandler(utiltests.FakeSensuURL, "aaaaaqqqqqwwww", "default", sensuHandler)
	expected := 1
	if utiltests.AddHandlersCalls != expected {
		t.Fatalf("Invalid 12.4 TestAddHandler %d", utiltests.AddHandlersCalls)
	}
}

func TestAddFilter(t *testing.T) {
	repo := utiltests.SensuRepositoryMock{}
	appcontext.Current.Add(appcontext.SensuRepository, repo)
	sensuFilter := &v2.EventFilter{
		Action:      "allow",
		Expressions: []string{"event.entity.labels['environment'] == 'production'"},
		ObjectMeta: v2.ObjectMeta{
			Name:      "test",
			Namespace: "default",
		},
	}
	_ = AddFilter(utiltests.FakeSensuURL, "aaaaaqqqqqwwww", "default", sensuFilter)
	expected := 1
	if utiltests.AddFiltersCalls != expected {
		t.Fatalf("Invalid 12.5 TestAddFilter %d", utiltests.AddFiltersCalls)
	}
}

func TestAddMutator(t *testing.T) {
	repo := utiltests.SensuRepositoryMock{}
	appcontext.Current.Add(appcontext.SensuRepository, repo)
	sensuMutator := &v2.Mutator{
		Command: "echo",
		ObjectMeta: v2.ObjectMeta{
			Name:      "test",
			Namespace: "default",
		},
	}
	_ = AddMutator(utiltests.FakeSensuURL, "aaaaaqqqqqwwww", "default", sensuMutator)
	expected := 1
	if utiltests.AddMutatorsCalls != expected {
		t.Fatalf("Invalid 12.6 TestAddMutator %d", utiltests.AddMutatorsCalls)
	}
}

func TestDeleteAsset(t *testing.T) {
	repo := utiltests.SensuRepositoryMock{}
	appcontext.Current.Add(appcontext.SensuRepository, repo)
	_ = DeleteAsset(utiltests.FakeSensuURL, "aaaaaqqqqqwwww", "default", "test")
	expected := 1
	if utiltests.DeleteAssetsCalls != expected {
		t.Fatalf("Invalid 13.1 TestDeleteAsset %d", utiltests.DeleteAssetsCalls)
	}
}

func TestDeleteNamespace(t *testing.T) {
	repo := utiltests.SensuRepositoryMock{}
	appcontext.Current.Add(appcontext.SensuRepository, repo)
	_ = DeleteNamespace(utiltests.FakeSensuURL, "aaaaaqqqqqwwww", "default")
	expected := 1
	if utiltests.DeleteNamespacesCalls != expected {
		t.Fatalf("Invalid 13.2 TestDeleteNamespace %d", utiltests.DeleteNamespacesCalls)
	}
}

func TestDeleteCheck(t *testing.T) {
	repo := utiltests.SensuRepositoryMock{}
	appcontext.Current.Add(appcontext.SensuRepository, repo)
	_ = DeleteCheck(utiltests.FakeSensuURL, "aaaaaqqqqqwwww", "default", "test")
	expected := 1
	if utiltests.DeleteChecksCalls != expected {
		t.Fatalf("Invalid 13.3 TestDeleteCheck %d", utiltests.DeleteChecksCalls)
	}
}

func TestDeleteHandler(t *testing.T) {
	repo := utiltests.SensuRepositoryMock{}
	appcontext.Current.Add(appcontext.SensuRepository, repo)
	_ = DeleteHandler(utiltests.FakeSensuURL, "aaaaaqqqqqwwww", "default", "test")
	expected := 1
	if utiltests.DeleteHandlersCalls != expected {
		t.Fatalf("Invalid 13.4 TestDeleteHandler %d", utiltests.DeleteHandlersCalls)
	}
}

func TestDeleteFilter(t *testing.T) {
	repo := utiltests.SensuRepositoryMock{}
	appcontext.Current.Add(appcontext.SensuRepository, repo)
	_ = DeleteFilter(utiltests.FakeSensuURL, "aaaaaqqqqqwwww", "default", "test")
	expected := 1
	if utiltests.DeleteFiltersCalls != expected {
		t.Fatalf("Invalid 13.5 TestDeleteFilter %d", utiltests.DeleteFiltersCalls)
	}
}

func TestDeleteMutator(t *testing.T) {
	repo := utiltests.SensuRepositoryMock{}
	appcontext.Current.Add(appcontext.SensuRepository, repo)
	_ = DeleteMutator(utiltests.FakeSensuURL, "aaaaaqqqqqwwww", "default", "test")
	expected := 1
	if utiltests.DeleteMutatorsCalls != expected {
		t.Fatalf("Invalid 13.6 TestDeleteMutator %d", utiltests.DeleteMutatorsCalls)
	}
}
