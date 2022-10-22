package main

import (
	"encoding/base64"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/miekg/dns"
	"github.com/patrickmn/go-cache"
	"log"
	"net/url"
	"os"
	"time"
)

//main function
func main() {
	//new Fiber Instance
	app := fiber.New(fiber.Config{
		Prefork:      os.Getenv("PREFORK") == "true",
		ServerHeader: "HttpDNS",
	})

	//new DNS Client Instance (default 2s timeout)
	dnsClient := new(dns.Client)

	//new Cache Instance (default 5m ttl)
	dnsCache := cache.New(5*time.Minute, 10*time.Minute)

	//get DNS Server from env
	dnsServer := os.Getenv("DNS_SERVER")
	if dnsServer == "" {
		dnsServer = "4.2.2.4:53"
	}

	//app default route
	app.Get("/:class/:type/:domain", func(c *fiber.Ctx) error {
		responseType := c.Query("type", "http")
		domain := c.Params("domain")
		domainBytes, err := base64.StdEncoding.DecodeString(domain)
		if err != nil {
			return sendError(c, fiber.StatusBadRequest, err, responseType)
		}
		domain = string(domainBytes)
		dnsType := c.Params("type")
		dnsClass := c.Params("class")
		_, err = url.Parse(domain)
		if err != nil {
			return sendError(c, fiber.StatusBadRequest, err, responseType)
		}
		cacheKey := fmt.Sprintf("%s/%s/%s/%s", domain, dnsType, dnsClass, responseType)
		value, exists := dnsCache.Get(cacheKey)
		if exists {
			msg := value.(*dns.Msg)
			return sendData(c, msg, responseType)
		}
		m := &dns.Msg{
			MsgHdr: dns.MsgHdr{
				RecursionDesired: true,
			},
			Question: make([]dns.Question, 1),
		}
		m.Question[0] = dns.Question{
			Name:   domain,
			Qtype:  dns.StringToType[dnsType],
			Qclass: dns.StringToClass[dnsClass],
		}

		resp, _, err := dnsClient.Exchange(m, dnsServer)
		if err != nil {
			return sendError(c, fiber.StatusInternalServerError, err, responseType)
		}
		if len(resp.Answer) == 0 {
			return sendError(c, fiber.StatusNotFound, err, responseType)
		}
		dnsCache.Set(cacheKey, resp, cache.DefaultExpiration)
		return sendData(c, resp, responseType)
	})

	log.Fatal(app.Listen(":3000"))
}

// send error to client based on responseType
func sendError(c *fiber.Ctx, request int, err error, responseType string) error {
	if responseType == "dns" {
		return c.SendStatus(request)
	} else {
		return c.Status(request).JSON(fiber.Map{
			"status":  "error",
			"message": err.Error(),
		})
	}

}

// send data to client based on responseType
func sendData(ctx *fiber.Ctx, resp *dns.Msg, responseType string) error {
	if responseType == "dns" {
		wire, err := resp.Pack()
		if err != nil {
			return sendError(ctx, fiber.StatusInternalServerError, err, responseType)
		}
		return ctx.Send(wire)
	} else {
		return ctx.JSON(fiber.Map{
			"status": "success",
			"data":   parseHttp(resp),
			"ttl":    resp.Answer[0].Header().Ttl,
		})
	}
}

//translate dns response to json format
func parseHttp(resp *dns.Msg) interface{} {
	var values []any
	for _, a := range resp.Answer {
		switch a.(type) {
		case *dns.A:
			values = append(values, a.(*dns.A).A.String())
		case *dns.AAAA:
			values = append(values, a.(*dns.AAAA).AAAA.String())
		case *dns.CNAME:
			values = append(values, a.(*dns.CNAME).Target)
		case *dns.MX:
			values = append(values, a.(*dns.MX).Mx)
		case *dns.NS:
			values = append(values, a.(*dns.NS).Ns)
		case *dns.PTR:
			values = append(values, a.(*dns.PTR).Ptr)
		case *dns.SOA:
			values = append(values, a.(*dns.SOA))
		case *dns.SRV:
			values = append(values, a.(*dns.SRV).Target)
		case *dns.TXT:
			values = append(values, a.(*dns.TXT).Txt[0])
		}
	}
	return values
}
