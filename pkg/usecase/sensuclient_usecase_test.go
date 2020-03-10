package usecase

import (
	"testing"

	"github.com/betorvs/sensu-operator/pkg/appcontext"
	v2 "github.com/sensu/sensu-go/api/core/v2"
)

const (
	fakeSensuURL = "http://sensu-api.svc.cluster.local:8080"
)

var (
	SensuBackendHealthCalls           int
	SensuVersionCalls                 int
	GetSensuAPITokenCalls             int
	CreateOperatorUserGetTokenCalls   int
	GetOperatorUserSensuAPITokenCalls int
	SensuTestTokenCalls               int
	GetClusterIDCalls                 int
	CheckMemberExistCalls             int
	AddNewMemberCalls                 int
	RemoveMemberCalls                 int
	AddResourceCalls                  int
	AddAssetsCalls                    int
	AddNamespacesCalls                int
	AddChecksCalls                    int
	AddHandlersCalls                  int
	AddFiltersCalls                   int
	AddMutatorsCalls                  int
	CheckResourceExistCalls           int
	CheckAssetsExistCalls             int
	CheckNamespacesExistCalls         int
	ChecksExistCalls                  int
	CheckHandlersExistCalls           int
	CheckFiltersExistCalls            int
	CheckMutatorsExistCalls           int
	DeleteResourceCalls               int
	DeleteAssetsCalls                 int
	DeleteNamespacesCalls             int
	DeleteChecksCalls                 int
	DeleteHandlersCalls               int
	DeleteFiltersCalls                int
	DeleteMutatorsCalls               int
)

type SensuRepositoryMock struct {
}

func (repo SensuRepositoryMock) SensuBackendHealth(sensuurl string) bool {
	SensuBackendHealthCalls++
	return true
}
func (repo SensuRepositoryMock) SensuVersion(sensuurl string, version string) bool {
	SensuVersionCalls++
	return true
}
func (repo SensuRepositoryMock) GetSensuAPIToken(sensuurl string) string {
	GetSensuAPITokenCalls++
	return ""
}
func (repo SensuRepositoryMock) CreateOperatorUserGetToken(sensuurl string) string {
	CreateOperatorUserGetTokenCalls++
	return ""
}
func (repo SensuRepositoryMock) GetOperatorUserSensuAPIToken(sensuurl string) string {
	GetOperatorUserSensuAPITokenCalls++
	return ""
}
func (repo SensuRepositoryMock) SensuTestToken(sensuurl string, token string) bool {
	SensuTestTokenCalls++
	return true
}
func (repo SensuRepositoryMock) GetClusterID(sensuurl string, token string) string {
	GetClusterIDCalls++
	return ""
}
func (repo SensuRepositoryMock) CheckMemberExist(sensuurl string, token string, member string) bool {
	CheckMemberExistCalls++
	return true
}
func (repo SensuRepositoryMock) AddNewMember(sensuurl string, token string, member string) error {
	AddNewMemberCalls++
	return nil
}
func (repo SensuRepositoryMock) RemoveMember(sensuurl string, token string, member string) error {
	RemoveMemberCalls++
	return nil
}
func (repo SensuRepositoryMock) CheckResourceExist(sensuurl string, token string, resource string, namespace string, name string) bool {
	switch resource {
	case "assets":
		CheckAssetsExistCalls++
	case "namespaces":
		CheckNamespacesExistCalls++
	case "checks":
		ChecksExistCalls++
	case "handlers":
		CheckHandlersExistCalls++
	case "filters":
		CheckFiltersExistCalls++
	case "mutators":
		CheckMutatorsExistCalls++
	default:
		CheckResourceExistCalls++
	}

	return true
}
func (repo SensuRepositoryMock) AddResource(sensuurl string, token string, resource string, namespace string, bodymarshal []byte) error {

	switch resource {
	case "assets":
		AddAssetsCalls++
	case "namespaces":
		AddNamespacesCalls++
	case "checks":
		AddChecksCalls++
	case "handlers":
		AddHandlersCalls++
	case "filters":
		AddFiltersCalls++
	case "mutators":
		AddMutatorsCalls++
	default:
		AddResourceCalls++
	}
	return nil
}
func (repo SensuRepositoryMock) DeleteResource(sensuurl string, token string, namespace string, resource string, name string) error {
	switch resource {
	case "assets":
		DeleteAssetsCalls++
	case "namespaces":
		DeleteNamespacesCalls++
	case "checks":
		DeleteChecksCalls++
	case "handlers":
		DeleteHandlersCalls++
	case "filters":
		DeleteFiltersCalls++
	case "mutators":
		DeleteMutatorsCalls++
	default:
		DeleteResourceCalls++
	}

	return nil
}

func TestSensuHealth(t *testing.T) {
	repo := SensuRepositoryMock{}
	appcontext.Current.Add(appcontext.SensuRepository, repo)
	_ = SensuHealth(fakeSensuURL)
	expected := 1
	if SensuBackendHealthCalls != expected {
		t.Fatalf("Invalid 1.1 TestSensuHealth %d", SensuBackendHealthCalls)
	}
}
func TestSensuVersion(t *testing.T) {
	repo := SensuRepositoryMock{}
	appcontext.Current.Add(appcontext.SensuRepository, repo)
	_ = SensuVersion(fakeSensuURL, "5.18.0")
	expected := 1
	if SensuVersionCalls != expected {
		t.Fatalf("Invalid 2.1 TestSensuVersion %d", SensuVersionCalls)
	}
}
func TestGetSensuAPIToken(t *testing.T) {
	repo := SensuRepositoryMock{}
	appcontext.Current.Add(appcontext.SensuRepository, repo)
	_ = GetSensuAPIToken(fakeSensuURL)
	expected := 1
	if GetSensuAPITokenCalls != expected {
		t.Fatalf("Invalid 3.1 TestGetSensuAPIToken %d", GetSensuAPITokenCalls)
	}
}
func TestCreateOperatorUserGetToken(t *testing.T) {
	repo := SensuRepositoryMock{}
	appcontext.Current.Add(appcontext.SensuRepository, repo)
	_ = CreateOperatorUserGetToken(fakeSensuURL)
	expected := 1
	if CreateOperatorUserGetTokenCalls != expected {
		t.Fatalf("Invalid 4.1 TestCreateOperatorUserGetToke %d", CreateOperatorUserGetTokenCalls)
	}
}
func TestGetOperatorUserSensuAPIToken(t *testing.T) {
	repo := SensuRepositoryMock{}
	appcontext.Current.Add(appcontext.SensuRepository, repo)
	_ = GetOperatorUserSensuAPIToken(fakeSensuURL)
	expected := 1
	if GetOperatorUserSensuAPITokenCalls != expected {
		t.Fatalf("Invalid 5.1 TestGetOperatorUserSensuAPIToken %d", GetOperatorUserSensuAPITokenCalls)
	}
}
func TestSensuTestToken(t *testing.T) {
	repo := SensuRepositoryMock{}
	appcontext.Current.Add(appcontext.SensuRepository, repo)
	_ = SensuTestToken(fakeSensuURL, "aaaaaqqqqqwwww")
	expected := 1
	if SensuTestTokenCalls != expected {
		t.Fatalf("Invalid 6.1 TestSensuTestToken %d", SensuTestTokenCalls)
	}
}
func TestGetClusterID(t *testing.T) {
	repo := SensuRepositoryMock{}
	appcontext.Current.Add(appcontext.SensuRepository, repo)
	_ = GetClusterID(fakeSensuURL, "aaaaaqqqqqwwww")
	expected := 1
	if GetClusterIDCalls != expected {
		t.Fatalf("Invalid 7.1 TestGetClusterID %d", GetClusterIDCalls)
	}
}
func TestCheckMemberExist(t *testing.T) {
	repo := SensuRepositoryMock{}
	appcontext.Current.Add(appcontext.SensuRepository, repo)
	_ = CheckMemberExist(fakeSensuURL, "aaaaaqqqqqwwww", "sensubackend-3")
	expected := 1
	if CheckMemberExistCalls != expected {
		t.Fatalf("Invalid 8.1 TestCheckMemberExist %d", CheckMemberExistCalls)
	}
}

func TestAddNewMember(t *testing.T) {
	repo := SensuRepositoryMock{}
	appcontext.Current.Add(appcontext.SensuRepository, repo)
	_ = AddNewMember(fakeSensuURL, "aaaaaqqqqqwwww", "sensubackend-3")
	expected := 1
	if AddNewMemberCalls != expected {
		t.Fatalf("Invalid 9.1 TestAddNewMember %d", AddNewMemberCalls)
	}
}

func TestRemoveMember(t *testing.T) {
	repo := SensuRepositoryMock{}
	appcontext.Current.Add(appcontext.SensuRepository, repo)
	_ = RemoveMember(fakeSensuURL, "aaaaaqqqqqwwww", "sensubackend-3")
	expected := 1
	if RemoveMemberCalls != expected {
		t.Fatalf("Invalid 10.1 TestRemoveMember %d", RemoveMemberCalls)
	}
}

func TestCheckAssetExist(t *testing.T) {
	repo := SensuRepositoryMock{}
	appcontext.Current.Add(appcontext.SensuRepository, repo)
	_ = CheckAssetExist(fakeSensuURL, "aaaaaqqqqqwwww", "default", "test")
	expected := 1
	if CheckAssetsExistCalls != expected {
		t.Fatalf("Invalid 11.1 TestCheckAssetExist %d", CheckAssetsExistCalls)
	}
}

func TestCheckExist(t *testing.T) {
	repo := SensuRepositoryMock{}
	appcontext.Current.Add(appcontext.SensuRepository, repo)
	_ = CheckExist(fakeSensuURL, "aaaaaqqqqqwwww", "default", "test")
	expected := 1
	if ChecksExistCalls != expected {
		t.Fatalf("Invalid 11.2 TestCheckExist %d", ChecksExistCalls)
	}
}

func TestCheckNamespaceExist(t *testing.T) {
	repo := SensuRepositoryMock{}
	appcontext.Current.Add(appcontext.SensuRepository, repo)
	_ = CheckNamespaceExist(fakeSensuURL, "aaaaaqqqqqwwww", "default")
	expected := 1
	if CheckNamespacesExistCalls != expected {
		t.Fatalf("Invalid 11.3 TestCheckNamespaceExist %d", CheckNamespacesExistCalls)
	}
}

func TestCheckHandlerExist(t *testing.T) {
	repo := SensuRepositoryMock{}
	appcontext.Current.Add(appcontext.SensuRepository, repo)
	_ = CheckHandlerExist(fakeSensuURL, "aaaaaqqqqqwwww", "default", "test")
	expected := 1
	if CheckHandlersExistCalls != expected {
		t.Fatalf("Invalid 11.4 TestHandlerAssetExist %d", CheckHandlersExistCalls)
	}
}

func TestCheckFilterExist(t *testing.T) {
	repo := SensuRepositoryMock{}
	appcontext.Current.Add(appcontext.SensuRepository, repo)
	_ = CheckFilterExist(fakeSensuURL, "aaaaaqqqqqwwww", "default", "test")
	expected := 1
	if CheckFiltersExistCalls != expected {
		t.Fatalf("Invalid 11.5 TestCheckFilterExist %d", CheckFiltersExistCalls)
	}
}

func TestCheckMutatorExist(t *testing.T) {
	repo := SensuRepositoryMock{}
	appcontext.Current.Add(appcontext.SensuRepository, repo)
	_ = CheckMutatorExist(fakeSensuURL, "aaaaaqqqqqwwww", "default", "test")
	expected := 1
	if CheckMutatorsExistCalls != expected {
		t.Fatalf("Invalid 11.6 TestCheckMutatorExist %d", CheckMutatorsExistCalls)
	}
}

func TestAddAsset(t *testing.T) {
	repo := SensuRepositoryMock{}
	appcontext.Current.Add(appcontext.SensuRepository, repo)
	sensuAsset := &v2.Asset{
		URL:    "https://bonsai.sensu.io",
		Sha512: "qqqaaazzzqqqaazzz",
		ObjectMeta: v2.ObjectMeta{
			Name:      "test",
			Namespace: "default",
		},
	}
	_ = AddAsset(fakeSensuURL, "aaaaaqqqqqwwww", "default", sensuAsset)
	expected := 1
	if AddAssetsCalls != expected {
		t.Fatalf("Invalid 12.1 TestAddAsset %d", AddAssetsCalls)
	}
}

func TestAddNamespace(t *testing.T) {
	repo := SensuRepositoryMock{}
	appcontext.Current.Add(appcontext.SensuRepository, repo)
	_ = AddNamespace(fakeSensuURL, "aaaaaqqqqqwwww", "test")
	expected := 1
	if AddNamespacesCalls != expected {
		t.Fatalf("Invalid 12.2 TestAddNamespace %d", AddNamespacesCalls)
	}
}

func TestAddCheck(t *testing.T) {
	repo := SensuRepositoryMock{}
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
	_ = AddCheck(fakeSensuURL, "aaaaaqqqqqwwww", "default", sensuCheck)
	expected := 1
	if AddChecksCalls != expected {
		t.Fatalf("Invalid 12.3 TestAddChecks %d", AddChecksCalls)
	}
}

func TestAddHandler(t *testing.T) {
	repo := SensuRepositoryMock{}
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
	_ = AddHandler(fakeSensuURL, "aaaaaqqqqqwwww", "default", sensuHandler)
	expected := 1
	if AddHandlersCalls != expected {
		t.Fatalf("Invalid 12.4 TestAddHandler %d", AddHandlersCalls)
	}
}

func TestAddFilter(t *testing.T) {
	repo := SensuRepositoryMock{}
	appcontext.Current.Add(appcontext.SensuRepository, repo)
	sensuFilter := &v2.EventFilter{
		Action:      "allow",
		Expressions: []string{"event.entity.labels['environment'] == 'production'"},
		ObjectMeta: v2.ObjectMeta{
			Name:      "test",
			Namespace: "default",
		},
	}
	_ = AddFilter(fakeSensuURL, "aaaaaqqqqqwwww", "default", sensuFilter)
	expected := 1
	if AddFiltersCalls != expected {
		t.Fatalf("Invalid 12.5 TestAddFilter %d", AddFiltersCalls)
	}
}

func TestAddMutator(t *testing.T) {
	repo := SensuRepositoryMock{}
	appcontext.Current.Add(appcontext.SensuRepository, repo)
	sensuMutator := &v2.Mutator{
		Command: "echo",
		ObjectMeta: v2.ObjectMeta{
			Name:      "test",
			Namespace: "default",
		},
	}
	_ = AddMutator(fakeSensuURL, "aaaaaqqqqqwwww", "default", sensuMutator)
	expected := 1
	if AddMutatorsCalls != expected {
		t.Fatalf("Invalid 12.6 TestAddMutator %d", AddMutatorsCalls)
	}
}

func TestDeleteAsset(t *testing.T) {
	repo := SensuRepositoryMock{}
	appcontext.Current.Add(appcontext.SensuRepository, repo)
	_ = DeleteAsset(fakeSensuURL, "aaaaaqqqqqwwww", "default", "test")
	expected := 1
	if DeleteAssetsCalls != expected {
		t.Fatalf("Invalid 13.1 TestDeleteAsset %d", DeleteAssetsCalls)
	}
}

func TestDeleteNamespace(t *testing.T) {
	repo := SensuRepositoryMock{}
	appcontext.Current.Add(appcontext.SensuRepository, repo)
	_ = DeleteNamespace(fakeSensuURL, "aaaaaqqqqqwwww", "default")
	expected := 1
	if DeleteNamespacesCalls != expected {
		t.Fatalf("Invalid 13.2 TestDeleteNamespace %d", DeleteNamespacesCalls)
	}
}

func TestDeleteCheck(t *testing.T) {
	repo := SensuRepositoryMock{}
	appcontext.Current.Add(appcontext.SensuRepository, repo)
	_ = DeleteCheck(fakeSensuURL, "aaaaaqqqqqwwww", "default", "test")
	expected := 1
	if DeleteChecksCalls != expected {
		t.Fatalf("Invalid 13.3 TestDeleteCheck %d", DeleteChecksCalls)
	}
}

func TestDeleteHandler(t *testing.T) {
	repo := SensuRepositoryMock{}
	appcontext.Current.Add(appcontext.SensuRepository, repo)
	_ = DeleteHandler(fakeSensuURL, "aaaaaqqqqqwwww", "default", "test")
	expected := 1
	if DeleteHandlersCalls != expected {
		t.Fatalf("Invalid 13.4 TestDeleteHandler %d", DeleteHandlersCalls)
	}
}

func TestDeleteFilter(t *testing.T) {
	repo := SensuRepositoryMock{}
	appcontext.Current.Add(appcontext.SensuRepository, repo)
	_ = DeleteFilter(fakeSensuURL, "aaaaaqqqqqwwww", "default", "test")
	expected := 1
	if DeleteFiltersCalls != expected {
		t.Fatalf("Invalid 13.5 TestDeleteFilter %d", DeleteFiltersCalls)
	}
}

func TestDeleteMutator(t *testing.T) {
	repo := SensuRepositoryMock{}
	appcontext.Current.Add(appcontext.SensuRepository, repo)
	_ = DeleteMutator(fakeSensuURL, "aaaaaqqqqqwwww", "default", "test")
	expected := 1
	if DeleteMutatorsCalls != expected {
		t.Fatalf("Invalid 13.6 TestDeleteMutator %d", DeleteMutatorsCalls)
	}
}
