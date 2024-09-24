package main

import (
	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/aws-cdk-go/awscdk/v2/assertions"
	"github.com/aws/jsii-runtime-go"
	"testing"
)

// example tests. To run these tests, uncomment this file along with the
// example resource in aws-cdk-setup_test.go
func TestAwsCdkSetupStack(t *testing.T) {
	// GIVEN
	app := awscdk.NewApp(nil)

	// WHEN
	stack := NewGoGatorCdkProjectStack(app, "MyStack", &AwsSetupStackProps{
		awscdk.StackProps{
			Env: env(),
		},
	})

	// THEN
	template := assertions.Template_FromStack(stack, nil)

	template.HasResourceProperties(jsii.String("AWS::EC2::VPC"), map[string]interface{}{
		"VisibilityTimeout": 300,
	})
}
