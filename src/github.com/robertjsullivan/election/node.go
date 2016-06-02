package election

import (
	"net/http"
	"github.com/gorilla/websocket"
	"math/rand"
	"time"
	"encoding/json"
	"strconv"
	"sync"
)

type Node struct{
	sync.Mutex
	Id               int
	WebsocketAddress string
	term             int
	Voted            bool
	Members          []Member
	timeout          int
	Heartbeat        bool
	votes            int
	Run              bool
}

type Member struct{
	Id int
	WebsocketAddress string
}

type Vote struct{
	Sender_id int
	Vote_id int
	Term int
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func NewNode(id int, address string, members []Member) *Node {
	return &Node{Id: id, WebsocketAddress: address, term: 0, Voted: false, Members: members, timeout: GenerateRandomTimeout(), Heartbeat: false, votes: 0, Run: true}
}


func (n *Node) VoteHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		panic(err)
	}


	switch {
	case r.Form.Get("startElection") == "true":
		id, err :=  strconv.Atoi(r.Form.Get("id"))
		if err != nil {
			panic(err)
		}
		term, err := strconv.Atoi(r.Form.Get("term"))
		if err != nil {
			panic(err)
		}
		myVote := Vote {
			Sender_id: n.Id,
			Vote_id: -1,
			Term: term,
		}
		n.Lock()
		if term > n.term || n.Voted == false {
			n.term = term
			n.Voted = true
			myVote.Vote_id = id
		}
		n.Unlock()

		jsonArr, err := json.Marshal(myVote)
		if err != nil {
			panic(err)
		}
		w.Write(jsonArr)
	case r.Form.Get("heartbeat") == "true":
		n.Heartbeat = true
		w.Write([]byte{'o','k'})
	}
}

func (n *Node) Start(nodeService NodeService){

	go nodeService.SetupHttpEndpoint(n)

	for n.Run{
		n.votes = 0
		time.Sleep(time.Duration(n.timeout) * time.Millisecond)

		if !n.Voted && !n.Heartbeat {
			leader := nodeService.InitiateVoting(n)
			if leader {
				nodeService.Heartbeater(n)
			}
		}
		n.Heartbeat = false
	}
}

func GenerateRandomTimeout() int {
	return rand.Intn(200) + 150
	//return rand.Intn(20)+1
}