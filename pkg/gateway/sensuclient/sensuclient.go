package sensuclient

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/betorvs/sensu-operator/pkg/appcontext"
	"github.com/betorvs/sensu-operator/pkg/config"
	"github.com/betorvs/sensu-operator/pkg/domain"
	v2 "github.com/sensu/sensu-go/api/core/v2"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
)

// Const
const (
	notFound = "Not Found"
)

// SensuRepository type
type SensuRepository struct {
	Client *http.Client
}

// var currentToken = domain.SensuToken{}
var currentToken = v2.Tokens{}

var logs = logf.Log.WithName("gateway_sensuclient")

func printLog(err error, function string, action string) {
	if config.GatewayDebug == "true" {
		reqLogger := logs.WithValues("Request.Function", function, "Request.Action", action)
		reqLogger.Info("Logging Requests")
		if err != nil {
			log.Printf("gateway_sensuclient %s, action %s, error %s", function, action, err)
			reqLogger.Error(err, "Error.Request.Function", function, "Request.Action", action)
		}
	}

}

// basicAuth func go to sensu api and get a new token
func (repo SensuRepository) basicAuth(basicURL string, user string, password string) (*v2.Tokens, error) {
	if currentToken.Access != "" &&
		currentToken.ExpiresAt > time.Now().Unix() {
		return &currentToken, nil
	}

	sensuURL := fmt.Sprintf("%s/auth", basicURL)
	req, err := http.NewRequest("GET", sensuURL, nil)
	if err != nil {
		printLog(err, "basicAuth", "newrequest")
		return &v2.Tokens{}, err
	}
	req.SetBasicAuth(user, password)
	resp, err := repo.Client.Do(req)
	if err != nil {
		printLog(err, "basicAuth", "Client.Do")
		return &v2.Tokens{}, err
	}
	bodyText, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		printLog(err, "basicAuth", "readbody")
		return &v2.Tokens{}, err
	}
	currentToken := new(v2.Tokens)
	_ = json.Unmarshal(bodyText, &currentToken)
	defer resp.Body.Close()
	return currentToken, nil
}

// sensuCreateAPIToken func
func (repo SensuRepository) sensuCreateAPIToken(url string, user string, password string) (string, error) {
	token, errToken := repo.basicAuth(url, user, password)
	if errToken != nil {
		action := fmt.Sprintf("basicAuth with user %s", user)
		printLog(errToken, "sensuCreateAPIToken", action)
		return "", errToken
	}
	var bearer = fmt.Sprintf("Bearer %s", token.Access)
	formPost := domain.Payload{
		Username: user,
	}
	bodymarshal, _ := json.Marshal(&formPost)
	sensuURL := fmt.Sprintf("%s/api/core/v2/apikeys", url)
	req, err := http.NewRequest("POST", sensuURL, bytes.NewBuffer(bodymarshal))
	if err != nil {
		printLog(err, "sensuCreateAPIToken", "newRequest")
		return "", err
	}
	req.Header.Add("Authorization", bearer)
	req.Header.Set("Content-Type", "application/json")
	resp, err := repo.Client.Do(req)
	if err != nil {
		printLog(errToken, "sensuCreateAPIToken", "Client.Do")
		return "", err
	}
	apitoken := resp.Header.Get("location")
	defer resp.Body.Close()
	return apitoken, nil
}

// sensuGet func
func (repo SensuRepository) sensuGet(token string, url string) ([]byte, error) {

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		printLog(err, "sensuGet", "NewRequest")
		return []byte{}, err
	}
	var bearer = "Key " + token
	req.Header.Add("Authorization", bearer)
	resp, err := repo.Client.Do(req)
	if err != nil {
		printLog(err, "sensuGet", "Client.Do")
		return []byte{}, err
	}
	bodyText, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		printLog(err, "sensuGet", "bodyText")
		return []byte{}, err
	}
	printLog(nil, "sensuGet", resp.Status)
	if resp.StatusCode <= 204 {
		return bodyText, nil
	}
	defer resp.Body.Close()
	return []byte{}, nil
}

// sensuPost func
func (repo SensuRepository) sensuPost(token string, sensuurl string, body []byte) ([]byte, error) {
	var req *http.Request
	var err error
	if body != nil {
		req, err = http.NewRequest("POST", sensuurl, bytes.NewBuffer(body))
		if err != nil {
			printLog(err, "sensuPost", "newRequest")
			return []byte{}, err
		}
	} else {
		req, err = http.NewRequest("POST", sensuurl, nil)
		if err != nil {
			printLog(err, "sensuPost", "newRequest")
			return []byte{}, err
		}
	}

	var bearer = fmt.Sprintf("Key %s", token)
	req.Header.Add("Authorization", bearer)
	req.Header.Set("Content-Type", "application/json")
	resp, err := repo.Client.Do(req)
	if err != nil {
		printLog(err, "sensuPost", "newRequest")
		return []byte{}, err
	}
	bodyText, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		printLog(err, "sensuPost", "newRequest")
		return []byte{}, err
	}
	printLog(nil, "sensuPost", resp.Status)
	if resp.StatusCode <= 204 {
		return bodyText, nil
	}
	defer resp.Body.Close()
	return []byte{}, nil
}

// sensuDelete func
func (repo SensuRepository) sensuDelete(token string, sensuurl string) error {
	req, err := http.NewRequest("DELETE", sensuurl, nil)
	if err != nil {
		printLog(err, "sensuDelete", "newRequest")
		return err
	}
	var bearer = fmt.Sprintf("Key %s", token)
	req.Header.Add("Authorization", bearer)
	req.Header.Set("Content-Type", "application/json")
	resp, err := repo.Client.Do(req)
	if err != nil {
		printLog(err, "sensuDelete", "Client.Do")
		return err
	}
	if resp.StatusCode <= 204 {
		return nil
	}
	defer resp.Body.Close()
	returnError := fmt.Errorf("Cannot delete")
	return returnError
}

// getMemberID func
func (repo SensuRepository) getMemberID(sensuurl string, token string, member string) string {
	sensuURL := fmt.Sprintf("%s/api/core/v2/cluster/members", sensuurl)
	body, err := repo.sensuGet(token, sensuURL)
	if err != nil {
		action := fmt.Sprintf("Failed to call Sensu API %s to get a members list", sensuurl)
		printLog(err, "getMemberID", action)
		return notFound
	}
	result := new(domain.ClusterMembers)
	_ = json.Unmarshal(body, &result)
	for _, n := range result.Members {
		if member == n.Name {
			hex := fmt.Sprintf("%x", n.ID)
			return hex
		}
	}
	return notFound
}

// SensuBackendHealth func
func (repo SensuRepository) SensuBackendHealth(sensuurl string) bool {
	sensuURL := fmt.Sprintf("%s/health", sensuurl)
	req, err := http.NewRequest("GET", sensuURL, nil)
	if err != nil {
		printLog(err, "SensuHealth", "newRequest")
		return false
	}
	resp, err := repo.Client.Do(req)
	if err != nil {
		printLog(err, "SensuHealth", "Client.Do")
		return false
	}
	defer resp.Body.Close()
	return true
}

// SensuVersion func
func (repo SensuRepository) SensuVersion(sensuurl string, version string) bool {
	sensuURL := fmt.Sprintf("%s/version", sensuurl)
	req, err := http.NewRequest("GET", sensuURL, nil)
	if err != nil {
		printLog(err, "SensuVersion", "newRequest")
		return false
	}
	resp, err := repo.Client.Do(req)
	if err != nil {
		printLog(err, "SensuVersion", "Client.Do")
		return false
	}
	bodyText, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		printLog(err, "SensuVersion", "Read Body")
		return false
	}
	result := new(v2.Version)
	_ = json.Unmarshal(bodyText, &result)
	if result.SensuBackend == version {
		return true
	}
	defer resp.Body.Close()
	return false
}

// GetSensuAPIToken func
func (repo SensuRepository) GetSensuAPIToken(sensuurl string) string {
	var token string
	token, err := repo.sensuCreateAPIToken(sensuurl, config.DefaultUser, config.DefaultPassword)
	if err != nil {
		action := fmt.Sprintf("Failed to call Sensu API %s to get a new token", sensuurl)
		printLog(err, "GetSensuAPIToken", action)
		return ""
	}
	valueToken := token[strings.LastIndex(token, "/")+1:]
	return valueToken
}

// sensuCreateOperatorUser func is used by sensu backend controller to create sensu-operator user inside sensu API
func (repo SensuRepository) sensuCreateOperatorUser(sensuurl string) error {
	token, err := repo.basicAuth(sensuurl, config.DefaultUser, config.DefaultPassword)
	if err != nil {
		action := fmt.Sprintf("Cannot authetication in Sensu API %s", sensuurl)
		printLog(err, "sensuCreateOperatorUser", action)
		return err
	}
	var bearer = fmt.Sprintf("Bearer %s", token.Access)
	userForm := v2.User{
		Username: config.OperatorSensuUser,
		Groups:   []string{"cluster-admins"},
		Disabled: false,
		Password: config.OperatorSensuPassword,
	}
	userBody, _ := json.Marshal(&userForm)
	sensuUserURL := fmt.Sprintf("%s/api/core/v2/users", sensuurl)
	reqUser, errUser := http.NewRequest("POST", sensuUserURL, bytes.NewBuffer(userBody))
	if errUser != nil {
		printLog(errUser, "sensuCreateOperatorUser", "newrequest")
		return errUser
	}
	reqUser.Header.Add("Authorization", bearer)
	reqUser.Header.Set("Content-Type", "application/json")
	resp, errClient := repo.Client.Do(reqUser)
	if errClient != nil {
		printLog(errClient, "sensuCreateOperatorUser", "Client.Do")
		return errClient
	}
	defer resp.Body.Close()
	return nil
}

//sensuURLGenerator func create final sensu url
func sensuURLGenerator(sensuurl string, resource string, namespace string, name string) string {

	basicURI := fmt.Sprintf("api/core/v2/namespaces/%s/%s/%s", namespace, resource, name)
	if resource == "namespaces" {
		basicURI = fmt.Sprintf("api/core/v2/namespaces/%s", name)
	}

	sensuBase := fmt.Sprintf("%s/%s", sensuurl, basicURI)
	return sensuBase
}

//sensuPostURLGenerator func create final sensu url
func sensuPostURLGenerator(sensuurl string, resource string, namespace string) string {
	basicURI := "api/core/v2/namespaces"
	if resource != "namespaces" {
		basicURI = fmt.Sprintf("api/core/v2/namespaces/%s/%s", namespace, resource)
	}

	sensuBase := fmt.Sprintf("%s/%s", sensuurl, basicURI)
	return sensuBase
}

// CreateOperatorUserGetToken func used by sensu backend controller
func (repo SensuRepository) CreateOperatorUserGetToken(sensuurl string) string {
	err := repo.sensuCreateOperatorUser(sensuurl)
	if err != nil {
		action := fmt.Sprintf("Failed to call Sensu API %s to create new operator user", sensuurl)
		printLog(err, "CreateOperatorUserGetToken", action)
		return ""
	}
	token, errToken := repo.sensuCreateAPIToken(sensuurl, config.OperatorSensuUser, config.OperatorSensuPassword)
	if errToken != nil {
		action := fmt.Sprintf("Failed to call Sensu API %s to get a new token", sensuurl)
		printLog(errToken, "CreateOperatorUserGetToken", action)
		return ""
	}
	valueToken := token[strings.LastIndex(token, "/")+1:]
	return valueToken
}

// GetOperatorUserSensuAPIToken func used for others controllers: asset, checks and go on.
func (repo SensuRepository) GetOperatorUserSensuAPIToken(sensuurl string) string {
	token, errToken := repo.sensuCreateAPIToken(sensuurl, config.OperatorSensuUser, config.OperatorSensuPassword)
	if errToken != nil {
		action := fmt.Sprintf("Failed to call Sensu API %s to get a new token", sensuurl)
		printLog(errToken, "GetOperatorUserSensuAPIToken", action)
		return ""
	}
	valueToken := token[strings.LastIndex(token, "/")+1:]
	return valueToken
}

// SensuTestToken func to test if a token still valid
func (repo SensuRepository) SensuTestToken(sensuurl string, token string) bool {
	accessToken, errToken := repo.basicAuth(sensuurl, config.OperatorSensuUser, config.OperatorSensuPassword)
	if errToken != nil {
		action := fmt.Sprintf("Cannot autheticated in Sensu API %s using %s", sensuurl, config.OperatorSensuUser)
		printLog(errToken, "SensuTestToken", action)
		return false
	}
	var bearer = fmt.Sprintf("Bearer %s", accessToken.Access)
	sensuURL := fmt.Sprintf("%s/api/core/v2/apikeys/%s", sensuurl, token)
	req, err := http.NewRequest("GET", sensuURL, nil)
	if err != nil {
		printLog(err, "SensuTestToken", "newRequest")
		return false
	}
	req.Header.Add("Authorization", bearer)
	resp, err := repo.Client.Do(req)
	if err != nil {
		printLog(err, "SensuTestToken", "Client.Do")
		return false
	}
	if resp.StatusCode == 404 {
		return false
	}
	defer resp.Body.Close()
	return true
}

// GetClusterID func
func (repo SensuRepository) GetClusterID(sensuurl string, token string) string {
	sensuURL := fmt.Sprintf("%s/api/core/v2/cluster/id", sensuurl)
	clusterID, err := repo.sensuGet(token, sensuURL)
	if err != nil {
		action := fmt.Sprintf("Failed to call Sensu API %s to get Cluster ID", sensuURL)
		printLog(err, "GetClusterID", action)
		return ""
	}
	result := string(clusterID)
	reg, err := regexp.Compile("[^a-zA-Z0-9]+")
	if err != nil {
		printLog(err, "GetClusterID", "regexp.Compile")
		return ""
	}
	processedString := reg.ReplaceAllString(result, "")
	return processedString
}

// CheckMemberExist func
func (repo SensuRepository) CheckMemberExist(sensuurl string, token string, member string) bool {
	sensuURL := fmt.Sprintf("%s/api/core/v2/cluster/members", sensuurl)
	body, err := repo.sensuGet(token, sensuURL)
	if err != nil {
		action := fmt.Sprintf("Failed to call Sensu API %s to get a members list", sensuURL)
		printLog(err, "CheckMemberExist", action)
		return false
	}
	result := new(domain.ClusterMembers)
	_ = json.Unmarshal(body, &result)
	for _, n := range result.Members {
		if member == n.Name {
			return true
		}
	}
	return false
}

// AddNewMember func
func (repo SensuRepository) AddNewMember(sensuurl string, token string, member string) error {
	memberURL := fmt.Sprintf("http://%s:2380", member)
	sensuURL := fmt.Sprintf("%s/api/core/v2/cluster/members?peer-addrs=%s", sensuurl, memberURL)
	_, err := repo.sensuPost(token, sensuURL, nil)
	if err != nil {
		action := fmt.Sprintf("Failed to add New Member to Sensu Cluster: %s", sensuURL)
		printLog(err, "AddNewMember", action)
		return err
	}
	return nil
}

// RemoveMember func
func (repo SensuRepository) RemoveMember(sensuurl string, token string, member string) error {
	hex := repo.getMemberID(sensuurl, token, member)
	if hex != notFound {
		sensuURL := fmt.Sprintf("%s/api/core/v2/cluster/members/%s", sensuurl, hex)
		err := repo.sensuDelete(token, sensuURL)
		if err != nil {
			action := fmt.Sprintf("Failed to delete Member from Sensu Cluster: %s", sensuURL)
			printLog(err, "RemoveMember", action)
			return err
		}
	}
	return nil
}

// CheckResourceExist func
func (repo SensuRepository) CheckResourceExist(sensuurl string, token string, resource string, namespace string, name string) bool {
	sensuURL := sensuURLGenerator(sensuurl, resource, namespace, name)
	_, err := repo.sensuGet(token, sensuURL)
	if err != nil {
		action := fmt.Sprintf("Failed to call Sensu API to get a %s list: %s", resource, sensuURL)
		printLog(err, "CheckResourceExist", action)
		return false
	}
	return true
}

// AddResource func
func (repo SensuRepository) AddResource(sensuurl string, token string, resource string, namespace string, bodymarshal []byte) error {
	sensuURL := sensuPostURLGenerator(sensuurl, resource, namespace)
	_, err := repo.sensuPost(token, sensuURL, bodymarshal)
	if err != nil {
		value := fmt.Sprintf("sensuPost %s", resource)
		printLog(err, "AddResource", value)
		return err
	}
	return nil
}

// DeleteResource func
func (repo SensuRepository) DeleteResource(sensuurl string, token string, namespace string, resource string, name string) error {
	sensuURL := sensuURLGenerator(sensuurl, resource, namespace, name)
	err := repo.sensuDelete(token, sensuURL)
	if err != nil {
		action := fmt.Sprintf("Failed to delete %s from Sensu: %s", resource, sensuURL)
		printLog(err, "DeleteResource", action)
		return err
	}
	return nil
}

// init func to start http client
func init() {
	if config.GetEnv("TESTRUN", "false") == "true" {
		return
	}
	client := http.Client{
		Timeout: time.Second * 10,
	}
	if config.CACertificate != "Absent" {
		caCert, err := ioutil.ReadFile(config.CACertificate)
		if err != nil {
			log.Fatal(err)
		}
		caCertPool := x509.NewCertPool()
		caCertPool.AppendCertsFromPEM(caCert)

		client = http.Client{
			Timeout: time.Second * 10,
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{
					RootCAs: caCertPool,
				},
			},
		}
	}
	appcontext.Current.Add(appcontext.SensuRepository, SensuRepository{Client: &client})
}
