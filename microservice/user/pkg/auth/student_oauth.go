package auth

import (
	"crypto/sha1"
	"encoding/hex"
	"errors"
	"net/url"
	"regexp"
	"strings"

	"forum-user/util"

	"github.com/spf13/viper"
)

const (
	StudentLoginProviderLegacy = "legacy"
	StudentLoginProviderOAuth  = "oauth"

	defaultStudentOAuthCASLoginURL         = "https://account.ccnu.edu.cn/cas/login"
	defaultStudentOAuthCASCallbackURL      = "https://pass.muxixyz.com/auth/api/oauth/cas/callback"
	defaultStudentOAuthTokenPath           = "/auth/api/oauth/token"
	defaultStudentOAuthUserInfoPath        = "/auth/api/user"
	defaultStudentOAuthBusinessCallbackURL = "http://localhost:3000/login/student-oauth"
)

type StudentOAuthConfig struct {
	ClientID            string
	ClientSecret        string
	CASLoginURL         string
	OAuthCASCallbackURL string
	TokenURL            string
	UserInfoURL         string
	BusinessCallbackURL string
}

type StudentOAuthUserInfo struct {
	StudentID string `json:"student_id"`
	Email     string `json:"email"`
	Username  string `json:"username"`
	Name      string `json:"name"`
	AvatarURL string `json:"avatar_url"`
}

type studentOAuthTokenResponse struct {
	util.BasicResponse
	Data TokenItem `json:"data"`
}

type studentOAuthUserInfoResponse struct {
	util.BasicResponse
	Data StudentOAuthUserInfo `json:"data"`
}

func LoadStudentOAuthConfig() StudentOAuthConfig {
	host := strings.TrimSpace(viper.GetString("auth_server.host"))
	basicURL := ""
	if host != "" {
		basicURL = withHTTP(host)
	}

	tokenURL := strings.TrimSpace(viper.GetString("student_oauth.token_url"))
	if tokenURL == "" && basicURL != "" {
		tokenURL = basicURL + defaultStudentOAuthTokenPath
	}

	userInfoURL := strings.TrimSpace(viper.GetString("student_oauth.userinfo_url"))
	if userInfoURL == "" && basicURL != "" {
		userInfoURL = basicURL + defaultStudentOAuthUserInfoPath
	}

	return StudentOAuthConfig{
		ClientID:            strings.TrimSpace(viper.GetString("student_oauth.client_id")),
		ClientSecret:        strings.TrimSpace(viper.GetString("student_oauth.client_secret")),
		CASLoginURL:         stringWithDefault(viper.GetString("student_oauth.cas_login_url"), defaultStudentOAuthCASLoginURL),
		OAuthCASCallbackURL: stringWithDefault(viper.GetString("student_oauth.oauth_cas_callback_url"), defaultStudentOAuthCASCallbackURL),
		TokenURL:            tokenURL,
		UserInfoURL:         userInfoURL,
		BusinessCallbackURL: stringWithDefault(viper.GetString("student_oauth.business_callback_url"), defaultStudentOAuthBusinessCallbackURL),
	}
}

func BuildStudentOAuthLoginURL(cfg StudentOAuthConfig, callbackURL string) (string, error) {
	if cfg.ClientID == "" {
		return "", errors.New("student oauth client_id is blank")
	}
	if cfg.CASLoginURL == "" {
		return "", errors.New("student oauth cas_login_url is blank")
	}
	if cfg.OAuthCASCallbackURL == "" {
		return "", errors.New("student oauth cas callback url is blank")
	}

	if strings.TrimSpace(callbackURL) == "" {
		callbackURL = cfg.BusinessCallbackURL
	}

	inner := url.Values{}
	inner.Set("callback_url", callbackURL)
	inner.Set("client_id", cfg.ClientID)
	serviceURL := cfg.OAuthCASCallbackURL + "?" + inner.Encode()

	outer := url.Values{}
	outer.Set("service", serviceURL)
	return cfg.CASLoginURL + "?" + outer.Encode(), nil
}

func ExchangeStudentOAuthCode(cfg StudentOAuthConfig, code, callbackURL string) (*StudentOAuthUserInfo, error) {
	if cfg.ClientID == "" || cfg.ClientSecret == "" {
		return nil, errors.New("student oauth client info is blank")
	}
	if cfg.TokenURL == "" {
		return nil, errors.New("student oauth token_url is blank")
	}
	if cfg.UserInfoURL == "" {
		return nil, errors.New("student oauth userinfo_url is blank")
	}
	if strings.TrimSpace(callbackURL) == "" {
		callbackURL = cfg.BusinessCallbackURL
	}

	token, err := getStudentOAuthToken(code, cfg.ClientID, cfg.ClientSecret, cfg.TokenURL, callbackURL)
	if err != nil {
		return nil, err
	}

	return getStudentOAuthUserInfo(token.AccessToken, cfg.UserInfoURL)
}

func getStudentOAuthToken(code, clientID, clientSecret, tokenURL, callbackURL string) (*TokenItem, error) {
	query := map[string]string{
		"client_id":     clientID,
		"response_type": "token",
		"grant_type":    "authorization_code",
	}
	bodyData := map[string]string{
		"code":          code,
		"client_secret": clientSecret,
		"redirect_uri":  callbackURL,
	}

	body, err := util.SendHTTPRequest(tokenURL, "POST", &util.RequestData{
		Query:       query,
		BodyData:    bodyData,
		ContentType: util.FormData,
	})
	if err != nil {
		return nil, err
	}

	var rp studentOAuthTokenResponse
	if err := util.UnmarshalBodyForCustomData(body, &rp); err != nil {
		return nil, err
	}
	if rp.Code != 0 {
		return nil, errors.New(rp.Message)
	}
	if rp.Data.AccessToken == "" {
		return nil, errors.New("student oauth access token is blank")
	}
	return &rp.Data, nil
}

func getStudentOAuthUserInfo(accessToken, userInfoURL string) (*StudentOAuthUserInfo, error) {
	body, err := util.SendHTTPRequest(userInfoURL, "GET", &util.RequestData{
		Header: map[string]string{"token": accessToken},
	})
	if err != nil {
		return nil, err
	}

	var rp studentOAuthUserInfoResponse
	if err := util.UnmarshalBodyForCustomData(body, &rp); err != nil {
		return nil, err
	}
	if rp.Code != 0 {
		return nil, errors.New(rp.Message)
	}

	rp.Data.StudentID = strings.TrimSpace(rp.Data.StudentID)
	rp.Data.Email = strings.TrimSpace(rp.Data.Email)
	rp.Data.Username = strings.TrimSpace(rp.Data.Username)
	rp.Data.Name = strings.TrimSpace(rp.Data.Name)
	if rp.Data.StudentID == "" {
		if studentID, ok := studentIDFromCASUsername(rp.Data.Username); ok {
			rp.Data.StudentID = studentID
		}
	}
	if rp.Data.Username == "" {
		rp.Data.Username = rp.Data.Name
	}
	if rp.Data.StudentID == "" && rp.Data.Email == "" {
		return nil, errors.New("student oauth user info missing student_id and email")
	}
	return &rp.Data, nil
}

// casUsernamePattern 对应通行证cas影子账号的命名规则。
var casUsernamePattern = regexp.MustCompile(`^cas_(.+)_([0-9a-f]{8})$`)

func studentIDFromCASUsername(username string) (string, bool) {
	matches := casUsernamePattern.FindStringSubmatch(username)
	if matches == nil {
		return "", false
	}

	base, hash := matches[1], matches[2]
	if base == "" || base == "user" {
		// 上游对无法规范化的用户名兜底成 "user"，原值已不可还原。
		return "", false
	}

	for _, candidate := range []string{base, strings.ToUpper(base)} {
		sum := sha1.Sum([]byte(candidate))
		if hex.EncodeToString(sum[:])[:8] == hash {
			return candidate, true
		}
	}

	return "", false
}

func stringWithDefault(value, fallback string) string {
	value = strings.TrimSpace(value)
	if value == "" {
		return fallback
	}
	return value
}
