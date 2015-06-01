package controllers_test

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"

	"github.com/astaxie/beego"
	"github.com/fsouza/go-dockerclient"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "netstatd"
	. "netstatd/api/server/controllers"
	. "netstatd/namespace"
	. "netstatd/namespace/discovery"
	. "netstatd/test_helper"
)

var _ = Describe("Netstats", func() {
	Describe("ShowAll", func() {
		BeforeEach(func() {
			namespace := NewNamespace(CURRENT_NAMESPACE_PID, "host")
			err := D.AddNameSpace(namespace)
			Expect(err).ShouldNot(HaveOccurred())
		})

		AfterEach(func() {
			D.Stop()
		})

		It("lists all interfaces's netstats", func() {
			r, _ := http.NewRequest("GET", "/v1/netstats", nil)
			w := httptest.NewRecorder()
			beego.BeeApp.Handlers.ServeHTTP(w, r)

			Expect(w.Code).Should(Equal(200))

			var netstats []NetStatJson
			body, _ := ioutil.ReadAll(w.Body)
			err := json.Unmarshal(body, &netstats)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(len(netstats)).Should(BeNumerically(">", 0))
		})
	})

	Describe("Show", func() {
		BeforeEach(func() {
			namespace := NewNamespace(CURRENT_NAMESPACE_PID, "host")
			err := D.AddNameSpace(namespace)
			Expect(err).ShouldNot(HaveOccurred())
		})

		AfterEach(func() {
			D.Stop()
		})

		It("gets a interface's netstat", func() {
			r, _ := http.NewRequest("GET", "/v1/netstats/eth0", nil)
			w := httptest.NewRecorder()
			beego.BeeApp.Handlers.ServeHTTP(w, r)

			Expect(w.Code).Should(Equal(200))

			var netstat NetStatJson
			body, _ := ioutil.ReadAll(w.Body)
			json.Unmarshal(body, &netstat)
			Expect(netstat.Namespace.Tag).Should(Equal("host"))
			Expect(netstat.Iface.Name).Should(Equal("eth0"))

		})

		Context("Docker", func() {
			var container *docker.Container

			BeforeEach(func() {
				container = DockerRunContainer()
				dockerDiscovery := NewDockerDiscovery()
				namespace, err := dockerDiscovery.GetNamespace(container.ID)
				Expect(err).ShouldNot(HaveOccurred())

				err = D.AddNameSpace(namespace)
				Expect(err).ShouldNot(HaveOccurred())

			})

			AfterEach(func() {
				DockerCleanAllContainers()
				D.Stop()
			})

			It("gets a interface's netstat", func() {
				r, _ := http.NewRequest("GET", "/v1/netstats/eth0?docker_container_id="+container.ID, nil)
				w := httptest.NewRecorder()
				beego.BeeApp.Handlers.ServeHTTP(w, r)

				Expect(w.Code).Should(Equal(200))

				var netstat NetStatJson
				body, _ := ioutil.ReadAll(w.Body)
				json.Unmarshal(body, &netstat)
				Expect(netstat.Namespace.Tag).Should(Equal("docker<" + container.ID[0:12] + ">"))
				Expect(netstat.Iface.Name).Should(Equal("eth0"))
			})
		})
	})
})
