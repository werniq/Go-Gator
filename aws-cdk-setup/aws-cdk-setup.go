package main

import (
	awscdk "github.com/aws/aws-cdk-go/awscdk/v2"
	ec2 "github.com/aws/aws-cdk-go/awscdk/v2/awsec2"
	eks "github.com/aws/aws-cdk-go/awscdk/v2/awseks"
	iam "github.com/aws/aws-cdk-go/awscdk/v2/awsiam"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
)

type AwsSetupStackProps struct {
	awscdk.StackProps
}

const (
	// CoreDnsAddonVersion is the version of the CoreDNS addon to install.
	CoreDnsAddonVersion = "v1.11.3-eksbuild.1"

	// KubeProxyAddonVersion is the version of the kube-proxy addon to install.
	KubeProxyAddonVersion = "v1.30.0-eksbuild.3"

	// AmazonVpcCniAddonVersion is the version of the Amazon VPC CNI addon to install.
	AmazonVpcCniAddonVersion = "v1.18.1-eksbuild.3"

	// PodIdentityAddonVersion is the version of the EKS pod identity addon to install.
	PodIdentityAddonVersion = "v1.3.2-eksbuild.2"
)

func AddParameters(stack awscdk.Stack) {
	awscdk.NewCfnParameter(stack, jsii.String("NodeGroupInstanceType"), &awscdk.CfnParameterProps{
		Type:        jsii.String("String"),
		Description: jsii.String("The EC2 instance type for the node group."),
		Default:     jsii.String("t3.medium"),
		AllowedValues: &[]*string{
			jsii.String("t3.medium"),
			jsii.String("t2.medium"),
			jsii.String("t2.small"),
		},
	})

	_ = awscdk.NewCfnParameter(stack, jsii.String("MinimalNodeGroupSize"), &awscdk.CfnParameterProps{
		Type:        jsii.String("Number"),
		Description: jsii.String("The minimal number of nodes in the node group."),
		Default:     jsii.Number(1),
		MinValue:    jsii.Number(1),
		MaxValue:    jsii.Number(5),
	})

	_ = awscdk.NewCfnParameter(stack, jsii.String("NodeDiskSize"), &awscdk.CfnParameterProps{
		Type:        jsii.String("Number"),
		Description: jsii.String("The minimal number of nodes in the node group."),
		MinValue:    jsii.Number(1 * 1024 * 1024 * 1024),
		MaxValue:    jsii.Number(50 * 1024 * 1024 * 1024),
		Default:     jsii.Number(10 * 1024 * 1024 * 1024),
	})

	_ = awscdk.NewCfnParameter(stack, jsii.String("MaximalNodeGroupSize"), &awscdk.CfnParameterProps{
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
		//MinValue:    minimalNodeGroupSize.ValueAsNumber(),
		//MaxValue:    maximalNodeGroupSize.ValueAsNumber(),
	})

	awscdk.NewCfnParameter(stack, jsii.String("KubernetesVersion"), &awscdk.CfnParameterProps{
		Type:        jsii.String("String"),
		Description: jsii.String("The Kubernetes version for the cluster."),
		Default:     jsii.String("1.30"),
		AllowedValues: &[]*string{
			jsii.String("1.24"),
			jsii.String("1.25"),
			jsii.String("1.26"),
			jsii.String("1.27"),
			jsii.String("1.28"),
			jsii.String("1.29"),
			jsii.String("1.30"),
		},
	})
}

func NewGoGatorCdkProjectStack(scope constructs.Construct, id string, props *AwsSetupStackProps) awscdk.Stack {
	stack := awscdk.NewStack(scope, &id, &props.StackProps)

	AddParameters(stack)

	vpc := ec2.NewVpc(stack, jsii.String("GoGatorVpc"), &ec2.VpcProps{
		EnableDnsSupport:   jsii.Bool(true),
		EnableDnsHostnames: jsii.Bool(true),
		MaxAzs:             jsii.Number(2),
		SubnetConfiguration: &[]*ec2.SubnetConfiguration{
			{
				Name:       jsii.String("GoGatorSubnet1"),
				CidrMask:   jsii.Number(24),
				SubnetType: ec2.SubnetType_PRIVATE_WITH_EGRESS,
			},
			{
				Name:       jsii.String("GoGatorSubnet2"),
				CidrMask:   jsii.Number(24),
				SubnetType: ec2.SubnetType_PUBLIC,
			},
		},
	})

	sg := ec2.NewSecurityGroup(stack, jsii.String("GoGatorSecurityGroup"), &ec2.SecurityGroupProps{
		Vpc:               vpc,
		SecurityGroupName: jsii.String("Go-Gator-Security-Group"),
		Description:       jsii.String("Allow inbound traffic from port 22 and 443"),
		AllowAllOutbound:  jsii.Bool(true),
	})
	sg.AddIngressRule(ec2.Peer_AnyIpv4(),
		ec2.Port_Tcp(jsii.Number(22)),
		jsii.String("Allow SSH"),
		jsii.Bool(true))
	sg.AddIngressRule(ec2.Peer_AnyIpv4(),
		ec2.Port_Tcp(jsii.Number(443)),
		jsii.String("Allow HTTPS"),
		jsii.Bool(true))

	subnet2RouteTable := ec2.NewCfnRouteTable(stack, jsii.String("GoGatorRouteTable2"), &ec2.CfnRouteTableProps{
		VpcId: vpc.VpcId(),
	})
	ec2.NewCfnRoute(stack, jsii.String("GoGatorRoute2"), &ec2.CfnRouteProps{
		RouteTableId:         subnet2RouteTable.Ref(),
		DestinationCidrBlock: jsii.String("0.0.0.0/0"),
		GatewayId:            vpc.InternetGatewayId(),
	})

	subnets := vpc.PublicSubnets()
	subnet := (*subnets)[0]

	ec2.NewCfnSubnetRouteTableAssociation(stack, jsii.String("GoGatorSubnet2RouteTableAssociation"), &ec2.CfnSubnetRouteTableAssociationProps{
		SubnetId:     subnet.SubnetId(),
		RouteTableId: subnet2RouteTable.Ref(),
	})

	eksRole := iam.NewRole(stack, jsii.String("EksClusterRole"), &iam.RoleProps{
		AssumedBy: iam.NewServicePrincipal(jsii.String("eks.amazonaws.com"), nil),
		ManagedPolicies: &[]iam.IManagedPolicy{
			iam.ManagedPolicy_FromAwsManagedPolicyName(jsii.String("AmazonEKSClusterPolicy")),
			iam.ManagedPolicy_FromAwsManagedPolicyName(jsii.String("AmazonEKSServicePolicy")),
		},
	})

	cluster := eks.NewCluster(stack, jsii.String("NewsAggregatorCluster"), &eks.ClusterProps{
		ClusterName:     jsii.String("NewsGoGatorCluster"),
		Version:         eks.KubernetesVersion_V1_30(),
		Vpc:             vpc,
		Role:            eksRole,
		DefaultCapacity: jsii.Number(0),
		EndpointAccess:  eks.EndpointAccess_PUBLIC_AND_PRIVATE(),
	})

	userName := "oleksandr"
	iamUserArn := "arn:aws:iam::406477933661:user/" + userName
	cluster.AwsAuth().AddUserMapping(iam.User_FromUserArn(stack, jsii.String(userName), jsii.String(iamUserArn)), &eks.AwsAuthMapping{
		Username: jsii.String(userName),
		Groups: &[]*string{
			jsii.String("system:masters"),
		},
	})

	eks.NewCfnAddon(stack, jsii.String("GoGatorCoreDnsAddon"), &eks.CfnAddonProps{
		ClusterName:  cluster.ClusterName(),
		AddonName:    jsii.String("coredns"),
		AddonVersion: jsii.String(CoreDnsAddonVersion),
	})

	eks.NewCfnAddon(stack, jsii.String("GoGatorKubeProxyAddon"), &eks.CfnAddonProps{
		ClusterName:  cluster.ClusterName(),
		AddonName:    jsii.String("kube-proxy"),
		AddonVersion: jsii.String(KubeProxyAddonVersion),
	})

	eks.NewCfnAddon(stack, jsii.String("GoGatorVpcCniAddon"), &eks.CfnAddonProps{
		ClusterName:  cluster.ClusterName(),
		AddonName:    jsii.String("vpc-cni"),
		AddonVersion: jsii.String(AmazonVpcCniAddonVersion),
	})

	eks.NewCfnAddon(stack, jsii.String("GoGatorEksPodIdentityAddon"), &eks.CfnAddonProps{
		AddonName:    jsii.String("eks-pod-identity-agent"),
		ClusterName:  cluster.ClusterName(),
		AddonVersion: jsii.String(PodIdentityAddonVersion),
	})

	nodeGroupRole := iam.NewRole(stack, jsii.String("node-group-role"), &iam.RoleProps{
		AssumedBy:   iam.NewServicePrincipal(jsii.String("ec2.amazonaws.com"), nil),
		Description: jsii.String("Role for EKS Node Group"),
		ManagedPolicies: &[]iam.IManagedPolicy{
			iam.ManagedPolicy_FromAwsManagedPolicyName(jsii.String("AmazonEC2ContainerRegistryReadOnly")),
			iam.ManagedPolicy_FromAwsManagedPolicyName(jsii.String("AmazonEC2FullAccess")),
			iam.ManagedPolicy_FromAwsManagedPolicyName(jsii.String("AmazonEKS_CNI_Policy")),
			iam.ManagedPolicy_FromAwsManagedPolicyName(jsii.String("AmazonEKSWorkerNodePolicy")),
		},
		RoleName: jsii.String("GoGatorNodeGroupRoleAdmin"),
	})

	cluster.AddNodegroupCapacity(jsii.String("GoGatorNodeGroup"), &eks.NodegroupOptions{
		DesiredSize:   jsii.Number(1),
		MaxSize:       jsii.Number(5),
		MinSize:       jsii.Number(1),
		NodegroupName: jsii.String("GoGatorNodeGroup"),
		NodeRole:      nodeGroupRole,
		RemoteAccess: &eks.NodegroupRemoteAccess{
			SshKeyName: jsii.String("Go-Gator-Client-Keys"),
		},
		Subnets: &ec2.SubnetSelection{
			Subnets: subnets,
		},
	})

	//cluster.AddHelmChart(jsii.String("GoGatorHelmChart-Go-Gator"), &eks.HelmChartOptions{
	//	Repository:      jsii.String("406477933661.dkr.ecr.us-east-2.amazonaws.com"),
	//	CreateNamespace: jsii.Bool(true),
	//	Namespace:       jsii.String("go-gator"),
	//	Chart:           jsii.String("qniw984/go-gator:0.1.0"),
	//})
	//
	//cluster.AddHelmChart(jsii.String("GoGatorHelmChart-Day-0"), &eks.HelmChartOptions{
	//	Repository: jsii.String("406477933661.dkr.ecr.us-east-2.amazonaws.com"),
	//	Chart:      jsii.String("qniw984/day-0:0.1.0"),
	//})

	return stack
}

func main() {
	defer jsii.Close()

	app := awscdk.NewApp(nil)

	NewGoGatorCdkProjectStack(app, "NewsAggregatorAwsSetupStack", &AwsSetupStackProps{
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
