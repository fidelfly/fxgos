module github.com/fidelfly/fxgos

go 1.12

require (
	github.com/BurntSushi/toml v0.3.1
	github.com/codegangsta/inject v0.0.0-20140425184007-37d7f8432a3e
	github.com/fidelfly/fxgo v0.0.0
	github.com/go-sql-driver/mysql v1.3.0
	github.com/go-xorm/builder v0.3.0
	github.com/go-xorm/core v0.6.0
	github.com/go-xorm/xorm v0.7.0
	github.com/gorilla/context v1.1.1
	github.com/gorilla/mux v1.7.2
	github.com/gorilla/websocket v1.4.0
	github.com/natefinch/lumberjack v2.0.0+incompatible
	github.com/patrickmn/go-cache v2.1.0+incompatible
	github.com/satori/go.uuid v1.2.0
	github.com/sirupsen/logrus v1.4.2
	github.com/tidwall/btree v0.0.0-20170113224114-9876f1454cf0
	github.com/tidwall/buntdb v1.0.0
	github.com/tidwall/gjson v1.1.3
	github.com/tidwall/grect v0.0.0-20161006141115-ba9a043346eb
	github.com/tidwall/match v1.0.1
	github.com/tidwall/rtree v0.0.0-20180113144539-6cd427091e0e
	github.com/tidwall/tinyqueue v0.0.0-20180302190814-1e39f5511563
	gopkg.in/oauth2.v3 v3.10.0
)

replace github.com/fidelfly/fxgo => ../fxgo
