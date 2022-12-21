package automation

import (
	"errors"
	"fmt"
	"golang.org/x/time/rate"
	"strconv"
)

type Attribute interface {
	Name() string
	Value() any
	String() string
	Validate() error
}

type attribute struct {
	name  string
	value any
}

func NewAttribute(name string, value any) Attribute {
	return &attribute{name: name, value: value}
}

func nilAttribute(name string) Attribute {
	return NewAttribute(name, nil)
}

func (a *attribute) Name() string {
	return a.name
}

func (a *attribute) Value() any {
	return a.value
}

func (a *attribute) String() string {
	if a.Value() == nil {
		return "nil"
	}
	if val, ok := a.Value().(int); ok {
		return strconv.Itoa(val)
	}
	if val, ok := a.Value().(bool); ok {
		return strconv.FormatBool(val)
	}
	if val, ok := a.Value().(rate.Limit); ok {
		if val == rate.Inf {
			return InfValue
		}
		return fmt.Sprintf("%v", val)
	}
	if val, ok := a.Value().(string); ok {
		return val
	}
	return "nil"
}

func (a *attribute) Validate() error {
	if a.Name() == "" {
		return errors.New("invalid attribute name : name is empty")
	}
	if a.Value() == nil {
		return errors.New(fmt.Sprintf("invalid attribute value: value is nil for [%v]", a.Name()))
	}
	return nil
}
