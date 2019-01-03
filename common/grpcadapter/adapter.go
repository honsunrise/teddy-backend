package grpcadapter

import (
	"context"
	"github.com/casbin/casbin/model"
	"github.com/casbin/casbin/util"
	"github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/grpc"
	"strings"
)

type Adapter struct {
	target string
	client PolicyAdapterClient
}

func NewAdapter(target string) (*Adapter, error) {
	a := Adapter{}
	a.target = target
	conn, err := grpc.Dial(a.target, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	a.client = NewPolicyAdapterClient(conn)
	return &a, nil
}

func loadPolicyLine(line string, model model.Model) {
	if line == "" {
		return
	}

	tokens := strings.Split(line, ", ")

	key := tokens[0]
	sec := key[:1]
	model[sec][key].Policy = append(model[sec][key].Policy, tokens[1:])
}

func (a *Adapter) LoadPolicy(model model.Model) error {
	resp, err := a.client.LoadPolicy(context.Background(), &empty.Empty{})
	if err != nil {
		return err
	}
	for _, line := range resp.Rules {
		loadPolicyLine(line, model)
	}
	return nil
}

func (a *Adapter) SavePolicy(model model.Model) error {
	var rules []string
	for ptype, ast := range model["p"] {
		for _, rule := range ast.Policy {
			tmp := ptype + ", " + util.ArrayToString(rule)
			rules = append(rules, tmp)
		}
	}

	for ptype, ast := range model["g"] {
		for _, rule := range ast.Policy {
			tmp := ptype + ", " + util.ArrayToString(rule)
			rules = append(rules, tmp)
		}
	}

	_, err := a.client.SavePolicy(context.Background(), &Policy{
		Rules: rules,
	})
	return err
}

func (a *Adapter) AddPolicy(sec string, ptype string, rule []string) error {
	_, err := a.client.AddPolicy(context.Background(), &AddPolicyReq{
		Sec:   sec,
		Ptype: ptype,
		Rule:  rule,
	})
	return err
}

func (a *Adapter) RemovePolicy(sec string, ptype string, rule []string) error {
	_, err := a.client.RemovePolicy(context.Background(), &RemovePolicyReq{
		Sec:   sec,
		Ptype: ptype,
		Rule:  rule,
	})
	return err
}

func (a *Adapter) RemoveFilteredPolicy(sec string, ptype string, fieldIndex int, fieldValues ...string) error {
	_, err := a.client.RemoveFilteredPolicy(context.Background(), &RemoveFilteredPolicyReq{
		Sec:         sec,
		Ptype:       ptype,
		FieldIndex:  int64(fieldIndex),
		FieldValues: fieldValues,
	})
	return err
}
