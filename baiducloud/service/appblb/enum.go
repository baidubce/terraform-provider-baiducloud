package appblb

const (
	ProtocolTCP   = "TCP"
	ProtocolUDP   = "UDP"
	ProtocolHTTP  = "HTTP"
	ProtocolHTTPS = "HTTPS"
	ProtocolICMP  = "ICMP"
)

var IpGroupPolicyTypes = []string{
	ProtocolTCP,
	ProtocolUDP,
	ProtocolHTTP,
	ProtocolHTTPS,
}

var IpGroupHealthCheckTypes = []string{
	ProtocolTCP,
	ProtocolUDP,
	ProtocolHTTP,
	ProtocolHTTPS,
	ProtocolICMP,
}
