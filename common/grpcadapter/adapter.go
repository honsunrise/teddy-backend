package grpcadapter

import (
	"context"
	"github.com/casbin/casbin/model"
	"github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/grpc"
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

func loadPolicyRule(policy *Policy, model model.Model) {
	if policy == nil {
		return
	}

	key := policy.Ptype
	sec := key[:1]
	model[sec][key].Policy = append(model[sec][key].Policy, policy.Rule)
}

func (a *Adapter) LoadPolicy(model model.Model) error {
	resp, err := a.client.LoadPolicy(context.Background(), &empty.Empty{})
	if err != nil {
		return err
	}
	for _, line := range resp.Policies {
		loadPolicyRule(line, model)
	}
	return nil
}

func (a *Adapter) SavePolicy(model model.Model) error {
	var policies []*Policy
	for ptype, ast := range model["p"] {
		for _, rule := range ast.Policy {
			policies = append(policies, &Policy{
				Ptype: ptype,
				Rule:  rule,
			})
		}
	}

	for ptype, ast := range model["g"] {
		for _, rule := range ast.Policy {
			policies = append(policies, &Policy{
				Ptype: ptype,
				Rule:  rule,
			})
		}
	}

	_, err := a.client.SavePolicy(context.Background(), &Policies{
		Policies: policies,
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
