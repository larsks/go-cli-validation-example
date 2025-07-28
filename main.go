// viper-validator-example -- example of command line processing and validation
// Copyright (C) 2025 Lars Kellogg-Stedman
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.
package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

type Config struct {
	Color             string `mapstructure:"color" validate:"required,oneof=red green blue yellow"`
	Size              int    `mapstructure:"size" validate:"required,min=1,max=100"`
	Count             int    `mapstructure:"count" validate:"min=1,max=1000"`
	IncludeCupHolders bool   `mapstructure:"include-cupholders"`
}

// This is terrible error handling just to make my linter
// stop yelling at me.
func must(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	var config Config

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

	pflag.String("color", "", "Color option (red, green, blue, yellow)")
	pflag.Int("size", 10, "Size value (1-100)")
	pflag.Int("count", 1, "Count value (1-1000)")
	pflag.Bool("include-cupholders", false, "Whether or not to include cupholders")
	pflag.Parse()

	must(viper.BindPFlags(pflag.CommandLine))

	// Look for YAML config file `config.yaml` in the current directory.
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			fmt.Println("No config file found, using defaults and command line flags")
		} else {
			log.Fatalf("Error reading config file: %v", err)
		}
	}

	// Look for environment variables prefixed with `EXAMPLE_`, such as
	// EXAMPLE_COLOR.
	viper.SetEnvPrefix("EXAMPLE")

	// We need to  translate between environment variable names like
	// EXAMPLE_INCLUDE_CUPHOLDERS and option names like --include-cupholders.
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	viper.AutomaticEnv()

	if err := viper.Unmarshal(&config); err != nil {
		log.Fatalf("Unable to decode into struct: %v", err)
	}

	validate := validator.New()
	if err := validate.Struct(&config); err != nil {
		fmt.Printf("Validation failed:\n")

		// This loop produces "nicer" error messages that we get by default, but
		// you could dispense with this and just print out the string value of err:
		//
		// fmt.Printf("Validation failed: %s\n", err)
		for _, err := range err.(validator.ValidationErrors) {
			switch err.Tag() {
			case "required":
				fmt.Printf("- %s is required\n", err.Field())
			case "oneof":
				fmt.Printf("- %s must be one of: red, green, blue, yellow\n", err.Field())
			case "min":
				fmt.Printf("- %s must be at least %s\n", err.Field(), err.Param())
			case "max":
				fmt.Printf("- %s must be at most %s\n", err.Field(), err.Param())
			default:
				fmt.Printf("- %s failed validation: %s\n", err.Field(), err.Tag())
			}
		}
		os.Exit(1)
	}

	fmt.Printf("Configuration loaded successfully:\n")
	fmt.Printf("Color: %s\n", config.Color)
	fmt.Printf("Size: %d\n", config.Size)
	fmt.Printf("Count: %d\n", config.Count)
	fmt.Printf("Include cup holders: %v\n", config.IncludeCupHolders)
}
