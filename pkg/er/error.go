package er

import (
	"fmt"
	"os"
)

func CheckErr(msg string, err error) {
	if err != nil {
		fmt.Printf("Error message: %s", msg)
		fmt.Printf("Error: %v", err)
		os.Exit(1)
	}
}
