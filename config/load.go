package config

import (
	"reflect"
	"strings"

	"github.com/spf13/viper"
)

const (
	viperDefaultDelimiter = "."
	defaultTagName        = "default"
	squashTagValue        = ",squash"
	mapStructureTagName   = "mapstructure"
)

func Read() (Config, error) {
	v := viper.New()
	v.SetEnvKeyReplacer(strings.NewReplacer(viperDefaultDelimiter, "_")) // replace default viper delimiter for env vars
	v.AutomaticEnv()
	v.SetTypeByDefaultValue(true)

	var cfg Config

	err := setDefaults("", v, &reflect.StructField{}, reflect.ValueOf(&cfg).Elem())
	if err != nil {
		return cfg, err
	}

	err = v.Unmarshal(&cfg)
	if err != nil {
		return cfg, err
	}

	return cfg, nil
}

// setDefaults sets default values for struct fields based using value from default tag.
func setDefaults(parentName string, vip *viper.Viper, t *reflect.StructField, v reflect.Value) error {
	if v.Kind() == reflect.Struct {
		value, ok := t.Tag.Lookup(mapStructureTagName)
		if ok && value != squashTagValue {
			if parentName != "" {
				parentName += viperDefaultDelimiter
			}

			parentName += strings.ToUpper(value)
		}

		for i := 0; i < v.NumField(); i++ {
			f := v.Type().Field(i)
			if err := setDefaults(parentName, vip, &f, v.Field(i)); err != nil {
				return err
			}
		}

		return nil
	}

	value, _ := t.Tag.Lookup(defaultTagName)
	fieldName, ok := t.Tag.Lookup(mapStructureTagName)

	if ok && fieldName != squashTagValue {
		if parentName != "" {
			fieldName = parentName + viperDefaultDelimiter + strings.ToUpper(fieldName)
		}

		vip.SetDefault(strings.ToUpper(fieldName), value)
	}

	return nil
}
