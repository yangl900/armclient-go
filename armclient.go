package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/urfave/cli"
)

const (
	flagVerbose = "verbose"
)

func main() {
	app := cli.NewApp()

	app.Name = "armclient"
	app.Usage = "Command line client for Azure Resource Manager APIs."
	app.Version = "0.1"
	app.Description = "This is a Go implementation of original windows version ARMClient (https://github.com/projectkudu/ARMClient/). " +
		"I intend to keep commands same as original, and now you can enjoy the useful tool on Linux. " +
		"Additionally in MSI environment like Azure Cloud Shell (https://shell.azure.com/), login is handled automatically. It just works."

	app.Action = func(c *cli.Context) error {
		fmt.Println("no verb specified!")
		cli.ShowAppHelp(c)
		return nil
	}

	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:  flagVerbose,
			Usage: "output verbose messages like request Uri, headers etc.",
		},
	}

	app.Commands = []cli.Command{
		{
			Name:   "get",
			Action: doRequest,
			Usage:  "Makes a GET request to ARM endpoint.",
		},
		{
			Name:   "head",
			Action: doRequest,
			Usage:  "Makes a HEAD request to ARM endpoint.",
		},
		{
			Name:   "put",
			Action: doRequest,
			Usage:  "Makes a PUT request to ARM endpoint.",
		},
		{
			Name:   "patch",
			Action: doRequest,
			Usage:  "Makes a PUT request to ARM endpoint.",
		},
		{
			Name:   "delete",
			Action: doRequest,
			Usage:  "Makes a DELETE request to ARM endpoint.",
		},
		{
			Name:   "post",
			Action: doRequest,
			Usage:  "Makes a POST request to ARM endpoint.",
		},
	}

	app.CustomAppHelpTemplate = cli.AppHelpTemplate

	err := app.Run(os.Args)
	if err != nil {
		fmt.Println(err)
	}
}

func doRequest(c *cli.Context) error {
	if len(c.Args()) == 0 {
		return errors.New("No path specified")
	}

	url, err := getRequestURL(c.Args().First())
	if err != nil {
		return err
	}

	client := &http.Client{}
	req, _ := http.NewRequest(strings.ToUpper(c.Command.Name), url, nil)

	token, err := acquireAuthToken()
	if err != nil {
		return errors.New("Failed to acquire auth token: " + err.Error())
	}

	req.Header.Set("Authorization", token)

	response, err := client.Do(req)
	if err != nil {
		return errors.New("Request failed: " + err.Error())
	}

	defer response.Body.Close()
	buf, err := ioutil.ReadAll(response.Body)

	if err != nil {
		return errors.New("Request failed: " + err.Error())
	}

	if c.GlobalBool(flagVerbose) {
		fmt.Println(responseDetail(response))
	}

	fmt.Println(prettyJSON(buf))
	return nil
}
