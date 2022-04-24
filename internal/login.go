package internal

import (
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"regexp"
)

var header = map[string]string{
	"User-Agent": "Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/63.0.3239.26 Safari/537.36",
}

var (
	ErrMatchIPFail    = errors.New("Can not match ip from page")
	ErrMatchTokenFail = errors.New("Can not match token from api")
)

type LoginManager struct {
	username           string
	password           string
	loginPageURL       string
	getChallengeApiURL string

	ip    string
	token string

	client *http.Client
}

type NewLoginManagerParams struct {
	Username string
	Password string

	LoginPageURL       string
	GetChallengeApiURL string
}

func NewLoginManager(params NewLoginManagerParams) *LoginManager {
	return &LoginManager{
		username:           params.Username,
		password:           params.Password,
		loginPageURL:       params.LoginPageURL,
		getChallengeApiURL: params.GetChallengeApiURL,
		client:             &http.Client{},
	}
}

func (l *LoginManager) getIP() error {
	log.Println("Step1: Get local ip returned from srun server.")

	// Get login page
	req, err := http.NewRequest(http.MethodGet, l.loginPageURL, nil)
	if err != nil {
		return err
	}
	for k, v := range header {
		req.Header.Add(k, v)
	}
	response, err := l.client.Do(req)
	if err != nil {
		return err
	}
	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return err
	}

	// Get IP from login page with regex
	r := regexp.MustCompile("id=\"user_ip\" value=\"(.*?)\"")
	match := r.FindAllSubmatch(body, -1)
	if len(match) == 0 {
		return ErrMatchIPFail
	}
	l.ip = string(match[0][1])

	log.Println("The IP is: ", l.ip)
	log.Println("----------------")
	return nil
}

func (l *LoginManager) getToken() error {
	log.Println("Step2: Get token by resolving challenge result.")

	// Get challenge
	params := url.Values{}
	parseURL, err := url.Parse(l.getChallengeApiURL)
	if err != nil {
		return err
	}
	params.Set("callback", "jQuery15815616146")
	params.Set("username", l.username)
	params.Set("ip", l.ip)
	parseURL.RawQuery = params.Encode()

	req, err := http.NewRequest(http.MethodGet, parseURL.String(), nil)
	if err != nil {
		return err
	}
	for k, v := range header {
		req.Header.Add(k, v)
	}
	response, err := l.client.Do(req)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return err
	}

	// Get token from challenge api with regex
	r := regexp.MustCompile("\"challenge\":\"(.*?)\"")
	match := r.FindAllSubmatch(body, -1)
	log.Printf("%q\n", match)
	if len(match) == 0 {
		return ErrMatchTokenFail
	}
	l.token = string(match[0][1])

	log.Println("The token is: ", l.token)
	log.Println("----------------")
	return nil
}

func (l *LoginManager) Login() error {
	err := l.getIP()
	if err != nil {
		return err
	}

	err = l.getToken()
	if err != nil {
		return err
	}

	return nil
}
