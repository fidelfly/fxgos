package mail

import (
	"html/template"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/gomail.v2"
)

var mailDialer *Broker

type Broker struct {
	*gomail.Dialer
	From string
}

func (b *Broker) DialAndSend(m ...*gomail.Message) error {
	if len(b.From) > 0 {
		for _, msg := range m {
			if len(msg.GetHeader("From")) == 0 {
				msg.SetHeader("From", b.From)
			}
		}
	}
	return b.Dialer.DialAndSend(m...)
}

var mailTemplates map[string]*template.Template

//export
func InitMailDelegator(config Config) {
	mailDialer = &Broker{
		gomail.NewDialer(config.Host, config.Port, config.User, config.Password),
		config.From}
	initTemplate(config.Templates...)
}

func initTemplate(templates ...Template) {
	if mailTemplates == nil {
		mailTemplates = make(map[string]*template.Template)
	}
	for _, tempCfg := range templates {
		if len(tempCfg.Source) == 0 {
			continue
		}
		ns := "default"
		if len(tempCfg.Namespace) > 0 {
			ns = tempCfg.Namespace
		}
		temp, ok := mailTemplates[ns]
		if !ok {
			temp = template.New(ns)
			mailTemplates[ns] = temp
		}
		for _, s := range tempCfg.Source {
			if len(s) > 0 {
				if strings.HasSuffix(s, ".tpl") {
					addTemplate(temp, s)
				} else {
					addTemplateS(temp, tempCfg.Recursive, s)
				}
			}
		}
	}
}

func addTemplate(t *template.Template, files ...string) {
	tempFiles := make([]string, 0)
	for _, f := range files {
		if strings.HasSuffix(f, ".tpl") {
			if _, err := os.Stat(f); err == nil {
				tempFiles = append(tempFiles, f)
			}
		}
	}
	_, _ = t.ParseFiles(tempFiles...)
}

func addTemplateS(t *template.Template, recursive bool, sources ...string) {
	files := make([]string, 0)
	for _, s := range sources {
		if f, err := os.Stat(s); err == nil {
			if f.IsDir() {
				if recursive {
					_ = filepath.Walk(s, func(path string, info os.FileInfo, err error) error {
						if info.IsDir() {
							return nil
						}
						if strings.HasSuffix(info.Name(), ".tpl") {
							files = append(files, path)
						}
						return nil
					})
				} else {
					if fis, err := ioutil.ReadDir(s); err == nil {
						for _, info := range fis {
							if info.IsDir() {
								continue
							}
							if strings.HasSuffix(info.Name(), ".tpl") {
								files = append(files, filepath.Join(s, info.Name()))
							}
						}
					}
				}
			}
		}
	}
	if len(files) > 0 {
		_, _ = t.ParseFiles(files...)
	}
}

func GetTemplate(ns string, name string) *template.Template {
	if len(ns) == 0 {
		ns = "default"
	} else if len(name) == 0 {
		return mailTemplates[ns]
	}
	tempNS := mailTemplates[ns]
	if tempNS == nil {
		tempNS = template.New(ns)
		if fi, err := os.Stat(name); err == nil {
			if fi.IsDir() {
				addTemplateS(tempNS, false, name)
				mailTemplates[ns] = tempNS
				return tempNS
			} else {
				if strings.HasSuffix(name, ".tpl") {
					addTemplate(tempNS, name)
					mailTemplates[ns] = tempNS
					return tempNS.Lookup(filepath.Base(name))
				}
			}
		}
		return nil
	}

	t := tempNS.Lookup(name)
	if t != nil {
		return t
	}

	if fi, err := os.Stat(name); err == nil {
		if !fi.IsDir() {
			if strings.HasSuffix(name, ".tpl") {
				addTemplate(tempNS, name)
				return tempNS.Lookup(filepath.Base(name))
			}
		}
	}

	return nil
}
