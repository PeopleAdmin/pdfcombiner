package gomon

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestGomon(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Gomon Suite")
}
