package config

func (bc *BlockerConfig) UpdateBlockList(blockList []string) {
	bc.BlockList = blockList
}

func (bc *BlockerConfig) AddHostToBlockList(blockHost string) {
	bc.BlockList = append(bc.BlockList, blockHost)
}
