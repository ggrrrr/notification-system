package config

import (
	"fmt"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

var (
	v viper.Viper
)

func init() {
	v = *viper.New()
	v.AddConfigPath("./configs")
	v.SetConfigType("yaml")
	// v.SetConfigFile("app.yaml")
	v.SetConfigFile("configs/app.yaml")
	v.AutomaticEnv()
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	readFile()
}

func readFile() {
	err := v.ReadInConfig()
	if err != nil {
		logrus.Errorf("unale to read file: %v", err)
	}
}
func GetInt(prefix, k string, val int) int {
	out := v.GetInt(key(prefix, k))
	if out == 0 {
		return val
	}
	return out
}

// Get config string with default value
func GetStringOrDefault(prefix, k, val string) string {
	v := strings.TrimSpace(v.GetString(key(prefix, k)))
	if v == "" {
		return val
	}
	return v
}

// Get config string with default value
func GetBool(prefix, k string) bool {
	return v.GetBool(key(prefix, k))
}

// Get config string with default value
func GetString(prefix, k string) string {
	key := key(prefix, k)
	val := v.GetString(key)
	return val
}

// Get time.Duration with default value
func GetDuration(prefix, k string, val time.Duration) time.Duration {
	str := v.GetString(key(prefix, k))
	out, err := time.ParseDuration(str)
	if err != nil {
		return val
	}
	return out
}

func key(prefix, key string) string {
	if prefix == "" {
		return key
	}
	return fmt.Sprintf("%v.%v", prefix, key)
}
