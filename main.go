package main

import (
	"log"
	"os"

	"github.com/urfave/cli/v2"
)

const (
	MAX_FILES_BY_SCAN = 20
)

var gg = NewGitGuardianClient()

func init() {
	gg.SetDebug(false)
	if token := os.Getenv("GITGUARDIAN_API_KEY"); token != "" {
		gg.LoginWithToken(token)
		return
	}
}
func main() {
	app := &cli.App{
		Name:    "GGscan",
		Usage:   "Scan your files for secrets with GitGuardian API",
		Version: "0.1.0",
		Commands: []*cli.Command{
			{
				Name:    "health",
				Aliases: []string{"he"},
				Usage:   "Ping /health endpoint and check everything's ready",
				Action:  healthCMD,
			},
			{
				Name:    "quotas",
				Aliases: []string{"q"},
				Usage:   "Get your current quotas status",
				Action:  quotasCMD,
			},
			{
				Name:    "scan",
				Aliases: []string{"s"},
				Usage:   "Scan directories for secrets",
				Action:  scanCMD,
			},
		},
	}
	app.Suggest = true
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
