package config

func (bc *BlockerConfig) UpdateBlockList(blockList []string) {
	bc.BlockList = blockList
}

func (bc *BlockerConfig) AddHostToBlockList(blockHost string) {
	bc.BlockList = append(bc.BlockList, blockHost)
}

func (bc *BlockerConfig) ToggleBlocker() {
	bc.Enabled = !bc.Enabled
}

// Sends query event to channel if channel is not nil
func (bc *BlockerConfig) sendQueryEvent(hostname string, blocked bool) {
	if bc.QueryChannel != nil {
		event := QueryEvent{
			hostname: hostname,
			blocked:  blocked,
		}
		bc.QueryChannel <- event
	}
}
