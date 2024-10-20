package general

import (
  // system packages
	"log"
	"time"
	"math/rand"
)

func CountCheckedUsers(m map[string]bool) int {
  cnt := 0

  for _, checked := range m {
    if checked {
      cnt++
    }
  }

  return cnt
}

func Contains(slice []string, value string) bool {
	for _, v := range slice {
		if v == value {
			return true
		}
	}
	return false
}

func GetRandomUser(userList []string) string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	// Generate a random number between 0 and 100
  // Intn(n) returns a random integer from 0 to n-1, so 101 gives 0 to 100
	randomNumber := r.Intn(len(userList)) 

  // debug
  log.Printf("Selected random user: list[%d] = %s\n", randomNumber, userList[randomNumber])

	return userList[randomNumber]
}
