package mongo_grpcadapter

import (
	"context"
	"fmt"
	"github.com/casbin/casbin/model"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/mongo"
	"teddy-backend/pkg/grpcadapter"
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

// adapter represents the MongoDB adapter for policy storage.
type adapter struct {
	client     *mongo.Client
	collection *mongo.Collection
}

// NewAdapter is the constructor for Adapter. If database name is not provided
// in the Mongo URL, 'casbin' will be used as database name.
func NewServer(client *mongo.Client, database string, collection string) grpcadapter.PolicyAdapterServer {
	a := &adapter{
		client:     client,
		collection: client.Database(database).Collection(collection),
	}
	return a
}

func loadPolicyLine(line CasbinRule, model model.Model) {
	key := line.PType
	sec := key[:1]

	var tokens []string
	if line.V0 != "" {
		tokens = append(tokens, line.V0)
	} else {
		goto LineEnd
	}

	if line.V1 != "" {
		tokens = append(tokens, line.V1)
	} else {
		goto LineEnd
	}

	if line.V2 != "" {
		tokens = append(tokens, line.V2)
	} else {
		goto LineEnd
	}

	if line.V3 != "" {
		tokens = append(tokens, line.V3)
	} else {
		goto LineEnd
	}

	if line.V4 != "" {
		tokens = append(tokens, line.V4)
	} else {
		goto LineEnd
	}

	if line.V5 != "" {
		tokens = append(tokens, line.V5)
	} else {
		goto LineEnd
	}

LineEnd:
	model[sec][key].Policy = append(model[sec][key].Policy, tokens)
}

func savePolicyLine(ptype string, rule []string) CasbinRule {
	line := CasbinRule{
		PType: ptype,
	}

	if len(rule) > 0 {
		line.V0 = rule[0]
	}
	if len(rule) > 1 {
		line.V1 = rule[1]
	}
	if len(rule) > 2 {
		line.V2 = rule[2]
	}
	if len(rule) > 3 {
		line.V3 = rule[3]
	}
	if len(rule) > 4 {
		line.V4 = rule[4]
	}
	if len(rule) > 5 {
		line.V5 = rule[5]
	}

	return line
}

func (a *adapter) LoadPolicy(ctx context.Context, req *empty.Empty) (*grpcadapter.Policies, error) {
	var policies []*grpcadapter.Policy

	cur, err := a.collection.Find(ctx, bson.D{})
	if err != nil {
		return nil, err
	}
	for cur.Next(ctx) {
		rule := CasbinRule{}
		err = cur.Decode(&rule)
		if err != nil {
			return nil, err
		}
		policies = append(policies, &grpcadapter.Policy{
			Ptype: rule.PType,
			Rule: []string{
				rule.V0,
				rule.V1,
				rule.V2,
				rule.V3,
				rule.V4,
				rule.V5,
			},
		})
	}
	return &grpcadapter.Policies{
		Policies: policies,
	}, nil
}

func (a *adapter) SavePolicy(ctx context.Context, req *grpcadapter.Policies) (*empty.Empty, error) {
	if err := a.collection.Drop(ctx); err != nil {
		return nil, err
	}

	var lines []interface{}

	for _, policy := range req.Policies {
		line := savePolicyLine(policy.Ptype, policy.Rule)
		lines = append(lines, &line)
	}

	_, err := a.collection.InsertMany(ctx, lines)
	if err != nil {
		return nil, err
	}
	return &empty.Empty{}, nil
}

func (a *adapter) AddPolicy(ctx context.Context, req *grpcadapter.AddPolicyReq) (*empty.Empty, error) {
	line := savePolicyLine(req.Ptype, req.Rule)
	_, err := a.collection.InsertOne(ctx, line)
	if err != nil {
		return nil, err
	}
	return &empty.Empty{}, nil
}

func (a *adapter) RemovePolicy(ctx context.Context, req *grpcadapter.RemovePolicyReq) (*empty.Empty, error) {
	filter := bson.D{
		{"ptype", req.Ptype},
	}

	for i, v := range req.Rule {
		filter = append(filter, bson.E{Key: fmt.Sprintf("v%d", i), Value: v})
	}

	_, err := a.collection.DeleteOne(ctx, filter)
	if err != nil {
		switch err {
		case mongo.ErrNoDocuments:
			return &empty.Empty{}, nil
		default:
			return nil, err
		}
	}
	return &empty.Empty{}, nil
}

func (a *adapter) RemoveFilteredPolicy(ctx context.Context, req *grpcadapter.RemoveFilteredPolicyReq) (*empty.Empty, error) {
	fieldIndex := req.FieldIndex
	fieldValues := req.FieldValues
	filter := bson.D{
		{"ptype", req.Ptype},
	}

	if fieldIndex <= 0 && 0 < fieldIndex+int64(len(fieldValues)) {
		if fieldValues[0-fieldIndex] != "" {
			filter = append(filter, bson.E{Key: "v0", Value: fieldValues[0-fieldIndex]})
		}
	}
	if fieldIndex <= 1 && 1 < fieldIndex+int64(len(fieldValues)) {
		if fieldValues[1-fieldIndex] != "" {
			filter = append(filter, bson.E{Key: "v1", Value: fieldValues[0-fieldIndex]})
		}
	}
	if fieldIndex <= 2 && 2 < fieldIndex+int64(len(fieldValues)) {
		if fieldValues[2-fieldIndex] != "" {
			filter = append(filter, bson.E{Key: "v2", Value: fieldValues[0-fieldIndex]})
		}
	}
	if fieldIndex <= 3 && 3 < fieldIndex+int64(len(fieldValues)) {
		if fieldValues[3-fieldIndex] != "" {
			filter = append(filter, bson.E{Key: "v3", Value: fieldValues[0-fieldIndex]})
		}
	}
	if fieldIndex <= 4 && 4 < fieldIndex+int64(len(fieldValues)) {
		if fieldValues[4-fieldIndex] != "" {
			filter = append(filter, bson.E{Key: "v4", Value: fieldValues[0-fieldIndex]})
		}
	}
	if fieldIndex <= 5 && 5 < fieldIndex+int64(len(fieldValues)) {
		if fieldValues[5-fieldIndex] != "" {
			filter = append(filter, bson.E{Key: "v5", Value: fieldValues[0-fieldIndex]})
		}
	}

	_, err := a.collection.DeleteMany(context.Background(), filter)
	if err != nil {
		return nil, err
	}
	return &empty.Empty{}, nil
}
