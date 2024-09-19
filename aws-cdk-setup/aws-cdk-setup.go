package main

import (
	awscdk "github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsec2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awseks"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
)

type AwsSetupStackProps struct {
	awscdk.StackProps
}

const (
	// CoreDnsAddonVersion is the version of the CoreDNS addon to install.
	CoreDnsAddonVersion = "v1.11.1-eksbuild.8"

	// KubeProxyAddonVersion is the version of the kube-proxy addon to install.
	KubeProxyAddonVersion = "v1.30.0-eksbuild.3"

	// AmazonVpcCniAddonVersion is the version of the Amazon VPC CNI addon to install.
	AmazonVpcCniAddonVersion = "v1.18.1-eksbuild.3"

	// PodIdentityAddonVersion is the version of the EKS pod identity addon to install.
	PodIdentityAddonVersion = "v1.3.2-eksbuild.2"
)

//func AddParameters(stack awscdk.Stack) {
//	awscdk.NewCfnParameter(stack, jsii.String("CidrBlock"), &awscdk.CfnParameterProps{
//		Type:        jsii.String("String"),
//		Description: jsii.String("The CIDR block for the VPC. This should be a valid private (RFC 1918) CIDR block."),
//		Default:     jsii.String("10.0.1.0/24"),
//	})
//
//	awscdk.NewCfnParameter(stack, jsii.String("NodeGroupInstanceType"), &awscdk.CfnParameterProps{
//		Type:        jsii.String("String"),
//		Description: jsii.String("The EC2 instance type for the node group."),
//		Default:     jsii.String("t3.medium"),
//		AllowedValues: &[]*string{
//			jsii.String("t3.medium"),
//			jsii.String("t2.medium"),
//			jsii.String("t2.small"),
//		},
//	})
//
//	minimalNodeGroupSize := awscdk.NewCfnParameter(stack, jsii.String("MinimalNodeGroupSize"), &awscdk.CfnParameterProps{
//		Type:        jsii.String("Number"),
//		Description: jsii.String("The minimal number of nodes in the node group."),
//		Default:     jsii.Number(1),
//		MinValue:    jsii.Number(1),
//		MaxValue:    jsii.Number(5),
//	})
//
//	maximalNodeGroupSize := awscdk.NewCfnParameter(stack, jsii.String("MaximalNodeGroupSize"), &awscdk.CfnParameterProps{
//		Type:        jsii.String("Number"),
//		Description: jsii.String("The maximal number of nodes in the node group."),
//		Default:     jsii.Number(5),
//		MinValue:    jsii.Number(1),
//		MaxValue:    jsii.Number(10),
//	})
//
//	awscdk.NewCfnParameter(stack, jsii.String("NodeGroupDesiredCapacity"), &awscdk.CfnParameterProps{
//		Type:        jsii.String("Number"),
//		Description: jsii.String("The desired number of nodes in the node group."),
//		Default:     jsii.Number(1),
//		MinValue:    minimalNodeGroupSize.ValueAsNumber(),
//		MaxValue:    maximalNodeGroupSize.ValueAsNumber(),
//	})
//
//	awscdk.NewCfnParameter(stack, jsii.String("KubernetesVersion"), &awscdk.CfnParameterProps{
//		Type:        jsii.String("String"),
//		Description: jsii.String("The Kubernetes version for the cluster."),
//		Default:     jsii.String("1.30"),
//		AllowedValues: &[]*string{
//			jsii.String("1.24"),
//			jsii.String("1.25"),
//			jsii.String("1.26"),
//			jsii.String("1.27"),
//			jsii.String("1.28"),
//			jsii.String("1.29"),
//			jsii.String("1.30"),
//		},
//	})
//}

func NewGoGatorCdkProjectStack(scope constructs.Construct, id string, props *AwsSetupStackProps) awscdk.Stack {
	stack := awscdk.NewStack(scope, &id, &props.StackProps)

	//AddParameters(stack)

	vpc := awsec2.NewVpc(stack, jsii.String("GoGatorVpc"), &awsec2.VpcProps{
		EnableDnsSupport:   jsii.Bool(true),
		EnableDnsHostnames: jsii.Bool(true),
		MaxAzs:             jsii.Number(2),
		SubnetConfiguration: &[]*awsec2.SubnetConfiguration{
			{
				Name:       jsii.String("GoGatorSubnet1"),
				CidrMask:   jsii.Number(24),
				SubnetType: awsec2.SubnetType_PRIVATE_WITH_EGRESS,
			},
			{
				Name:       jsii.String("GoGatorSubnet2"),
				CidrMask:   jsii.Number(24),
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

	//clusterCreationRole := awsiam.NewRole(stack, jsii.String("GoGatorRole"), &awsiam.RoleProps{
	//	AssumedBy: awsiam.NewServicePrincipal(jsii.String("eks.amazonaws.com"), &awsiam.ServicePrincipalOpts{
	//		Region: jsii.String("us-east-2"),
	//	}),
	//	Description: jsii.String("Role to create EKS cluster"),
	//})
	//clusterCreationRole.AddManagedPolicy(awsiam.ManagedPolicy_FromManagedPolicyName(stack,
	//	jsii.String("GoGatorAmazonEKSClusterPolicy"),
	//	jsii.String("AmazonEKSClusterPolicy")))

	subnet1RouteTable := awsec2.NewCfnRouteTable(stack, jsii.String("GoGatorRouteTable1"), &awsec2.CfnRouteTableProps{
		VpcId: vpc.VpcId(),
	})
	awsec2.NewCfnRoute(stack, jsii.String("GoGatorRoute1"), &awsec2.CfnRouteProps{
		RouteTableId:         subnet1RouteTable.Ref(),
		DestinationCidrBlock: jsii.String("0.0.0.0/0"),
		GatewayId:            vpc.InternetGatewayId(),
	})
	subnet2RouteTable := awsec2.NewCfnRouteTable(stack, jsii.String("GoGatorRouteTable2"), &awsec2.CfnRouteTableProps{
		VpcId: vpc.VpcId(),
	})
	awsec2.NewCfnRoute(stack, jsii.String("GoGatorRoute2"), &awsec2.CfnRouteProps{
		RouteTableId:         subnet2RouteTable.Ref(),
		DestinationCidrBlock: jsii.String("0.0.0.0/0"),
		GatewayId:            vpc.InternetGatewayId(),
	})

	subnets := vpc.PublicSubnets()
	subnet1 := (*subnets)[0]
	subnet2 := (*subnets)[1]

	awsec2.NewCfnSubnetRouteTableAssociation(stack, jsii.String("GoGatorSubnet1RouteTableAssociation"), &awsec2.CfnSubnetRouteTableAssociationProps{
		SubnetId:     subnet1.SubnetId(),
		RouteTableId: subnet1RouteTable.Ref(),
	})
	awsec2.NewCfnSubnetRouteTableAssociation(stack, jsii.String("GoGatorSubnet2RouteTableAssociation"), &awsec2.CfnSubnetRouteTableAssociationProps{
		SubnetId:     subnet2.SubnetId(),
		RouteTableId: subnet2RouteTable.Ref(),
	})

	cluster := awseks.NewCluster(stack, jsii.String("GoGatorCluster"), &awseks.ClusterProps{
		Version:     awseks.KubernetesVersion_V1_30(),
		ClusterName: jsii.String("GoGatorCluster"),
		//Role:          clusterCreationRole,
		SecurityGroup: sg,
		Vpc:           vpc,
	})

	awseks.NewCfnAddon(stack, jsii.String("GoGatorCoreDnsAddon"), &awseks.CfnAddonProps{
		ClusterName:  cluster.ClusterName(),
		AddonName:    jsii.String("coredns"),
		AddonVersion: jsii.String(CoreDnsAddonVersion),
	})

	awseks.NewCfnAddon(stack, jsii.String("GoGatorKubeProxyAddon"), &awseks.CfnAddonProps{
		ClusterName:  cluster.ClusterName(),
		AddonName:    jsii.String("kube-proxy"),
		AddonVersion: jsii.String(KubeProxyAddonVersion),
	})

	awseks.NewCfnAddon(stack, jsii.String("GoGatorVpcCniAddon"), &awseks.CfnAddonProps{
		ClusterName:  cluster.ClusterName(),
		AddonName:    jsii.String("vpc-cni"),
		AddonVersion: jsii.String(AmazonVpcCniAddonVersion),
	})

	awseks.NewCfnAddon(stack, jsii.String("GoGatorEksPodIdentityAddon"), &awseks.CfnAddonProps{
		AddonName:    jsii.String("eks-pod-identity-agent"),
		ClusterName:  cluster.ClusterName(),
		AddonVersion: jsii.String(PodIdentityAddonVersion),
	})

	cluster.AddAutoScalingGroupCapacity(jsii.String("GoGatorNodeGroup"), &awseks.AutoScalingGroupCapacityOptions{
		DesiredCapacity: jsii.Number(1),
		MaxCapacity:     jsii.Number(5),
		MinCapacity:     jsii.Number(1),
		VpcSubnets: &awsec2.SubnetSelection{
			AvailabilityZones: jsii.Strings(*subnet1.AvailabilityZone(), *subnet2.AvailabilityZone()),
		},
		InstanceType: awsec2.NewInstanceType(jsii.String("t2.micro")),
		MapRole:      jsii.Bool(true),
	})

	return stack
}

func main() {
	defer jsii.Close()

	app := awscdk.NewApp(nil)

	NewGoGatorCdkProjectStack(app, "NewGoGatorAwsSetupStack", &AwsSetupStackProps{
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
