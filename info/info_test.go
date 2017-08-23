package info_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/ghttp"
	. "github.com/jhamon/guac/info"
)

var _ = Describe("Info", func() {
	var (
		server *ghttp.Server
		context UaaContext
		infoResponse string
	)

	BeforeEach(func() {
		server = ghttp.NewServer()
		context = UaaContext{}
		context.BaseUrl = server.URL()
		infoResponse = `{
		  "app": {
			"version": "4.5.0"
		  },
		  "links": {
			"uaa": "https://uaa.run.pivotal.io",
			"passwd": "https://account.run.pivotal.io/forgot-password",
			"login": "https://login.run.pivotal.io",
			"register": "https://account.run.pivotal.io/sign-up"
		  },
		  "zone_name": "uaa",
		  "entityID": "login.run.pivotal.io",
		  "commit_id": "df80f63",
		  "idpDefinitions": {},
		  "prompts": {
			"username": [
			  "text",
			  "Email"
			],
			"password": [
			  "password",
			  "Password"
			]
		  },
		  "timestamp": "2017-07-21T22:45:01+0000"
		}`
	})

	AfterEach(func() {
		server.Close()
	})

	It("calls the /info endpoint", func() {
		server.RouteToHandler("GET", "/info", ghttp.CombineHandlers(
			ghttp.RespondWith(200, infoResponse),
			ghttp.VerifyRequest("GET", "/info"),
			ghttp.VerifyHeaderKV("Accept", "application/json"),
		))

		infoResponse, _ := Info(context)

		Expect(server.ReceivedRequests()).To(HaveLen(1))
		Expect(infoResponse.App.Version).To(Equal("4.5.0"))
		Expect(infoResponse.Links.ForgotPassword).To(Equal("https://account.run.pivotal.io/forgot-password"))
		Expect(infoResponse.Links.Uaa).To(Equal("https://uaa.run.pivotal.io"))
		Expect(infoResponse.Links.Registration).To(Equal("https://account.run.pivotal.io/sign-up"))
		Expect(infoResponse.Links.Login).To(Equal("https://login.run.pivotal.io"))
		Expect(infoResponse.ZoneName).To(Equal("uaa"))
		Expect(infoResponse.EntityId).To(Equal("login.run.pivotal.io"))
		Expect(infoResponse.CommitId).To(Equal("df80f63"))
		Expect(infoResponse.Prompts.Username).To(Equal([]string{"text", "Email"}))
		Expect(infoResponse.Prompts.Password).To(Equal([]string{"password", "Password"}))
		Expect(infoResponse.Timestamp).To(Equal("2017-07-21T22:45:01+0000"))
	})

	It("returns helpful error when /info request fails", func() {
		server.RouteToHandler("GET", "/info", ghttp.CombineHandlers(
			ghttp.RespondWith(500, ""),
			ghttp.VerifyRequest("GET", "/info"),
			ghttp.VerifyHeaderKV("Accept", "application/json"),
		))

		_, err := Info(context)

		Expect(err).NotTo(BeNil())
		Expect(err.Error()).To(ContainSubstring("An unknown error occurred while calling"))
	})

	It("returns helpful error when /info response can't be parsed", func() {
		server.RouteToHandler("GET", "/info", ghttp.CombineHandlers(
			ghttp.RespondWith(200, "{unparsable-json-response}"),
			ghttp.VerifyRequest("GET", "/info"),
			ghttp.VerifyHeaderKV("Accept", "application/json"),
		))

		_, err := Info(context)

		Expect(err).NotTo(BeNil())
		Expect(err.Error()).To(ContainSubstring("An unknown error occurred while parsing response from"))
		Expect(err.Error()).To(ContainSubstring("Response was {unparsable-json-response}"))
	})
})
