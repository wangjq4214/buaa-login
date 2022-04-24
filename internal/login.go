package internal

import (
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"regexp"
	"strings"

	"github.com/wangjq4214/buaa-login/internal/encryption"
)

var header = map[string]string{
	"User-Agent": "Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/63.0.3239.26 Safari/537.36",
}

var (
	ErrMatchIPFail    = errors.New("Can not match ip from page")
	ErrMatchTokenFail = errors.New("Can not match token from api")
	ErrLoginFail      = errors.New("Can not login from api")
)

func getRequestURL(s string, p map[string]string) (string, error) {
	params := url.Values{}
	parseURL, err := url.Parse(s)
	if err != nil {
		return "", err
	}
	for k, v := range p {
		params.Set(k, v)
	}
	parseURL.RawQuery = params.Encode()

	return parseURL.String(), nil
}

type LoginManager struct {
	username           string
	password           string
	n                  string
	acID               string
	enc                string
	vType              string
	loginPageURL       string
	getChallengeApiURL string
	loginApiURL        string

	ip    string
	token string

	client *http.Client
}

type NewLoginManagerParams struct {
	Username string
	Password string
	N        string
	AcID     string
	Enc      string
	VType    string

	LoginPageURL       string
	GetChallengeApiURL string
	LoginApiURL        string
}

func NewLoginManager(params NewLoginManagerParams) *LoginManager {
	return &LoginManager{
		username:           params.Username,
		password:           params.Password,
		n:                  params.N,
		acID:               params.AcID,
		enc:                params.Enc,
		vType:              params.VType,
		loginPageURL:       params.LoginPageURL,
		getChallengeApiURL: params.GetChallengeApiURL,
		loginApiURL:        params.LoginApiURL,
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
	u, err := getRequestURL(l.getChallengeApiURL, map[string]string{
		"callback": "jQuery15815616146",
		"username": l.username,
		"ip":       l.ip,
	})
	req, err := http.NewRequest(http.MethodGet, u, nil)
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
	if len(match) == 0 {
		return ErrMatchTokenFail
	}
	l.token = string(match[0][1])

	log.Println("The token is: ", l.token)
	log.Println("----------------")
	return nil
}

func (l *LoginManager) encryptInfo() (string, error) {
	info := [][]string{
		{"username", l.username},
		{"password", l.password},
		{"ip", l.ip},
		{"acid", l.acID},
		{"enc_ver", l.enc},
	}

	infoLength := len(info)
	builder := strings.Builder{}
	builder.WriteString("{")
	for i, v := range info {
		builder.WriteString("\"")
		builder.WriteString(strings.Trim(v[0], " "))
		builder.WriteString("\":\"")
		builder.WriteString(strings.Trim(v[1], " "))
		builder.WriteString("\"")
		if i < infoLength-1 {
			builder.WriteString(",")
		}
	}
	builder.WriteString("}")

	return "{SRBX1}" + encryption.GetBase64(encryption.GetXencode(builder.String(), l.token)), nil
}

func (l *LoginManager) encryptMD5() string {
	return "{MD5}" + encryption.GetMD5(l.token, "")
}

func (l *LoginManager) encryptChecksum(md5, info string) string {
	builder := strings.Builder{}

	builder.WriteString(l.token)
	builder.WriteString(l.username)

	builder.WriteString(l.token)
	builder.WriteString(encryption.GetMD5(l.token, ""))

	builder.WriteString(l.token)
	builder.WriteString(l.acID)

	builder.WriteString(l.token)
	builder.WriteString(l.ip)

	builder.WriteString(l.token)
	builder.WriteString(l.n)

	builder.WriteString(l.token)
	builder.WriteString(l.vType)

	builder.WriteString(l.token)
	builder.WriteString(info)

	return encryption.GetSHA1(builder.String())
}

func (l *LoginManager) login() error {
	log.Println("Step3: Login and resolve response.")
	info, err := l.encryptInfo()
	if err != nil {
		return err
	}
	md5 := l.encryptMD5()
	checksum := l.encryptChecksum(md5, info)

	u, err := getRequestURL(l.loginApiURL, map[string]string{
		"callback": "jsonp1583251661368",
		"action":   "login",
		"username": l.username,
		"password": md5,
		"ac_id":    l.acID,
		"ip":       l.ip,
		"info":     info,
		"chksum":   checksum,
		"n":        l.n,
		"type":     l.vType,
	})
	req, err := http.NewRequest(http.MethodGet, u, nil)
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

	log.Println(string(body))

	r := regexp.MustCompile("\"suc_msg\":\"(.*?)\"")
	match := r.FindAllSubmatch(body, -1)
	if len(match) == 0 {
		return ErrLoginFail
	}

	log.Println("The success msg is: ", string(match[0][1]))
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

	err = l.login()
	if err != nil {
		return err
	}

	return nil
}
