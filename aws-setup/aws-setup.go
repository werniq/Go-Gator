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

const (
	defaultCidrBlock = "0.0.0.0/0"

	// CoreDnsVersion is the version of the CoreDNS add-on
	CoreDnsVersion = "v1.11.1-eksbuild.18"

	// KubeProxyVersion is the version of the kube-proxy add-on
	KubeProxyVersion = "v1.30.0-eksbuild.3"

	// AmazonVpcCniVersion is the version of the Amazon VPC CNI add-on
	AmazonVpcCniVersion = "v1.18.1-eksbuild.3"

	// PodIdentityVersion is the version of the pod identity add-on
	PodIdentityVersion = "v1.3.2-eksbuild.2"
)

func InternetGateway(stack awscdk.Stack) awsec2.CfnInternetGateway {
	return awsec2.NewCfnInternetGateway(stack,
		jsii.String("GoGatorInternetGateway"),
		&awsec2.CfnInternetGatewayProps{},
	)
}

func Cluster(stack awscdk.Stack, role awsiam.IRole, sg awsec2.SecurityGroup, vpc awsec2.Vpc) awseks.Cluster {
	return awseks.NewCluster(stack, jsii.String("GoGatorCluster"), &awseks.ClusterProps{
		Version:       awseks.KubernetesVersion_V1_30(),
		ClusterName:   jsii.String("GoGatorCluster"),
		Role:          role,
		SecurityGroup: sg,
		Vpc:           vpc,
	})
}

func EksAddon(stack awscdk.Stack, cluster awseks.Cluster, addonName string, addonVersion string, addonId string) {
	awseks.NewCfnAddon(stack, jsii.String(addonId), &awseks.CfnAddonProps{
		ClusterName:  cluster.ClusterName(),
		AddonName:    jsii.String(addonName),
		AddonVersion: jsii.String(addonVersion),
	})
}

func Route(stack awscdk.Stack, routeTable awsec2.CfnRouteTable, igw awsec2.CfnInternetGateway,
	destinationCidrBlock string) {
	awsec2.NewCfnRoute(stack, jsii.String("GoGatorRoute"), &awsec2.CfnRouteProps{
		RouteTableId:         routeTable.Ref(),
		DestinationCidrBlock: jsii.String(destinationCidrBlock),
		GatewayId:            igw.Ref(),
	})
}

func RouteTable(stack awscdk.Stack, vpc awsec2.Vpc) awsec2.CfnRouteTable {
	return awsec2.NewCfnRouteTable(stack, jsii.String("GoGatorRouteTable"), &awsec2.CfnRouteTableProps{
		VpcId: vpc.VpcId(),
	})
}

func Vpc(stack awscdk.Stack) awsec2.Vpc {
	return awsec2.NewVpc(stack, jsii.String("GoGatorVpc"), &awsec2.VpcProps{
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
}

func NewSecurityGroupSecurityGroup(stack awscdk.Stack, vpc awsec2.Vpc) awsec2.SecurityGroup {
	sg := awsec2.NewSecurityGroup(stack, jsii.String("GoGatorSecurityGroup"), &awsec2.SecurityGroupProps{
		Vpc:               vpc,
		SecurityGroupName: jsii.String("Go-Gator-Security-Group"),
		Description:       jsii.String("Allow inbound traffic from port 22, 443 and 80 (SSH, HTTP and HTTPS)"),
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
	sg.AddIngressRule(awsec2.Peer_AnyIpv4(),
		awsec2.Port_Tcp(jsii.Number(80)),
		jsii.String("Allow HTTP"),
		jsii.Bool(true))

	return sg
}

// AddParameters adds the parameters to the stack
func AddParameters(stack awscdk.Stack) {
	awscdk.NewCfnParameter(stack, jsii.String("VpcName"), &awscdk.CfnParameterProps{
		Type:        jsii.String("String"),
		Description: jsii.String("The name of the VPC."),
		Default:     jsii.String("GoGatorVpc"),
	})

	awscdk.NewCfnParameter(stack, jsii.String("CidrBlock"), &awscdk.CfnParameterProps{
		Type:        jsii.String("String"),
		Description: jsii.String("The CIDR block for the VPC."),
		Default:     jsii.String("10.0.1.0/24"),
	})

	awscdk.NewCfnParameter(stack, jsii.String("NodeGroupInstanceType"), &awscdk.CfnParameterProps{
		Type:        jsii.String("String"),
		Description: jsii.String("The EC2 instance type for the node group."),
		Default:     jsii.String("t3.medium"),
		AllowedValues: jsii.Strings(
			"t3.medium",
			"t2.medium",
			"t2.small",
		),
	})

	awscdk.NewCfnParameter(stack, jsii.String("MinimalNodeGroupSize"), &awscdk.CfnParameterProps{
		Type:        jsii.String("Number"),
		Description: jsii.String("The minimal number of nodes in the node group."),
		Default:     jsii.Number(1),
		MinValue:    jsii.Number(1),
		MaxValue:    jsii.Number(5),
	})

	awscdk.NewCfnParameter(stack, jsii.String("MaximalNodeGroupSize"), &awscdk.CfnParameterProps{
		Type:        jsii.String("Number"),
		Description: jsii.String("The maximal number of nodes in the node group."),
		Default:     jsii.Number(5),
		MinValue:    jsii.Number(1),
		MaxValue:    jsii.Number(10),
	})

	awscdk.NewCfnParameter(stack, jsii.String("NodeGroupDesiredCapacity"), &awscdk.CfnParameterProps{
		Type:        jsii.String("Number"),
		Description: jsii.String("The desired number of nodes in the node group."),
		Default:     jsii.Number(1),
	})

	awscdk.NewCfnParameter(stack, jsii.String("KubernetesVersion"), &awscdk.CfnParameterProps{
		Type:        jsii.String("String"),
		Description: jsii.String("The Kubernetes version for the cluster."),
		Default:     jsii.String("1.30"),
		AllowedValues: jsii.Strings(
			"1.24",
			"1.25",
			"1.26",
			"1.27",
			"1.28",
			"1.29",
			"1.30",
		),
	})
}

func NewGoGatorCdkProjectStack(scope constructs.Construct, id string, props *AwsSetupStackProps) awscdk.Stack {
	stack := awscdk.NewStack(scope, &id, &props.StackProps)

	AddParameters(stack)

	vpc := Vpc(stack)
	sg := NewSecurityGroupSecurityGroup(stack, vpc)

	igw := InternetGateway(stack)
	awsec2.NewCfnVPCGatewayAttachment(stack, jsii.String("GoGatorAttachGateway"), &awsec2.CfnVPCGatewayAttachmentProps{
		VpcId:             vpc.VpcId(),
		InternetGatewayId: igw.Ref(),
	})

	routeTable := RouteTable(stack, vpc)
	Route(stack, routeTable, igw, defaultCidrBlock)

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

	cluster := Cluster(stack, role, sg, vpc)

	EksAddon(stack, cluster, "coredns", CoreDnsVersion, "GoGatorCodeDnsAddon")
	EksAddon(stack, cluster, "kube-proxy", KubeProxyVersion, "GoGatorKubeProxyAddon")
	EksAddon(stack, cluster, "vpc-cni", AmazonVpcCniVersion, "GoGatorVpcCniAddon")
	EksAddon(stack, cluster, "eks-pod-identity-agent", PodIdentityVersion, "GoGatorPodIdentityAddon")

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

	NewGoGatorCdkProjectStack(app, "GoGatorAwsSetupStack", &AwsSetupStackProps{
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
