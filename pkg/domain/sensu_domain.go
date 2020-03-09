package domain

import (
	"github.com/betorvs/sensu-operator/pkg/appcontext"
)

// ClusterMembers struct
type ClusterMembers struct {
	Header  Header    `json:"header"`
	Members []Members `json:"members"`
}

// Header of ClusterMembers
type Header struct {
	ClusterID int64 `json:"cluster_id"`
	MemberID  int64 `json:"member_id"`
	RaftTerm  int   `json:"raft_term"`
}

// Members of ClusterMembers
type Members struct {
	ID         int64    `json:"ID"`
	Name       string   `json:"name"`
	PeerURLs   []string `json:"peerURLs"`
	ClientURLs []string `json:"clientURLs"`
}

// Payload struct
type Payload struct {
	Username string `json:"username"`
}

// SensuRepository interface
type SensuRepository interface {
	appcontext.Component
	SensuBackendHealth(sensuurl string) bool
	SensuVersion(sensuurl string, version string) bool
	GetSensuAPIToken(sensuurl string) string
	CreateOperatorUserGetToken(sensuurl string) string
	GetOperatorUserSensuAPIToken(sensuurl string) string
	SensuTestToken(sensuurl string, token string) bool
	GetClusterID(sensuurl string, token string) string
	CheckMemberExist(sensuurl string, token string, member string) bool
	AddNewMember(sensuurl string, token string, member string) error
	RemoveMember(sensuurl string, token string, member string) error
	CheckResourceExist(sensuurl string, token string, resource string, namespace string, name string) bool
	AddResource(sensuurl string, token string, resource string, namespace string, bodymarshal []byte) error
	DeleteResource(sensuurl string, token string, namespace string, resource string, name string) error
}

// GetSensuRepository func return SensuRepository interface
func GetSensuRepository() SensuRepository {
	return appcontext.Current.Get(appcontext.SensuRepository).(SensuRepository)
}
