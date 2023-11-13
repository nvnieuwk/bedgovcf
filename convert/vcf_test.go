package bedgovcf

import (
	"testing"
)

func TestSetVersion(t *testing.T) {
	header := Header{}
	header.setVersion("4.2")
	if header.Version != "4.2" {
		t.Fatalf("Expected version to be 4.2, got %s", header.Version)
	}
}

func TestSetSample(t *testing.T) {
	header := Header{}
	header.setSample("test")
	if header.Sample != "test" {
		t.Fatalf("Expected sample to be test, got %s", header.Sample)
	}
}

func TestSetContigs(t *testing.T) {
	header := Header{}
	header.setContigs("../test_data/test.fai")
	if len(header.HeaderLines) != 2 {
		t.Fatalf("Expected 2 contigs, got %d", len(header.HeaderLines))
	}
	contig0 := HeaderLine{
		Category: "contig",
		Id:       "chr1",
		Length:   "248956422",
	}

	contig1 := HeaderLine{
		Category: "contig",
		Id:       "chr2",
		Length:   "242193529",
	}

	if header.HeaderLines[0] != contig0 {
		t.Fatalf("Expected contig 0 to be %v, got %v", contig0, header.HeaderLines[0])
	}
	if header.HeaderLines[1] != contig1 {
		t.Fatalf("Expected contig 1 to be %v, got %v", contig1, header.HeaderLines[1])
	}
}

func TestSetHeaderLines(t *testing.T) {
	header := Header{}
	config := Config{
		Header: map[string]ConfigHeaderStruct{
			"test": {
				Name:  "test",
				Value: "test",
			},
		},
		Alt: ConfigStandardFieldStruct{
			Value: "$test",
			Options: []ConfigHeaderStruct{
				{
					Name:        "DUP",
					Description: "Duplication",
				},
			},
		},
		Filter: ConfigStandardFieldStruct{
			Value: "$test",
			Options: []ConfigHeaderStruct{
				{
					Name:        "PASS",
					Description: "Passed filters",
				},
			},
		},
		Info: MapConfigInfoFormatStruct{
			"SVLEN": ConfigInfoFormatStruct{
				Value:       "$test",
				Description: "The length of the SV",
				Number:      "1",
				Type:        "Integer",
			},
		},
		Format: MapConfigInfoFormatStruct{
			"GT": ConfigInfoFormatStruct{
				Value:       "$test",
				Description: "Genotype",
				Number:      "1",
				Type:        "String",
			},
		},
	}
	header.setHeaderLines(config)
	if len(header.HeaderLines) != 5 {
		t.Fatalf("Expected 5 header lines, got %d", len(header.HeaderLines))
	}
	headerLine0 := HeaderLine{
		Category: "test",
		Content:  "test",
	}
	headerLine1 := HeaderLine{
		Category:    "ALT",
		Id:          "DUP",
		Description: "Duplication",
	}
	headerLine2 := HeaderLine{
		Category:    "INFO",
		Id:          "SVLEN",
		Description: "The length of the SV",
		Number:      "1",
		Type:        "Integer",
	}
	headerLine3 := HeaderLine{
		Category:    "FORMAT",
		Id:          "GT",
		Description: "Genotype",
		Number:      "1",
		Type:        "String",
	}
	headerLine4 := HeaderLine{
		Category:    "FILTER",
		Id:          "PASS",
		Description: "Passed filters",
	}
	if header.HeaderLines[0] != headerLine0 {
		t.Fatalf("Expected header line 0 to be %v, got %v", headerLine0, header.HeaderLines[0])
	}
	if header.HeaderLines[1] != headerLine1 {
		t.Fatalf("Expected header line 1 to be %v, got %v", headerLine1, header.HeaderLines[1])
	}
	if header.HeaderLines[2] != headerLine2 {
		t.Fatalf("Expected header line 2 to be %v, got %v", headerLine2, header.HeaderLines[2])
	}
	if header.HeaderLines[3] != headerLine3 {
		t.Fatalf("Expected header line 3 to be %v, got %v", headerLine3, header.HeaderLines[3])
	}
	if header.HeaderLines[4] != headerLine4 {
		t.Fatalf("Expected header line 4 to be %v, got %v", headerLine4, header.HeaderLines[4])
	}

}

func TestStandardGetValue(t *testing.T) {
	config := ConfigStandardFieldStruct{
		Value: "$test",
	}
	header := []string{"test", "test2"}
	values := []string{"value", "I don't want this"}
	value := config.getValue(values, header)
	if value != "value" {
		t.Fatalf("Expected value to be 'value', got %s", value)
	}

	config = ConfigStandardFieldStruct{
		Value:  "test",
		Prefix: "hello_",
	}
	value = config.getValue(values, header)
	if value != "hello_test" {
		t.Fatalf("Expected value to be 'hello_test', got %s", value)
	}

	config = ConfigStandardFieldStruct{
		Value: "$2",
	}
	header = []string{"0", "1", "2", "3"}
	values = []string{"value", "I don't want this", "this is the one", "definitely not this"}
	value = config.getValue(values, header)
	if value != "this is the one" {
		t.Fatalf("Expected value to be 'this is the one', got %s", value)
	}
}

func TestInfoFormatGetValue(t *testing.T) {
	config := ConfigInfoFormatStruct{
		Value: "$test",
	}
	header := []string{"test", "test2"}
	values := []string{"value", "I don't want this"}
	value := config.getValue(values, header)
	if value != "value" {
		t.Fatalf("Expected value to be 'value', got %s", value)
	}

	config = ConfigInfoFormatStruct{
		Value:  "test",
		Prefix: "hello_",
	}
	value = config.getValue(values, header)
	if value != "hello_test" {
		t.Fatalf("Expected value to be 'hello_test', got %s", value)
	}

	config = ConfigInfoFormatStruct{
		Value: "$2",
	}
	header = []string{"0", "1", "2", "3"}
	values = []string{"value", "I don't want this", "this is the one", "definitely not this"}
	value = config.getValue(values, header)
	if value != "this is the one" {
		t.Fatalf("Expected value to be 'this is the one', got %s", value)
	}
}

func TestVariantString(t *testing.T) {
	variant := Variant{
		Chrom:  "chr1",
		Pos:    "123",
		Id:     "test",
		Ref:    "A",
		Alt:    "C",
		Qual:   "100",
		Filter: "PASS",
		Info: MapVariantInfoFormat{
			"SVLEN": VariantInfoFormat{
				Number: "1",
				Type:   "Integer",
				Value:  "100",
			},
		},
		Format: MapVariantInfoFormat{
			"GT": VariantInfoFormat{
				Number: "1",
				Type:   "String",
				Value:  "0/1",
			},
		},
	}

	if variant.String(1) != "chr1\t123\ttest1\tA\tC\t100\tPASS\tSVLEN=100\tGT\t0/1\n" {
		t.Fatalf("Expected variant string to be 'chr1\t123\ttest\tA\tC\t100\tPASS\tSVLEN=100\tGT\t0/1\n', got '%s'", variant.String(1))
	}

	variant = Variant{
		Chrom:  "chr1",
		Pos:    "123",
		Id:     "test",
		Ref:    "A",
		Alt:    "C",
		Qual:   "100",
		Filter: "PASS",
		Info: MapVariantInfoFormat{
			"SVLEN": VariantInfoFormat{
				Number: "1",
				Type:   "Integer",
				Value:  "100",
			},
			"SVTYPE": VariantInfoFormat{
				Number: "1",
				Type:   "String",
				Value:  "DEL",
			},
		},
		Format: MapVariantInfoFormat{
			"GT": VariantInfoFormat{
				Number: "1",
				Type:   "String",
				Value:  "0/1",
			},
			"CN": VariantInfoFormat{
				Number: "1",
				Type:   "Integer",
				Value:  "2",
			},
		},
	}

	if variant.String(1) != "chr1	123	test1	A	C	100	PASS	SVLEN=100;SVTYPE=DEL	GT:CN	0/1:2\n" {
		t.Fatalf("Expected variant string to be 'chr1	123	test	A	C	100	PASS	SVLEN=100;SVTYPE=DEL	GT:CN	0/1:2\n', got '%s'", variant.String(1))
	}

}

func TestHeader(t *testing.T) {
	header := Header{
		Version: "4.2",
		Sample:  "test",
		HeaderLines: []HeaderLine{
			{
				Category: "test",
				Content:  "test",
			},
			{
				Category:    "ALT",
				Id:          "DEL",
				Description: "Deletion",
			},
			{
				Category: "contig",
				Id:       "chr1",
				Length:   "123",
			},
			{
				Category:    "INFO",
				Id:          "SVLEN",
				Number:      "1",
				Type:        "Integer",
				Description: "The length of the structural variant",
			},
		},
	}

	testString := "##fileformat=VCFv4.2\n" +
		"##test=test\n" +
		"##ALT=<ID=DEL,Description=\"Deletion\">\n" +
		"##contig=<ID=chr1,length=123>\n" +
		"##INFO=<ID=SVLEN,Number=1,Type=Integer,Description=\"The length of the structural variant\">\n" +
		"#CHROM\tPOS\tID\tREF\tALT\tQUAL\tFILTER\tINFO\tFORMAT\ttest\n"

	if header.String() != testString {
		t.Fatalf("Expected header string to be '%s', got '%s'", testString, header.String())
	}
}
