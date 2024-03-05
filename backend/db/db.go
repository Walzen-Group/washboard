package db

import (
	"context"
	"sync"
	"time"
	"washboard/state"
	"washboard/types"

	"go.mongodb.org/mongo-driver/bson"
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
	collection := conn.db.Collection(types.DbStackSettingsCollection)
	if conErr != nil {
		return conErr
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_, err := collection.InsertOne(ctx, stackSettings)
	return err
}

// GetStackSettings retrieves a stack settings document from the database by ID
func GetStackSettings(id int) (*types.StackSettings, error) {
	conn, conErr := GetConnection()
	if conErr != nil {
		return nil, conErr
	}
	collection := conn.db.Collection(types.DbStackSettingsCollection)
	var stackSettings types.StackSettings

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err := collection.FindOne(ctx, bson.M{"stackId": id}).Decode(&stackSettings)
	if err != nil {
		return nil, err
	}
	return &stackSettings, nil
}

// UpdateStackSettings updates a stack settings document in the database
func UpdateStackSettings(stackSettings *types.StackSettings) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	conn, conErr := GetConnection()
	if conErr != nil {
		return conErr
	}
	collection := conn.db.Collection(types.DbStackSettingsCollection)
	_, err := collection.ReplaceOne(ctx, bson.M{"stackId": stackSettings.StackId}, stackSettings)
	return err
}

// DeleteStackSettings deletes a stack settings document from the database by ID
func DeleteStackSettings(id string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	conn, conErr := GetConnection()
	if conErr != nil {
		return conErr
	}
	collection := conn.db.Collection(types.DbStackSettingsCollection)
	_, err := collection.DeleteOne(ctx, bson.M{"_id": id})
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

	_, err := collection.InsertOne(ctx, groupSettings)
	return err
}

// GetGroupSettings retrieves a group settings document from the database by ID
func GetGroupSettings(id string) (*types.GroupSettings, error) {
	conn, conErr := GetConnection()
	if conErr != nil {
		return nil, conErr
	}
	collection := conn.db.Collection(types.DbGroupSettingsCollection)
	var groupSettings types.GroupSettings
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err := collection.FindOne(ctx, bson.M{"_id": id}).Decode(&groupSettings)
	if err != nil {
		return nil, err
	}
	return &groupSettings, nil
}

// UpdateGroupSettings updates a group settings document in the database
func UpdateGroupSettings(groupSettings *types.GroupSettings) error {
	conn, conErr := GetConnection()
	if conErr != nil {
		return conErr
	}
	collection := conn.db.Collection(types.DbGroupSettingsCollection)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_, err := collection.ReplaceOne(ctx, bson.M{"_id": groupSettings.GroupName}, groupSettings)
	return err
}

// DeleteGroupSettings deletes a group settings document from the database by ID
func DeleteGroupSettings(id string) error {
	conn, conErr := GetConnection()
	if conErr != nil {
		return conErr
	}
	collection := conn.db.Collection(types.DbGroupSettingsCollection)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_, err := collection.DeleteOne(ctx, bson.M{"_id": id})
	return err
}
