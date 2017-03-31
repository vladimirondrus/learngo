package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"regexp"
)

func main() {
	var (
		cmdOut []byte
		err    error
	)

	url := flag.String("url", "empty", "a string")
	flag.Parse()
	fmt.Println("url: ", *url)

	data := map[string]string{}
	cmdName := "curl"
	cmdArgs := []string{*url, "--insecure", "--cert", "/etc/idbt/ssl/consul/certificate.crt", "--key", "/etc/idbt/ssl/consul/private_key.key"}
	if cmdOut, err = exec.Command(cmdName, cmdArgs...).Output(); err != nil {
		fmt.Fprintln(os.Stderr, "There was an error running", cmdName, ":", err)
		os.Exit(1)
	}

	err = json.Unmarshal(cmdOut, &data)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(cmdOut))
	problem := false
	for _, v := range data {
		//fmt.Printf("key[%s] value[%s]\n", k, v)
		matched, _ := regexp.MatchString("^CRITICAL|WARNING", v)
		if matched {
			problem = true
		}
	}

	if problem {
		os.Exit(1)
	}
}
