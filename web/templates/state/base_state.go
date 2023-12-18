package state

import (
	"github.com/a-h/templ"
	"github.com/guilycst/guigoes/internal/core/domain"
)

type State struct {
	Language string
}

type BaseState struct {
	State
	Title string
	Body  templ.Component
}

type IndexState struct {
	State
	Posts []*domain.Post
}
