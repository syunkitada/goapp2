package str_utils

import (
	"encoding/json"
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

func OutputByFormat(data interface{}, format string) {
	switch format {
	case "yaml":
		bytes, err := yaml.Marshal(&data)
		if err != nil {
			fmt.Println("Failed yaml.Marshal: err=%s", err.Error())
			os.Exit(1)
		}
		fmt.Println(string(bytes))
	case "json":
		bytes, err := json.Marshal(&data)
		if err != nil {
			fmt.Println("Failed yaml.Marshal: err=%s", err.Error())
			os.Exit(1)
		}
		fmt.Println(string(bytes))
	default:
		fmt.Println(data)
	}
}

func AppendOutputByFormat(outputs []string, data interface{}, format string) []string {
	var bytes []byte
	var err error
	switch format {
	case "yaml":
		bytes, err = yaml.Marshal(&data)
	case "json":
		bytes, err = json.Marshal(&data)
	}
	if err != nil {
		fmt.Println("Failed yaml.Marshal: err=%s", err.Error())
		os.Exit(1)
	}
	return append(outputs, string(bytes))
}
