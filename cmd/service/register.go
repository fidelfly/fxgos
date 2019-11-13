package service

var serviceMap = make(map[string]interface{})

func Register(name string, instance interface{}) {
	serviceMap[name] = instance
}

func GetService(name string) (interface{}, bool) {
	if instance, ok := serviceMap[name]; ok {
		return instance, ok
	}
	return nil, false
}
