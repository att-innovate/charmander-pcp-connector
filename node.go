package pcp

type Node struct {
	domains map[string]*Domain
	client  *Client
}

func NewNode(host string, port uint16) (*Node, error) {
	context := NewContext("local")
	node := &Node{
		client: NewClient(host, port, context),
	}
	err := node.client.RefreshContext()

	if err != nil {
		return nil, err
	}

	node.domains = make(map[string]*Domain)
	return node, nil
}

func (n *Node) Domain(name string) *Domain {
	var domain *Domain

	domain, ok := n.domains[name]
	if !ok {

	}

	return domain
}

func (n *Node) SetLogLevel(level int) {
	n.client.SetLogLevel(level)
}

type Domain struct {
	Name      string
	instances []*MetricInstance
}

type Instance struct {
	metrics []*Metric
}

func (d *Domain) Instance(name string) *Instance {
	instance := Instance{}
	return &instance
}

func (i *Instance) Metric(name string) *Metric {
	return &Metric{}
}
