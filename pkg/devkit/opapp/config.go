// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package opapp

import (
	"fmt"
	"strings"

	"openpitrix.io/openpitrix/pkg/util/jsonutil"
)

const (
	TypeArray    = "array"
	TypeString   = "string"
	TypeInteger  = "integer"
	TypeNumber   = "number"
	TypePassword = "password"
	TypeBoolean  = "boolean"
)

type ConfigTemplate struct {
	Config
	Raw string
}

type Config struct {
	Type        string    `json:"type,omitempty"`
	Properties  []*Config `json:"properties,omitempty"`
	Key         string    `json:"key,omitempty"`
	Description string    `json:"description,omitempty"`
	Required    bool      `json:"required,omitempty"`

	Default interface{} `json:"default,omitempty"`
	Pattern *string     `json:"pattern,omitempty"`

	Limits            map[string][]string `json:"limits,omitempty"`
	AllowedOperations []string            `json:"allowed_operations,omitempty"`

	Port *int `json:"port,omitempty"`

	Range []interface{} `json:"range,omitempty"`
	Min   *float64      `json:"min,omitempty"`
	Max   *float64      `json:"max,omitempty"`
	Step  *int64        `json:"step,omitempty"`

	Changeable  *bool   `json:"changeable,omitempty"`
	Separator   *string `json:"separator,omitempty"`
	Multichoice *bool   `json:"multichoice,omitempty"`
}

func (c *Config) SpecificConfig(key string) {
	if c.Key == "" && c.Type == TypeArray {
		var properties []*Config
		for _, p := range c.Properties {
			if p.Key == key {
				properties = append(properties, p)
				break
			}
		}
		c.Properties = properties
	}
}

func (c *Config) FillInDefaultConfig(defaultConfig jsonutil.Json) {
	c.fillInDefaultConfig(defaultConfig)
}

func (c *Config) fillInDefaultConfig(defaultConfig jsonutil.Json) {
	if c.Type == TypeArray {
		for _, p := range c.Properties {
			if c.Key != "" {
				p.fillInDefaultConfig(defaultConfig.Get(c.Key))
			} else {
				p.fillInDefaultConfig(defaultConfig)
			}
		}
	} else {
		if defaultConfig.Get(c.Key) != nil {
			c.Default = defaultConfig.Get(c.Key).Interface()
		}
	}
}

func (c *Config) GetDefaultConfig() jsonutil.Json {
	// FIXME: need improve performance
	conf := c.getDefaultConfig()
	return jsonutil.ToJson(conf)
}

func (c *Config) getDefaultConfig() interface{} {
	if c.Type == TypeArray {
		defaultConfig := make(map[string]interface{})
		for _, p := range c.Properties {
			defaultConfig[p.Key] = p.getDefaultConfig()
		}
		return defaultConfig
	}
	if c.Default != nil {
		return c.Default
	}
	switch c.Type {
	case TypeInteger:
		return 0
	case TypeNumber:
		return 0.0
	case TypePassword, TypeString:
		return ""
	case TypeBoolean:
		return false
	default:
		return ""
	}
}

// input user defined config, output rendered config with default values
func (c *Config) GetRenderedConfig(config jsonutil.Json) jsonutil.Json {
	v := config
	if v.Interface() == nil && c.Default != nil {
		v = jsonutil.ToJson(c.Default)
	}
	return v
}

func getParent(parent []string) string {
	return strings.Join(parent, ".")
}

// validate the rendered config
func (c *Config) Validate(config jsonutil.Json, parent ...string) error {
	if c.Key != "" {
		parent = append(parent, c.Key)
	}
	//logger.Debug(nil, "%+v %s", v.Interface(), c.Key)
	v := c.GetRenderedConfig(config)
	switch c.Type {
	case TypeInteger:
		value, err := v.Int64()
		if err != nil {
			return fmt.Errorf("[%s] is not [%s] type", getParent(parent), TypeInteger)
		}
		if c.Max != nil && value > int64(*c.Max) {
			return fmt.Errorf("[%s] large than max [%.f]", getParent(parent), *c.Max)
		}
		if c.Min != nil && value < int64(*c.Min) {
			return fmt.Errorf("[%s] less than min [%.f]", getParent(parent), *c.Min)
		}
		if c.Step != nil && value%*c.Step != 0 {
			return fmt.Errorf("[%s] can not exact division with [%d]", getParent(parent), *c.Step)
		}
	case TypeNumber:
		value, err := v.Float64()
		if err != nil {
			return fmt.Errorf("[%s] is not [%s] type", getParent(parent), TypeNumber)
		}
		if c.Max != nil && value > (*c.Max) {
			return fmt.Errorf("[%s] large than max [%f]", getParent(parent), *c.Max)
		}
		if c.Min != nil && value < (*c.Min) {
			return fmt.Errorf("[%s] less than min [%f]", getParent(parent), *c.Min)
		}
	case TypePassword, TypeString:
		value, err := v.String()
		if err != nil {
			return fmt.Errorf("[%s] is not [%s] type", getParent(parent), TypeString)
		}
		if c.Required && len(value) == 0 {
			// subnet will be validated in cluster manager
			if c.Key != "subnet" {
				return fmt.Errorf("[%s] is required", getParent(parent))
			}
		}
	case TypeBoolean:
		_, err := v.Bool()
		if err != nil {
			return fmt.Errorf("[%s] is not [%s] type", getParent(parent), TypeBoolean)
		}
	case TypeArray:
		if len(c.Properties) == 0 {
			return fmt.Errorf("properties must not be empty")
		}
		for _, p := range c.Properties {
			v = config.Get(p.Key)
			err := p.Validate(v)
			if err != nil {
				return err
			}
		}
	default:
		return nil
	}
	return nil
}
