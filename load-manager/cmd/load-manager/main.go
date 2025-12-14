package main 

import ( 
	"fmt"
)

func main() {
	n := 5
	for i := range(n) {
		form := fmt.Sprintf("%d", i)
		fmt.Println(form)
	}
}

