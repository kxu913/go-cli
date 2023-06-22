package operator

import "os"

func GetEnvWithDefault(key string, defValue string) string {
	val, err := os.LookupEnv(key)
	if !err {
		return defValue
	}
	return val

}

var ConfigFile = GetEnvWithDefault("kubeconfig", "C:/Users/kevin/.kube/config")
