package tools

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

type CertItem struct {
	Product        string `yaml:"product"`
	Domain         string `yaml:"domain"`
	PublicKeyPath  string `yaml:"publicKeyPath"`
	PrivateKeyPath string `yaml:"privateKeyPath"`
	Alias          string `yaml:"alias"`
}

type CvmItem struct {
	Region string `yaml:"region"`
}

type qcloudConfig struct {
	SecretId     string              `yaml:"secretId"`
	SecretKey    string              `yaml:"secretKey"`
	Certificates map[string]CertItem `yaml:"certificates"`
	Cvms         map[string]CvmItem  `yaml:"cvms"`
}

var config *qcloudConfig
var once sync.Once

func NewQcloudConfig(file string) *qcloudConfig {
	if file == "" {
		rootDir, _ := os.Executable()
		rootDir = filepath.Dir(rootDir)
		file = rootDir + "/../config/qcloud.yaml"
	}

	once.Do(func() {
		content, err := ioutil.ReadFile(file)
		if err != nil {
			fmt.Printf("failed to read yaml file : %v\n", err)
			panic(err)
		}

		err = yaml.Unmarshal(content, &config)
		if err != nil {
			panic(err)
		}
	})

	return config
}

func (config *qcloudConfig) GetCvmItem(group string) CvmItem {
	cvmItem, found := config.Cvms[group]
	if !found {
		panic("配置不存在")
	}
	return cvmItem
}

func (config *qcloudConfig) GetCertItem(group string) CertItem {
	certItem, found := config.Certificates[group]
	if !found {
		panic("配置不存在")
	}
	return certItem
}

func (config *qcloudConfig) GetCertRequestParam(group string) string {

	certItem, found := config.Certificates[group]
	if !found {
		return ""
	}

	return GetCertRequestParam(certItem.Domain, certItem.PrivateKeyPath, certItem.PublicKeyPath)
}

func GetCertRequestParam(domain, privateKeyPath, publicKeyPath string) string {

	publicData, _ := ioutil.ReadFile(publicKeyPath)
	privateData, _ := ioutil.ReadFile(privateKeyPath)

	publicKeyData := strings.TrimSpace(string(publicData))
	publicKeyData = strings.ReplaceAll(publicKeyData, "\n", "\\n")

	privateKeyData := strings.TrimSpace(string(privateData))
	privateKeyData = strings.ReplaceAll(privateKeyData, "\n", "\\n")

	params := "{\"Domain\":\"%s\",\"Https\":{\"Switch\":\"on\",\"Http2\":\"on\",\"CertInfo\":{\"Certificate\":\"%s\",\"PrivateKey\":\"%s\",\"Message\":\"%s\"}}}"
	params = fmt.Sprintf(params, domain, publicKeyData, privateKeyData, time.Now().Format("2006-01-02"))

	return params

}
