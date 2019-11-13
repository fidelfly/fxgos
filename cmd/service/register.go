package service

var serviceMap = make(map[string]interface{})

func Register(name string, instance interface{}) error {
	serviceMap[name] = instance
	return nil
}

func GetService(name string) (interface{}, bool) {
	if instance, ok := serviceMap[name]; ok {
		return instance, ok
	}
	return nil, false
}
