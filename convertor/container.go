package convertor

import (
	"strings"
)

type Element interface {
	Format() string
}

func NewContainer() *Container {
	return &Container{}
}

type Container struct {
	list []Element
}

func (c *Container) Append(f Element) {
	c.list = append(c.list, f)
}

// Render 将container输出为string
func (c *Container) Render(lineSuffix string) string {
	var sb strings.Builder

	for _, e := range c.list {
		sb.WriteString(e.Format())
		if len(lineSuffix) > 0 {
			sb.WriteString(lineSuffix)
		}
	}
	return sb.String()
}
