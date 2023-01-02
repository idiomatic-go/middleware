package accesslog

import (
	"errors"
	"fmt"
	"github.com/idiomatic-go/middleware/accessdata"
	"strings"
)

var ingressEntries []accessdata.Operator
var egressEntries []accessdata.Operator

func CreateIngressEntries(config []accessdata.Operator) error {
	ingressEntries = []accessdata.Operator{}
	return CreateEntries(&ingressEntries, config)
}

func CreateEgressEntries(config []accessdata.Operator) error {
	egressEntries = []accessdata.Operator{}
	return CreateEntries(&egressEntries, config)
}

func CreateEntries(items *[]accessdata.Operator, config []accessdata.Operator) error {
	if items == nil {
		return errors.New("invalid configuration : entries are nil")
	}
	if len(config) == 0 {
		return errors.New("invalid configuration : configuration is empty")
	}
	dup := make(map[string]string)
	for _, op := range config {
		op2, err := createOperator(op)
		if err != nil {
			return err
		}
		if IsEmpty(op2.Value) {
			return errors.New(fmt.Sprintf("invalid operator : operator is invalid %v", op2.Value))
		}
		if IsEmpty(op2.Name) {
			return errors.New(fmt.Sprintf("invalid reference : name is empty %v", op2.Name))
		}
		if _, ok := dup[op2.Value]; ok {
			return errors.New(fmt.Sprintf("invalid reference : name is a duplicate [%v]", op2.Value))
		}
		dup[op2.Value] = op2.Value
		*items = append(*items, op2)
	}
	return nil
}

func createOperator(op accessdata.Operator) (accessdata.Operator, error) {
	if IsEmpty(op.Value) {
		return accessdata.Operator{}, errors.New(fmt.Sprintf("invalid operator: value is empty %v", op.Name))
	}
	if !strings.HasPrefix(op.Value, accessdata.OperatorPrefix) {
		if IsEmpty(op.Name) {
			return accessdata.Operator{}, errors.New(fmt.Sprintf("invalid operator : name is empty [%v]", op.Value))
		}
		return accessdata.Operator{Name: accessdata.CreateDirect(op.Name), Value: op.Value}, nil
	}
	if op2, ok := accessdata.Operators[op.Value]; ok {
		newOp := accessdata.Operator{Name: op2.Name, Value: op.Value}
		if !IsEmpty(op.Name) {
			newOp.Name = op.Name
		}
		return newOp, nil
	}
	if strings.HasPrefix(op.Value, accessdata.RequestReferencePrefix) {
		return createHeaderOperator(op), nil
	}
	return accessdata.Operator{}, errors.New(fmt.Sprintf("invalid operator : value not found %v", op.Value))
}

func createHeaderOperator(op accessdata.Operator) accessdata.Operator {
	if IsEmpty(op.Value) || !strings.HasPrefix(op.Value, accessdata.RequestReferencePrefix) || len(op.Value) <= len(accessdata.RequestReferencePrefix) {
		return accessdata.Operator{}
	}
	s := op.Value[len(accessdata.RequestReferencePrefix):]
	tokens := strings.Split(s, ")")
	if len(tokens) == 1 || tokens[0] == "" {
		return accessdata.Operator{}
	}
	op1 := fmt.Sprintf("%v:%v", accessdata.HeaderPrefix, tokens[0])
	if IsEmpty(op.Name) {
		return accessdata.Operator{Name: tokens[0], Value: op1}
	}
	return accessdata.Operator{Name: op.Name, Value: op1}
}
