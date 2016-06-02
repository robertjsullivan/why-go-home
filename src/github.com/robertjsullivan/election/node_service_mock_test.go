package election

type NodeServiceMock struct {
	InitiatedVoting bool
	HeartBeating bool
	SetupHttpEndpoints bool
}

func (ns *NodeServiceMock) InitiateVoting(node *Node) bool {
	ns.InitiatedVoting = true;
	return true;
}

func (ns *NodeServiceMock) Heartbeater(node *Node){
	ns.HeartBeating = true;
}

func (ns *NodeServiceMock) SetupHttpEndpoint(node *Node){
	ns.SetupHttpEndpoints = true;
}