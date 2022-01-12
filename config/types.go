package config

type BlockerConfig struct {
	UpstreamDns string
	BlockList   []string
	Addr        string
	Enabled     bool
}
