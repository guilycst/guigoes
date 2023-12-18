package state

import "github.com/a-h/templ"

type BaseState struct {
	Title string
	Body  templ.Component
}
