package mongoadapter

import (
	"context"
	"fmt"
	"github.com/casbin/casbin/model"
	"github.com/casbin/casbin/persist"
	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/mongo"
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
func NewAdapter(client *mongo.Client, database string, collection string) persist.Adapter {
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

func (a *adapter) LoadPolicy(model model.Model) error {
	cur, err := a.collection.Find(context.Background(), bson.D{})
	if err != nil {
		return err
	}
	for cur.Next(context.Background()) {
		rule := CasbinRule{}
		err = cur.Decode(&rule)
		if err != nil {
			return err
		}
		loadPolicyLine(rule, model)
	}
	return nil
}

func (a *adapter) SavePolicy(model model.Model) error {
	if err := a.collection.Drop(context.Background()); err != nil {
		return err
	}

	var lines []interface{}

	for ptype, ast := range model["p"] {
		for _, rule := range ast.Policy {
			line := savePolicyLine(ptype, rule)
			lines = append(lines, &line)
		}
	}

	for ptype, ast := range model["g"] {
		for _, rule := range ast.Policy {
			line := savePolicyLine(ptype, rule)
			lines = append(lines, &line)
		}
	}

	_, err := a.collection.InsertMany(context.Background(), lines)
	if err != nil {
		return err
	}
	return nil
}

func (a *adapter) AddPolicy(sec string, ptype string, rule []string) error {
	line := savePolicyLine(ptype, rule)
	_, err := a.collection.InsertOne(context.Background(), line)
	if err != nil {
		return err
	}
	return nil
}

func (a *adapter) RemovePolicy(sec string, ptype string, rule []string) error {
	filter := bson.D{
		{"ptype", ptype},
	}

	for i, v := range rule {
		filter = append(filter, bson.E{Key: fmt.Sprintf("v%d", i), Value: v})
	}

	_, err := a.collection.DeleteOne(context.Background(), filter)
	if err != nil {
		switch err {
		case mongo.ErrNoDocuments:
			return nil
		default:
			return err
		}
	}
	return nil
}

func (a *adapter) RemoveFilteredPolicy(sec string, ptype string, fieldIndex int, fieldValues ...string) error {
	filter := bson.D{
		{"ptype", ptype},
	}

	if fieldIndex <= 0 && 0 < fieldIndex+len(fieldValues) {
		if fieldValues[0-fieldIndex] != "" {
			filter = append(filter, bson.E{Key: "v0", Value: fieldValues[0-fieldIndex]})
		}
	}
	if fieldIndex <= 1 && 1 < fieldIndex+len(fieldValues) {
		if fieldValues[1-fieldIndex] != "" {
			filter = append(filter, bson.E{Key: "v1", Value: fieldValues[0-fieldIndex]})
		}
	}
	if fieldIndex <= 2 && 2 < fieldIndex+len(fieldValues) {
		if fieldValues[2-fieldIndex] != "" {
			filter = append(filter, bson.E{Key: "v2", Value: fieldValues[0-fieldIndex]})
		}
	}
	if fieldIndex <= 3 && 3 < fieldIndex+len(fieldValues) {
		if fieldValues[3-fieldIndex] != "" {
			filter = append(filter, bson.E{Key: "v3", Value: fieldValues[0-fieldIndex]})
		}
	}
	if fieldIndex <= 4 && 4 < fieldIndex+len(fieldValues) {
		if fieldValues[4-fieldIndex] != "" {
			filter = append(filter, bson.E{Key: "v4", Value: fieldValues[0-fieldIndex]})
		}
	}
	if fieldIndex <= 5 && 5 < fieldIndex+len(fieldValues) {
		if fieldValues[5-fieldIndex] != "" {
			filter = append(filter, bson.E{Key: "v5", Value: fieldValues[0-fieldIndex]})
		}
	}

	_, err := a.collection.DeleteMany(context.Background(), filter)
	if err != nil {
		return err
	}
	return nil
}
