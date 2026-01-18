package config

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/vo"
	"github.com/spf13/viper"
	"log"
	"net/url"
	"os"
	"strconv"
	"strings"
)

const NACOS_ENV = "NACOSDSN"

type Config struct {
	Name string
}

func Init(cfg string, prefix string) error {
	c := Config{
		Name: cfg,
	}

	// 初始化配置文件
	if err := c.initConfig(prefix); err != nil {
		return err
	}

	return nil
}

func (c *Config) initConfig(prefix string) error {
	viper.SetConfigType("yaml")
	viper.AutomaticEnv()
	viper.SetEnvPrefix(prefix)

	// 尝试从 Nacos 拉取配置
	content, err := GetConfigFromNacos(NACOS_ENV)
	if err != nil {
		log.Println("降级到本地配置文件拉取...")
		if c.Name != "" {
			viper.SetConfigFile(c.Name)
		} else {
			viper.AddConfigPath("./conf")
			viper.SetConfigName("config")
		}

		// viper 解析本地配置文件
		if err := viper.ReadInConfig(); err != nil {
			return fmt.Errorf("本地配置加载失败: %w", err)
		}
		return nil
	}

	// 尝试加载到 viper
	if err := viper.ReadConfig(bytes.NewBufferString(content)); err != nil {
		log.Printf("成功拉取 Nacos 配置，但加载到 viper 失败: %v", err)
	}

	return nil
}

func GetConfigFromNacos(env string) (string, error) {
	dsn := os.Getenv(env)
	if dsn == "" {
		log.Printf("%s 环境变量未设置", env)
		return "", errors.New("环境变量未设置")
	}

	server, port, namespace, user, pass, group, dataId := ParseNacosDSN(dsn)

	serverConfigs := []constant.ServerConfig{
		{
			IpAddr: server,
			Port:   port,
			Scheme: "http",
		},
	}

	clientConfig := constant.ClientConfig{
		NamespaceId:         namespace,
		Username:            user,
		Password:            pass,
		TimeoutMs:           5000,
		NotLoadCacheAtStart: true,
		CacheDir:            "/tmp/nacos/cache",
		LogDir:              "/tmp/nacos/log",
	}

	configClient, err := clients.CreateConfigClient(map[string]interface{}{
		"serverConfigs": serverConfigs,
		"clientConfig":  clientConfig,
	})
	if err != nil {
		log.Fatal("初始化失败:", err)
	}

	content, err := configClient.GetConfig(vo.ConfigParam{
		DataId: dataId,
		Group:  group,
	})
	if err != nil {
		log.Fatal("拉取配置失败:", err)
	}
	return content, nil
}

func ParseNacosDSN(dsn string) (server string, port uint64, ns, user, pass, group, dataId string) {

	parts := strings.SplitN(dsn, "?", 2)
	host := parts[0]
	params := url.Values{}

	if len(parts) == 2 {
		params, _ = url.ParseQuery(parts[1])
	}

	hostParts := strings.Split(host, ":")
	server = hostParts[0]
	if len(hostParts) > 1 {
		p, _ := strconv.Atoi(hostParts[1])
		port = uint64(p)
	} else {
		port = 8848
	}

	ns = params.Get("namespace")
	if ns == "" {
		ns = "public"
	}

	user = params.Get("username")
	pass = params.Get("password")
	group = params.Get("group")
	dataId = params.Get("dataId")
	return
}
