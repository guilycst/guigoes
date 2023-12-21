package main

import (
	"context"
	"path/filepath"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	ginadapter "github.com/awslabs/aws-lambda-go-api-proxy/gin"
	"github.com/guilycst/guigoes/internal/handlers"
	"github.com/guilycst/guigoes/internal/services"
	"github.com/guilycst/guigoes/pkg"
	"github.com/otiai10/copy"
)

var ginLambda *ginadapter.GinLambda

func init() {
	pkg.LoadEnvFromOS()
	idxTmp := "/tmp/" + filepath.Base(pkg.BLEVE_IDX_PATH) + "/"
	copy.Copy(pkg.BLEVE_IDX_PATH, idxTmp)
	pkg.BLEVE_IDX_PATH = idxTmp
	lps := services.NewLocalPostService()
	gr := handlers.NewGinRouter(lps)
	ginLambda = ginadapter.New(gr.Engine)
}

func Handler(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// If no name is provided in the HTTP request body, throw an error
	return ginLambda.ProxyWithContext(ctx, req)
}

func main() {
	lambda.Start(Handler)
}
