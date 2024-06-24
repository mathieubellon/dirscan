package main

import (
	"context"
	"fmt"

	"github.com/schollz/progressbar/v3"
	"github.com/urfave/cli/v2"
)

type quotasResponse struct {
	Content struct {
		Count     int    `json:"count"`
		Limit     int    `json:"limit"`
		Remaining int    `json:"remaining"`
		Since     string `json:"since"`
	} `json:"content"`
}

func (c *GitGuardianClient) GetQuotas(ctx context.Context) (resp *quotasResponse, err error) {
	err = c.Get("/quotas").
		Do(ctx).
		Into(&resp)
	return
}

func quotasCMD(Ctx *cli.Context) error {
	quotas, err := gg.GetQuotas(nil)
	if err != nil {
		fmt.Println(err)
	}
	bar := progressbar.NewOptions(quotas.Content.Limit,
		progressbar.OptionEnableColorCodes(true),
		progressbar.OptionShowBytes(false),
		progressbar.OptionShowCount(),
		progressbar.OptionSetDescription("API usage since "+quotas.Content.Since),
		progressbar.OptionClearOnFinish(),
		progressbar.OptionSetRenderBlankState(true),
		progressbar.OptionSetElapsedTime(false),
		progressbar.OptionSetPredictTime(false),
		progressbar.OptionSetTheme(progressbar.Theme{
			Saucer:        "=",
			SaucerHead:    ">",
			SaucerPadding: " ",
			BarStart:      "[",
			BarEnd:        "]",
		}))
	bar.Add(quotas.Content.Count)
	fmt.Println("")

	return nil
}
