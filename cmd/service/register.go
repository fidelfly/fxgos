package service

var serviceMap = make(map[string]interface{})

var deps = make([][2]string, 0)

type Service interface {
	Start() error
}

func Register(name string, instance interface{}, dependencies ...string) {
	serviceMap[name] = instance
	if len(dependencies) > 0 {
		for _, dep := range dependencies {
			deps = append(deps, [2]string{name, dep})
		}
	}
}

func GetService(name string) (interface{}, bool) {
	if instance, ok := serviceMap[name]; ok {
		return instance, ok
	}
	return nil, false
}

func Start() error {
	srvs := getServices()
	for _, srv := range srvs {
		server := serviceMap[srv]
		if myService, ok := server.(Service); ok {
			if err := myService.Start(); err != nil {
				return err
			}
		}
	}
	return nil
}

func getServices() []string {
	var srvs = make([]string, len(serviceMap))
	var j = 0
	for k := range serviceMap {
		srvs[j] = k
		j++
	}
	return sortService(srvs)
}

func sortService(srvs []string) []string {
	var result = make([]string, 0)
	var relates = deps[:]
	for {
		var tempO = make([]string, 0)
		var tempS = make([]string, 0)
		var tempR = make([][2]string, 0)
		for _, srv := range srvs {
			var find = false
			for _, r := range relates {
				if r[0] == srv {
					find = true
					break
				}
			}
			if find {
				tempS = append(tempS, srv)
			} else {
				tempO = append(tempO, srv)
			}
		}
		if len(tempO) == 0 {
			panic("there is service dependency loop")
		}
		result = append(result, tempO...)

		if len(tempS) == 0 {
			break
		}

		srvs = tempS
		relates = tempR
	}
	return result
}
