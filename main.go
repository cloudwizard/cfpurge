package cfpurge

import (
	"context"

	"github.com/cloudflare/cloudflare-go"
)

type cfpurge struct {
	ApiToken     string // Under MyProfile/API Tokens
	ZoneId       string // click domain that you want, on right side, Zone Id
	EmailAddress string // Email address for your cloudflare account
	api          *cloudflare.API
}

func (c cfpurge) New(apitoken string, zoneid string, eaddr string) (cfpurge, error) {
	api, err := cloudflare.New(apitoken, eaddr)
	if err != nil {
		return cfpurge{}, err
	}
	newcf := cfpurge{apitoken, zoneid, eaddr, api}
	return newcf, nil
}

func (c cfpurge) PurgeFile(fn string) (cloudflare.PurgeCacheResponse, error) {
	// fn is "https://example.com/alpha"
	// https://pkg.go.dev/github.com/pcaminog/cloudflare-go#API.PurgeCache
	f := []string{fn}

	return c.PurgeFiles(f)
}

func (c cfpurge) PurgeFiles(files []string) (cloudflare.PurgeCacheResponse, error) {
	// https://pkg.go.dev/github.com/pcaminog/cloudflare-go#API.PurgeCache
	// up to 30 files at once
	// 30K calls/24 hours
	ctx := context.Background()

	pcr := cloudflare.PurgeCacheRequest{Files: files}

	return c.api.PurgeCache(ctx, "9513f845fdf6e463566ef637d9ae7032", pcr)
}

func (c cfpurge) PurgeAll() (cloudflare.PurgeCacheResponse, error) {
	ctx := context.Background()
	return c.api.PurgeEverything(ctx, c.ZoneId)

}
