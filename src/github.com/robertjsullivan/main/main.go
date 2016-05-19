package main

import (
	"fmt"
	"github.com/robertjsullivan/election"
)


func main() {
	member_count := 5;
	fmt.Printf("Running Leader Election with %d members.\n", member_count)
	var members []election.Member
	for i := 0; i < member_count; i++ {
		port := i +8080
		address  := fmt.Sprintf(":%d", port)
		member := election.Member{i, address,}
		members = append(members, member)
	}

	for _ , member := range members {
		node := election.NewNode(member.Id, member.WebsocketAddress, members)
		go node.Start()
	}

	for true{

	}
}