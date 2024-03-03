package cfpurge

import (
	"context"
	"os"

	"github.com/cloudflare/cloudflare-go"
	"gopkg.in/yaml.v3"
)

// YAML file in secrets folder
//	cf = new(cfpurge)
//

type Cfpurge struct {
	ApiToken     string          `yaml:"apitoken"` // Under MyProfile/API Tokens
	ZoneId       string          `yaml:"zoneid"`   // click domain that you want, on right side, Zone Id
	EmailAddress string          `yaml:"email"`    // Email address for your cloudflare account
	Api          *cloudflare.API `yaml:"-"`
}

func (c *Cfpurge) Init(apitoken string, zoneid string, eaddr string) error {
	api, err := cloudflare.New(apitoken, eaddr)
	if err != nil {
		return err
	}
	c.Api = api
	c.ApiToken = apitoken
	c.ZoneId = zoneid
	c.EmailAddress = eaddr
	return nil
}

func (c *Cfpurge) InitFromYaml(filename string) error {
	// https://stackoverflow.com/questions/30947534/how-to-read-a-yaml-file
	buf, err := os.ReadFile(filename)
	if err != nil {
		return err
	}

	err = yaml.Unmarshal(buf, c)
	if err != nil {
		return err
	}
	c.Api, err = cloudflare.New(c.ApiToken, c.EmailAddress)
	return err
}

func (c *Cfpurge) PurgeFile(fn string) (cloudflare.PurgeCacheResponse, error) {
	// fn is "https://example.com/alpha"
	// https://pkg.go.dev/github.com/pcaminog/cloudflare-go#API.PurgeCache
	f := []string{fn}

	return c.PurgeFiles(f)
}

func (c *Cfpurge) PurgeFiles(files []string) (cloudflare.PurgeCacheResponse, error) {
	// https://pkg.go.dev/github.com/pcaminog/cloudflare-go#API.PurgeCache
	// up to 30 files at once
	// 30K calls/24 hours
	ctx := context.Background()

	pcr := cloudflare.PurgeCacheRequest{Files: files}

	return c.Api.PurgeCache(ctx, c.ZoneId, pcr)
}

func (c *Cfpurge) PurgeAll() (cloudflare.PurgeCacheResponse, error) {
	ctx := context.Background()
	return c.Api.PurgeEverything(ctx, c.ZoneId)

}
