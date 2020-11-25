package aliyun



type STSPolicy struct {
	Version string
	Statement []*STSPolicyStatement
}

type STSPolicyStatement struct {
	Effect string
	Action []string
	Resource []string
}