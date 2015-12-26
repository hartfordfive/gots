package main

import (
	b64 "encoding/base64"
	"fmt"
	"github.com/codegangsta/cli"
	"github.com/franela/goreq"
	"net"
	"net/url"
	"os"
	"reflect"
	"strconv"
	"strings"
	"time"
)

const (
	DEFAULT_API_BASE_URL  = "https://onetimesecret.com"
	DEFAULT_SETTINGS_FILE = "~/.gots"
	USER_AGENT_STRING     = "Go OTS client/v0.1.0"
)

var api_methods map[string]string

type OneTimeSecretClient struct {
	ApiUrl  string
	User    string
	ApiKey  string
	Methods map[string]string
	Debug   bool
}

func main() {

	app := cli.NewApp()
	app.Name = "Go One-Time-Secret"

	ots := NewOneTimeSecretClient()

	fmt.Println()

	app.Commands = []cli.Command{
		{
			Name:  "share",
			Usage: "gots share [secret] [passphrase] [ttl] [recipient_email]",
			Action: func(c *cli.Context) {
				ots.ShareSecret(c.Args())
			},
		},
		{
			Name:  "generate",
			Usage: "gots generate [passphrase] [ttl] [metadata_ttl] [secret_ttl] [recipient_email]",
			Action: func(c *cli.Context) {
				ots.GenerateSecret(c.Args())
			},
		},
		{
			Name:  "get",
			Usage: "gots get [secret_key] [passphrase]",
			Action: func(c *cli.Context) {
				ots.GetSecret(c.Args())
			},
		},
		{
			Name:  "getmeta",
			Usage: "gots getmeta [metadata_key]",
			Action: func(c *cli.Context) {
				ots.GetMetadata(c.Args())
			},
		},
		{
			Name:  "recentmeta",
			Usage: "gots recentmeta",
			Action: func(c *cli.Context) {
				ots.GetRecentMetadata(c.Args())
			},
		},
		{
			Name:  "status",
			Usage: "gots status",
			Action: func(c *cli.Context) {
				ots.GetApiStatus()
			},
		},
	}
	app.Run(os.Args)
}

func NewOneTimeSecretClient() *OneTimeSecretClient {
	conf, err := loadApiCreds(DEFAULT_SETTINGS_FILE)
	if err != nil {
		fmt.Println("Error: Could not load credentials file")
		fmt.Println(err)
		os.Exit(1)
	}
	debug, _ := strconv.ParseBool(conf["debug"])
	api_host := DEFAULT_API_BASE_URL
	if _, ok := conf["ots_host"]; ok {
		if debug {
			fmt.Println("[DEBUG] Using custom host:", conf["ots_host"])
		}
		u, err := url.Parse(conf["ots_host"])
		if err != nil {
			fmt.Println("Error: Invalid hostname specified")
			os.Exit(1)
		}
		_, err = net.LookupHost(u.Host)
		if err != nil {
			fmt.Printf("Error: Hostname %s is unreachable", u.Host)
			os.Exit(1)
		}
		api_host = conf["ots_host"]
	}

	return &OneTimeSecretClient{
		ApiUrl: api_host,
		User:   conf["user"],
		ApiKey: conf["api_key"],
		Debug:  debug,
	}
}

func (c *OneTimeSecretClient) ShareSecret(params []string) {

	payload := []string{}

	if len(params) >= 1 {
		payload = append(payload, "secret="+params[0])
	} else {
		errorAndExit("You must specify a secret")
	}

	if len(params) >= 2 {
		payload = append(payload, "passphrase="+params[1])
	}

	if len(params) >= 3 {
		payload = append(payload, "ttl="+params[2])
	} else {
		payload = append(payload, "ttl=3600")
	}

	if len(params) >= 4 {
		payload = append(payload, "recipient="+params[3])
	}

	sr, _, err := c.issueRequest("/api/v1/share", "POST", payload)
	if err != nil {
		errorAndExit(err.Error())
	}
	fmt.Println("Secret url:", "https://onetimesecret.com/secret/"+sr.SecretKey+"\n")
	fmt.Println("Metadata key:", sr.MetadataKey)
	fmt.Println("Secret TTL:", sr.SecretTtl)
	fmt.Println()
	os.Exit(0)
}

func (c *OneTimeSecretClient) GenerateSecret(params []string) {

	payload := []string{}
	if len(params) >= 1 {
		payload = append(payload, "passphrase="+params[0])
	}

	if len(params) >= 2 {
		payload = append(payload, "ttl="+params[1])
	} else {
		payload = append(payload, "ttl=3600")
	}

	if len(params) >= 3 {
		payload = append(payload, "metadata_ttl="+params[2])
	} else {
		payload = append(payload, "metadata_ttl=3600")
	}

	if len(params) >= 4 {
		payload = append(payload, "secret_ttl="+params[3])
	} else {
		payload = append(payload, "secret_ttl=3600")
	}

	if len(params) >= 5 {
		payload = append(payload, "recipient="+params[4])
	}

	sr, _, err := c.issueRequest("/api/v1/generate", "POST", payload)
	if err != nil {
		fmt.Println("Error: ", err, "\n")
		os.Exit(1)
	}
	fmt.Println("Secret url: ", "https://onetimesecret.com/secret/"+sr.SecretKey)
	fmt.Println("Metadata key:", sr.MetadataKey)
	fmt.Println("Expires in", sr.SecretTtl, "seconds")
	fmt.Println()
	os.Exit(0)
}

func (c *OneTimeSecretClient) GetSecret(params []string) {

	payload := []string{}
	if len(params) >= 2 {
		payload = append(payload, "passphrase="+params[1])
	} else {
		fmt.Println("Error: Missing secret key")
		os.Exit(1)
	}

	sr, _, err := c.issueRequest("/api/v1/secret/"+params[0], "POST", payload)
	if err != nil {
		fmt.Println("Error: ", err, "\n")
		os.Exit(1)
	}
	fmt.Println("Secret value: ", sr.Value+"\n")
	os.Exit(0)
}

func (c *OneTimeSecretClient) GetMetadata(params []string) {
	res, _, err := c.issueRequest("/api/v1/private/"+params[0], "POST", []string{})
	if err != nil {
		fmt.Println("Error: ", err)
		os.Exit(1)
	}
	fmt.Println("Secret Metadata Information")
	fmt.Println("----------------------------")

	printStruct(res)
	fmt.Println()

	os.Exit(0)
}

func (c *OneTimeSecretClient) GetRecentMetadata(params []string) {
	res, _, err := c.issueRequest("/api/v1/private/recent", "POST", []string{})
	if err != nil {
		fmt.Println("Error: ", err)
		os.Exit(1)
	}
	fmt.Println("Recent Metadata: ", res)
	fmt.Println()
	os.Exit(0)
}

func (c *OneTimeSecretClient) GetApiStatus() {

	_, status_code, err := c.issueRequest("/api/v1/status", "GET", []string{})
	if err != nil {
		fmt.Println("Error: ", err, "\n")
		os.Exit(1)
	}
	status := "UP"
	if status_code != 200 {
		status = "DOWN"
	}
	fmt.Println("API Status: ", status)
	fmt.Println()
	os.Exit(0)
}

func (c *OneTimeSecretClient) issueRequest(uri string, method string, params []string) (SecretResponse, int, error) {

	if c.Debug {
		fmt.Println("\n------------- REQUEST DEBUG DATA ------------")
	}

	res, err := goreq.Request{
		Method:    method,
		UserAgent: USER_AGENT_STRING,
		Uri:       c.ApiUrl + uri,
		ShowDebug: c.Debug,
		Body:      strings.Join(params, "&"),
		Timeout:   10000 * time.Millisecond,
	}.WithHeader(
		"Authorization", "Basic "+b64.StdEncoding.EncodeToString([]byte(c.User+":"+c.ApiKey)),
	).Do()

	if c.Debug {
		fmt.Println("----------------------------")
	}

	if err != nil {
		fmt.Println("\nError: ", err)
		return SecretResponse{}, res.StatusCode, err
	}

	if c.Debug {
		fmt.Println("[DEBUG] Response Status Code: ", res.StatusCode, "\n")
	}

	var sr SecretResponse
	res.Body.FromJsonTo(&sr)

	return sr, res.StatusCode, nil
}

func printStruct(res SecretResponse) {
	typ := reflect.TypeOf(res)
	v := reflect.ValueOf(res)

	for i := 0; i < v.NumField(); i++ {
		p := typ.Field(i)
		if !p.Anonymous {
			fmt.Println(p.Name, ":", v.Field(i).Interface())
		}
	}
}

func errorAndExit(msg string) {
	fmt.Println("Error:", msg)
	fmt.Println()
	os.Exit(1)
}
