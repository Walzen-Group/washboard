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

func CreateIgnoredImage(ignoredImage *types.IgnoredImage) error {
	conn, conErr := GetConnection()
	if conErr != nil {
		return conErr
	}
	collection := conn.db.Collection(types.DbIgnoredImagesCollection)
	indexModel := mongo.IndexModel{
		Keys:    bson.D{primitive.E{Key: "name", Value: 1}},
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
	_, err = collection.InsertOne(ctx, ignoredImage)
	if err != nil {
		err = werrors.NewCannotInsertError(err, err.(mongo.WriteException).WriteErrors[0].Message) // 🌯 the error
	}
	return err
}

func DeleteIgnoredImage(name string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	conn, conErr := GetConnection()
	if conErr != nil {
		return conErr
	}
	collection := conn.db.Collection(types.DbIgnoredImagesCollection)
	res, err := collection.DeleteOne(ctx, bson.M{"name": name})
	if err != nil {
		return err
	}
	if res.DeletedCount == 0 {
		return fmt.Errorf("No ignored image found with name %s", name)
	}
	return err
}

func GetAllIgnoredImages() ([]types.IgnoredImage, error) {
	conn, conErr := GetConnection()
	if conErr != nil {
		return nil, conErr
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := conn.db.Collection(types.DbIgnoredImagesCollection)
	var ignoredImages []types.IgnoredImage

	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	err = cursor.All(ctx, &ignoredImages)
	if err != nil {
		return nil, err
	}
	return ignoredImages, nil
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
			return nil, &werrors.DoesNotExistError{
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

func UpdateStackPriority(stackSettings *types.StackSettings) error {
	conn, conErr := GetConnection()
	if conErr != nil {
		return conErr
	}
	collection := conn.db.Collection(types.DbStackSettingsCollection)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// get all stack settings and move the priorities downwards from the inser position and upwards from where it was taken from
	allStackSettings, err := GetAllStackSettings()
	if err != nil {
		return err
	}
	oldStackSettings, err := GetStackSettings(stackSettings.StackName)
	if err != nil {
		return err
	}

	moveUp := false
	if oldStackSettings.Priority > stackSettings.Priority {
		moveUp = true
	}

	var updates []mongo.WriteModel
	for _, liveStackSetting := range allStackSettings {
		if liveStackSetting.StackName == stackSettings.StackName {
			liveStackSetting.Priority = stackSettings.Priority
		} else if moveUp {
			if liveStackSetting.Priority >= stackSettings.Priority && liveStackSetting.Priority < oldStackSettings.Priority {
				liveStackSetting.Priority++
			}
		} else {
			if liveStackSetting.Priority <= stackSettings.Priority && liveStackSetting.Priority > oldStackSettings.Priority {
				liveStackSetting.Priority--
			}
		}

		update := mongo.NewUpdateOneModel()
		update.SetFilter(bson.M{"stackName": liveStackSetting.StackName})
		update.SetUpdate(bson.M{"$set": bson.M{"priority": liveStackSetting.Priority}})
		update.SetUpsert(false)
		updates = append(updates, update)
	}

	_, err = collection.BulkWrite(ctx, updates)
	if err != nil {
		return err
	}

	return nil
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
