package viper

import (
	"reflect"
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func Load(confPath string, receiver interface{}) {
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	viper.AddConfigPath(".")
	viper.SetConfigName("config")
	viper.SetConfigType("yml")

	configType := "default"
	if confPath != "" {
		viper.SetConfigFile(confPath)
		configType = "supplied"
	}

	if err := viper.ReadInConfig(); err != nil {
		logrus.WithField("config_type", configType).WithError(err).Fatal("read config")
	}

	logrus.WithFields(logrus.Fields{
		"config":      viper.ConfigFileUsed(),
		"config_type": configType,
	}).Info("viper using config")

	bindEnvs(reflect.ValueOf(receiver))

	if err := viper.Unmarshal(receiver); err != nil {
		logrus.WithError(err).Fatal("unmarshal viper config file")
	}
}

//nolint:exhaustive // later
func bindEnvs(v reflect.Value, parts ...string) {
	if v.Kind() == reflect.Ptr {
		if v.IsNil() {
			return
		}

		bindEnvs(v.Elem(), parts...)
		return
	}

	ift := v.Type()
	for i := 0; i < ift.NumField(); i++ {
		val := v.Field(i)
		t := ift.Field(i)
		tv, ok := t.Tag.Lookup("mapstructure")
		if !ok {
			continue
		}

		switch val.Kind() {
		case reflect.Struct:
			bindEnvs(val, append(parts, tv)...)
		default:
			if err := viper.BindEnv(strings.Join(append(parts, tv), ".")); err != nil {
				logrus.WithError(err).Fatal("bind env")
			}
		}
	}
}
