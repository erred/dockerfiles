package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"sort"
	"strconv"
	"strings"

	"github.com/cloudflare/cloudflare-go"
)

const (
	EnvRecord = "RECORD"

	XAuthEmail = "X_AUTH_EMAIL"
	XAuthKey   = "X_AUTH_KEY"
)

type Config struct {
	cf         *cloudflare.API
	authEmail  string
	authKey    string
	zoneName   string
	zoneID     string
	recordName string
	recordType string
	arg        string
	argN       int
	noProxy    bool
}

func NewConfig() *Config {
	ss := strings.Split(os.Getenv(EnvRecord), ":")
	if len(ss) != 5 {
		log.Fatalf("expected %v to have 5 parts, got %v: %v", EnvRecord, len(ss), ss)
	}
	c := Config{
		authEmail:  os.Getenv(XAuthEmail),
		authKey:    os.Getenv(XAuthKey),
		zoneName:   ss[0],
		recordName: ss[2] + "." + ss[0],
		recordType: ss[3],
		arg:        ss[4],
	}

	c.argN, _ = strconv.Atoi(ss[4])
	switch strings.ToLower(ss[1]) {
	case "noproxy":
		c.noProxy = true
	}
	return &c
}

func (c *Config) Auth() {
	var err error
	c.cf, err = cloudflare.New(c.authKey, c.authEmail)
	if err != nil {
		log.Fatal("auth: ", err)
	}
}

func (c *Config) GetZone() {
	var err error
	c.zoneID, err = c.cf.ZoneIDByName(c.zoneName)
	if err != nil {
		log.Fatal("get zone ", c.zoneName, ": ", err)
	}
}

func (c *Config) GetRecords() []cloudflare.DNSRecord {
	var err error
	rs, err := c.cf.DNSRecords(c.zoneID, cloudflare.DNSRecord{
		Type: c.recordType,
		Name: c.recordName,
	})
	if err != nil {
		log.Fatal("get records: ", err)
	}
	return rs
}

func (c *Config) add(content string) error {
	res, err := c.cf.CreateDNSRecord(c.zoneID, cloudflare.DNSRecord{
		Type:    c.recordType,
		Name:    c.recordName,
		Proxied: !c.noProxy,
		Content: content,
	})
	if err != nil {
		log.Printf("add %s: %+v", c.recordType, res)
	}
	return err
}
func (c *Config) deleteRecord(rid string) error {
	return c.cf.DeleteDNSRecord(c.zoneID, rid)
}
func (c *Config) update(r cloudflare.DNSRecord, content string) error {
	return c.cf.UpdateDNSRecord(c.zoneID, r.ID, cloudflare.DNSRecord{
		Type:    c.recordType,
		Name:    c.recordName,
		Content: content,
		Proxied: !c.noProxy,
	})
}
func getIP() (string, error) {
	res, err := http.Get("https://api.ipify.org?format=json")
	if err != nil {
		return "", fmt.Errorf("get ip from ext: %v", err)
	}
	if res.StatusCode < 200 || res.StatusCode >= 300 {
		return "", fmt.Errorf("res status %v %v\n", res.StatusCode, res.Status)
	}
	defer res.Body.Close()

	type IP struct {
		IP string `json:"ip"`
	}
	var ip IP
	err = json.NewDecoder(res.Body).Decode(&ip)
	if err != nil {
		return "", fmt.Errorf("decode ip: %v", err)
	}
	return ip.IP, nil

}

// usage: cf-dns-update kluster seankhliao.com A 1
// usage: cf-dns-update badger seankhliao.com CNAME kluster.seankhliao.com
func main() {
	c := NewConfig()

	c.Auth()
	c.GetZone()
	rs := c.GetRecords()
	sort.Sort(Records(rs))

	switch c.recordType {
	case "A":
		if c.argN < 1 {
			log.Fatal("argN < 1: ", c.argN)
		}
		ip, err := getIP()
		if err != nil {
			log.Fatal("get IP from API: ", err)
		}
		for _, r := range rs {
			if r.Content == ip {
				fmt.Println("record already exists")
				os.Exit(0)
			}
		}
		if len(rs) < c.argN {
			err := c.add(ip)
			if err != nil {
				log.Fatal("add A record: ", err)
			}
			fmt.Println("A record added")
			os.Exit(0)
		}
		if len(rs) > c.argN {
			for i := 0; i < len(rs)-c.argN; i++ {
				err := c.deleteRecord(rs[0].ID)
				if err != nil {
					log.Println("delete A record: ", err)
				}
				rs = rs[1:]
			}
		}
		err = c.update(rs[0], ip)
		if err != nil {
			log.Fatal("update A record: ", err)
		}
		fmt.Println("A record updated")
		os.Exit(0)

	case "CNAME":
		if len(rs) == 0 {
			err := c.add(c.arg)
			if err != nil {
				log.Fatal("add CNAME: ", err)
			}
			fmt.Println("CNAME record added")
			os.Exit(0)
		}
		r := rs[0]
		if r.Content == c.arg {
			fmt.Println("record already exists")
			os.Exit(0)
		}
		err := c.update(r, c.arg)
		if err != nil {
			log.Fatal("update CNAME record: ", err)
		}
		fmt.Println("CNAME record updated")
		os.Exit(0)

	default:
		log.Fatal("unimplemented record type: ", c.recordType)
	}

}

type Records []cloudflare.DNSRecord

func (r Records) Len() int           { return len(r) }
func (r Records) Less(i, j int) bool { return r[i].ModifiedOn.Before(r[j].ModifiedOn) }
func (r Records) Swap(i, j int)      { r[i], r[j] = r[j], r[i] }
