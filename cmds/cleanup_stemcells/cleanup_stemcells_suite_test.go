package cleanup_stemcells_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestCleanupStemcells(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "CleanupStemcells Suite")
}
