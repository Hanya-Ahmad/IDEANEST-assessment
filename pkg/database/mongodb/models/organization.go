package models

import (
	"context"
	"errors"
	"time"

	"github.com/Hanya-Ahmad/IDEANEST-assessment/pkg/utils"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// ErrOrganizationExists is raised when a new organization is added with a name linked to an existing organization
var ErrOrganizationExists = errors.New("an organization with that name already exists")

// Organization struct contains Organization's table fields
type Organization struct{
	OrganizationID     primitive.ObjectID   `bson:"_id,omitempty" json:"organization_id"`
	Name               string   `bson:"name" json:"name" validate:"required"`
	Description        string   `bson:"description" json:"description"`
	OrganizationMembers []Member `bson:"organization_members" json:"organization_members"`
}
// Validate validates the organization struct using the validate tag
func (organization *Organization) Validate() error {
	validate := validator.New()
	return validate.Struct(organization)
}


// CreateNewOrganization creates a new organization in the database
func (db *DBClient) CreateNewOrganization(org Organization)(string,error){
	collection := db.Client.Collection("Organization")
	_, err :=db.GetOrganizationByName(org.Name)
	if err == mongo.ErrNoDocuments{
		if org.OrganizationMembers == nil {
            org.OrganizationMembers = []Member{}
        }
		ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	result, err := collection.InsertOne(ctx, org)
	if err != nil {
		return "", err
	}
	newID := result.InsertedID.(primitive.ObjectID).Hex()
	return newID, nil
	}else{
		return "", ErrOrganizationExists
	}
}

// GetOrganizationByName retrieves an organization from the database
func (db *DBClient) GetOrganizationByName (name string)(Organization, error){
collection := db.Client.Collection("Organization")
	filter := bson.M{"name": name}
	var checkedOrg Organization
    err := collection.FindOne(context.Background(), filter).Decode(&checkedOrg)
	return checkedOrg, err	
}

func (db *DBClient) GetOrganizationByID (id primitive.ObjectID)(Organization, error){
	collection := db.Client.Collection("Organization")
	filter := bson.M{"_id": id}
	var checkedOrg Organization
    err := collection.FindOne(context.Background(), filter).Decode(&checkedOrg)
	return checkedOrg, err	
}

func (db *DBClient) GetAllOrganizations ()([]Organization,error){
	collection := db.Client.Collection("Organization")
    var organizations []Organization
    cursor, err := collection.Find(context.Background(), bson.M{})
    if err != nil {
        return nil, err
    }
    defer cursor.Close(context.Background())
    if err := cursor.All(context.Background(), &organizations); err != nil {
        return nil, err
    }
    return organizations, nil

}

// UpdateOrganization updates an organization's information in the database
func (db *DBClient) UpdateOrganization(id primitive.ObjectID, update bson.M) (error) {
    _, err := db.GetOrganizationByID(id)
	if err!=nil{
		return err
	}
    collection := db.Client.Collection("Organization")
    filter := bson.M{"_id": id}
    updateOptions := options.Update().SetUpsert(true)
    _, err = collection.UpdateOne(context.Background(), filter, update, updateOptions)
    if err != nil {
        return err
    }
   
    return nil
}

// DeleteOrganization deletes an organization from the database
func (db *DBClient) DeleteOrganization(id primitive.ObjectID) error{
	  _, err := db.GetOrganizationByID(id)
	if err!=nil{
		return err
	}
    collection := db.Client.Collection("Organization")
    filter := bson.M{"_id": id}
    _, err = collection.DeleteOne(context.Background(), filter)
    if err != nil {
        return err
    }
    return nil
}

// SendMail checks that the organization exists in the database before sending the invitation email
func (db *DBClient) SendMail(email string, id primitive.ObjectID) error{
	org, err := db.GetOrganizationByID(id)
	if err!=nil{
		return err
	}
	_, err = db.GetUserByEmail(email)
	if err!=nil{
		return ErrUserNotFound
	}
	strID := id.Hex()
	err = utils.SendMail(email, org.Name, strID)
	return err
}