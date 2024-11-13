package config

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/fs"
	"log/slog"
	"os"
	"os/user"
	"path/filepath"
	"reflect"
	"strings"
)

type Config struct {
	DBUrl           string `json:"db_url"`
	CurrentUsername string `json:"current_username"`
	userConfPath    string
}

var DefaultPaths []string = []string{
	"/etc/gatorconfig.json",
	"~/.config/gatorconfig.json",
	"~/.gatorconfig.json",
}

func ReadConfig(paths []string) (Config, error) {
	conf := Config{}
	for _, path := range paths {
		path, err := expandPath(path)
		if err != nil {
			return Config{}, fmt.Errorf("Invalid path: " + path)
		}
		data, err := os.ReadFile(path)
		if err != nil {
			slog.Debug("Skipping file due to readfile err", "err", err)
			continue
		}
		decoder := json.NewDecoder(bytes.NewBuffer(data))
		var readConf Config
		err = decoder.Decode(&readConf)
		if err != nil {
			slog.Debug("Skipping file due to decode err", "err", err)
			continue
		}
		conf = combineConfigs(conf, readConf)
		conf.userConfPath = path
	}
	return conf, nil
}

func (conf Config) write() error {
	const perms fs.FileMode = 644
	if conf.userConfPath == "" {
		conf.userConfPath = DefaultPaths[len(DefaultPaths)-1]
	}
	dir := filepath.Dir(conf.userConfPath)
	jsonData, err := json.Marshal(conf)
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}
	slog.Debug("Writing config", "config", jsonData)

	err = os.MkdirAll(dir, perms)
	if err != nil {
		return fmt.Errorf("failed to instantiate dir: `%s`, err: %w", dir, err)
	}

	err = os.WriteFile(conf.userConfPath, jsonData, perms)
	if err != nil {
		return fmt.Errorf("failed to write file: `%s`, err: %w", conf.userConfPath, err)
	}

	return nil

}

func (conf *Config) SetUser(username string) error {
	conf.CurrentUsername = username
	return conf.write()
}

func combineConfigs(configA, configB Config) Config {
	outConfig := Config{}
	va := reflect.ValueOf(configA)
	vb := reflect.ValueOf(configB)
	vout := reflect.Indirect(reflect.ValueOf(&outConfig))
	slog.Debug("Mering configs", "configA", configA, "configB", configB)
	for i := 0; i < va.NumField(); i++ {
		fieldName := va.Type().Field(i).Name

		fieldA := va.Field(i)
		fieldB := vb.Field(i)
		fieldOut := vout.Field(i)

		slog.Debug("", "fieldName", fieldName, "fieldA", fieldA, "fieldB", fieldB)

		switch {
		case !fieldB.IsZero():
			slog.Debug("FieldB is non-default - merging it", "fieldB", fieldB, "fieldName", fieldName)
			fieldOut.Set(fieldB)
		case !fieldA.IsZero():
			slog.Debug("fieldB is default, fieldA is not - merging fieldA",
				"fieldA", fieldA, "fieldB", fieldB, "fieldName", fieldName)
			fieldOut.Set(fieldA)
		default:
			slog.Debug("fieldA & fieldB are default - ignoring", "fieldA", fieldA, "fieldB", fieldB, "fieldName", fieldName)
		}
	}
	slog.Debug("merged output", "outConfig", outConfig)
	return outConfig
}

func expandPath(p string) (string, error) {
	// Check if the path starts with ~
	if strings.HasPrefix(p, "~") {
		// Get the current user's home directory
		usr, err := user.Current()
		if err != nil {
			return "", err
		}

		// If path is just ~, return the home directory
		if p == "~" {
			return usr.HomeDir, nil
		}

		// Replace ~ with the user's home directory
		return filepath.Join(usr.HomeDir, p[1:]), nil
	}
	return p, nil
}
