package hub

var (
	SubChannel = make(chan [2]interface{}, 2)
	NewChannel = make(chan [2]interface{}, 2)
	PubChannel = make(chan [2]interface{}, 2)
	JoinChan   = make(chan User)
	Categories = make([]string, 0, 10)
)
