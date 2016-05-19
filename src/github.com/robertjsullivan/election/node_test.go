package election_test

import (
	. "github.com/robertjsullivan/election"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Election", func() {

	It("sets the correct values", func(){

	})

	It("generates a random timeout between 150 and 350 ms", func(){

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
