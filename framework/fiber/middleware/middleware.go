package middleware

import goauth "github.com/SwanHtetAungPhyo/go-auth"

type Maker struct {
	cfg goauth.Config
}

func NewMaker(cfg goauth.Config) *Maker {
	return &Maker{cfg: cfg}
}
