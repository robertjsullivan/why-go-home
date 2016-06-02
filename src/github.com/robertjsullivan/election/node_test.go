package election_test

import (
	. "github.com/robertjsullivan/election"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"net/http"
	"net/http/httptest"
)

var _ = Describe("Election", func() {

	It("sets the correct values", func(){
		var members []Member
		member := Member{1, "some-address",}
		members = append(members, member)
		member = Member{2, "some-address",}
		members = append(members, member)

		n := NewNode(3, "some-address", members)
		Expect(n.Id).To(Equal(3))
		Expect(n.WebsocketAddress).To(Equal("some-address"))
		Expect(len(n.Members)).To(Equal(len(members)))
	})

	It("handles startElection requests", func(){
		request, _ := http.NewRequest("GET", "http://example.com", nil)
		request.ParseForm()
		request.Form.Add("startElection", "true")
		request.Form.Add("id", "1")
		request.Form.Add("term", "1")
		recorder := httptest.NewRecorder()
		n := &Node{}
		n.VoteHandler(recorder, request)
		Expect(string(recorder.Body.Bytes())).To(Equal(`{"Sender_id":0,"Vote_id":1,"Term":1}`))
	})

	It("votes no if it has already voted in that term", func(){
		request, _ := http.NewRequest("GET", "http://example.com", nil)
		request.ParseForm()
		request.Form.Add("startElection", "true")
		request.Form.Add("id", "1")
		request.Form.Add("term", "0")
		recorder := httptest.NewRecorder()
		n := &Node{}
		n.Voted = true
		n.Id = 3
		n.VoteHandler(recorder, request)
		Expect(string(recorder.Body.Bytes())).To(Equal(`{"Sender_id":3,"Vote_id":-1,"Term":0}`))
	})

	It("handles heartbeat requests", func(){
		request, _ := http.NewRequest("GET", "http://example.com", nil)
		request.ParseForm()
		request.Form.Add("heartbeat", "true")
		recorder := httptest.NewRecorder()
		n := &Node{}
		n.VoteHandler(recorder, request)
		Expect(string(recorder.Body.Bytes())).To(Equal("ok"))
		Expect(n.Heartbeat).To(BeTrue())
	})


	It("generates a random timeout between 150 and 350 ms", func(){
		Expect(GenerateRandomTimeout()).To(BeNumerically(">=", 150))
		Expect(GenerateRandomTimeout()).To(BeNumerically("<=", 350))
	})

	It("initiates a vote when the timeout expires", func(){

	})

	It("does not initiate if already voted", func(){

	})

	It("does not initiate if it receives a heartbeat", func(){

	})

	It("votes if it hasn't already voted this term", func(){

	})

	It("it records a heartbeat", func(){

	})

	It("initiates a vote if it doesn't receive a heartbeat", func(){

	})


})
