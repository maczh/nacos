package nacos

import (
	"encoding/json"
	"fmt"
	"github.com/levigross/grequests"
	"math/rand"
	"time"
)

var clientInstance *Instance
var nacosServer *NameServer
var serviceCache = make(map[string][]Instance)

func NewNameServer(host string, port int) *NameServer {
	nacosServer = &NameServer{
		Host: host,
		Port: port,
	}
	return nacosServer
}

func GetNamingServer() *NameServer {
	return nacosServer
}

func (n *NameServer) GetNacosServerUrl() string {
	return fmt.Sprintf("http://%s:%d", n.Host, n.Port)
}

func NewInstance(serviceName, ip, groupName string, port int) *Instance {
	if groupName == "" {
		groupName = "DEFAULT_GROUP"
	}
	clientInstance = &Instance{
		Port:        port,
		Ip:          ip,
		Weight:      1.0,
		ServiceName: serviceName,
		GroupName:   groupName,
	}
	return clientInstance
}

func GetNamingInstance() *Instance {
	return clientInstance
}

func (instance *Instance) SetMetadataMap(meta map[string]string) *Instance {
	instance.Metadata = meta
	return instance
}

func (instance *Instance) SetMetadata(key, value string) *Instance {
	if instance.Metadata == nil {
		instance.Metadata = make(map[string]string)
	}
	instance.Metadata[key] = value
	return instance
}

func (instance *Instance) DelMetadata(key string) *Instance {
	if instance.Metadata == nil {
		return instance
	}
	delete(instance.Metadata, key)
	return instance
}

func (instance *Instance) GetMetadata(key string) string {
	if instance.Metadata == nil {
		return ""
	}
	return instance.Metadata[key]
}

func (instance *Instance) Register() error {
	server := GetNamingServer()
	if server == nil {
		return fmt.Errorf("nacos server unset")
	}
	logger.Debug(fmt.Sprintf("%s/nacos/v1/ns/instance", server.GetNacosServerUrl()))
	instance.Healthy = true
	data := toMap(*instance)
	j, _ := json.Marshal(data)
	logger.Debug(fmt.Sprintf("请求参数:%s", string(j)))
	resp, err := grequests.Post(fmt.Sprintf("%s/nacos/v1/ns/instance", server.GetNacosServerUrl()), &grequests.RequestOptions{
		Data: toMap(*instance),
	})
	if err != nil {
		logger.Error("register to nacos server failed: " + err.Error())
		return err
	}
	logger.Debug("Nacos server register response: " + resp.String())
	switch resp.StatusCode {
	case 200:
		return nil
	case 400:
		return fmt.Errorf("Bad Request")
	case 403:
		return fmt.Errorf("Forbidden")
	case 404:
		return fmt.Errorf("Not Found")
	case 500:
		return fmt.Errorf("Internal Server Error")
	}
	return nil
}

func (instance *Instance) Unregister() error {
	server := GetNamingServer()
	if server == nil {
		return fmt.Errorf("nacos server unset")
	}
	data := toMap(*instance)
	j, _ := json.Marshal(data)
	logger.Debug(fmt.Sprintf("请求参数:%s", string(j)))
	resp, err := grequests.Delete(fmt.Sprintf("%s/nacos/v1/ns/instance", server.GetNacosServerUrl()), &grequests.RequestOptions{
		Data: data,
	})
	if err != nil {
		logger.Error("register to nacos server failed: " + err.Error())
		return err
	}
	logger.Debug("Nacos server register response: " + resp.String())
	switch resp.StatusCode {
	case 200:
		return nil
	case 400:
		return fmt.Errorf("Bad Request")
	case 403:
		return fmt.Errorf("Forbidden")
	case 404:
		return fmt.Errorf("Not Found")
	case 500:
		return fmt.Errorf("Internal Server Error")
	}
	return nil
}

func (instance *Instance) queryServiceInstances(serviceName, groupName string) error {
	server := GetNamingServer()
	if server == nil {
		return fmt.Errorf("nacos server unset")
	}
	if groupName == "" {
		groupName = instance.GroupName
	}
	queryParams := SelectInstancesParam{
		ServiceName: serviceName,
		GroupName:   groupName,
		HealthyOnly: true,
	}
	resp, err := grequests.Get(fmt.Sprintf("%s/nacos/v1/ns/instance/list", server.GetNacosServerUrl()), &grequests.RequestOptions{
		Params: toMap(queryParams),
	})
	if err != nil {
		logger.Error("query nacos server service instance failed: " + err.Error())
		return err
	}
	if resp.StatusCode == 200 {
		var response SelectInstancesResponse
		err = json.Unmarshal(resp.Bytes(), &response)
		if err != nil {
			logger.Error("json unmarshal failed: " + err.Error())
			return err
		}
		if len(response.Hosts) > 0 {
			serviceCache[serviceName] = response.Hosts
		} else {
			if groupName != "DEFAULT_GROUP" {
				return instance.queryServiceInstances(serviceName, "DEFAULT_GROUP")
			}
		}
		return nil
	} else {
		return fmt.Errorf("nacos server return code %d", resp.StatusCode)
	}
}

func (instance *Instance) ListServiceInstances(serviceName, groupName string) []Instance {
	if serviceCache[serviceName] == nil {
		instance.queryServiceInstances(serviceName, groupName)
	}
	return serviceCache[serviceName]
}

func (instance *Instance) GetServiceInstance(serviceName string) Instance {
	instances := instance.ListServiceInstances(serviceName, "")
	if instances == nil || len(instances) == 0 {
		return Instance{}
	}
	if len(instances) == 1 {
		return instances[0]
	}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	i := 0
	l := len(instances)
	for {
		i = r.Intn(l)
		if instances[i].Healthy == true && instances[i].Metadata != nil && instances[i].Metadata["debug"] != "true" {
			break
		}
	}
	return instances[i]
}

func (instance *Instance) GetServiceUrl(serviceName string) string {
	ins := instance.GetServiceInstance(serviceName)
	if ins.Ip == "" {
		return ""
	}
	if ins.Metadata != nil && ins.Metadata["ssl"] == "true" {
		return fmt.Sprintf("https://%s:%d", ins.Ip, ins.Port)
	}
	return fmt.Sprintf("http://%s:%d", ins.Ip, ins.Port)
}

func (instance *Instance) HealthBeat() error {
	server := GetNamingServer()
	if server == nil {
		return fmt.Errorf("nacos server unset")
	}
	beatInfo := BeatInfo{
		Ip:          instance.Ip,
		Port:        instance.Port,
		Weight:      instance.Weight,
		ServiceName: fmt.Sprintf("%s@@%s", instance.GroupName, instance.ServiceName),
		Metadata:    instance.Metadata,
		Cluster:     "DEFAULT",
	}
	beatJson, _ := json.Marshal(beatInfo)
	beatParam := BeatParam{
		ServiceName: fmt.Sprintf("%s@@%s", instance.GroupName, instance.ServiceName),
		Beat:        string(beatJson),
	}
	resp, err := grequests.Put(fmt.Sprintf("%s/nacos/v1/ns/instance/beat", server.GetNacosServerUrl()), &grequests.RequestOptions{
		Data: toMap(beatParam),
	})
	if err != nil {
		logger.Error("Instance send health beat nacos server service instance failed: " + err.Error())
		return err
	}
	switch resp.StatusCode {
	case 200:
		return nil
	case 400:
		return fmt.Errorf("Bad Request")
	case 403:
		return fmt.Errorf("Forbidden")
	case 404:
		return fmt.Errorf("Not Found")
	case 500:
		return fmt.Errorf("Internal Server Error")
	}
	return nil
}

func (instance *Instance) refreshCache() {
	if serviceCache == nil || len(serviceCache) == 0 {
		return
	}
	for serviceName, cache := range serviceCache {
		delete(serviceCache, serviceName)
		instance.queryServiceInstances(serviceName, cache[0].GroupName)
	}
}

func Init(serverIp string, serverPort int, serviceName, groupName, serviceIp string, servicePort int, meta map[string]string) {
	NewNameServer(serverIp, serverPort)
	instance := NewInstance(serviceName, serviceIp, groupName, servicePort)
	if meta != nil {
		instance = instance.SetMetadataMap(meta)
	}
	instance.Register()
	//设置定时任务--心跳
	tickerBeat := time.NewTicker(time.Second * 20)
	go func() {
		for _ = range tickerBeat.C {
			instance.HealthBeat()
		}
	}()

	//设置定时任务--刷新缓存
	tickerRefreshCache := time.NewTicker(time.Minute * 5)
	go func() {
		for _ = range tickerRefreshCache.C {
			instance.refreshCache()
		}
	}()
}
