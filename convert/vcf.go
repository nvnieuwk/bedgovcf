package bedgovcf

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	cli "github.com/urfave/cli/v2"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

// Set the header of the VCF struct according to the config and fai
func (v *Vcf) SetHeader(cCtx *cli.Context, config Config) error {
	err := v.Header.setVersion("4.2")
	if err != nil {
		return err
	}

	if cCtx.String("sample") == "" {
		err = v.Header.setSample(strings.Split(filepath.Base(cCtx.String("bed")), ".")[0])
	} else {
		err = v.Header.setSample(cCtx.String("sample"))
	}
	if err != nil {
		return err
	}

	err = v.Header.setHeaderLines(config)
	if err != nil {
		return err
	}

	err = v.Header.setContigs(cCtx.String("fai"))
	if err != nil {
		return err
	}

	return nil
}

// Set the header lines of the VCF struct according to the config
func (h *Header) setHeaderLines(config Config) error {
	for _, v := range config.Header {
		h.HeaderLines = append(h.HeaderLines, HeaderLine{
			Category: v.Name,
			Content:  v.Content,
		})
	}

	for _, v := range config.Alt.Options {
		h.HeaderLines = append(h.HeaderLines, HeaderLine{
			Category:    "ALT",
			Id:          v.Name,
			Description: v.Description,
		})
	}

	for _, v := range config.Info {
		number := v.Number
		if number == "" {
			number = "."
		}
		typeField := v.Type
		if typeField == "" {
			typeField = "String"
		}
		h.HeaderLines = append(h.HeaderLines, HeaderLine{
			Category:    "INFO",
			Id:          v.Name,
			Number:      number,
			Type:        typeField,
			Description: v.Description,
		})
	}

	for _, v := range config.Format {
		number := v.Number
		if number == "" {
			number = "."
		}
		typeField := v.Type
		if typeField == "" {
			typeField = "String"
		}
		h.HeaderLines = append(h.HeaderLines, HeaderLine{
			Category:    "FORMAT",
			Id:          v.Name,
			Number:      number,
			Type:        typeField,
			Description: v.Description,
		})
	}

	for _, v := range config.Filter.Options {
		h.HeaderLines = append(h.HeaderLines, HeaderLine{
			Category:    "FILTER",
			Id:          v.Name,
			Description: v.Description,
		})
	}

	return nil
}

func (h *Header) setSample(sample string) error {
	h.Sample = sample
	return nil
}

func (h *Header) setVersion(version string) error {
	h.Version = version
	return nil
}

// Read the fasta index file and add the contigs to the VCF header
func (h *Header) setContigs(faidx string) error {

	file, err := os.Open(faidx)
	if err != nil {
		return fmt.Errorf("failed to open the fasta index file: %v", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.Split(scanner.Text(), "\t")
		h.HeaderLines = append(h.HeaderLines, HeaderLine{
			Category: "contig",
			Id:       line[0],
			Length:   line[1],
		})
	}
	return nil
}

// Read the BED file and add the variants to the VCF struct
func (v *Vcf) AddVariants(cCtx *cli.Context, config Config) error {
	file, err := os.Open(cCtx.String("bed"))
	if err != nil {
		return fmt.Errorf("failed to open the bed file: %v", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	header := []string{}
	var skipCount int64

	for scanner.Scan() {
		if skipCount < cCtx.Int64("skip") {
			skipCount++
			continue
		}
		line := strings.Split(scanner.Text(), "\t")

		if len(header) == 0 {
			if cCtx.Bool("header") {
				header = line
				continue
			} else {
				for k := range line {
					header = append(header, fmt.Sprintf("%v", k))
				}
			}
		}

		if len(line) != len(header) {
			return errors.New("the amount of columns in the BED file is not consistent\n check if there aren't any additional lines at the top of the bed file (and use --skip to tell bedgovcf to skip these lines)")
		}

		variant := Variant{}

		//Standard fields
		variant.Chrom, err = config.Chrom.getValue(line, header)
		if err != nil {
			return err
		}
		variant.Pos, err = config.Pos.getValue(line, header)
		if err != nil {
			return err
		}
		variant.Id, err = config.Id.getValue(line, header)
		if err != nil {
			return err
		}
		variant.Ref, err = config.Ref.getValue(line, header)
		if err != nil {
			return err
		}
		variant.Alt, err = config.Alt.getValue(line, header)
		if err != nil {
			return err
		}
		variant.Qual, err = config.Qual.getValue(line, header)
		if err != nil {
			return err
		}
		variant.Filter, err = config.Filter.getValue(line, header)
		if err != nil {
			return err
		}
		variant.Info, err = config.Info.getValues(line, header)
		if err != nil {
			return err
		}
		variant.Format, err = config.Format.getValues(line, header)
		if err != nil {
			return err
		}

		v.Variants = append(v.Variants, variant)
	}

	return nil
}

// Get the values of all info fields and transform them to a map
func (mcifs *SliceConfigInfoFormatStruct) getValues(values []string, header []string) (SliceVariantInfoFormat, error) {
	infoMap := SliceVariantInfoFormat{}
	for _, v := range *mcifs {
		value, err := v.getValue(values, header)
		if err != nil {
			return nil, err
		}
		infoMap = append(infoMap, VariantInfoFormat{
			Name:   v.Name,
			Number: v.Number,
			Type:   v.Type,
			Value:  value,
		})
	}
	return infoMap, nil
}

// Get the value for the given field based on the config
func (cifs *ConfigInfoFormatStruct) getValue(values []string, header []string) (string, error) {
	var prefix string
	if cifs.Prefix != "" {
		prefix = cifs.Prefix
	}

	value, err := resolveField(strings.Split(cifs.Value, " "), values, header)
	if err != nil {
		return "", err
	}

	return prefix + value, nil

}

// Get the value for the given field based on the config
func (csfs *ConfigStandardFieldStruct) getValue(values []string, header []string) (string, error) {
	var prefix string
	if csfs.Prefix != "" {
		prefix = csfs.Prefix
	}

	value, err := resolveField(strings.Split(csfs.Value, " "), values, header)
	if err != nil {
		return "", err
	}

	return prefix + value, nil

}

// Write the VCF struct to stdout or a file
func (v *Vcf) Write(cCtx *cli.Context) error {
	stdout := true
	if cCtx.String("output") != "" {
		stdout = false
	}

	if stdout {
		fmt.Print(v.Header.String())
		for count, variant := range v.Variants {
			fmt.Print(variant.String(count))
		}
	} else {
		file, err := os.Create(cCtx.String("output"))
		if err != nil {
			return fmt.Errorf("failed to create the output file: %v", err)
		}
		defer file.Close()
		file.WriteString(v.Header.String())
		for count, variant := range v.Variants {
			file.WriteString(variant.String(count))
		}
	}
	return nil
}

// Convert a variant to a string
func (v Variant) String(count int) string {

	id := fmt.Sprintf("%v%v", v.Id, count)

	variant := strings.Join([]string{
		v.Chrom,
		v.Pos,
		id,
		v.Ref,
		v.Alt,
		v.Qual,
		v.Filter,
		v.Info.infoString(),
		v.Format.formatString(),
	}, "\t")
	return variant + "\n"
}

// Convert the info map to a string
func (mvif SliceVariantInfoFormat) infoString() string {
	var infoSlice []string
	for _, v := range mvif {
		upperInfo := strings.ToUpper(v.Name)
		switch infoType := strings.ToLower(v.Type); infoType {
		case "flag":
			infoSlice = append(infoSlice, upperInfo)
		default:
			infoSlice = append(infoSlice, fmt.Sprintf("%v=%v", upperInfo, v.Value))
		}
	}

	return strings.Join(infoSlice, ";")
}

// Convert the format map to a string
func (mvif SliceVariantInfoFormat) formatString() string {
	var formatField []string
	var formatValues []string
	for _, v := range mvif {
		upperFormat := strings.ToUpper(v.Name)
		formatField = append(formatField, upperFormat)
		formatValues = append(formatValues, v.Value)
	}

	return strings.Join(formatField, ":") + "\t" + strings.Join(formatValues, ":")
}

// Convert the VCF header to a string
func (h Header) String() string {
	header := ""
	header += fmt.Sprintf("##fileformat=VCFv%v\n", h.Version)
	for _, v := range h.HeaderLines {
		header += fmt.Sprintf("%v\n", v.String())
	}
	header += fmt.Sprintf("#CHROM\tPOS\tID\tREF\tALT\tQUAL\tFILTER\tINFO\tFORMAT\t%v\n", h.Sample)
	return header
}

// Convert the VCF header line to a string
func (h HeaderLine) String() string {
	line := ""
	switch category := strings.ToLower(h.Category); category {
	case "contig":
		line = fmt.Sprintf("##%v=<ID=%v,length=%v>", strings.ToLower(h.Category), h.Id, h.Length)
	case "info", "format":
		lineType := cases.Title(language.English, cases.Compact).String(strings.ToLower(h.Type))
		line = fmt.Sprintf("##%v=<ID=%v,Number=%v,Type=%v,Description=\"%v\">", strings.ToUpper(h.Category), strings.ToUpper(h.Id), h.Number, lineType, h.Description)
	case "alt", "filter":
		line = fmt.Sprintf("##%v=<ID=%v,Description=\"%v\">", strings.ToUpper(h.Category), strings.ToUpper(h.Id), h.Description)
	default:
		line = fmt.Sprintf("##%v=%v", h.Category, h.Content)
	}

	return line
}
