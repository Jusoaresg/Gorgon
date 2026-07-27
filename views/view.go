package views

type View struct {
	Layout    string
	Default   string
	Templates map[string]string
	Styles    []string
	Data      any
}
