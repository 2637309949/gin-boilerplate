package viper

import (
	"github.com/spf13/viper"
)

var (
	Get                = viper.Get
	GetBool            = viper.GetBool
	GetFloat64         = viper.GetFloat64
	GetInt             = viper.GetInt
	GetIntSlice        = viper.GetIntSlice
	GetString          = viper.GetString
	GetStringMap       = viper.GetStringMap
	GetStringMapString = viper.GetStringMapString
	GetStringSlice     = viper.GetStringSlice
	GetTime            = viper.GetTime
	GetDuration        = viper.GetDuration
	IsSet              = viper.IsSet
	AllSettings        = viper.AllSettings
)

func init() {
	viper.SetConfigName("app")  // name of config file (without extension)
	viper.SetConfigType("yaml") // REQUIRED if the config file does not have the extension in the name
	viper.AddConfigPath(".")    // optionally look for config in the working directory
	err := viper.ReadInConfig() // Find and read the config file
	if err != nil {             // Handle errors reading the config file
		panic(err)
	}
}
