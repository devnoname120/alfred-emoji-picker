package scoring_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestScoring(t *testing.T) {
	RegisterFailHandler(Fail)
	_, reporterConfig := GinkgoConfiguration()
	reporterConfig.FullTrace = true
	reporterConfig.Verbose = true
	RunSpecs(t, "Scoring Suite", reporterConfig)
}
