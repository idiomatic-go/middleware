package accesslog

import (
	"errors"
	"fmt"
	"github.com/idiomatic-go/middleware/accessdata"
)

var ingressOperators []accessdata.Operator
var egressOperators []accessdata.Operator

func CreateOperators(items *[]accessdata.Operator, config []accessdata.Operator) error {
	if items == nil {
		return errors.New("invalid configuration: operators slice is nil")
	}
	if len(config) == 0 {
		return errors.New("invalid configuration: configuration slice is empty")
	}
	dup := make(map[string]string)
	for _, op := range config {
		op2, err := createOperator(op)
		if err != nil {
			return err
		}
		if _, ok := dup[op2.Name]; ok {
			return errors.New(fmt.Sprintf("invalid operator: name is a duplicate [%v]", op2.Name))
		}
		dup[op2.Name] = op2.Name
		*items = append(*items, op2)
	}
	return nil
}

func createOperator(op accessdata.Operator) (accessdata.Operator, error) {
	if IsEmpty(op.Value) {
		return accessdata.Operator{}, errors.New(fmt.Sprintf("invalid operator: value is empty %v", op.Name))
	}
	if accessdata.IsDirectOperator(op) {
		if IsEmpty(op.Name) {
			return accessdata.Operator{}, errors.New(fmt.Sprintf("invalid operator: name is empty [%v]", op.Value))
		}
		return accessdata.Operator{Name: op.Name, Value: op.Value}, nil
	}
	if op2, ok := accessdata.Operators[op.Value]; ok {
		newOp := accessdata.Operator{Name: op2.Name, Value: op.Value}
		if !IsEmpty(op.Name) {
			newOp.Name = op.Name
		}
		return newOp, nil
	}
	if accessdata.IsRequestOperator(op) {
		return accessdata.Operator{Name: accessdata.RequestOperatorHeaderName(op), Value: op.Value}, nil
	}
	return accessdata.Operator{}, errors.New(fmt.Sprintf("invalid operator: value not found or invalid %v", op.Value))
}

/*
func createRequestOperator(op accessdata.Operator) accessdata.Operator {
	if len(op.Value) <= len(accessdata.RequestReferencePrefix) {
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


*/
