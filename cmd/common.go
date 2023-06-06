package cmd

import "fmt"

func checkAuthorized() error {
	if TokenResp == nil {
		return fmt.Errorf("not authorized, execute `auth` command to authorize first")
	}
	return nil
}
