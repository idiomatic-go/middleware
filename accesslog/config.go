package accesslog

import (
	"errors"
	"fmt"
	"github.com/idiomatic-go/middleware/accessdata"
	"strings"
)

var ingressOperators []accessdata.Operator
var egressOperators []accessdata.Operator

func CreateIngressOperators(config []accessdata.Operator) error {
	ingressOperators = []accessdata.Operator{}
	return CreateOperators(&ingressOperators, config)
}

func CreateEgressOperators(config []accessdata.Operator) error {
	egressOperators = []accessdata.Operator{}
	return CreateOperators(&egressOperators, config)
}

func CreateOperators(items *[]accessdata.Operator, config []accessdata.Operator) error {
	if items == nil {
		return errors.New("invalid configuration : operators are nil")
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
		//if IsEmpty(op2.Value) {
		//	return errors.New(fmt.Sprintf("invalid operator : operator is invalid %v", op2.Value))
		//}
		//if IsEmpty(op2.Name) {
		//	return errors.New(fmt.Sprintf("invalid reference : name is empty %v", op2.Name))
		//}
		if _, ok := dup[op2.Value]; ok {
			return errors.New(fmt.Sprintf("invalid operator : value is a duplicate [%v]", op2.Value))
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
	if accessdata.IsDirectOperator(op) {
		if IsEmpty(op.Name) {
			return accessdata.Operator{}, errors.New(fmt.Sprintf("invalid operator : name is empty [%v]", op.Value))
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
		newOp, ok := accessdata.ParseRequestOperator(op)
		if !ok {
			return accessdata.Operator{}, errors.New(fmt.Sprintf("invalid operator : request is empty or invalid %v", op.Value))
		}
		return newOp, nil
	}
	return accessdata.Operator{}, errors.New(fmt.Sprintf("invalid operator : value not found or invalid %v", op.Value))
}

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
