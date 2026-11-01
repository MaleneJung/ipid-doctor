package ipid

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"os"
	"strings"
)

type IPIDConfig struct {
	ValidationKey string `json:"validationKey"`
}

func LoadIPIDConfig() IPIDConfig {
	config := IPIDConfig{
		ValidationKey: "remember-to-change",
	}
	bytes, err := os.ReadFile("ipid-doctor.json")
	if err != nil {
		return config
	}
	_ = json.Unmarshal(bytes, &config)
	return config
}

func GenerateIPID(ip string) string {
	config := LoadIPIDConfig()
	hasher := sha256.New()
	hasher.Write([]byte(ip))
	return config.ValidationKey + hex.EncodeToString(hasher.Sum(nil))
}

func ValidateIPID(ipid string) bool {
	config := LoadIPIDConfig()
	return strings.HasPrefix(ipid, config.ValidationKey)
}

func GetIPIDFromPath(path string) string {
	if strings.Count(path, "?") >= 1 {
		pathSplit := strings.SplitN(path, "?", 2)
		paramsSplit := strings.Split(pathSplit[1], "&")
		for i := range paramsSplit {
			cutIPID, foundIPID := strings.CutPrefix(paramsSplit[i], "ipid=")
			if foundIPID {
				return cutIPID
			}
		}
	}
	return ""
}

func GetIPIDFromRequest(r *http.Request) string {
	fromPath := GetIPIDFromPath(r.URL.Path)
	if ValidateIPID(fromPath) {
		return fromPath
	}
	return GenerateIPID(r.RemoteAddr)
}
