package main

import (
	"context"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/urfave/cli/v2"
)

type scanPayload struct {
	Filename string `json:"filename"`
	Document string `json:"document"`
}

type scanResponse struct {
	PolicyBreakCount int      `json:"policy_break_count"`
	Policies         []string `json:"policies"`
	PolicyBreaks     []struct {
		Type    string `json:"type"`
		Policy  string `json:"policy"`
		Matches []struct {
			Type  string `json:"type"`
			Match string `json:"match"`
		} `json:"matches"`
		Validity string `json:"validity,omitempty"`
	} `json:"policy_breaks"`
}

func displayResults(scan scanResponse) {
	if scan.PolicyBreakCount == 0 {
		return
	}
	for _, policyBreak := range scan.PolicyBreaks {
		fmt.Printf("Policy: %s\n", policyBreak.Policy)
		for _, match := range policyBreak.Matches {
			fmt.Printf("Match: %s\n", match.Match)
		}
	}
}

// buildPayload walks the path and builds the payload to send to the API.
// If dryRun is true, it will not read the file content (when we just want a files count)
func buildPayload(path string, dryRun bool) []scanPayload {
	var payload []scanPayload
	root := path
	err := filepath.Walk(root, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			if dryRun {
				payload = append(payload, scanPayload{Filename: filepath.Base(path), Document: ""})
			} else {
				fileContent, err := os.ReadFile(path)
				if err != nil {
					fmt.Println(err)
					return nil
				}
				payload = append(payload, scanPayload{Filename: filepath.Base(path), Document: string(fileContent)})
			}
		}
		return nil
	})
	if err != nil {
		panic(err)
	}
	return payload

}

func (c *GitGuardianClient) PostScan(ctx context.Context, payload *scanPayload) (resp *scanResponse, err error) {
	err = c.Post("/scan").
		SetBody(payload).
		Do(ctx).
		Into(&resp)
	return
}

func scanCMD(Ctx *cli.Context) error {

	selectedPath := Ctx.Args().Get(0)

	// if path is empty launch filepicker
	if selectedPath == "" {
		fmt.Println("Selectedpath is empty")
		fmt.Println("Trigger Filepicker here!")
		return nil
	}

	// We use /multiscan for 1 file or a dir w/ multiple files for sake of simplicity
	payload := buildPayload(selectedPath, false)

	fmt.Println(len(payload))

	for _, file := range payload {
		resp, err := gg.PostScan(context.Background(), &file)
		if err != nil {
			fmt.Println(err)
			continue
		}
		displayResults(*resp)
	}

	return nil
}
