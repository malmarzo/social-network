package utils

import ("math/rand"
"time"
"fmt"
// "log"
// "net/http"
// "social-network/pkg/db/queries"
)

func GenerateSessionID() string {
	rand.Seed(time.Now().UnixNano())
	return fmt.Sprintf("sess_%d", rand.Intn(1000000000))
}
