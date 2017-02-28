package main_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"

	. "github.com/scoiatael/archai"
	"github.com/scoiatael/archai/actions"
	"github.com/scoiatael/archai/util"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

const testingKeyspace = "archai_test"

var (
	config Config
)

var _ = BeforeSuite(func() {
	config = Config{}
	config.Hosts = []string{"127.0.0.1"}
	config.Keyspace = testingKeyspace
})

var _ = AfterSuite(func() {

})

var _ = Describe("Actions", func() {
	Describe("HttpServer", func() {
		var (
			action  actions.HttpServer
			port    int
			address string
			stream  string
			buf     io.Reader
		)
		port = 9080
		BeforeEach(func() {
			action.Addr = "127.0.0.1"
			port = port + 1 + util.RandomInt(1000)
			action.Port = port
		})
		JustBeforeEach(func() {
			go action.Run(config)
			address = fmt.Sprintf("http://127.0.0.1:%d", action.Port)
		})

		AfterEach(func() {
			action.Stop()
		})

		Describe("/stream/:id", func() {
			JustBeforeEach(func() {
				stream = util.RandomString(10)
				address = fmt.Sprintf("%s/stream/%s", address, stream)
				buf = bytes.NewBufferString(`{ "foo": "bar" }`)
			})

			It("allows writing events", func() {
				resp, err := http.Post(address, "application/json", buf)

				Expect(err).NotTo(HaveOccurred())
				body, err := ioutil.ReadAll(resp.Body)
				Expect(err).NotTo(HaveOccurred())
				Expect(string(body)).To(Equal(`"OK"`))
			})

			Context("After some event was written", func() {
				JustBeforeEach(func() {
					event := actions.WriteEvent{}
					event.Stream = stream
					event.Payload = []byte(`{ "foo": "bar"}`)
					event.Meta = make(map[string]string)

					event.Run(config)

				})
				It("allows reading events", func() {
					resp, err := http.Get(address)

					Expect(err).NotTo(HaveOccurred())
					body, err := ioutil.ReadAll(resp.Body)
					Expect(err).NotTo(HaveOccurred())
					js := make(map[string]interface{})
					err = json.Unmarshal(body, &js)
					Expect(err).NotTo(HaveOccurred())
					Expect(js["results"]).NotTo(BeEmpty())
				})
			})
		})
	})
})
