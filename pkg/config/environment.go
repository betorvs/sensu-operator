package config

import (
	"os"
	"strconv"
	"time"
)

var (
	// DefaultUser to access sensu api
	DefaultUser string
	// DefaultPassword to access sensu api
	DefaultPassword string
	// OperatorSensuUser string
	OperatorSensuUser string
	// OperatorSensuPassword string
	OperatorSensuPassword string
	// CACertificate string
	CACertificate string
	// GatewayDebug string
	GatewayDebug string
	// SensuImageTag string
	SensuImageTag string
	// RequeueTime time.Duration
	RequeueTime time.Duration
)

// GetEnv func return a default value if dont find a environment variable
func GetEnv(key, defaultValue string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return defaultValue
}

func init() {
	DefaultUser = GetEnv("SENSU_BACKEND_CLUSTER_ADMIN_USERNAME", "admin")
	DefaultPassword = GetEnv("SENSU_BACKEND_CLUSTER_ADMIN_PASSWORD", "P@ssw0rd!2GO")
	OperatorSensuUser = GetEnv("OPERATOR_SENSU_USER", "sensu-operator")
	OperatorSensuPassword = GetEnv("OPERATOR_SENSU_PASSWORD", "P@ssw0rd!2GO")
	CACertificate = GetEnv("SENSU_CA_CERTIFICATE", "ABSENT")
	SensuImageTag = GetEnv("SENSU_IMAGE_TAG", "5.18.0")
	GatewayDebug = GetEnv("OPERATOR_GATEWAY_DEBUG", "false")
	timeout := GetEnv("OPERATOR_REQUEUE_TIME", "300")
	tmpTimeout, err := strconv.Atoi(timeout)
	if err != nil {
		tmpTimeout = 300
	}
	RequeueTime = time.Second * time.Duration(tmpTimeout)
}
