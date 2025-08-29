package tadpoles

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"tadpoles-backup/internal/api"
	"tadpoles-backup/internal/utils"
	"time"

	"github.com/sirupsen/logrus"
)

func loginAdmit(client *http.Client, admitUrl *url.URL, cookieFile string) (expires *time.Time, err error) {
	logrus.Debug("Admit...")

	zone, _ := time.Now().Zone()

	resp, err := client.PostForm(
		admitUrl.String(),
		url.Values{
			"tz": {zone},
		},
	)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, utils.NewRequestError(resp, "tadpoles admit failed")
	}

	logrus.Debug("Admit successful")

	return api.SerializeResponseCookies(cookieFile, resp)
}

func login(
	client *http.Client,
	loginUrl *url.URL,
	email, password string,
) error {
	resp, err := client.PostForm(
		loginUrl.String(),
		url.Values{
			"email":    {email},
			"password": {password},
			"service":  {"tadpoles"},
		},
	)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		return utils.NewRequestError(resp, "tadpoles login failed")
	}

	return nil
}

type HostHeaderTransport struct {
	hostHeader string
}

func (t *HostHeaderTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	req.Header.Add("Host", t.hostHeader)
	return http.DefaultTransport.RoundTrip(req)
}

func requestPasswordReset(
	client *http.Client,
	resetUrl *url.URL,
	email string,
) error {
	resp, err := client.PostForm(
		resetUrl.String(),
		url.Values{
			"email":   {email},
			"app":     {"parent"},
			"service": {"tadpoles"},
		},
	)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		return utils.NewRequestError(resp, "tadpoles reset password request failed")
	}

	return nil
}

func getCookie(token string) (string, error) {
	url := fmt.Sprintf("https://www.tadpoles.com/auth/passwordtoken?token=%s&app=parent", token)
	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", fmt.Errorf("error creating request: %v", err)
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10.15; rv:142.0) Gecko/20100101 Firefox/142.0")
	resp, err := client.Do(req)
	if err != nil {
		if !strings.Contains(err.Error(), "stopped after") {
			return "", fmt.Errorf("error sending request: %v", err)
		}
	}
	defer resp.Body.Close()

	var cookieStrings []string
	for _, cookie := range resp.Cookies() {
		cookieStrings = append(cookieStrings, fmt.Sprintf("%s=%s", cookie.Name, cookie.Value))
	}
	fullCookieHeader := strings.Join(cookieStrings, "; ")
	return fullCookieHeader, nil
}
func resetPassword(
	client *http.Client,
	resetUrl *url.URL,
	resetCode, newPassword string,
) error {
	cookie, err := getCookie(resetCode)
	if err != nil {
		return fmt.Errorf("error getting cookie: %v", err)
	}
	data := url.Values{
		"action":       {"set"},
		"new_password": {newPassword},
		"token":        {resetCode},
		"client":       {"dashboard"},
	}
	req, err := http.NewRequest("POST", resetUrl.String(), nil)
	if err != nil {
		return fmt.Errorf("error creating request: %v", err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
	req.Body = ioutil.NopCloser(bytes.NewBufferString(data.Encode()))
	req.Header.Set("Cookie", cookie)

	cl := &http.Client{}
	resp, err := cl.Do(req)
	if err != nil {
		return fmt.Errorf("error sending request: %v", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("error reading response: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		return utils.NewRequestError(resp, fmt.Sprintf("tadpoles reset password failed: %s", string(body)))
	}
	if string(body) != `{"status":"ok"}` {
		return fmt.Errorf("unexpected response body: %s", string(body))
	}
	return nil
}
