package db

import (
	"context"
	"fmt"
	"sync"
	"time"
	"washboard/state"
	"washboard/types"
	"washboard/werrors"

	"github.com/kpango/glg"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var instance *DataStore
var once sync.Once

type DataStore struct {
	db     *mongo.Database
	client *mongo.Client
}

func (ds *DataStore) Db() *mongo.Database {
	return ds.db
}

func (ds *DataStore) Client() *mongo.Client {
	return ds.client
}

// Get establishes a connection to the database and returns the db handle
func GetConnection() (*DataStore, error) {
	var connectErr error
	once.Do(func() {
		instance = new(DataStore)
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		dbUrl := state.Instance().Config.DbUrl
		glg.Infof("Connecting to database at %s", dbUrl)
		client, err := mongo.Connect(ctx, options.Client().ApplyURI(dbUrl))
		if err != nil {
			connectErr = err
			return
		}
		instance.client = client
		instance.db = client.Database(types.DbName)
	})
	if connectErr != nil {
		return nil, connectErr
	}
	return instance, nil
}

// CreateStackSettings creates a new stack settings document in the database
func CreateStackSettings(stackSettings *types.StackSettings) error {
	conn, conErr := GetConnection()
	if conErr != nil {
		return conErr
	}
	collection := conn.db.Collection(types.DbStackSettingsCollection)
	indexModel := mongo.IndexModel{
		Keys:    bson.D{primitive.E{Key: "stackName", Value: 1}},
		Options: options.Index().SetUnique(true),
	}
	ctxOp2, cancelOp2 := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancelOp2()
	_, err := collection.Indexes().CreateOne(ctxOp2, indexModel)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_, err = collection.InsertOne(ctx, stackSettings)
	if err != nil {
		err = werrors.NewCannotInsertError(err, err.(mongo.WriteException).WriteErrors[0].Message) // 🌯 the error
	}
	return err
}

// GetStackSettings retrieves a stack settings document from the database by ID
func GetStackSettings(name string) (*types.StackSettings, error) {
	conn, conErr := GetConnection()
	if conErr != nil {
		return nil, conErr
	}
	collection := conn.db.Collection(types.DbStackSettingsCollection)
	var stackSettings types.StackSettings

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err := collection.FindOne(ctx, bson.M{"stackName": name}).Decode(&stackSettings)
	if err != nil {
		if err.Error() == "mongo: no documents in result" {
			return nil , &werrors.DoesNotExistError{
				Context: "empty response",
				Err:     err,
			}
		}
		return nil, err
	}
	return &stackSettings, nil
}

// GetAllStackSettings retrieves all stack settings from the database
func GetAllStackSettings() ([]types.StackSettings, error) {
	conn, conErr := GetConnection()
	if conErr != nil {
		return nil, conErr
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := conn.db.Collection(types.DbStackSettingsCollection)
	var stackSettings []types.StackSettings

	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	err = cursor.All(ctx, &stackSettings)
	if err != nil {
		return nil, err
	}
	return stackSettings, nil
}

// UpdateStackSettings updates a stack settings document in the database
func UpdateStackSettings(stackSettings *types.StackSettings, stackName string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	conn, conErr := GetConnection()
	if conErr != nil {
		return conErr
	}
	collection := conn.db.Collection(types.DbStackSettingsCollection)
	res, err := collection.ReplaceOne(ctx, bson.M{"stackName": stackName}, stackSettings)
	if err != nil {
		return err
	}
	if res.MatchedCount == 0 {
		return fmt.Errorf("No stack settings found with stack name %s", stackName)
	}
	return err
}

// DeleteStackSettings deletes a stack settings document from the database by ID
func DeleteStackSettings(stackName string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	conn, conErr := GetConnection()
	if conErr != nil {
		return conErr
	}
	collection := conn.db.Collection(types.DbStackSettingsCollection)
	res, err := collection.DeleteOne(ctx, bson.M{"stackName": stackName})
	if err != nil {
		return err
	}
	if res.DeletedCount == 0 {
		return fmt.Errorf("No stack settings found with stack name %s", stackName)
	}
	return err
}

// CreateGroupSettings creates a new group settings document in the database
func CreateGroupSettings(groupSettings *types.GroupSettings) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	conn, conErr := GetConnection()
	if conErr != nil {
		return conErr
	}
	collection := conn.db.Collection(types.DbGroupSettingsCollection)

	indexModel := mongo.IndexModel{
		Keys:    bson.D{primitive.E{Key: "groupName", Value: 1}},
		Options: options.Index().SetUnique(true),
	}
	ctxOp2, cancelOp2 := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancelOp2()
	_, err := collection.Indexes().CreateOne(ctxOp2, indexModel)
	if err != nil {
		return err
	}

	_, err = collection.InsertOne(ctx, groupSettings)
	return err
}

// GetGroupSettings retrieves a group settings document from the database by ID
func GetGroupSettings(groupName string) (*types.GroupSettings, error) {
	conn, conErr := GetConnection()
	if conErr != nil {
		return nil, conErr
	}
	collection := conn.db.Collection(types.DbGroupSettingsCollection)
	var groupSettings types.GroupSettings
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err := collection.FindOne(ctx, bson.M{"groupName": groupName}).Decode(&groupSettings)
	if err != nil {
		if err.Error() == "mongo: no documents in result" {
			return nil , &werrors.DoesNotExistError{
				Context: "empty response",
				Err:     err,
			}
		}
		return nil, err
	}
	return &groupSettings, nil
}

// GetGroupSettings retrieves a group settings document from the database by ID
func GetAllGroupSettings() ([]types.GroupSettings, error) {
	conn, conErr := GetConnection()
	if conErr != nil {
		return nil, conErr
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := conn.db.Collection(types.DbGroupSettingsCollection)
	var groupSettings []types.GroupSettings

	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	err = cursor.All(ctx, &groupSettings)
	if err != nil {
		return nil, err
	}
	return groupSettings, nil
}

// UpdateGroupSettings updates a group settings document in the database
func UpdateGroupSettings(groupSettings *types.GroupSettings, groupName string) error {
	conn, conErr := GetConnection()
	if conErr != nil {
		return conErr
	}
	collection := conn.db.Collection(types.DbGroupSettingsCollection)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	res, err := collection.ReplaceOne(ctx, bson.M{"groupName": groupName}, groupSettings)
	if err != nil {
		return err
	}
	if res.MatchedCount == 0 {
		return fmt.Errorf("No group settings found with name %s", groupName)
	}
	return err
}

// DeleteGroupSettings deletes a group settings document from the database by ID
func DeleteGroupSettings(groupName string) error {
	conn, conErr := GetConnection()
	if conErr != nil {
		return conErr
	}
	collection := conn.db.Collection(types.DbGroupSettingsCollection)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	res, err := collection.DeleteOne(ctx, bson.M{"groupName": groupName})
	if err != nil {
		return err
	}
	if res.DeletedCount == 0 {
		return fmt.Errorf("No group settings found with name %s", groupName)
	}
	return err
}
