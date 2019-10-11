package mail

type Server struct {
	Host     string
	Port     int
	User     string
	Password string
}

type Config struct {
	Server
	Templates []Template
	From      string
}

type Template struct {
	Namespace string
	Source    []string
	Recursive bool
}
