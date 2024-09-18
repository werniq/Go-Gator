package main

import (
	awscdk "github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsec2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awseks"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsiam"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
)

type AwsSetupStackProps struct {
	awscdk.StackProps
}

func NewGoGatorCdkProjectStack(scope constructs.Construct, id string, props *AwsSetupStackProps) awscdk.Stack {
	stack := awscdk.NewStack(scope, &id, &props.StackProps)

	vpc := awsec2.NewVpc(stack, jsii.String("GoGatorVpc"), &awsec2.VpcProps{
		EnableDnsSupport:   jsii.Bool(true),
		EnableDnsHostnames: jsii.Bool(true),
		MaxAzs:             jsii.Number(2),
		SubnetConfiguration: &[]*awsec2.SubnetConfiguration{
			{
				Name:       jsii.String("GoGatorSubnet1"),
				CidrMask:   jsii.Number(26),
				SubnetType: awsec2.SubnetType_PRIVATE_WITH_EGRESS,
			},
			{
				Name:       jsii.String("GoGatorSubnet2"),
				CidrMask:   jsii.Number(26),
				SubnetType: awsec2.SubnetType_PUBLIC,
			},
		},
	})

	sg := awsec2.NewSecurityGroup(stack, jsii.String("GoGatorSecurityGroup"), &awsec2.SecurityGroupProps{
		Vpc:               vpc,
		SecurityGroupName: jsii.String("Go-Gator-Security-Group"),
		Description:       jsii.String("Allow inbound traffic from port 22 and 443"),
		AllowAllOutbound:  jsii.Bool(true),
	})
	sg.AddIngressRule(awsec2.Peer_AnyIpv4(),
		awsec2.Port_Tcp(jsii.Number(22)),
		jsii.String("Allow SSH"),
		jsii.Bool(true))
	sg.AddIngressRule(awsec2.Peer_AnyIpv4(),
		awsec2.Port_Tcp(jsii.Number(443)),
		jsii.String("Allow HTTPS"),
		jsii.Bool(true))

	routeTable := awsec2.NewCfnRouteTable(stack, jsii.String("GoGatorRouteTable"), &awsec2.CfnRouteTableProps{
		VpcId: vpc.VpcId(),
	})
	awsec2.NewCfnRoute(stack, jsii.String("GoGatorRoute"), &awsec2.CfnRouteProps{
		RouteTableId:         routeTable.Ref(),
		DestinationCidrBlock: jsii.String("0.0.0.0/0"),
		GatewayId:            vpc.InternetGatewayId(),
	})

	subnets := vpc.PublicSubnets()
	subnet1 := (*subnets)[0]
	subnet2 := (*subnets)[1]

	awsec2.NewCfnSubnetRouteTableAssociation(stack, jsii.String("GoGatorSubnet1RouteTableAssociation"), &awsec2.CfnSubnetRouteTableAssociationProps{
		SubnetId:     subnet1.SubnetId(),
		RouteTableId: routeTable.Ref(),
	})
	awsec2.NewCfnSubnetRouteTableAssociation(stack, jsii.String("GoGatorSubnet2RouteTableAssociation"), &awsec2.CfnSubnetRouteTableAssociationProps{
		SubnetId:     subnet2.SubnetId(),
		RouteTableId: routeTable.Ref(),
	})

	role := awsiam.Role_FromRoleArn(stack, jsii.String("EksClusterRole"), jsii.String("arn:aws:iam::406477933661:role/Oleksandr"), nil)

	cluster := awseks.NewCluster(stack, jsii.String("GoGatorCluster"), &awseks.ClusterProps{
		Version:       awseks.KubernetesVersion_V1_30(),
		ClusterName:   jsii.String("GoGatorCluster"),
		Role:          role,
		SecurityGroup: sg,
		Vpc:           vpc,
	})

	awseks.NewCfnAddon(stack, jsii.String("GoGatorCoreDnsAddon"), &awseks.CfnAddonProps{
		ClusterName: cluster.ClusterName(),
		AddonName:   jsii.String("coredns"),
	})

	awseks.NewCfnAddon(stack, jsii.String("GoGatorKubeProxyAddon"), &awseks.CfnAddonProps{
		ClusterName: cluster.ClusterName(),
		AddonName:   jsii.String("kube-proxy"),
	})

	awseks.NewCfnAddon(stack, jsii.String("GoGatorVpcCniAddon"), &awseks.CfnAddonProps{
		ClusterName: cluster.ClusterName(),
		AddonName:   jsii.String("vpc-cni"),
	})

	awseks.NewCfnAddon(stack, jsii.String("GoGatorEksPodIdentityAddon"), &awseks.CfnAddonProps{
		AddonName:   jsii.String("eks-pod-identity-agent"),
		ClusterName: cluster.ClusterName(),
	})

	cluster.AddAutoScalingGroupCapacity(jsii.String("GoGatorNodeGroup"), &awseks.AutoScalingGroupCapacityOptions{
		InstanceType:    awsec2.NewInstanceType(jsii.String("t2.micro")),
		DesiredCapacity: jsii.Number(1),
		MinCapacity:     jsii.Number(1),
		MaxCapacity:     jsii.Number(5),
		MapRole:         jsii.Bool(true),
		VpcSubnets: &awsec2.SubnetSelection{
			AvailabilityZones: jsii.Strings(*subnet1.AvailabilityZone(), *subnet2.AvailabilityZone()),
		},
	})

	return stack
}

func main() {
	defer jsii.Close()

	app := awscdk.NewApp(nil)

	NewGoGatorCdkProjectStack(app, "GoGatorAwsSetupStackV2", &AwsSetupStackProps{
		awscdk.StackProps{
			Env: env(),
		},
	})

	app.Synth(nil)
}

// env determines the AWS environment (account+region) in which our stack is to
// be deployed. For more information see: https://docs.aws.amazon.com/cdk/latest/guide/environments.html
func env() *awscdk.Environment {
	return &awscdk.Environment{
		Account: jsii.String("406477933661"),
		Region:  jsii.String("us-east-2"),
	}
}
