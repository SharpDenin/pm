package config

import (
	"encoding/json"
	"io/ioutil"
	"strings"

	"gopkg.in/yaml.v2"
)

// Dependency represents a package dependency with name and version constraint
type Dependency struct {
	Name string `json:"name" yaml:"name"`
	Ver  string `json:"ver" yaml:"ver"` // e.g., ">=1.10", "<=2.0", or exact "1.10"
}

// Target represents a file target for packing
type Target struct {
	Path    string `json:"path,omitempty" yaml:"path,omitempty"`
	Exclude string `json:"exclude,omitempty" yaml:"exclude,omitempty"`
}

// PacketConfig for packing
type PacketConfig struct {
	Name    string        `json:"name" yaml:"name"`
	Ver     string        `json:"ver" yaml:"ver"`
	Targets []interface{} `json:"targets" yaml:"targets"` // array of string or Target
	Packets []Dependency  `json:"packets" yaml:"packets"`
}

// PackagesConfig for updating
type PackagesConfig struct {
	Packages []Dependency `json:"packages" yaml:"packages"`
}

// LoadPacketConfig loads and normalizes packet configuration
func LoadPacketConfig(file string) (*PacketConfig, error) {
	data, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}

	var config PacketConfig
	if strings.HasSuffix(file, ".yaml") || strings.HasSuffix(file, ".yml") {
		err = yaml.Unmarshal(data, &config)
	} else {
		err = json.Unmarshal(data, &config)
	}
	if err != nil {
		return nil, err
	}

	// Normalize targets: convert strings to Target structs
	for i, t := range config.Targets {
		switch v := t.(type) {
		case string:
			config.Targets[i] = Target{Path: v}
		case map[string]interface{}:
			path := v["path"].(string)
			exclude := ""
			if ex, ok := v["exclude"]; ok {
				exclude = ex.(string)
			}
			config.Targets[i] = Target{Path: path, Exclude: exclude}
		}
	}
	return &config, nil
}

// LoadPackagesConfig loads packages configuration
func LoadPackagesConfig(file string) (*PackagesConfig, error) {
	data, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}

	var config PackagesConfig
	if strings.HasSuffix(file, ".yaml") || strings.HasSuffix(file, ".yml") {
		err = yaml.Unmarshal(data, &config)
	} else {
		err = json.Unmarshal(data, &config)
	}
	if err != nil {
		return nil, err
	}
	return &config, nil
}
