package chat

import (
	// "github.com/phuslu/log"
	"fmt"
	"os"
	"path"
	"reflect"
	"strings"
	"zirc/helpers"

	"github.com/phuslu/log"
	"gopkg.in/yaml.v3"
)

type Config struct {
	Server struct {
		Server_version             string `yaml:"server_version"`
		Super_admin_password_file  string `yaml:"super_admin_password" validate:"secret"`
		Default_server_name        string `yaml:"default_server_name"`
		Dns_name                   string `yaml:"dns_name" default:"localhost"`
		Capabilities               string `yaml:"capabilities"`
		Supported_auth_types       string `yaml:"supported_auth_types"`
		User_modes                 string `yaml:"user_modes"`
		Channel_modes              string `yaml:"channel_modes"`
		Password_file              string `yaml:"password" validate:"secret"`
		Init_worker_count          uint64 `yaml:"init_worker_count"`
		Max_buffer_size            int    `yaml:"max_buffer_size"`
		Max_user_channels          string `yaml:"max_user_channels"`
		Server_manager_port        string `yaml:"server_manager_port"`
		Host                       string `yaml:"host"`
		Port                       string `yaml:"port"`
		Tls_port                   string `yaml:"tls_port"`
		Tls_cert_path              string `yaml:"tls_cert_path"`
		Tls_key_path               string `yaml:"tls_key_path"`
		Enable_tls                 bool   `yaml:"enable_tls"`
		Irc_verison                string `yaml:"irc_verison"`
		Use_user_pass              string `yaml:"use_user_pass"`
		Chat_server_log_level      string `yaml:"chat_server_log_level"`
		Ping_pong_timeout          string `yaml:"ping_pong_timeout"`
		Ping_pong_timeout_duration string `yaml:"ping_pong_timeout_duration"`
	}

	DB struct {
		Host                 string `yaml:"host" default:""`
		Port                 string `yaml:"port" default:""`
		User                 string `yaml:"user" default:""`
		Password_file        string `yaml:"password" validate:"secret"`
		Database             string `yaml:"database" default:""`
		Application_name     string `yaml:"application_name" default:""`
		Enable_tls           bool   `yaml:"enable_tls" default:""`
		Tls_cert_path        string `yaml:"tls_cert_path" default:""`
		Tls_key_path         string `yaml:"tls_key_path" default:""`
		Healthcheck_interval string `yaml:"healthcheck_interval" default:""`
	}

	S2S struct {
		Password_file   string `yaml:"password" validate:"secret"`
		Port            string `yaml:"port" default:"7000"`
		Tls_port        string `yaml:"tls_port" default:"7001"`
		Enable_tls      bool   `yaml:"enable_tls" default:"false"`
		Tls_cert_path   string `yaml:"tls_cert_path" default:"./.tls/s2s.crt"`
		Tls_key_path    string `yaml:"tls_key_path"  default:"./.tls/s2s.crt"`
		Max_buffer_size int    `yaml:"max_buffer_size" default:"262144"`
		Config          []struct {
			Host          string `yaml:"host,flow"`
			Port          string `yaml:"port,flow"`
			Password_file string `yaml:"password,flow" validate:"secret"`
		} `yaml:"config,flow,omitempty"`
	}
}

func parseTags(v interface{}) error {
	value := reflect.ValueOf(v)
	if value.Kind() != reflect.Ptr || value.Elem().Kind() != reflect.Struct {
		return fmt.Errorf("ParseTags requires a pointer to a struct")
	}

	value = value.Elem() // Dereference to get struct
	t := value.Type()    // Get struct type

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		tag := field.Tag.Get("validate")
		// defaultTag := field.Tag.Get("default")

		if tag == "secret" && strings.Contains(field.Name, "_file") {
			fieldValue := value.Field(i)

			// Check if the field is a settable string
			if fieldValue.Kind() != reflect.String || !fieldValue.CanSet() {
				return fmt.Errorf("%s field must be a settable string", field.Name)
			}

			strValue := fieldValue.String()
			log.Debug().Msgf("Reading secret from file %s", strValue)

			secret_dir := helpers.GetEnv("SECRET_BASE_PATH", "/run/secrets")
			secret_path := path.Join(secret_dir, strValue)
			_secret, err := os.ReadFile(secret_path)
			secret := strings.TrimSuffix(string(_secret), "\n")

			if err != nil || len(secret) == 0 {
				log.Warn().Msgf("Config tried reading secret %s from file, but failed... attempting to get env var", field.Name)

				// Fallback: Try getting from env variable
				i := strings.Index(field.Name, "_file")
				envVar := os.Getenv(strings.ToUpper(field.Name[:i]))
				if envVar == "" {
					log.Error().Msgf("Config couldn't find value for %s", field.Name)
					// return fmt.Errorf("%s field is required", field.Name)
				}

				// TODO
				// can we just fall back to the default?
				// Assign environment variable to field
				fieldValue.SetString(envVar)
			} else {
				// Assign file content to field
				fieldValue.SetString(secret)
			}
		}
	}
	return nil
}

// Reload reloads the the Config struct dynamically
// by calling load_config again
func (c *Config) Reload() {
	load_config()
}

var G_Config *Config

// load_config reads a config file from disk
// and populates the Config struct
func load_config() {
	config_path := helpers.GetEnv("CONFIG_FILE", "/chat_server/config.yaml")
	config_data, err := os.ReadFile(config_path)
	if err != nil {
		log.Error().Msgf("Error reading config file %v", err)
		panic(err)
	}

	err = yaml.Unmarshal([]byte(config_data), &G_Config)
	if err != nil {
		log.Error().Msgf("Couldn't unmarshal config: %v", err)
		panic(err)
	}

	err = parseTags(G_Config)
	if err != nil {
		log.Error().Msg(err.Error())
		panic(err)
	}

	log.Info().Msgf("Successfully loaded config: %+v", G_Config)
}

func init() {
	load_config()
}
