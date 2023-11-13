package bedgovcf

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strings"

	"gopkg.in/yaml.v3"
)

// Read the configuration file, cast it to its struct and validate
func ReadConfig(configString string) (error, Config) {
	configFile, err := os.ReadFile(configString)
	if err != nil {
		return errors.New(fmt.Sprintf("Failed to open the config file: %v", err)), Config{}
	}

	var config Config

	if err := yaml.Unmarshal(configFile, &config); err != nil {
		return errors.New(fmt.Sprintf("Failed to open the config file: %v", err)), Config{}
	}

	config.validate()
	return nil, config
}

// Validate the config
func (c *Config) validate() {
	logger := log.New(os.Stderr, "", 0)
	if c.Chrom.Value == "" {
		logger.Printf("No value defined for CHROM, defaulting to the column 0")
		c.Chrom.Value = "$0"
	}

	if c.Pos.Value == "" {
		logger.Printf("No value defined for POS, defaulting to the column 1")
		c.Pos.Value = "$1"
	}

	if c.Id.Value == "" && c.Id.Prefix == "" {
		logger.Printf("No value or prefix specified for the ID, defaulting to prefix 'id_")
		c.Id.Prefix = "id_"
	}

	if c.Ref.Value == "" {
		logger.Printf("No value specified for the REF, defaulting to value 'N")
		c.Ref.Value = "N"
	}

	if c.Alt.Value == "" {
		logger.Printf("No value specified for the ALT, defaulting to value '<CNV>")
		c.Alt.Value = "<CNV>"
	}

	if c.Qual.Value == "" {
		logger.Printf("No value specified for the QUAL, defaulting to value '.'")
		c.Qual.Value = "."
	}

	if c.Filter.Value == "" {
		logger.Printf("No value specified for the FILTER, defaulting to value 'PASS'")
		c.Filter.Value = "PASS"
	}

	if len(c.Info) != 0 {
		for _, v := range c.Info {
			if v.Value == "" {
				logger.Printf("No value specified for the INFO/%v", strings.ToUpper(v.Name))
			}
		}
	}

	if len(c.Format) != 0 {
		for _, v := range c.Format {
			if v.Value == "" {
				logger.Printf("No value specified for the FORMAT/%v", strings.ToUpper(v.Name))
			}
		}
	}

}
