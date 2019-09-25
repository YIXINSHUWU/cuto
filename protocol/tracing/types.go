package tracing


// Configuration for memory dump. Used only when "memory-infra" category is enabled.
type MemoryDumpConfig struct {
	
}	

// 
type TraceConfig struct {
	
	// Controls how the trace buffer stores data.
	// Possible value: recordUntilFull,recordContinuously,recordAsMuchAsPossible,echoToConsole,
	RecordMode	string	`json:"recordMode"`
	
	// Turns on JavaScript stack sampling.
	
	EnableSampling	bool	`json:"enableSampling"`
	
	// Turns on system tracing.
	
	EnableSystrace	bool	`json:"enableSystrace"`
	
	// Turns on argument filter.
	
	EnableArgumentFilter	bool	`json:"enableArgumentFilter"`
	
	// Included category filters.
	
	IncludedCategories	[]string	`json:"includedCategories"`
	
	// Excluded category filters.
	
	ExcludedCategories	[]string	`json:"excludedCategories"`
	
	// Configuration to synthesize the delays in tracing.
	
	SyntheticDelays	[]string	`json:"syntheticDelays"`
	
	// Configuration for memory dump triggers. Used only when "memory-infra" category is enabled.
	
	MemoryDumpConfig	MemoryDumpConfig	`json:"memoryDumpConfig"`
	
}	

// Data format of a trace. Can be either the legacy JSON format or the
	// protocol buffer format. Note that the JSON format will be deprecated soon.
type StreamFormat string	

// Compression type to use for traces returned via streams.
type StreamCompression string	

