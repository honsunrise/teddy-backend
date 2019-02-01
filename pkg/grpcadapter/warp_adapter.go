package grpcadapter

import (
	"context"
	"github.com/casbin/casbin/model"
	"github.com/casbin/casbin/persist"
	"github.com/golang/protobuf/ptypes/empty"
)

// CasbinRule represents a rule in Casbin.
type CasbinRule struct {
	PType string `bson:"ptype"`
	V0    string `bson:"v0"`
	V1    string `bson:"v1"`
	V2    string `bson:"v2"`
	V3    string `bson:"v3"`
	V4    string `bson:"v4"`
	V5    string `bson:"v5"`
}

// warpAdapter represents the MongoDB warpAdapter for policy storage.
type warpAdapter struct {
	adapter persist.Adapter
}

// NewAdapter is the constructor for Adapter. If database name is not provided
// in the Mongo URL, 'casbin' will be used as database name.
func NewServer(adapter persist.Adapter) PolicyAdapterServer {
	a := &warpAdapter{
		adapter: adapter,
	}
	return a
}

func (a *warpAdapter) LoadPolicy(ctx context.Context, req *empty.Empty) (*Policies, error) {
	var policies []*Policy

	m := model.Model{}
	a.adapter.LoadPolicy(m)

	for ptype, ast := range m["p"] {
		for _, rule := range ast.Policy {
			policies = append(policies, &Policy{
				Ptype: ptype,
				Rule:  rule,
			})
		}
	}

	for ptype, ast := range m["g"] {
		for _, rule := range ast.Policy {
			policies = append(policies, &Policy{
				Ptype: ptype,
				Rule:  rule,
			})
		}
	}

	return &Policies{
		Policies: policies,
	}, nil
}

func (a *warpAdapter) SavePolicy(ctx context.Context, req *Policies) (*empty.Empty, error) {
	m := model.Model{}

	for _, policy := range req.Policies {
		key := policy.Ptype
		sec := key[:1]
		m[sec][key].Policy = append(m[sec][key].Policy, policy.Rule)
	}

	a.adapter.SavePolicy(m)
	return &empty.Empty{}, nil
}

func (a *warpAdapter) AddPolicy(ctx context.Context, req *AddPolicyReq) (*empty.Empty, error) {
	err := a.adapter.AddPolicy(req.Sec, req.Ptype, req.Rule)
	if err != nil {
		return nil, err
	}
	return &empty.Empty{}, nil
}

func (a *warpAdapter) RemovePolicy(ctx context.Context, req *RemovePolicyReq) (*empty.Empty, error) {
	err := a.adapter.RemovePolicy(req.Sec, req.Ptype, req.Rule)

	if err != nil {
		return nil, err
	}
	return &empty.Empty{}, nil
}

func (a *warpAdapter) RemoveFilteredPolicy(ctx context.Context, req *RemoveFilteredPolicyReq) (*empty.Empty, error) {
	a.adapter.RemoveFilteredPolicy(req.Sec, req.Ptype, int(req.FieldIndex), req.FieldValues...)
	return &empty.Empty{}, nil
}
