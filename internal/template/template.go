package template

import (
	"bytes"
	"check-id-api/internal/logger"
	"text/template"
)

func GenerateTemplateMail(param map[string]string) (string, error) {
	bf := &bytes.Buffer{}
	tpl := &template.Template{}

	tpl = template.Must(template.New("").ParseGlob("templates/*.gohtml"))
	err := tpl.ExecuteTemplate(bf, param["TEMPLATE-PATH"], &param)
	if err != nil {
		logger.Error.Printf("couldn't generate template body mail: %v", err)
		return "", err
	}
	return bf.String(), nil
}
