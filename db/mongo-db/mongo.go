package mongo_db

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"

	"rnd-git.valsun.cn/ebike-server/go-common/logs"

	"go-gpt-assistant/config"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var Client *mongo.Client // 包级别的 Client 变量

type Document struct {
	MetaData map[string]any `bson:"metadata"`
	Content  string         `bson:"content"`
}

func InitMongo() {
	clientOptions := options.Client().ApplyURI(fmt.Sprintf("mongodb://%s:%s@%s", config.GetAppConf().MongoDbUser,
		config.GetAppConf().MongoDbPwd, config.GetAppConf().MongoDbAddr))
	var err error
	Client, err = mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		logs.Errorf("Failed to connect to MongoDB: %v", err)
	}
	err = Client.Ping(context.Background(), nil)
	if err != nil {
		logs.Errorf("Failed to ping MongoDB: %v", err)
	}
	logs.Infof("Connected to MongoDB")
}

func Disconnect() {
	Client.Disconnect(context.TODO())
	Client = nil
	log.Println("Disconnected from MongoDB!")
	return
}

func GetCollection() *mongo.Collection {
	return Client.Database(config.GetAppConf().MongoDbDatabase).Collection(config.GetAppConf().MongoDbCollection)
}

func InsertOne(data interface{}) error {
	collection := GetCollection()
	result, err := collection.InsertOne(context.TODO(), data)
	if err != nil {
		return err
	}
	logs.Infof("Inserted a single document: %v", result.InsertedID)
	return nil
}

func Searchtext(text string) ([]Document, error) {
	collection := GetCollection()
	filter := bson.D{{"content", bson.D{{"$regex", text}}}}
	ctx, _ := context.WithTimeout(context.Background(), 2*time.Second)
	cur, err := collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)
	var results []Document
	for cur.Next(ctx) {
		var result Document
		err := cur.Decode(&result)
		if err != nil {
			return nil, err
		}
		results = append(results, result)
	}
	if err := cur.Err(); err != nil {
		return nil, err
	}
	return results, nil
}
