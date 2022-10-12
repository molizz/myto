package convertor

import (
	"strings"
)

type Element interface {
	Format() string
	AppendClient(e Element)
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
func (c *Container) Render() string {
	var sb strings.Builder

	for _, e := range c.list {
		sb.WriteString(e.Format())
	}
	return sb.String()
}

type defaultElement struct {
}

func (d *defaultElement) AppendClient(e Element) {}
