package utiltests

const (
	// FakeSensuURL const to be used in tests
	FakeSensuURL = "http://sensu-api.svc.cluster.local:8080"
)

// All variables are int and used to Mock SensuRepository
// When each method was called, this increment this int by 1
var (
	// SensuBackendHealthCalls int
	SensuBackendHealthCalls int
	// SensuVersionCalls int
	SensuVersionCalls int
	// GetSensuAPITokenCalls int
	GetSensuAPITokenCalls int
	// CreateOperatorUserGetTokenCalls int
	CreateOperatorUserGetTokenCalls int
	// GetOperatorUserSensuAPITokenCalls int
	GetOperatorUserSensuAPITokenCalls int
	// SensuTestTokenCalls int
	SensuTestTokenCalls int
	// GetClusterIDCalls int
	GetClusterIDCalls int
	// CheckMemberExistCalls int
	CheckMemberExistCalls int
	// AddNewMemberCalls int
	AddNewMemberCalls int
	// RemoveMemberCalls int
	RemoveMemberCalls int
	// AddResourceCalls int
	AddResourceCalls int
	// AddAssetsCalls int
	AddAssetsCalls int
	// AddNamespacesCalls int
	AddNamespacesCalls int
	// AddChecksCalls
	AddChecksCalls int
	// AddHandlersCalls int
	AddHandlersCalls int
	// AddFiltersCalls int
	AddFiltersCalls int
	// AddMutatorsCalls int
	AddMutatorsCalls int
	// CheckResourceExistCalls int
	CheckResourceExistCalls int
	// CheckAssetsExistCalls int
	CheckAssetsExistCalls int
	// CheckNamespacesExistCalls int
	CheckNamespacesExistCalls int
	// ChecksExistCalls int
	ChecksExistCalls int
	// CheckHandlersExistCalls int
	CheckHandlersExistCalls int
	// CheckFiltersExistCalls int
	CheckFiltersExistCalls int
	// CheckMutatorsExistCalls int
	CheckMutatorsExistCalls int
	// DeleteResourceCalls int
	DeleteResourceCalls int
	// DeleteAssetsCalls int
	DeleteAssetsCalls int
	// DeleteNamespacesCalls int
	DeleteNamespacesCalls int
	// DeleteChecksCalls int
	DeleteChecksCalls int
	// DeleteHandlersCalls int
	DeleteHandlersCalls int
	// DeleteFiltersCalls int
	DeleteFiltersCalls int
	// DeleteMutatorsCalls int
	DeleteMutatorsCalls int
)

// SensuRepositoryMock struct is used for Mock SensuRepository requests
type SensuRepositoryMock struct {
}

// SensuBackendHealth func
func (repo SensuRepositoryMock) SensuBackendHealth(sensuurl string) bool {
	SensuBackendHealthCalls++
	return true
}

// SensuVersion func
func (repo SensuRepositoryMock) SensuVersion(sensuurl string, version string) bool {
	SensuVersionCalls++
	return true
}

// GetSensuAPIToken func
func (repo SensuRepositoryMock) GetSensuAPIToken(sensuurl string) string {
	GetSensuAPITokenCalls++
	return ""
}

// CreateOperatorUserGetToken func
func (repo SensuRepositoryMock) CreateOperatorUserGetToken(sensuurl string) string {
	CreateOperatorUserGetTokenCalls++
	return ""
}

// GetOperatorUserSensuAPIToken func
func (repo SensuRepositoryMock) GetOperatorUserSensuAPIToken(sensuurl string) string {
	GetOperatorUserSensuAPITokenCalls++
	return ""
}

// SensuTestToken func
func (repo SensuRepositoryMock) SensuTestToken(sensuurl string, token string) bool {
	SensuTestTokenCalls++
	return true
}

// GetClusterID func
func (repo SensuRepositoryMock) GetClusterID(sensuurl string, token string) string {
	GetClusterIDCalls++
	return ""
}

// CheckMemberExist func
func (repo SensuRepositoryMock) CheckMemberExist(sensuurl string, token string, member string) bool {
	CheckMemberExistCalls++
	return true
}

// AddNewMember func
func (repo SensuRepositoryMock) AddNewMember(sensuurl string, token string, member string) error {
	AddNewMemberCalls++
	return nil
}

// RemoveMember func
func (repo SensuRepositoryMock) RemoveMember(sensuurl string, token string, member string) error {
	RemoveMemberCalls++
	return nil
}

// CheckResourceExist func
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

// AddResource func
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

// DeleteResource func
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
