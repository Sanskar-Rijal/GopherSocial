package env

import "os"

func GetString(key, fallback string) string {
	value, exists :=os.LookupEnv(key)
	if !exists {
		return fallback
	}
	return value 
}