package cfg

import "fmt"

var (
	user     = "gastonpalomeque"
	password = "ghpalo21"
	host     = "localhost"
	port     = "6000"
	name     = "joblib"
	sslmode  = "disable"
)

// URL specifies db connection url
var URL = fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s", user, password, host, port, name, sslmode)
