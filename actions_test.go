package main_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"time"

	. "github.com/scoiatael/archai"
	"github.com/scoiatael/archai/actions"
	"github.com/scoiatael/archai/simplejson"
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
	config.StatsdAddr = "dd-agent.service.consul:8125"
	err := config.Init()
	if err != nil {
		panic(err)
	}
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
			time.Sleep(10 * time.Millisecond)
			address = fmt.Sprintf("http://127.0.0.1:%d", action.Port)
		})

		AfterEach(func() {
			action.Stop()
		})

		Describe("/bulk/stream/:id", func() {
			JustBeforeEach(func() {
				stream = util.RandomString(10)
				address = fmt.Sprintf("%s/bulk/stream/%s", address, stream)
				buf = bytes.NewBufferString(`{ "data": [["foo",1,2,3], ["bar",4,5,6]], "schema": ["name", "likes", "shares", "comments"] }`)
			})

			It("allows writing events", func() {
				resp, err := http.Post(address, "application/json", buf)

				Expect(err).NotTo(HaveOccurred())
				body, err := ioutil.ReadAll(resp.Body)
				Expect(err).NotTo(HaveOccurred())
				Expect(string(body)).To(Equal(`"OK"`))

				time.Sleep(20 * time.Millisecond)

				action := actions.ReadEvents{}
				action.Amount = 5
				action.Stream = stream
				err = action.Run(config)
				Expect(err).NotTo(HaveOccurred())
				Expect(action.Events).NotTo(BeEmpty())
				Expect(action.Events).To(HaveLen(2))

				js, err := simplejson.Read(action.Events[0].Blob)
				Expect(err).NotTo(HaveOccurred())
				Expect(js["name"]).To(Equal("foo"))
				Expect(js["likes"]).To(Equal(1.0))
				Expect(js["shares"]).To(Equal(2.0))
				Expect(js["comments"]).To(Equal(3.0))
			})
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

				get := func(query string) interface{} {
					resp, err := http.Get(address + query)

					Expect(err).NotTo(HaveOccurred())
					body, err := ioutil.ReadAll(resp.Body)
					Expect(err).NotTo(HaveOccurred())
					js := make(map[string]interface{})
					err = json.Unmarshal(body, &js)
					Expect(err).NotTo(HaveOccurred())
					return js["results"]
				}

				It("allows reading events", func() {
					results := get("")
					Expect(results).NotTo(BeEmpty())
					Expect(results).To(HaveLen(1))
				})
				It("allows reading events with cursor", func() {
					cursor := get("").([]interface{})[0].(map[string]interface{})["ID"].(string)

					results := get("?cursor=" + cursor)
					Expect(results).To(BeEmpty())
				})
			})
		})
	})
})
