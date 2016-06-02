package election

import (
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"net/url"
	"strconv"
	"io/ioutil"
	"time"
)

type NodeService interface {
	InitiateVoting(node *Node) bool
	Heartbeater(node *Node)
	SetupHttpEndpoint(node *Node)
}

type NodeServiceImpl struct {

}

func (ns *NodeServiceImpl) InitiateVoting(node *Node) bool {
	fmt.Printf("initiating voting for id: %d\n",node.Id)
	node.term++;
	for _, member := range node.Members {
		uri := fmt.Sprintf("http://127.0.0.1%s/%d", member.WebsocketAddress, member.Id)
		resp, err := http.PostForm(uri,
			url.Values{"startElection": {"true"}, "id": {strconv.Itoa(node.Id)}, "term": {strconv.Itoa(node.term)}})
		if err != nil {
			panic(err)
		}
		vote := &Vote{}
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			panic(err)
		}
		if err != nil {
			panic(err)
		}
		err = json.Unmarshal(body, vote)
		if err != nil {
			panic(err)
		}
		node.Lock()
		if vote.Vote_id == node.Id {
			fmt.Printf("id: %d received vote from %d term: %d\n",node.Id,vote.Sender_id,vote.Term)
			node.votes++
		}else{
			fmt.Printf("id: %d denied vote from %d term: %d\n",node.Id,vote.Sender_id,vote.Term)
		}
		node.Unlock()



	}

	if float64(node.votes) > math.Ceil(float64(len(node.Members) / 2)){
		fmt.Printf("id: %d has won the election.  time to heartbeat.\n", node.Id)
		return true;
	}

	return false;
}

func (ns *NodeServiceImpl) Heartbeater(node *Node){
	fmt.Printf("heart beating\n");
	for true {
		for _, member := range node.Members {
			uri := fmt.Sprintf("http://127.0.0.1%s/%d", member.WebsocketAddress, member.Id)
			_, err := http.PostForm(uri,
				url.Values{"heartbeat": {"true"}})

			if err != nil {
				panic(err)
			}
		}

		time.Sleep(50 * time.Millisecond)

	}
}

func (ns *NodeServiceImpl) SetupHttpEndpoint(node *Node){
	uri := fmt.Sprintf("/%d", node.Id)
	http.HandleFunc(uri, node.VoteHandler)
	fmt.Printf("created endpoint for uri"+node.WebsocketAddress +uri+"\n")
	err := http.ListenAndServe(node.WebsocketAddress, nil)

	if err != nil {
		panic("ListenAndServe: " + err.Error())
	}
}