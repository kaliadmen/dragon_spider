package main

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/fatih/color"
	"io/ioutil"
	"os"
	"regexp"
	"strings"
)

func makeKey() error {
	randKey := ds.RandomString(32)
	color.Yellow("Random 32 character encryption key: %s", randKey)
	err := autoAdd(randKey)
	if err != nil {
		gracefulExit(err)
	}
	return nil
}

func autoAdd(key string) error {
	fmt.Println("Do you want to set this key in your env file? (y|n)")
	reader := bufio.NewReader(os.Stdin)
	res, _ := reader.ReadString('\n')
	// convert CRLF to LF
	res = strings.Replace(res, "\n", "", -1)

	switch res {

	case "y", "Y", "yes", "Yes":
		err := addKeyToEnv(key)
		if err != nil {
			gracefulExit(err)
		}
		return nil

	case "n", "N", "no", "No":
		return nil

	default:
		color.Yellow("You may have to add the key manually! Key: %s", key)
		return nil
	}
}

func addKeyToEnv(key string) error {
	filePath := ds.RootPath + "/.env"

	if !fileExists(filePath) {
		return errors.New("no .env file found")
	}

	data, err := os.ReadFile(filePath)
	if err != nil {
		gracefulExit(err)
	}

	env := string(data)
	keyRegex := regexp.MustCompile(`KEY=[a-zA-Z\d!@#$%?^&*()-_+{}]*`)
	env = keyRegex.ReplaceAllString(env, "KEY="+key)

	err = ioutil.WriteFile(filePath, []byte(env), 0644)
	if err != nil {
		gracefulExit(err)
	}

	return nil
}
