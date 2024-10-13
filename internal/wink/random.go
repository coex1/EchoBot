package wink

// system packages
import (
	"log"
	"time"
	"math/rand"
)

func GetRandomUser(userList []string) string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	// Generate a random number between 0 and 100
  // Intn(n) returns a random integer from 0 to n-1, so 101 gives 0 to 100
	randomNumber := r.Intn(len(userList)) 

  // debug
  log.Println("Selected random user: list[%d] = %s", randomNumber, userList[randomNumber])

	return userList[randomNumber]
}
