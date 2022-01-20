package config

type BlockerConfig struct {
	UpstreamDns  string
	BlockList    []string
	Addr         string
	Enabled      bool
	QueryChannel chan QueryEvent
}

type QueryEvent struct {
	Hostname string
	Blocked  bool
}
