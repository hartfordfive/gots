package main

import (
	"io/ioutil"
	"os/user"
	"strings"
)

func loadApiCreds(file string) (map[string]string, error) {

	usr, _ := user.Current()
	dir := usr.HomeDir

	if file[:2] == "~/" {
		file = strings.Replace(file, "~/", dir+"/", 1)
	}

	contents, err := ioutil.ReadFile(file)
	conf := map[string]string{}
	if err == nil {
		lines := strings.Split(string(contents), "\n")
		for _, l := range lines {
			if strings.TrimSpace(l) == "" || string(strings.TrimSpace(l)[0:0]) == "#" {
				continue
			}
			parts := strings.Split(l, "=")
			conf[strings.Trim(parts[0], " ")] = strings.Trim(parts[1], " ")
		}
		return conf, nil
	}
	return conf, err
}
