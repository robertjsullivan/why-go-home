package election

import (
	"net/http"
	"github.com/gorilla/websocket"
	"math/rand"
	"time"
	"fmt"
	"net/url"
	"encoding/json"
	"strconv"
	"io/ioutil"
	"sync"
	"math"
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
	return &Node{Id: id, WebsocketAddress: address, term: 0, Voted: false, Members: members, timeout: GenerateRandomTimeout(), Heartbeat: false, votes: 0}
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

func (n *Node) Start(){


	go n.setupHttpEndpoint()

	for true{
		n.Heartbeat = false
		n.votes = 0
		time.Sleep(time.Duration(n.timeout) * time.Millisecond)

		if !n.Voted && !n.Heartbeat {
			leader := n.initiateVoting()
			if leader {
				n.heartbeater()
			}
		}
	}
}

func GenerateRandomTimeout() int {
	return rand.Intn(200) + 150
	//return rand.Intn(20)+1
}

func (n *Node) initiateVoting() bool {
	fmt.Printf("initiating voting for id: %d\n",n.Id)
	n.term++;
	for _, member := range n.Members {
		uri := fmt.Sprintf("http://127.0.0.1%s/%d", member.WebsocketAddress, member.Id)
		resp, err := http.PostForm(uri,
			url.Values{"startElection": {"true"}, "id": {strconv.Itoa(n.Id)}, "term": {strconv.Itoa(n.term)}})
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
		n.Lock()
		if vote.Vote_id == n.Id {
			fmt.Printf("id: %d received vote from %d term: %d\n",n.Id,vote.Sender_id,vote.Term)
			n.votes++
		}else{
			fmt.Printf("id: %d denied vote from %d term: %d\n",n.Id,vote.Sender_id,vote.Term)
		}
		n.Unlock()



	}

	if float64(n.votes) > math.Ceil(float64(len(n.Members) / 2)){
		fmt.Printf("id: %d has won the election.  time to heartbeat.\n", n.Id)
		return true;
	}

	return false;
}

func (n *Node) heartbeater(){
	fmt.Printf("heart beating\n");
	for true {
		for _, member := range n.Members {
			uri := fmt.Sprintf("http://127.0.0.1%s/%d", member.WebsocketAddress, member.Id)
			//http.PostForm(url,)
			_, err := http.PostForm(uri,
				url.Values{"heartbeat": {"true"}})

			//c, _, err := websocket.DefaultDialer.Dial(uri, nil)
			if err != nil {
			 panic(err)
			}
			//defer c.Close()
			//c.WriteJSON("{bob:true}")

		}

		time.Sleep(50 * time.Millisecond)

	}
}

func (n *Node) setupHttpEndpoint(){


	uri := fmt.Sprintf("/%d", n.Id)
	http.HandleFunc(uri, n.VoteHandler)
	fmt.Printf("created endpoint for uri"+n.WebsocketAddress +uri+"\n")
	err := http.ListenAndServe(n.WebsocketAddress, nil)

	if err != nil {
		panic("ListenAndServe: " + err.Error())
	}
}