package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/test/pkg/models"
)

type Response struct {
	Message any `json:"message"`
}

func CreateUserHandler(req events.APIGatewayProxyRequest, tableName string, dynoClient *dynamodb.DynamoDB) (*events.APIGatewayProxyResponse, error) {
	user := &models.User{}

	if err := json.Unmarshal([]byte(req.Body), &user); err != nil {
		return &events.APIGatewayProxyResponse{
			StatusCode: http.StatusBadRequest,
			Body:       "invalid Data",
		}, nil
	}

	res, err := user.Create(tableName, dynoClient)
	if err != nil {
		return &events.APIGatewayProxyResponse{
			StatusCode: http.StatusOK,
			Headers: map[string]string{
				"Content-Type": "application/json",
			},
			Body: err.Error(),
		}, nil
	}
	responseBody := Response{
		Message: res,
	}

	jsonResponse, err := json.Marshal(responseBody)
	if err != nil {
		return &events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       err.Error(),
		}, nil
	}

	return &events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
		Body: string(jsonResponse),
	}, nil
}

func GetUserByEmail(req events.APIGatewayProxyRequest, tableName string, dynoCLient *dynamodb.DynamoDB) (*events.APIGatewayProxyResponse, error) {

	email := req.QueryStringParameters["email"]

	user, err := models.GetByEmail(email, tableName, dynoCLient)
	if err != nil {
		return &events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       err.Error(),
		}, nil
	}
	res := &Response{
		Message: user,
	}

	data, err := json.Marshal(res)
	if err != nil {
		return &events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       "marshing error ",
		}, nil
	}

	return &events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Body:       string(data),
	}, nil
}

func DeleteUserHandler(req events.APIGatewayProxyRequest, tableName string, dynoCLient *dynamodb.DynamoDB) (*events.APIGatewayProxyResponse, error) {

	email := req.QueryStringParameters["email"]
	res, err := models.DeleteUser(email, tableName, dynoCLient)
	if err != nil {
		return &events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       err.Error(),
		}, nil
	}
	return &events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Body:       string(res),
	}, nil
}

func UpdateUserHandler(req events.APIGatewayProxyRequest, tableName string, dynoClient *dynamodb.DynamoDB) (*events.APIGatewayProxyResponse, error) {

	user := &models.User{}
	if err := json.Unmarshal([]byte(req.Body), &user); err != nil {
		return &events.APIGatewayProxyResponse{
			StatusCode: http.StatusBadRequest,
			Body:       "invalid request",
		}, nil
	}

	usr, err := user.UpdateUser(tableName, dynoClient)
	if err != nil {
		return &events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       err.Error(),
		}, nil
	}

	res := &Response{
		Message: usr,
	}

	updatedUser, err := json.Marshal(res)
	if err != nil {
		return &events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       err.Error(),
		}, nil
	}
	return &events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Body:       string(updatedUser),
	}, nil

}
