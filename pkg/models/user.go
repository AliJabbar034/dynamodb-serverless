package models

import (
	"errors"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

type User struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (u *User) Create(tableName string, dynoClient *dynamodb.DynamoDB) (string, error) {

	user, err := dynamodbattribute.MarshalMap(u)
	if err != nil {
		return "", errors.New("marshing user data error")
	}

	input := &dynamodb.PutItemInput{
		Item:      user,
		TableName: aws.String(tableName),
	}

	_, err = dynoClient.PutItem(input)
	if err != nil {
		return "", err
	}

	return "user created successfuly", nil

}

func GetByEmail(email, tableName string, dynoClient *dynamodb.DynamoDB) (*User, error) {

	var user User
	input := &dynamodb.GetItemInput{
		TableName: aws.String(tableName),
		Key: map[string]*dynamodb.AttributeValue{
			"email": {
				S: aws.String(email),
			},
		},
	}

	result, err := dynoClient.GetItem(input)
	if err != nil {
		return nil, err
	}
	if result.Item == nil {
		return nil, err
	}

	if err := dynamodbattribute.UnmarshalMap(result.Item, &user); err != nil {
		return nil, err
	}
	return &user, nil

}

func DeleteUser(email, tableName string, dynoClient *dynamodb.DynamoDB) (string, error) {

	input := &dynamodb.DeleteItemInput{
		TableName: aws.String(tableName),
		Key: map[string]*dynamodb.AttributeValue{
			"email": {
				S: aws.String(email),
			},
		},
	}

	_, err := dynoClient.DeleteItem(input)
	if err != nil {
		return "", nil

	}

	return "Deleted Successfully", nil
}

func (u *User) UpdateUser(tableName string, dynoCLient *dynamodb.DynamoDB) (*User, error) {

	input := &dynamodb.UpdateItemInput{
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":n": {
				S: aws.String(u.Name),
			},
		},
		TableName: aws.String(tableName),
		Key: map[string]*dynamodb.AttributeValue{
			"email": {
				S: aws.String(u.Email),
			},
		},
		ReturnValues:     aws.String("UPDATED_NEW"),
		UpdateExpression: aws.String("SET #n = :n"),
		ExpressionAttributeNames: map[string]*string{
			"#n": aws.String("name"),
		},
	}

	_, err := dynoCLient.UpdateItem(input)
	if err != nil {
		return nil, err
	}
	return u, nil
}
