package views

import (
	"bytes"
	"fmt"
	"html/template"
	"io"
	"strings"

	episodeModel "github.com/jusoaresg/gorgon/internal/episode/model"
	"github.com/labstack/echo/v4"
)

type Template struct {
	templates *template.Template
}

type PageData struct {
	TemplateName string
	Data         any
	Styles       []string
}

func (t *Template) Render(w io.Writer, name string, data any, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

func NewTemplate() *Template {
	tmpl := template.New("")

	tmpl.Funcs(template.FuncMap{
		"toLower": func(text string) string {
			return strings.ToLower(text)
		},
		"render": func(name string, data any) (template.HTML, error) {
			var buf bytes.Buffer
			err := tmpl.ExecuteTemplate(&buf, name, data)
			return template.HTML(buf.String()), err
		},
		"safeHtml": func(s string) template.HTML {
			return template.HTML(s)
		},
		"derefFloat": func(f *float64) float64 {
			if f == nil {
				return 0
			}
			return *f
		},
		"countDownloaded": func(episodes []episodeModel.Episode) int {
			count := 0
			for _, ep := range episodes {
				if ep.Tracking == episodeModel.TrackingDownloaded {
					count++
				}
			}
			return count
		},
		"calcPercentage": func(part, total int) int {
			if total == 0 {
				return 0
			}
			return int((float64(part) / float64(total)) * 100)
		},
	})

	tmpl, err := tmpl.ParseFS(FrontFS, "ui/*.html", "ui/components/*.html", "ui/components/settings/*.html")
	if err != nil {
		panic(fmt.Errorf("error while loading embedded front %w", err))
	}

	return &Template{
		templates: tmpl,
	}
}
