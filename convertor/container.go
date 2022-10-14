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

func NewContainerWithSuffix(suffix string, ignore bool) *Container {
	c := NewContainer()
	c.lineSuffix = suffix
	c.ignoreLastLineSuffix = ignore
	return c
}

type Container struct {
	list                 []Element
	lineSuffix           string
	ignoreLastLineSuffix bool
}

func (c *Container) Append(f Element) {
	c.list = append(c.list, f)
}

// Render 将container输出为string
func (c *Container) Render() string {
	var sb strings.Builder

	for i, e := range c.list {
		sb.WriteString(e.Format())
		if len(c.lineSuffix) > 0 {
			if c.ignoreLastLineSuffix && i == len(c.list)-1 {
				break
			}
			sb.WriteString(c.lineSuffix)
		}
	}
	return sb.String()
}
