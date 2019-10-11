module github.com/fidelfly/fxgos

go 1.12

require (
	github.com/fidelfly/fxgo v0.0.0
	github.com/go-sql-driver/mysql v1.4.1
	github.com/go-xorm/xorm v0.7.6
	github.com/lib/pq v1.1.1 // indirect
	github.com/sirupsen/logrus v1.4.2
	gopkg.in/alexcesaro/quotedprintable.v3 v3.0.0-20150716171945-2caba252f4dc // indirect
	gopkg.in/gomail.v2 v2.0.0-20160411212932-81ebce5c23df
	xorm.io/core v0.7.2
)

replace github.com/fidelfly/fxgo => ../fxgo
