package nacos

type NameServer struct {
	Host string `json:"host"`
	Port int    `json:"port"`
}

type Instance struct {
	InstanceId  string            `json:"instanceId"`
	Port        int               `json:"port"`
	Ip          string            `json:"ip"`
	Weight      float64           `json:"weight"`
	Metadata    map[string]string `json:"metadata"`
	ClusterName string            `json:"clusterName"`
	ServiceName string            `json:"serviceName"`
	GroupName   string            `json:"groupName"` //optional,default:DEFAULT_GROUP
	Healthy     bool              `json:"healthy"`
}

type RegisterInstanceParam struct {
	Ip          string            `json:"ip"`          //required
	Port        int               `json:"port"`        //required
	Weight      float64           `json:"weight"`      //required,it must be lager than 0
	Enable      bool              `json:"enabled"`     //required,the instance can be access or not
	Healthy     bool              `json:"healthy"`     //required,the instance is health or not
	Metadata    map[string]string `json:"metadata"`    //optional
	ClusterName string            `json:"clusterName"` //optional,default:DEFAULT
	ServiceName string            `json:"serviceName"` //required
	GroupName   string            `json:"groupName"`   //optional,default:DEFAULT_GROUP
}

type DeregisterInstanceParam struct {
	Ip          string `json:"ip"`          //required
	Port        int    `json:"port"`        //required
	Cluster     string `json:"cluster"`     //optional,default:DEFAULT
	ServiceName string `json:"serviceName"` //required
	GroupName   string `json:"groupName"`   //optional,default:DEFAULT_GROUP
	Ephemeral   bool   `json:"ephemeral"`   //optional
}

type UpdateInstanceParam struct {
	Ip          string            `json:"ip"`          // required
	Port        int               `json:"port"`        // required
	ClusterName string            `json:"cluster"`     // optional,default:DEFAULT
	ServiceName string            `json:"serviceName"` // required
	GroupName   string            `json:"groupName"`   // optional,default:DEFAULT_GROUP
	Ephemeral   bool              `json:"ephemeral"`   // optional
	Weight      float64           `json:"weight"`      // required,it must be lager than 0
	Enable      bool              `json:"enabled"`     // required,the instance can be access or not
	Metadata    map[string]string `json:"metadata"`    // optional
}

type GetServiceParam struct {
	Clusters    []string `json:"clusters"`    //optional,default:DEFAULT
	ServiceName string   `json:"serviceName"` //required
	GroupName   string   `json:"groupName"`   //optional,default:DEFAULT_GROUP
}

type GetAllServiceInfoParam struct {
	NameSpace string `json:"nameSpace"` //optional,default:public
	GroupName string `json:"groupName"` //optional,default:DEFAULT_GROUP
	PageNo    uint32 `json:"pageNo"`    //optional,default:1
	PageSize  uint32 `json:"pageSize"`  //optional,default:10
}

type SelectAllInstancesParam struct {
	Clusters    []string `json:"clusters"`    //optional,default:DEFAULT
	ServiceName string   `json:"serviceName"` //required
	GroupName   string   `json:"groupName"`   //optional,default:DEFAULT_GROUP
}

type SelectInstancesParam struct {
	ServiceName string `json:"serviceName"` //required
	GroupName   string `json:"groupName"`   //optional,default:DEFAULT_GROUP
	HealthyOnly bool   `json:"healthyOnly"` //optional,return only healthy instance
}

type SelectOneHealthInstanceParam struct {
	Clusters    []string `json:"clusters"`    //optional,default:DEFAULT
	ServiceName string   `json:"serviceName"` //required
	GroupName   string   `json:"groupName"`   //optional,default:DEFAULT_GROUP
}

type SelectInstancesResponse struct {
	Name                     string     `json:"name"`
	GroupName                string     `json:"groupName"`
	Clusters                 string     `json:"clusters"`
	CacheMillis              int        `json:"cacheMillis"`
	Hosts                    []Instance `json:"hosts"`
	LastRefTime              int64      `json:"lastRefTime"`
	Checksum                 string     `json:"checksum"`
	AllIPs                   bool       `json:"allIPs"`
	ReachProtectionThreshold bool       `json:"reachProtectionThreshold"`
	Valid                    bool       `json:"valid"`
}

type SubscribeService struct {
	ClusterName string            `json:"clusterName"`
	Enable      bool              `json:"enable"`
	InstanceId  string            `json:"instanceId"`
	Ip          string            `json:"ip"`
	Metadata    map[string]string `json:"metadata"`
	Port        int               `json:"port"`
	ServiceName string            `json:"serviceName"`
	Valid       bool              `json:"valid"`
	Weight      float64           `json:"weight"`
}

type BeatInfo struct {
	Ip          string            `json:"ip"`
	Port        int               `json:"port"`
	Weight      float64           `json:"weight"`
	ServiceName string            `json:"serviceName"`
	Cluster     string            `json:"cluster"`
	Metadata    map[string]string `json:"metadata"`
}

type BeatParam struct {
	ServiceName string `json:"serviceName"`
	Beat        string `json:"beat"`
}
