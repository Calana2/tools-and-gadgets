Module for brute-force attacks 

Example of usage

``` go
package main

import (
	"fmt"
	"main/cracker"
	"sync"
)

func main (){
	c, err := cracker.NewCombiner("", 1, 3)
	if err != nil {
		fmt.Println("Error, can't create combinator")
    return
	}

	passPipe := make(chan string, 1)
	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		wg.Wait()
		close(passPipe)
	}()

	go c.GenerateToPipe(passPipe, &wg)
	for p := range passPipe {
		fmt.Println(p)
	}
}

```
