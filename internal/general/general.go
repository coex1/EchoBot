package general

import (
  // system packages
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

// return a random number between 'min' and 'max'
func Random(min int, max int) (r int) {
  if min > max {
    return -1
  }

	random := rand.New(rand.NewSource(time.Now().UnixNano()))

  // Intn(n) -> (0 to n-1),
  // Ex: 101 gives 0 to 100
	r = random.Intn(max - min + 1) + min

  return
}
