package election_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestElection(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "election Suite")
}
