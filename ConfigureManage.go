package nacos

import (
	"fmt"
	"github.com/levigross/grequests"
	"strings"
)

type ConfigServer struct {
	Url string `json:"url"`
	Group string `json:"group"`
}

func NewConfigServer(url,group string) *ConfigServer {
	if group == "" {
		group = "DEFAULT_GROUP"
	}
	if !strings.HasSuffix(url,"/") {
		url += "/"
	}
	return &ConfigServer{
		Url: url,
		Group: group,
	}
}

func (c *ConfigServer) GetConfig(configName string) string {
	confUrl := fmt.Sprintf("%sv1/cs/configs?group=%s&dataId=%s", c.Url, c.Group,configName)
	resp,_ := grequests.Get(confUrl,nil)
	return resp.String()
}