package spell

import (
	"errors"
	"github.com/evanw/esbuild/pkg/api"
)

func (e *Engine) BuildAssets() error {
	r := api.Build(e.EsbuildOptions)

	if len(r.Errors) > 0 {
		return errors.New(r.Errors[0].Text)
	}

	return nil
}
