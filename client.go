package main

import (
	"crypto/sha1"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const (
	clientVersion = "1.2.25"
	clientBuild   = 96

	clientAPIHost   = "rpc.hentaiathome.net"
	clientAPIScheme = "http"
	clientAPIPath   = "clientapi.php"

	argClientBuild    = "clientbuild"
	argAction         = "act"
	argClientID       = "cid"
	argActionArgument = "add"
	argActionKey      = "actkey"
	argTime           = "acttime"

	actionKeyStart     = "hentai@home"
	actionKeyDelimiter = "-"
	actionStart        = "client_start"

	httpGET = "get"
)

type APIResponce struct {
	Success bool
	Message string
	Data    string
}

// Client is api for hath rpc
type Client struct {
	id         int64
	key        string
	httpClient *http.Client
}

func (c Client) getURL(args ...string) *url.URL {
	// preparing time and arguments
	if len(args) == 0 {
		panic("bad arguments lenght in getServerAPIURL")
	}
	unixTime := time.Now().Unix()
	action := args[0]
	argument := ""
	if len(args) > 1 {
		argument = args[1]
	}

	// building action key
	toHash := strings.Join([]string{actionKeyStart,
		action,
		argument,
		fmt.Sprint(c.id),
		fmt.Sprint(unixTime),
		c.key},
		actionKeyDelimiter,
	)
	h := sha1.New()
	fmt.Fprint(h, toHash)
	actionKey := fmt.Sprintf("%x", h.Sum(nil))

	// building url
	u := &url.URL{Scheme: clientAPIScheme, Path: clientAPIPath, Host: clientAPIHost}
	values := make(url.Values)
	values.Add(argClientBuild, fmt.Sprint(clientBuild))
	values.Add(argAction, action)
	values.Add(argActionArgument, argument)
	values.Add(argActionKey, actionKey)
	values.Add(argTime, fmt.Sprint(unixTime))
	values.Add(argClientID, fmt.Sprint(c.id))
	u.RawQuery = values.Encode()
	return u
}

// ActionURL - get url for action
func (c Client) ActionURL(action string) *url.URL {
	return c.getURL(action)
}

func (c Client) getResponce(action string) (r APIResponce, err error) {
	return
}

func (c Client) printRequest(action string) {
	u := c.ActionURL(action)
	log.Println(httpGET, u)
	res, err := c.httpClient.Get(u.String())
	if err != nil {
		log.Print(err)
	}
	defer res.Body.Close()
	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Print(err)
	}
	log.Println(string(data))
}

// NewClient creates new client for api
func NewClient(id int64, key string) *Client {
	c := new(Client)
	c.id = id
	c.key = key
	c.httpClient = http.DefaultClient
	return c
}
