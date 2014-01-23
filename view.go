package views

import z "github.com/zaiuz/zaiuz"

type View interface {
	Subview(filenames ...string) View
	Render(c *z.Context, data interface{}) error
}
