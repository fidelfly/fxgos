module github.com/fidelfly/fxgos

go 1.12

require (
	github.com/BurntSushi/toml v0.3.1
	github.com/casbin/casbin v1.9.1
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/fidelfly/gostool v0.0.0
	github.com/fidelfly/gox v0.0.0
	github.com/go-sql-driver/mysql v1.4.1
	github.com/go-xorm/xorm v0.7.9
	github.com/golang/protobuf v1.3.2
	github.com/lib/pq v1.1.1 // indirect
	github.com/sirupsen/logrus v1.4.2
	github.com/tidwall/buntdb v1.0.0
	golang.org/x/crypto v0.0.0-20190605123033-f99c8df09eb5 // indirect
	golang.org/x/net v0.0.0-20190603091049-60506f45cf65 // indirect
	golang.org/x/sys v0.0.0-20190602015325-4c4f7f33c9ed // indirect
	golang.org/x/text v0.3.2 // indirect
	golang.org/x/tools v0.0.0-20190606050223-4d9ae51c2468 // indirect
	google.golang.org/grpc v1.19.0
	gopkg.in/gomail.v2 v2.0.0-20160411212932-81ebce5c23df
	xorm.io/core v0.7.2
)

replace (
	github.com/fidelfly/gostool => ./gospkg
	github.com/fidelfly/gox => ../gox
)
