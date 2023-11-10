package bedgovcf

import (
	"bufio"
	"fmt"
	"log"
	"math"
	"os"
	"path/filepath"
	"slices"
	"strconv"
	"strings"

	cli "github.com/urfave/cli/v2"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

// Set the header of the VCF struct according to the config and fai
func (v *Vcf) SetHeader(cCtx *cli.Context, config Config) {
	v.Header.setVersion("4.2")

	if cCtx.String("sample") == "" {
		v.Header.setSample(strings.Split(filepath.Base(cCtx.String("bed")), ".")[0])
	} else {
		v.Header.setSample(cCtx.String("sample"))
	}

	v.Header.setHeaderLines(config)
	v.Header.setContigs(cCtx.String("fai"))

}

// Set the header lines of the VCF struct according to the config
func (h *Header) setHeaderLines(config Config) {
	for k, v := range config.Header {
		h.HeaderLines = append(h.HeaderLines, HeaderLine{
			Category: k,
			Content:  v.Value,
		})
	}

	for _, v := range config.Alt.Options {
		h.HeaderLines = append(h.HeaderLines, HeaderLine{
			Category:    "ALT",
			Id:          v.Name,
			Description: v.Description,
		})
	}

	for k, v := range config.Info {
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
			Id:          k,
			Number:      number,
			Type:        typeField,
			Description: v.Description,
		})
	}

	for k, v := range config.Format {
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
			Id:          k,
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

}

func (h *Header) setSample(sample string) {
	h.Sample = sample
}

func (h *Header) setVersion(version string) {
	h.Version = version
}

// Read the fasta index file and add the contigs to the VCF header
func (h *Header) setContigs(faidx string) {

	file, err := os.Open(faidx)
	if err != nil {
		log.Fatalf("Failed to open the fasta index file: %v", err)
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
}

// Read the BED file and add the variants to the VCF struct
func (v *Vcf) AddVariants(cCtx *cli.Context, config Config) {
	file, err := os.Open(cCtx.String("bed"))
	if err != nil {
		log.Fatalf("Failed to open the bed file: %v", err)
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
			log.Fatal("The amount of columns in the BED file is not consistent.\n Check if there aren't any additional lines at the top of the bed file (and use --skip to tell bedgovcf to skip these lines).")
		}

		variant := Variant{}

		//Standard fields
		variant.Chrom = config.Chrom.getValue(line, header)
		variant.Pos = config.Pos.getValue(line, header)
		variant.Id = config.Id.getValue(line, header)
		variant.Ref = config.Ref.getValue(line, header)
		variant.Alt = config.Alt.getValue(line, header)
		variant.Qual = config.Qual.getValue(line, header)
		variant.Filter = config.Filter.getValue(line, header)
		variant.Info = config.Info.getValues(line, header)
		variant.Format = config.Format.getValues(line, header)

		v.Variants = append(v.Variants, variant)
	}
}

// Get the values of all info fields and transform them to a map
func (mcifs *MapConfigInfoFormatStruct) getValues(values []string, header []string) MapVariantInfoFormat {
	infoMap := MapVariantInfoFormat{}
	for k, v := range *mcifs {
		infoMap[k] = VariantInfoFormat{
			Number: v.Number,
			Type:   v.Type,
			Value:  v.getValue(values, header),
		}
	}
	return infoMap
}

// Get the value for the given field based on the config
func (cifs *ConfigInfoFormatStruct) getValue(values []string, header []string) string {
	var prefix string
	if cifs.Prefix != "" {
		prefix = cifs.Prefix
	}

	// TODO write a resolve function for conditional fields
	if cifs.Field != "" {
		var headerIndex int
		for i, v := range header {
			if v == cifs.Field {
				headerIndex = i
				break
			}
		}
		return prefix + values[headerIndex]
	} else if cifs.Value != "" {
		return prefix + resolveField(cifs.Value, values, header)
	}
	return ""
}

// Get the value for the given field based on the config
func (csfs *ConfigStandardFieldStruct) getValue(values []string, header []string) string {
	var prefix string
	if csfs.Prefix != "" {
		prefix = csfs.Prefix
	}

	// TODO write a resolve function for conditional fields
	if csfs.Field != "" {
		var headerIndex int
		for i, v := range header {
			if v == csfs.Field {
				headerIndex = i
				break
			}
		}
		return prefix + values[headerIndex]
	} else if csfs.Value != "" {
		return prefix + resolveField(csfs.Value, values, header)
	}
	return ""
}

func resolveField(value string, values []string, header []string) string {
	rawInput := strings.Split(value, " ")
	function := ""
	if strings.HasPrefix(rawInput[0], "~") {
		function = rawInput[0][1:]
	} else {
		return value
	}

	input := []string{}
	for _, v := range rawInput {
		if strings.HasPrefix(v, "$") {
			var headerIndex int
			for j, w := range header {
				if w == v[1:] {
					headerIndex = j
					break
				}
			}
			input = append(input, values[headerIndex])
			continue
		}
		input = append(input, v)
	}

	switch function {
	case "round":
		// ~round <value>
		float, err := strconv.ParseFloat(input[1], 64)
		if err != nil {
			log.Fatalf("Failed to parse the value (%v) to a float: %v", input[1], err)
		}
		return fmt.Sprintf("%v", math.Round(float))
	case "sum":
		// ~sum <value1> <value2> ...
		var sum float64
		for _, v := range input[1:] {
			float, err := strconv.ParseFloat(v, 64)
			if err != nil {
				log.Fatalf("Failed to parse the value (%v) to a float: %v", v, err)
			}
			sum += float
		}

		return strconv.FormatFloat(sum, 'g', -1, 64)
	case "min":
		// ~min <startValue> <valueToSubstract1> <valueToSubstract2> ...
		min, err := strconv.ParseFloat(input[1], 64)
		if err != nil {
			log.Fatalf("Failed to parse the value (%v) to a float: %v", input[1], err)
		}
		for _, v := range input[2:] {
			float, err := strconv.ParseFloat(v, 64)
			if err != nil {
				log.Fatalf("Failed to parse the value (%v) to a float: %v", v, err)
			}
			min -= float
		}
		return strconv.FormatFloat(min, 'g', -1, 64)
	case "if":
		// ~if <value1> <operator> <value2> <value_if_true> <value_if_false>
		// supported operators: > < >= <= ==
		v1 := input[1]
		operator := input[2]
		v2 := input[3]
		vTrue := input[4]
		vFalse := input[5]

		floatV1, err1 := strconv.ParseFloat(v1, 64)
		floatV2, err2 := strconv.ParseFloat(v2, 64)

		floatOperators := []string{"<", ">", "<=", ">="}
		if slices.Contains(floatOperators, operator) && (err1 != nil || err2 != nil) {
			log.Fatalf("Failed to parse the values (%v and %v) to a float: %v and %v", v1, v2, err1, err2)
		}

		switch operator {
		case "<":
			if floatV1 < floatV2 {
				return vTrue
			} else {
				return vFalse
			}
		case ">":
			if floatV1 > floatV2 {
				return vTrue
			} else {
				return vFalse
			}
		case ">=":
			if floatV1 >= floatV2 {
				return vTrue
			} else {
				return vFalse
			}
		case "<=":
			if floatV1 <= floatV2 {
				return vTrue
			} else {
				return vFalse
			}
		case "==":
			if err1 == nil && err2 == nil {
				if floatV1 == floatV2 {
					return vTrue
				} else {
					return vFalse
				}
			} else {
				if v1 == v2 {
					return vTrue
				} else {
					return vFalse
				}
			}
		}
	}

	return ""
}

// Write the VCF struct to stdout or a file
func (v *Vcf) Write(cCtx *cli.Context) {
	stdout := true
	if cCtx.String("output") != "" {
		stdout = false
	}

	if stdout {
		fmt.Print(v.Header.String())
		for _, variant := range v.Variants {
			fmt.Print(variant.String())
		}
	} else {
		file, err := os.Create(cCtx.String("output"))
		if err != nil {
			log.Fatalf("Failed to create the output file: %v", err)
		}
		defer file.Close()
		file.WriteString(v.Header.String())
		for _, variant := range v.Variants {
			file.WriteString(variant.String())
		}
	}
}

// Convert a variant to a string
func (v Variant) String() string {

	variant := strings.Join([]string{
		v.Chrom,
		v.Pos,
		v.Id,
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
func (mvif MapVariantInfoFormat) infoString() string {
	var infoSlice []string
	for k, v := range mvif {
		upperInfo := strings.ToUpper(k)
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
func (mcifs MapVariantInfoFormat) formatString() string {
	var formatField []string
	var formatValues []string
	for k, v := range mcifs {
		upperFormat := strings.ToUpper(k)
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
