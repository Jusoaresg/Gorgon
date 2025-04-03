package templates

import (
	"github.com/a-h/templ"
	"github.com/labstack/echo/v4"
)

func Render(c echo.Context, component templ.Component) error {
	return component.Render(c.Request().Context(), c.Response())
}

// var templates map[string]*template.Template
// var componentTemplates map[string]*template.Template
//
// func LoadTemplates() {
// 	templates = make(map[string]*template.Template)
// 	componentTemplates = make(map[string]*template.Template)
//
// 	templateDirs := []string{
// 		"assets/templates/components/*.html",
// 		"assets/templates/layouts/*.html",
// 	}
//
// 	var layoutFiles []string
// 	for _, dir := range templateDirs {
// 		files, err := filepath.Glob(dir)
// 		if err != nil {
// 			fmt.Printf("Error while loading template %s: %v\n", dir, err)
// 			continue
// 		}
// 		layoutFiles = append(layoutFiles, files...)
// 	}
//
// 	pages, err := filepath.Glob("assets/templates/pages/*.html")
// 	if err != nil {
// 		panic(err)
// 	}
//
// 	for _, page := range pages {
// 		name := filepath.Base(page)
//
// 		tmpl := template.Must(template.New(name).ParseFiles(append(layoutFiles, page)...))
//
// 		templates[name[:len(name)-5]] = tmpl // Remove ".html" from key
// 	}
//
// 	allTemplates, err := filepath.Glob("assets/templates/**/*.html")
// 	if err != nil {
// 		panic(err)
// 	}
//
// 	// Criar um template global para os componentes
// 	globalTemplate := template.New("")
//
// 	for _, file := range allTemplates {
// 		_, err := globalTemplate.ParseFiles(file)
// 		if err != nil {
// 			fmt.Printf("Error parsing template %s: %v\n", file, err)
// 		}
// 	}
//
// 	// Adicionar cada componente ao mapa
// 	for _, file := range allTemplates {
// 		name := filepath.Base(file)
// 		componentTemplates[name[:len(name)-5]] = globalTemplate.Lookup(name)
// 		fmt.Println("Loaded component:", name[:len(name)-5]) // 🔍 Debug
// 	}
//
// }
//
// func Render(c echo.Context, templateName string, data gin.H) {
// 	if templates == nil {
// 		LoadTemplates()
// 	}
//
// 	template, exists := templates[templateName]
// 	if !exists {
// 		c.String(http.StatusInternalServerError, "Template not found: "+templateName)
// 		return
// 	}
//
// 	c.Response().Header().Set(echo.HeaderContentType, echo.MIMETextHTMLCharsetUTF8)
// 	if err := template.ExecuteTemplate(c.Response().Writer, "base", data); err != nil {
// 		c.String(http.StatusInternalServerError, "Template error: "+err.Error())
// 	}
// }
//
// func RenderComponent(c echo.Context, componentName string, data gin.H) {
// 	if componentTemplates == nil {
// 		LoadTemplates()
// 	}
//
// 	template, exists := componentTemplates[componentName]
// 	if !exists {
// 		c.String(http.StatusInternalServerError, "Component not found: "+componentName)
// 		return
// 	}
//
// 	c.Response().Header().Set(echo.HeaderContentType, echo.MIMETextHTMLCharsetUTF8)
//
// 	if err := template.Execute(c.Response().Writer, data); err != nil {
// 		c.String(http.StatusInternalServerError, "Component render error: "+err.Error())
// 		return
// 	}
// }
