package accesslog

import "strings"

type Directory map[string]*Entry

// Reference - configuration of logging entries
type Reference struct {
	Operator string
	Name     string
}

type Entry struct {
	Operator    string
	Name        string
	Value       string
	StringValue bool
}

func (e Entry) IsHeader() bool {
	return strings.HasPrefix(e.Operator, headerPrefix)
}

func (e Entry) IsDirect() bool {
	return e.Operator == directOperator
}

func NewEntry(operator, name, value string, stringValue bool) Entry {
	return Entry{Operator: operator, Name: name, Value: value, StringValue: stringValue}
}
