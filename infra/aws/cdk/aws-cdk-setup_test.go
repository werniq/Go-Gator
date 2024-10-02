package main

import (
	"testing"

	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/aws-cdk-go/awscdk/v2/assertions"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsec2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awseks"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsiam"
)

func TestGoGatorCdkProjectStack(t *testing.T) {
	app := awscdk.NewApp(nil)

	stack := NewGoGatorCdkProjectStack(app, "TestStack", &AwsSetupStackProps{
		awscdk.StackProps{
			Env: env(),
		},
	})
	template := assertions.Template_FromStack(stack, nil)

	template.HasResourceProperties(awsec2.CfnVPC_CFN_RESOURCE_TYPE_NAME(), map[string]interface{}{
		"EnableDnsSupport":   true,
		"EnableDnsHostnames": true,
	})

	template.HasResourceProperties(awsec2.CfnSecurityGroup_CFN_RESOURCE_TYPE_NAME(), map[string]interface{}{
		"SecurityGroupIngress": []interface{}{
			map[string]interface{}{
				"IpProtocol": "tcp",
				"FromPort":   22,
				"ToPort":     22,
				"CidrIp":     "0.0.0.0/0",
			},
			map[string]interface{}{
				"IpProtocol": "tcp",
				"FromPort":   443,
				"ToPort":     443,
				"CidrIp":     "0.0.0.0/0",
			},
		},
	})

	template.HasResource(awsiam.CfnRole_CFN_RESOURCE_TYPE_NAME(), map[string]interface{}{})

	template.HasResource(awseks.CfnAddon_CFN_RESOURCE_TYPE_NAME(), map[string]interface{}{
		"Properties": map[string]interface{}{
			"AddonName":    "coredns",
			"AddonVersion": "v1.11.3-eksbuild.1",
			"ClusterName": map[string]string{
				"Ref": "NewsAggregatorClusterDDE6C3E5",
			},
		},
	})

	template.HasResource(awseks.CfnAddon_CFN_RESOURCE_TYPE_NAME(), map[string]interface{}{
		"Properties": map[string]interface{}{
			"AddonName":    "kube-proxy",
			"AddonVersion": "v1.30.0-eksbuild.3",
			"ClusterName": map[string]string{
				"Ref": "NewsAggregatorClusterDDE6C3E5",
			},
		},
	})

	template.HasResource(awseks.CfnAddon_CFN_RESOURCE_TYPE_NAME(), map[string]interface{}{
		"Properties": map[string]interface{}{
			"AddonName":    "eks-pod-identity-agent",
			"AddonVersion": "v1.3.2-eksbuild.2",
			"ClusterName": map[string]string{
				"Ref": "NewsAggregatorClusterDDE6C3E5",
			},
		},
	})

	template.HasResource(awseks.CfnAddon_CFN_RESOURCE_TYPE_NAME(), map[string]interface{}{
		"Properties": map[string]interface{}{
			"AddonName":    "vpc-cni",
			"AddonVersion": "v1.18.1-eksbuild.3",
			"ClusterName": map[string]string{
				"Ref": "NewsAggregatorClusterDDE6C3E5",
			},
		},
	})

	template.HasResourceProperties(awseks.CfnNodegroup_CFN_RESOURCE_TYPE_NAME(), map[string]interface{}{
		"NodegroupName": "GoGatorNodeGroup",
		"ScalingConfig": map[string]interface{}{
			"DesiredSize": 1,
			"MinSize":     1,
			"MaxSize":     5,
		},
	})
}
