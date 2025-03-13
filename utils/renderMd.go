package forum

import (
	"html/template"
        "github.com/yuin/goldmark"
        "bytes"
)

func RenderMarkdown(md string) template.HTML {
	var buf bytes.Buffer
	if err := goldmark.Convert([]byte(md), &buf); err != nil {
		return template.HTML(md) 
	}
	return template.HTML(buf.String())
}
