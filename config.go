package main

import (
	"fmt"
	"io/ioutil"
	"reflect"
	"strconv"
	"strings"
)

type Config struct {
	Server  string
	Channel string
	Nr      int
	IPv6    bool
}

func loadConfig(file string) (config Config, err error) {
	configStr, err := ioutil.ReadFile(file)
	if err != nil {
		return config, err
	}

	configReflect := reflect.ValueOf(&config).Elem()

	for _, line := range strings.Split(string(configStr), "\n") {
		lineParts := strings.Split(line, "\"")
		line = ""
		for k, v := range lineParts {
			if k%2 == 0 {
				if i := strings.Index(v, "#"); i != -1 {
					line += v[:i]
					break
				}
			}
			line += v
		}

		if strings.TrimSpace(line) == "" {
			continue
		}

		keyVal := strings.SplitN(line, "=", 2)
		if len(keyVal) < 2 {
			return config, fmt.Errorf("Config line must contain \"=\": %s", line)
		}
		key := strings.TrimSpace(keyVal[0])
		value := strings.TrimSpace(keyVal[1])

		field := configReflect.FieldByName(key)
		if !field.IsValid() {
			return config, fmt.Errorf("Config key is not valid: %s", key)
		}

		switch field.Kind() {
		case reflect.String:
			field.SetString(value)
		case reflect.Int:
			i, err := strconv.ParseInt(value, 10, 64)
			if err != nil {
				return config, fmt.Errorf("Invalid int \"%s\" in key \"%s\": %s", value, key, err)
			}
			field.SetInt(i)
		case reflect.Bool:
			if value != "true" && value != "false" {
				return config, fmt.Errorf("Invalid bool \"%s\" in key \"%s\": must be true or false", value, key)
			}
			field.SetBool(value == "true")
		}
	}

	return config, nil
}

func main() {
	config, err := loadConfig("itkbot.conf")
	if err != nil {
		fmt.Println("Couldn't load config:", err)
		return
	}

	fmt.Println("Server:", config.Server)
	fmt.Println("Channel:", config.Channel)
	fmt.Println("Nr:", config.Nr)
	fmt.Println("IPv6:", config.IPv6)
}
