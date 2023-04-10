package utils

import "os"

func GetEnvVariableOrDefault(identifier string, defaultAddress string) string {
	address := os.Getenv(identifier)
	if address == "" {
		address = defaultAddress
	}
	return address
}
