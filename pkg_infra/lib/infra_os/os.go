package infra_os

import "os"

func Exit(disable bool, code int) {
	if !disable {
		os.Exit(code)
	}
}
