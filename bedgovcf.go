package main

import (
	"log"
	"os"

	bedgovcf "github.com/nvnieuwk/bedgovcf/convert"
	cli "github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:            "bedgovcf",
		Usage:           "Convert a BED file to a VCF file according to a YAML config",
		HideHelpCommand: true,
		Version:         "0.1.1",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "output",
				Aliases:  []string{"o"},
				Usage:    "The location to the output VCF file, defaults to stdout",
				Category: "Optional",
			},
			&cli.StringFlag{
				Name:     "sample",
				Aliases:  []string{"s"},
				Usage:    "The name of the sample to use in the VCF file, defaults to the basename of the BED file",
				Category: "Optional",
			},
			&cli.Int64Flag{
				Name:     "skip",
				Aliases:  []string{"k"},
				Usage:    "The amount of lines to skip in the BED file",
				Category: "Optional",
			},
			&cli.BoolFlag{
				Name:     "header",
				Aliases:  []string{"l"},
				Usage:    "The BED file contains a header line",
				Category: "Optional",
			},
			&cli.StringFlag{
				Name:     "config",
				Aliases:  []string{"c"},
				Usage:    "Configuration file to use for the conversion in YAML format",
				Required: true,
				Category: "Required",
			},
			&cli.StringFlag{
				Name:     "bed",
				Aliases:  []string{"b"},
				Usage:    "The input BED file",
				Required: true,
				Category: "Required",
			},
			&cli.StringFlag{
				Name:     "fai",
				Aliases:  []string{"f"},
				Usage:    "The location to the fasta index file",
				Required: true,
				Category: "Required",
			},
		},
		Action: func(c *cli.Context) error {
			logger := log.New(os.Stderr, "", 0)
			config, err := bedgovcf.ReadConfig(c.String("config"))
			if err != nil {
				logger.Fatal(err)
			}
			vcf := bedgovcf.Vcf{}
			err = vcf.SetHeader(c, config)
			if err != nil {
				logger.Fatal(err)
			}
			err = vcf.AddVariants(c, config)
			if err != nil {
				logger.Fatal(err)
			}
			err = vcf.Write(c)
			if err != nil {
				logger.Fatal(err)
			}

			return nil
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
