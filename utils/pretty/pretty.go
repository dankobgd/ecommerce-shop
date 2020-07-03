package pretty

import (
	"encoding/json"
	"fmt"
)

// PrintJSON prints json string
func PrintJSON(i interface{}) {
	s, _ := json.MarshalIndent(i, "", "  ")
	fmt.Println(string(s))
}
