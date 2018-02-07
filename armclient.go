package main

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/urfave/cli"
)

const (
	appVersion   = "0.1"
	flagVerbose  = "verbose"
	userAgentStr = "github.com/yangl900/armclient-go"
)

func main() {
	app := cli.NewApp()

	app.Name = "armclient"
	app.Usage = "Command line client for Azure Resource Manager APIs."
	app.Version = appVersion
	app.Description = `
		This is a Go implementation of original windows version ARMClient (https://github.com/projectkudu/ARMClient/).
		Commands are kept same as much as possible, and now you can enjoy the useful tool on Linux & Mac.
		Additionally in Azure Cloud Shell (https://shell.azure.com/), login is handled automatically. It just works.`

	app.Action = func(c *cli.Context) error {
		cli.ShowAppHelp(c)
		return nil
	}

	log.SetOutput(ioutil.Discard)

	verboseFlag := cli.BoolFlag{
		Name:  flagVerbose,
		Usage: "output verbose messages like request Uri, headers etc.",
	}

	app.Flags = []cli.Flag{verboseFlag}

	app.Commands = []cli.Command{
		{
			Name:   "get",
			Action: doRequest,
			Usage:  "Makes a GET request to ARM endpoint.",
			Flags:  []cli.Flag{verboseFlag},
		},
		{
			Name:   "head",
			Action: doRequest,
			Usage:  "Makes a HEAD request to ARM endpoint.",
			Flags:  []cli.Flag{verboseFlag},
		},
		{
			Name:   "put",
			Action: doRequest,
			Usage:  "Makes a PUT request to ARM endpoint.",
			Flags:  []cli.Flag{verboseFlag},
		},
		{
			Name:   "patch",
			Action: doRequest,
			Usage:  "Makes a PUT request to ARM endpoint.",
			Flags:  []cli.Flag{verboseFlag},
		},
		{
			Name:   "delete",
			Action: doRequest,
			Usage:  "Makes a DELETE request to ARM endpoint.",
			Flags:  []cli.Flag{verboseFlag},
		},
		{
			Name:   "post",
			Action: doRequest,
			Usage:  "Makes a POST request to ARM endpoint.",
			Flags:  []cli.Flag{verboseFlag},
		},
	}

	app.CustomAppHelpTemplate = cli.AppHelpTemplate

	err := app.Run(os.Args)
	if err != nil {
		fmt.Println(err)
	}
}

func isWriteVerb(verb string) bool {
	v := strings.ToUpper(verb)
	return v == "PUT" || v == "POST" || v == "PATCH"
}

func doRequest(c *cli.Context) error {
	if len(c.Args()) == 0 {
		return errors.New("No path specified")
	}

	url, err := getRequestURL(c.Args().First())
	if err != nil {
		return err
	}

	var reqBody string
	if isWriteVerb(c.Command.Name) && c.NArg() > 1 {
		reqBody = c.Args().Get(1)

		if strings.HasPrefix(reqBody, "@") {
			filePath, _ := filepath.Abs(strings.TrimSuffix(strings.TrimPrefix(strings.TrimPrefix(reqBody, "@"), "'"), "'"))

			if _, err := os.Stat(filePath); err != nil {
				return errors.New("File not found: " + filePath)
			}

			buffer, err := ioutil.ReadFile(filePath)
			if err != nil {
				return errors.New("Failed to read file: " + filePath)
			}

			reqBody = prettyJSON(buffer)
		} else {
			reqBody = prettyJSON([]byte(reqBody))
			fmt.Println(reqBody)
		}
	}

	client := &http.Client{}
	req, _ := http.NewRequest(strings.ToUpper(c.Command.Name), url, bytes.NewReader([]byte(reqBody)))

	token, err := acquireAuthToken()
	if err != nil {
		return errors.New("Failed to acquire auth token: " + err.Error())
	}

	req.Header.Set("Authorization", token)
	req.Header.Set("User-Agent", userAgentStr)
	req.Header.Set("x-ms-client-request-id", newUUID())
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")

	start := time.Now()
	response, err := client.Do(req)
	if err != nil {
		return errors.New("Request failed: " + err.Error())
	}

	defer response.Body.Close()
	buf, err := ioutil.ReadAll(response.Body)

	if err != nil {
		return errors.New("Request failed: " + err.Error())
	}

	if c.GlobalBool(flagVerbose) || c.Bool(flagVerbose) {
		fmt.Println(responseDetail(response, time.Now().Sub(start)))
	}

	fmt.Println(prettyJSON(buf))
	return nil
}
