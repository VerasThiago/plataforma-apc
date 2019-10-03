package project

import (
	"context"

	"github.com/apc-unb/apc-api/web/components/admin"

	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/bson/primitive"
	"github.com/mongodb/mongo-go-driver/mongo"
	"github.com/mongodb/mongo-go-driver/mongo/options"
)

// Here's an array in which you can store the decoded documents
func GetProjects(db *mongo.Client, studentID primitive.ObjectID, databaseName, collectionName string) ([]Project, error) {

	collection := db.Database(databaseName).Collection(collectionName)

	studentProjects := []Project{}

	// Passing bson.D{{}} as the filter matches all documents in the collection
	cursor, err := collection.Find(
		context.TODO(),
		bson.M{
			"studentID": studentID,
		},
		options.Find(),
	)

	if err != nil {
		return nil, err
	}

	// Finding multiple documents returns a cursor
	// Iterating through the cursor allows us to decode documents one at a time
	for cursor.Next(context.TODO()) {

		// create a value into which the single document can be decoded
		var elem Project

		// Checks if decoding method didn't return any errors
		if err := cursor.Decode(&elem); err != nil {
			return nil, err
		}

		// Push school class inside student array
		studentProjects = append(studentProjects, elem)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	// Close the cursor once finished
	cursor.Close(context.TODO())

	return studentProjects, nil
}

func CreateProject(db *mongo.Client, projectInfo Project, databaseName string) (admin.Admin, error) {

	//
	// GETTING RANDOM MONITOR BASED ON
	// ON LEAST PROJECTS ASSIGNED TO HIS NAME
	//
	var monitorInfo admin.Admin

	collection := db.Database(databaseName).Collection("admin")

	projection := bson.D{
		{"firstname", 1},
		{"lastname", 1},
		{"email", 1},
	}

	sortMethod := bson.D{
		{"projects", 1},
	}

	// Passing bson.D{{}} as the filter matches all documents in the collection
	cursor, err := collection.Find(
		context.TODO(),
		bson.M{
			"classid": projectInfo.ClassID,
		},
		options.Find().SetProjection(projection).SetSort(sortMethod).SetLimit(1),
	)

	if err != nil {
		return monitorInfo, err
	}

	// Finding multiple documents returns a cursor
	// Iterating through the cursor allows us to decode documents one at a time
	for cursor.Next(context.TODO()) {

		// Checks if decoding method didn't return any errors
		if err := cursor.Decode(&monitorInfo); err != nil {
			return monitorInfo, err
		}

		break
	}

	if err := cursor.Err(); err != nil {
		return monitorInfo, err
	}

	// Close the cursor once finished
	cursor.Close(context.TODO())

	// Assign current project to the monitor
	projectInfo.MonitorID = monitorInfo.ID

	// Updating amountto projects
	update := bson.D{
		{"$inc", bson.D{
			{"projects", 1},
		}},
	}

	filter := bson.M{
		"_id": monitorInfo.ID,
	}

	if _, err := collection.UpdateOne(context.TODO(), filter, update, nil); err != nil {
		return monitorInfo, err
	}

	//
	// Inserting project into DB
	//

	collection = db.Database(databaseName).Collection("projects")

	if _, err := collection.InsertOne(context.TODO(), projectInfo); err != nil {
		return monitorInfo, err
	}

	return monitorInfo, nil

}