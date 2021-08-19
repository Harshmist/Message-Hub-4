package variables

var (
	SubChannel = make(chan [2]interface{}, 2)
	NewChannel = make(chan [2]interface{}, 2)
	PubChannel = make(chan [2]interface{}, 2)
	JoinChan   = make(chan User)
)

type User struct {
	Name    string
	Address chan string
}
