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
	DeleteResourceCalls               int
	CheckResourceExistCalls           int
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
	CheckResourceExistCalls++
	return true
}
func (repo SensuRepositoryMock) AddResource(sensuurl string, token string, resource string, namespace string, bodymarshal []byte) error {
	AddResourceCalls++
	return nil
}
func (repo SensuRepositoryMock) DeleteResource(sensuurl string, token string, namespace string, resource string, name string) error {
	DeleteResourceCalls++
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
	if CheckResourceExistCalls != expected {
		t.Fatalf("Invalid 11.1 TestCheckAssetExist %d", CheckResourceExistCalls)
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
	if AddResourceCalls != expected {
		t.Fatalf("Invalid 12.1 TestAddAsset %d", AddResourceCalls)
	}
}

func TestDeleteAsset(t *testing.T) {
	repo := SensuRepositoryMock{}
	appcontext.Current.Add(appcontext.SensuRepository, repo)
	_ = DeleteAsset(fakeSensuURL, "aaaaaqqqqqwwww", "default", "test")
	expected := 1
	if DeleteResourceCalls != expected {
		t.Fatalf("Invalid 13.1 TestDeleteAsset %d", DeleteResourceCalls)
	}
}
