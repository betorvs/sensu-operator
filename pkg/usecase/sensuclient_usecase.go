package usecase

import (
	"encoding/json"

	"github.com/betorvs/sensu-operator/pkg/domain"
	v2 "github.com/sensu/sensu-go/api/core/v2"
)

var sensuRepository = domain.GetSensuRepository()

// SensuHealth func used to test if Sensu API are up and running
func SensuHealth(sensuurl string) bool {
	return sensuRepository.SensuHealth(sensuurl)
}

// SensuVersion func used to test if Sensu API are up and running
func SensuVersion(sensuurl string, version string) bool {
	return sensuRepository.SensuVersion(sensuurl, version)
}

// GetSensuAPIToken func is used by sensu backend to get admin token using Default credentials
func GetSensuAPIToken(sensuurl string) string {
	return sensuRepository.GetSensuAPIToken(sensuurl)
}

// CreateOperatorUserGetToken func used by sensu backend controller to create operator user and return his token
func CreateOperatorUserGetToken(sensuurl string) string {
	return sensuRepository.CreateOperatorUserGetToken(sensuurl)
}

// GetOperatorUserSensuAPIToken func used by others controllers to get sensu operator user token
func GetOperatorUserSensuAPIToken(sensuurl string) string {
	return sensuRepository.GetOperatorUserSensuAPIToken(sensuurl)
}

// SensuTestToken func is used by sensu backend controller to test admin and operator token
func SensuTestToken(sensuurl string, token string) bool {
	return sensuRepository.SensuTestToken(sensuurl, token)
}

// GetClusterID func returns Cluster ID from sensu api
func GetClusterID(sensuurl string, token string) string {
	return sensuRepository.GetClusterID(sensuurl, token)
}

// CheckMemberExist func check if that member exist in that sensu cluster
func CheckMemberExist(sensuurl string, token string, member string) bool {
	return sensuRepository.CheckMemberExist(sensuurl, token, member)
}

// CheckNamespaceExist func check if that namespace already exist in sensu api
func CheckNamespaceExist(sensuurl string, token string, namespace string) bool {
	// return sensuRepository.CheckNamespaceExist(sensuurl, token, namespace)
	return sensuRepository.CheckResourceExist(sensuurl, token, "namespaces", namespace, namespace)
}

// CheckAssetExist func check if that asset already exist in sensu api
func CheckAssetExist(sensuurl string, token string, namespace string, asset string) bool {
	return sensuRepository.CheckResourceExist(sensuurl, token, "assets", namespace, asset)
}

// CheckExist func check if that check already exist in sensu api
func CheckExist(sensuurl string, token string, namespace string, check string) bool {
	return sensuRepository.CheckResourceExist(sensuurl, token, "checks", namespace, check)
}

// CheckHandlerExist func check if that handler already exist in sensu api
func CheckHandlerExist(sensuurl string, token string, namespace string, handler string) bool {
	return sensuRepository.CheckResourceExist(sensuurl, token, "handlers", namespace, handler)
}

// CheckFilterExist func check if that filter already exist in sensu api
func CheckFilterExist(sensuurl string, token string, namespace string, filter string) bool {
	return sensuRepository.CheckResourceExist(sensuurl, token, "filters", namespace, filter)
}

// CheckMutatorExist func check if that mutator already exist in sensu api
func CheckMutatorExist(sensuurl string, token string, namespace string, mutator string) bool {
	return sensuRepository.CheckResourceExist(sensuurl, token, "mutators", namespace, mutator)
}

// AddNewMember func add a new member to a cluster. Used by sensu backend controller.
func AddNewMember(sensuurl string, token string, member string) error {
	return sensuRepository.AddNewMember(sensuurl, token, member)
}

// AddNamespace func add a new namespace in sensu api
func AddNamespace(sensuurl string, token string, namespace string) error {
	metadata := v2.Namespace{
		Name: namespace,
	}
	bodymarshal, err := json.Marshal(&metadata)
	if err != nil {
		return err
	}
	return sensuRepository.AddResource(sensuurl, token, "namespaces", namespace, bodymarshal)
}

// AddAsset func add a new asset in sensu api
func AddAsset(sensuurl string, token string, namespace string, asset *v2.Asset) error {
	bodymarshal, err := json.Marshal(&asset)
	if err != nil {
		return err
	}
	return sensuRepository.AddResource(sensuurl, token, "assets", namespace, bodymarshal)
}

// AddCheck func add a new check in sensu api
func AddCheck(sensuurl string, token string, namespace string, check *v2.Check) error {
	bodymarshal, err := json.Marshal(&check)
	if err != nil {
		return err
	}
	return sensuRepository.AddResource(sensuurl, token, "checks", namespace, bodymarshal)
}

// AddHandler func add a new handler in sensu api
func AddHandler(sensuurl string, token string, namespace string, handler *v2.Handler) error {
	bodymarshal, err := json.Marshal(&handler)
	if err != nil {
		return err
	}
	return sensuRepository.AddResource(sensuurl, token, "handlers", namespace, bodymarshal)
}

// AddFilter func add a new filter in sensu api
func AddFilter(sensuurl string, token string, namespace string, filter *v2.EventFilter) error {
	bodymarshal, err := json.Marshal(&filter)
	if err != nil {
		return err
	}
	return sensuRepository.AddResource(sensuurl, token, "filters", namespace, bodymarshal)
}

// AddMutator func add a new mutator in sensu api
func AddMutator(sensuurl string, token string, namespace string, mutator *v2.Mutator) error {
	bodymarshal, err := json.Marshal(&mutator)
	if err != nil {
		return err
	}
	return sensuRepository.AddResource(sensuurl, token, "mutators", namespace, bodymarshal)
}

// RemoveMember func remove a member from cluster to downsize an sensu cluster. Used by sensu backend controller
func RemoveMember(sensuurl string, token string, member string) error {
	return sensuRepository.RemoveMember(sensuurl, token, member)
}

// DeleteNamespace func delete a namespace from sensu api
func DeleteNamespace(sensuurl string, token string, namespace string) error {
	return sensuRepository.DeleteResource(sensuurl, token, namespace, "namespaces", namespace)
}

// DeleteAsset func delete a asset from sensu api
func DeleteAsset(sensuurl string, token string, namespace string, asset string) error {
	return sensuRepository.DeleteResource(sensuurl, token, namespace, "assets", asset)
}

// DeleteCheck func delete a check from sensu api
func DeleteCheck(sensuurl string, token string, namespace string, check string) error {
	return sensuRepository.DeleteResource(sensuurl, token, namespace, "checks", check)
}

// DeleteHandler func delete a handler from sensu api
func DeleteHandler(sensuurl string, token string, namespace string, handler string) error {
	return sensuRepository.DeleteResource(sensuurl, token, namespace, "handlers", handler)
}

// DeleteFilter func delete a filter from sensu api
func DeleteFilter(sensuurl string, token string, namespace string, filter string) error {
	return sensuRepository.DeleteResource(sensuurl, token, namespace, "filters", filter)
}

// DeleteMutator func delete a mutator from sensu api
func DeleteMutator(sensuurl string, token string, namespace string, mutator string) error {
	return sensuRepository.DeleteResource(sensuurl, token, namespace, "mutators", mutator)
}
