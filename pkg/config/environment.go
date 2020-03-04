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

func getEnv(key, defaultValue string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return defaultValue
}

func init() {
	DefaultUser = getEnv("SENSU_BACKEND_CLUSTER_ADMIN_USERNAME", "admin")
	DefaultPassword = getEnv("SENSU_BACKEND_CLUSTER_ADMIN_PASSWORD", "P@ssw0rd!2GO")
	OperatorSensuUser = getEnv("OPERATOR_SENSU_USER", "sensu-operator")
	OperatorSensuPassword = getEnv("OPERATOR_SENSU_PASSWORD", "P@ssw0rd!2GO")
	CACertificate = getEnv("SENSU_CA_CERTIFICATE", "ABSENT")
	SensuImageTag = getEnv("SENSU_IMAGE_TAG", "5.18.0")
	GatewayDebug = getEnv("OPERATOR_GATEWAY_DEBUG", "false")
	timeout := getEnv("OPERATOR_REQUEUE_TIME", "300")
	tmpTimeout, err := strconv.Atoi(timeout)
	if err != nil {
		tmpTimeout = 300
	}
	RequeueTime = time.Second * time.Duration(tmpTimeout)
}
