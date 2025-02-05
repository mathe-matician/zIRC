package chat

import (
	// "github.com/phuslu/log"

	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"
	"zirc/helpers"

	"github.com/phuslu/log"
	"gopkg.in/yaml.v3"
)

const (
	envVarPrefix = "IRC_"
)

type Config struct {
	Server struct {
		Server_version              string `yaml:"server_version" default:"v99.99.99+default"`
		Server_role                 string `yaml:"server_role" default:"leaf"`
		Super_admin_password_file   string `yaml:"super_admin_password_file" validate:"secret" default:"password"`
		Default_server_name         string `yaml:"default_server_name" default:"Z.IRC"`
		Dns_name                    string `yaml:"dns_name" default:"localhost"`
		Capabilities                string `yaml:"capabilities" default:"sasl account-registration"`
		Supported_auth_types        string `yaml:"supported_auth_types" default:"PLAIN,SCRAM-SHA-256,OAUTHBEARER,EXTERNAL"`
		User_modes                  string `yaml:"user_modes" default:"oiws"`
		Channel_modes               string `yaml:"channel_modes" default:"beiklmnopPst"`
		Password_file               string `yaml:"password_file" validate:"secret" default:""`
		Init_worker_count           uint64 `yaml:"init_worker_count" default:"3"`
		Max_buffer_size             int    `yaml:"max_buffer_size" default:"8192"`
		Max_user_channels           int    `yaml:"max_user_channels" default:"20"`
		Host                        string `yaml:"host" default:"0.0.0.0"`
		Port                        string `yaml:"port" default:"6667"`
		Tls_port                    string `yaml:"tls_port" default:"6677"`
		Tls_cert_path               string `yaml:"tls_cert_path" default:"./.tls/server.crt"`
		Tls_key_path                string `yaml:"tls_key_path" default:"./.tls/server.key"`
		Enable_tls                  bool   `yaml:"enable_tls" default:"false"`
		Irc_verison                 string `yaml:"irc_verison" default:"302"`
		Supported_protocol_versions string `yaml:"supported_protocol_versions" default:"302"`
		Use_user_pass               string `yaml:"use_user_pass" default:"true"`      // When true, will not ask the user to authenticate again to use chat, but will log them in via their password automatically
		Chat_server_log_level       string `yaml:"chat_server_log_level" default:"3"` // info https://pkg.go.dev/github.com/phuslu/log@v1.0.110#Level
		Ping_pong_timeout           int    `yaml:"ping_pong_timeout" default:"2"`
		Ping_pong_timeout_duration  string `yaml:"ping_pong_timeout_duration" default:"minute"`
	}

	DB struct {
		Enabled              bool   `yaml:"enabled" default:"false"`
		Host                 string `yaml:"host" default:"zirc_db"`
		Port                 string `yaml:"port" default:"5432"`
		User                 string `yaml:"user" default:"postgres"`
		Password_file        string `yaml:"password_file" validate:"secret" default:""`
		Database             string `yaml:"database" default:"postgres"`
		Application_name     string `yaml:"application_name" default:"zirc_server"`
		Enable_tls           bool   `yaml:"enable_tls" default:"false"`
		Tls_cert_path        string `yaml:"tls_cert_path" default:"./.tls/db.crt"`
		Tls_key_path         string `yaml:"tls_key_path" default:"./.tls/db.key"`
		Healthcheck_interval int    `yaml:"healthcheck_interval" default:"2"`
	}

	S2S struct {
		Password_file   string `yaml:"password_file" validate:"secret"`
		Port            string `yaml:"port" default:"7000"`
		Tls_port        string `yaml:"tls_port" default:"7001"`
		Enable_tls      bool   `yaml:"enable_tls" default:"false"`
		Tls_cert_path   string `yaml:"tls_cert_path" default:"./.tls/s2s.crt"`
		Tls_key_path    string `yaml:"tls_key_path"  default:"./.tls/s2s.crt"`
		Max_buffer_size int    `yaml:"max_buffer_size" default:"262144"`
		// Config          []struct {
		// 	Host          string `yaml:"host,flow"`
		// 	Port          string `yaml:"port,flow"`
		// 	Password_file string `yaml:"password,flow" validate:"secret"`
		// } `yaml:"config,flow,omitempty"`
		// shouldn't have Config here
		// if we can dynamically reload the config
		// this is because init servers may not be good to reload mid execution
		// if they have been removed or etc
	}
}

// convertToEnvName converts a field name in config struct to an env var
// using known global prefix + parentStructName + structFieldName
// If it is a secret, the `_file` will be removed from the field name
//
// E.g. if:
//
//	DB struct {
//		Password_file: bool
//		Max_connections: string
//	}
//
// Then if those values aren't set, it will attempt to find env vars:
// IRC_DB_PASSWORD
// IRC_DB_MAX_CONNECTIONS
func convertToEnvName(fieldName string, parentStructName string, isSecret bool) string {
	if isSecret {
		i := strings.Index(fieldName, "_file")
		fieldName = fieldName[:i]
	}

	envVarName := strings.ToUpper(envVarPrefix + parentStructName + "_" + fieldName)

	return envVarName
}

func parseTags(v interface{}) error {
	value := reflect.ValueOf(v)
	if value.Kind() != reflect.Ptr || value.Elem().Kind() != reflect.Struct {
		return fmt.Errorf("ParseTags requires a pointer to a struct")
	}

	value = value.Elem() // Dereference to get struct
	t := value.Type()    // Get struct type

	for j := 0; j < t.NumField(); j++ {
		parentStructName := t.Field(j).Name // Get parent struct field name (DB)
		parentStructValue := value.Field(j) // Get parent struct value (DB struct)

		if parentStructValue.Kind() == reflect.Struct {
			parentType := parentStructValue.Type()

			for i := 0; i < parentType.NumField(); i++ {
				field := parentType.Field(i)
				tag := field.Tag.Get("validate")
				defaultTagValue := field.Tag.Get("default")

				fieldValue := parentStructValue.Field(i)

				if tag == "secret" && strings.Contains(field.Name, "_file") {
					// Check if the field is a settable string
					if fieldValue.Kind() != reflect.String || !fieldValue.CanSet() {
						return fmt.Errorf("%s field must be a settable string", field.Name)
					}

					strValue := fieldValue.Interface().(string)
					log.Info().Msgf("Searching for secret at path: %s", strValue)
					_secret, err := os.ReadFile(strValue)
					secret := strings.TrimSuffix(string(_secret), "\n")

					if err != nil || len(secret) == 0 {
						envVarName := convertToEnvName(field.Name, parentStructName, true)
						// envVarName := strings.ToUpper(envVarPrefix + parentStructName + "_" + field.Name[:i])
						log.Warn().Msgf("Config tried reading secret from file %s, but failed... attempting to get env var %s", strValue, envVarName)
						envVar := os.Getenv(envVarName)

						if envVar == "" {
							log.Error().Msgf("Config couldnt find value for %s", field.Name)
							// TODO
							// check if it is required and fail if it is
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
				} else {
					if fieldValue.String() == "" {
						envVarName := convertToEnvName(field.Name, parentStructName, false)
						log.Warn().Msgf("Config tried reading %s from config, but failed... attempting to get env var", envVarName)
						envVar := os.Getenv(envVarName)

						if envVar == "" {
							log.Error().Msgf("Config couldn't find value for %s", field.Name)
							// TODO
							// check if it is required and fail if it is
							// return fmt.Errorf("%s field is required", field.Name)
						}
					}
				}

				if value.Field(j).String() == "" && defaultTagValue != "" && fieldValue.CanSet() {
					switch fieldValue.Kind() {
					case reflect.String:
						fieldValue.SetString(defaultTagValue)
					case reflect.Bool:
						boolValue, err := strconv.ParseBool(defaultTagValue)
						if err != nil {
							log.Error().Msgf("CONFIG: Error parsing bool default value: %s", err.Error())
							boolValue = false
						}
						fieldValue.SetBool(boolValue)
					case reflect.Int:
						intValue, err := strconv.ParseInt(defaultTagValue, 10, 64)
						if err != nil {
							log.Error().Msgf("CONFIG: Error parsing int64 default value: %s", err.Error())
							// since we don't have any way of knowing what a sane default
							// for any int typed field is
							// we just panic.
							panic(err)
						}
						fieldValue.SetInt(intValue)
						// case reflect.Struct:
					}
				}
			}
		}
		// TODO
		// handle non struct config values
		// i.e. top level keys that aren't structs (don't _have_ to, can just force it to be all nested)
	}

	// for i := 0; i < t.NumField(); i++ {
	// 	field := t.Field(i)
	// 	tag := field.Tag.Get("validate")
	// 	defaultTagValue := field.Tag.Get("default")

	// 	fieldValue := value.Field(i)

	// 	if tag == "secret" && strings.Contains(field.Name, "_file") {
	// 		// Check if the field is a settable string
	// 		if fieldValue.Kind() != reflect.String || !fieldValue.CanSet() {
	// 			return fmt.Errorf("%s field must be a settable string", field.Name)
	// 		}

	// 		strValue := fieldValue.String()
	// 		log.Debug().Msgf("Reading secret from file %s", strValue)

	// 		secret_dir := helpers.GetEnv("SECRET_BASE_PATH", "/run/secrets")
	// 		secret_path := path.Join(secret_dir, strValue)
	// 		log.Info().Msgf("Searching for secret at path: %s", secret_path)
	// 		_secret, err := os.ReadFile(secret_path)
	// 		secret := strings.TrimSuffix(string(_secret), "\n")

	// 		if err != nil || len(secret) == 0 {
	// 			log.Warn().Msgf("Config tried reading secret %s from file, but failed... attempting to get env var", field.Name)

	// 			// Fallback: Try getting from env variable
	// 			i := strings.Index(field.Name, "_file")
	// 			envVar := os.Getenv(strings.ToUpper(field.Name[:i]))
	// 			if envVar == "" {
	// 				log.Error().Msgf("Config couldn't find value for %s", field.Name)
	// 				// return fmt.Errorf("%s field is required", field.Name)
	// 			}

	// 			// TODO
	// 			// can we just fall back to the default?
	// 			// Assign environment variable to field
	// 			fieldValue.SetString(envVar)
	// 		} else {
	// 			// Assign file content to field
	// 			fieldValue.SetString(secret)
	// 		}
	// 	}

	// 	if value.Field(i).String() == "" && defaultTagValue != "" && fieldValue.CanSet() {
	// 		switch fieldValue.Kind() {
	// 		case reflect.String:
	// 			fieldValue.SetString(defaultTagValue)
	// 		case reflect.Bool:
	// 			boolValue, err := strconv.ParseBool(defaultTagValue)
	// 			if err != nil {
	// 				log.Error().Msgf("CONFIG: Error parsing bool default value: %s", err.Error())
	// 				boolValue = false
	// 			}
	// 			fieldValue.SetBool(boolValue)
	// 		case reflect.Int:
	// 			intValue, err := strconv.ParseInt(defaultTagValue, 10, 64)
	// 			if err != nil {
	// 				log.Error().Msgf("CONFIG: Error parsing int64 default value: %s", err.Error())
	// 				// since we don't have any way of knowing what a sane default
	// 				// for any int typed field is
	// 				// we just panic.
	// 				panic(err)
	// 			}
	// 			fieldValue.SetInt(intValue)
	// 			// case reflect.Struct:
	// 		}
	// 	}
	// }
	return nil
}

// Reload reloads the the Config struct dynamically
// by calling load_config again
func (c Config) Reload() {
	load_config()
}

var G_Config Config

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

	err = parseTags(&G_Config)
	if err != nil {
		log.Error().Msg(err.Error())
		panic(err)
	}

	// TODO
	// redact sensitive configurations
	// overriding String() doesn't seem to work for whatever reason
	log.Info().Msgf("Successfully loaded config: %+v", G_Config)
}

func init() {
	load_config()
}
