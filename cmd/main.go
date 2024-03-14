package main

import (
	"context"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/test/pkg/handlers"
)

const tableName = "user"

var dynoClient *dynamodb.DynamoDB

func handler(ctx context.Context, request events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {

	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	dynoClient = dynamodb.New(sess)
	switch request.HTTPMethod {
	case "POST":
		return handlers.CreateUserHandler(request, tableName, dynoClient)
	case "GET":
		return handlers.GetUserByEmail(request, tableName, dynoClient)
	case "PUT":
		return handlers.UpdateUserHandler(request, tableName, dynoClient)

	case "DELETE":
		return handlers.DeleteUserHandler(request, tableName, dynoClient)
	default:
		return &events.APIGatewayProxyResponse{
			StatusCode: http.StatusMethodNotAllowed,
			Body:       "Method not allowed",
		}, nil
	}

}

func main() {
	lambda.Start(handler)
}
