package main_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/scoiatael/archai/actions"
	"github.com/scoiatael/archai/simplejson"
	"github.com/scoiatael/archai/util"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

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

				write := <-config.BackgroundJobs()
				write.Run(config)

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
				buf = bytes.NewBufferString(`{ "foo": "bar", "baz": 2 }`)
			})

			It("allows writing events", func() {
				resp, err := http.Post(address, "application/json", buf)

				Expect(err).NotTo(HaveOccurred())
				body, err := ioutil.ReadAll(resp.Body)
				Expect(err).NotTo(HaveOccurred())
				Expect(string(body)).To(Equal(`"OK"`))

				write := <-config.BackgroundJobs()
				write.Run(config)

				action := actions.ReadEvents{}
				action.Amount = 5
				action.Stream = stream
				err = action.Run(config)
				Expect(err).NotTo(HaveOccurred())
				Expect(action.Events).NotTo(BeEmpty())
				Expect(action.Events).To(HaveLen(1))

				js, err := simplejson.Read(action.Events[0].Blob)
				Expect(err).NotTo(HaveOccurred())
				Expect(js["foo"]).To(Equal("bar"))
				Expect(js["baz"]).To(Equal(2.0))
			})

			Context("After some event was written", func() {
				JustBeforeEach(func() {
					event := actions.WriteEvent{}
					event.Stream = stream
					event.Payload = []byte(`{ "foo": "bar"}`)
					event.Meta = make(map[string]string)

					event.Run(config)

				})

				get := func(query string) map[string]interface{} {
					resp, err := http.Get(address + query)

					Expect(err).NotTo(HaveOccurred())
					body, err := ioutil.ReadAll(resp.Body)
					Expect(err).NotTo(HaveOccurred())
					js := make(map[string]interface{})
					err = json.Unmarshal(body, &js)
					Expect(err).NotTo(HaveOccurred())
					return js
				}

				It("allows reading events", func() {
					results := get("")["results"]
					Expect(results).NotTo(BeEmpty())
					Expect(results).To(HaveLen(1))
				})
				It("allows reading events with cursor", func() {
					cursor := get("")["cursor"].(map[string]interface{})["next"].(string)

					results := get("?cursor=" + cursor)["results"]
					Expect(results).To(BeEmpty())
				})
			})
		})
	})
})
