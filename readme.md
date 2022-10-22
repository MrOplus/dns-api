# DNS API

## Currently Supported Types: 
- A
- AAAA
- CNAME
- MX
- NS
- PTR
- SRV
- TXT
- SOA

## Response Types 
- JSON
- WIRE
## Envronment Variables
- `DNS_SERVER` - The DNS server to use. Defaults to 4.2.2.4:53`
- `PREFORK` - Prefork Enabled. Defaults to false
## Usage
```
localhost:3000/class/type/base64encodedUrl?type=[dns,http]
```
Example : 
```
localhost:3000/IN/A/Z29vZ2xlLmNvbS4=?type=dns
```

Powered By : 
- [miekg DNS](github.com/miekg/dns)
- [Fiber](github.com/gofiber/fiber)

