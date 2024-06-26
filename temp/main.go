package main

import "fmt"
import "strings"

func main() {
	links := []string{"0-1", "0-5", "5-2", "2-1", "0-6","6-1" }
	start := "0"
	end := "1"

	for i := 0; i < len(links); i++ {
		findAllPath(links, start, end, i)
	}
}


func findAllPath(links []string, start,end string, index int)  {

	first:= strings.Split(links[index], "-")[0]
	second := strings.Split(links[index], "-")[1]


	if first == start{
		fmt.Print(first+"-")

		
	if second == end{
		fmt.Println(end)
	}

		index++
		if index>=len(links){
			return
		}
		findAllPath(links, second, end, index)
	}




}

