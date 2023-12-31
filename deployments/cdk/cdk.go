package main

import (
	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsapigateway"
	"github.com/aws/aws-cdk-go/awscdk/v2/awslambda"
	"github.com/aws/aws-cdk-go/awscdk/v2/awss3assets"
	"github.com/aws/aws-cdk-go/awscdklambdagoalpha/v2"

	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
)

type CdkStackProps struct {
	awscdk.StackProps
}

func GuigoesCdkStack(scope constructs.Construct, id string, props *CdkStackProps) awscdk.Stack {
	var sprops awscdk.StackProps
	if props != nil {
		sprops = props.StackProps
	}

	stack := awscdk.NewStack(scope, &id, &sprops)

	lambda := awscdklambdagoalpha.NewGoFunction(stack, sptr("GuigoesLambda"), &awscdklambdagoalpha.GoFunctionProps{
		Runtime: awslambda.Runtime_GO_1_X(),
		Entry:   sptr("../../cmd/lambda/main.go"),
		Environment: &map[string]*string{
			"POSTS_PATH":     sptr("/opt/posts/"),
			"DIST_PATH":      sptr("/opt/web/dist"),
			"BLEVE_IDX_PATH": sptr("/opt/blog.bleve"),
			"GIN_MODE":       sptr("release"),
		},
	})

	postsLayer := awslambda.NewLayerVersion(stack, sptr("GuigoesLayer"), &awslambda.LayerVersionProps{
		Code: awslambda.AssetCode_FromAsset(sptr("../../"), &awss3assets.AssetOptions{
			Exclude: &[]*string{
				sptr("cmd"),
				sptr("deployments"),
				sptr("internal"),
				sptr("pkg"),
				sptr("tmp"),
				sptr("web/templates"),
				sptr("web/css"),
				sptr("*.mod"),
				sptr("*.sum"),
				sptr("*.work"),
				sptr(".env"),
				sptr(".gitignore"),
				sptr("*.go"),
				sptr(".air.toml"),
				sptr("*.sh"),
				sptr("*.js"),
				sptr("makefile"),
			},
		}),
	})

	lambda.AddLayers(postsLayer)

	api := awsapigateway.NewLambdaRestApi(stack, sptr("GuigoesApi"), &awsapigateway.LambdaRestApiProps{
		Handler:          lambda,
		BinaryMediaTypes: &[]*string{sptr("*/*")},
	})

	awscdk.NewCfnOutput(stack, sptr("api-gateway-endpoint"),
		&awscdk.CfnOutputProps{
			ExportName: sptr("API-Gateway-Endpoint"),
			Value:      api.Url()})

	return stack
}

type numberType interface {
	~float32 | ~float64 |
		~int | ~int8 | ~int16 | ~int32 | ~int64 |
		~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr
}

func nptr[T numberType](v T) *float64 {
	return jsii.Number(v)
}

func boolptr(b bool) *bool {
	return jsii.Bool(b)
}

func sptr(s string) *string {
	return jsii.String(s)
}

func main() {
	defer jsii.Close()

	app := awscdk.NewApp(nil)

	GuigoesCdkStack(app, "GuigoesStack", &CdkStackProps{
		awscdk.StackProps{
			Env: env(),
		},
	})

	app.Synth(nil)
}

// env determines the AWS environment (account+region) in which our stack is to
// be deployed. For more information see: https://docs.aws.amazon.com/cdk/latest/guide/environments.html
func env() *awscdk.Environment {
	// If unspecified, this stack will be "environment-agnostic".
	// Account/Region-dependent features and context lookups will not work, but a
	// single synthesized template can be deployed anywhere.
	//---------------------------------------------------------------------------
	return nil

	// Uncomment if you know exactly what account and region you want to deploy
	// the stack to. This is the recommendation for production stacks.
	//---------------------------------------------------------------------------
	// return &awscdk.Environment{
	//  Account: jsii.String("123456789012"),
	//  Region:  jsii.String("us-east-1"),
	// }

	// Uncomment to specialize this stack for the AWS Account and Region that are
	// implied by the current CLI configuration. This is recommended for dev
	// stacks.
	//---------------------------------------------------------------------------
	// return &awscdk.Environment{
	//  Account: jsii.String(os.Getenv("CDK_DEFAULT_ACCOUNT")),
	//  Region:  jsii.String(os.Getenv("CDK_DEFAULT_REGION")),
	// }
}
