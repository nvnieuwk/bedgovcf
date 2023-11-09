package bedgovcf

import (
	"fmt"
	"log"
	"os"
	"strings"

	"gopkg.in/yaml.v3"
)

func ReadConfig(configString string) Config {
	configFile, err := os.ReadFile(configString)
	if err != nil {
		log.Fatalf("Failed to open the config file: %v", err)
	}

	var config Config

	if err := yaml.Unmarshal(configFile, &config); err != nil {
		log.Fatalf("Failed to parse the config file: %v", err)
	}

	config.validate()
	return config
}

func (c *Config) validate() {
	if c.Chrom.Field == "" {
		log.Printf("No field defined for CHROM, defaulting to the column 0")
		c.Chrom.Field = "0"
	}

	if c.Pos.Field == "" {
		log.Printf("No field defined for POS, defaulting to the column 1")
		c.Pos.Field = "1"
	}

	if c.Id.Field == "" && c.Id.Prefix == "" {
		log.Printf("No field or prefix specified for the ID, defaulting to prefix 'id_")
		c.Id.Prefix = "id_"
	}

	if c.Ref.Field == "" && c.Ref.Value == "" {
		log.Printf("No field or value specified for the REF, defaulting to value 'N")
		c.Ref.Value = "N"
	}

	if c.Alt.Field == "" && c.Alt.Value == "" {
		log.Printf("No field or value specified for the ALT, defaulting to value '<CNV>")
		c.Alt.Value = "<CNV>"
	}

	if c.Qual.Field == "" && c.Qual.Value == "" {
		log.Printf("No field or value specified for the QUAL, defaulting to value '.'")
		c.Qual.Value = "."
	}

	if c.Filter.Field == "" && c.Filter.Value == "" {
		log.Printf("No field or value specified for the FILTER, defaulting to value 'PASS'")
		c.Filter.Value = "PASS"
	}

	errStrings := []string{}

	if len(c.Info) != 0 {
		for k, v := range c.Info {
			if v.Field == "" && v.Value == "" {
				s := fmt.Sprintf("No field or value specified for the INFO/%v", strings.ToUpper(k))
				errStrings = append(errStrings, s)
			}
		}
	}

	if len(c.Format) != 0 {
		for k, v := range c.Format {
			if v.Field == "" && v.Value == "" {
				s := fmt.Sprintf("No field or value specified for the FORMAT/%v", strings.ToUpper(k))
				errStrings = append(errStrings, s)
			}
		}
	}

	if len(errStrings) != 0 {
		log.Fatalf("The following errors were found in the config file:\n%v", strings.Join(errStrings, "\n"))
	}
}
