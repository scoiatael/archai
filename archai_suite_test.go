package main_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"

	. "github.com/scoiatael/archai"
)

const testingKeyspace = "archai_test"

var (
	config Config
)

var _ = BeforeSuite(func() {
	config = Config{}
	config.Features = make(map[string]bool)
	config.Hosts = []string{"127.0.0.1"}
	// NOTE: makes it far easier to spot panics, but throws a log of noise otherwise
	// config.Features["dev_logger"] = true
	config.Keyspace = testingKeyspace
	config.StatsdAddr = "dd-agent.service.consul:8125"
	err := config.Init()
	if err != nil {
		panic(err)
	}
})

var _ = AfterSuite(func() {

})

func TestArchai(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Archai Suite")
}
