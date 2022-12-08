package accesslog

import "strings"

type Directory map[string]*Entry

// Reference - configuration of logging entries
type Reference struct {
	Operator string
	Name     string
}
type Entry struct {
	Ref         Reference
	Value       string
	StringValue bool
}

func (e Entry) IsHeader() bool {
	return strings.HasPrefix(e.Ref.Operator, headerPrefix)
}

func (e Entry) IsDirect() bool {
	return e.Ref.Operator == directOperator
}

func (e Entry) Operator() string {
	return e.Ref.Operator
}

func (e Entry) Name() string {
	return e.Ref.Name
}

func NewEntry(operator, name, value string, stringValue bool) Entry {
	return Entry{Ref: Reference{Operator: operator, Name: name}, Value: value, StringValue: stringValue}
}
