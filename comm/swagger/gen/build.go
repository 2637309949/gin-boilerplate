package gen

// Build builds swagger json file  for given searchDir and mainAPIFile. Returns json
func Build(config *Config) {
	go New().Build(config)
}
