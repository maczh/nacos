package nacos

import (
	"encoding/json"
	"fmt"
	"github.com/sadlil/gologger"
	"reflect"
)

var logger = gologger.GetLogger()

func toMap(obj interface{}) map[string]string {
	refVal := reflect.ValueOf(obj)
	refType := reflect.TypeOf(obj)
	m := make(map[string]string)
	for i := 0; i < refVal.NumField(); i++ {
		key := refType.Field(i).Tag.Get("json")
		if refVal.Field(i).Interface() == nil {
			continue
		}
		field := refVal.Field(i)
		valueType := field.Kind()
		switch valueType {
		case reflect.String:
			m[key] = field.String()
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Uint, reflect.Uint32, reflect.Uint64:
			m[key] = fmt.Sprintf("%d", field.Int())
		case reflect.Float32, reflect.Float64, reflect.Complex64, reflect.Complex128:
			m[key] = fmt.Sprintf("%f", field.Float())
		case reflect.Bool:
			m[key] = fmt.Sprintf("%v", field.Bool())
		case reflect.Array, reflect.Map, reflect.Struct, reflect.Slice:
			j, err := json.Marshal(field.Interface())
			if err != nil {
				logger.Error("json Marshal failed: " + err.Error())
				continue
			}
			m[key] = string(j)
		}
	}
	//js, _ := json.Marshal(m)
	//logger.Debug("转换结果: " + string(js))
	return m
}
