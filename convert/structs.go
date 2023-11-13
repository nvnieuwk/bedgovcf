package bedgovcf

//
// CONFIG
//

// The main config struct
type Config struct {
	Header []ConfigHeaderStruct        // Additional headers to add to the VCF
	Chrom  ConfigStandardFieldStruct   // The chromosome field
	Pos    ConfigStandardFieldStruct   // The position field
	Id     ConfigStandardFieldStruct   // The ID field
	Ref    ConfigStandardFieldStruct   // The reference field
	Alt    ConfigStandardFieldStruct   // The alt field
	Qual   ConfigStandardFieldStruct   // The quality field
	Filter ConfigStandardFieldStruct   // The filter field
	Info   SliceConfigInfoFormatStruct // The info fields
	Format SliceConfigInfoFormatStruct // The format fields
}

// The struct for the additional headers
type ConfigHeaderStruct struct {
	Name        string // The name of the header line
	Content     string // The content of the header line
	Description string // The description of the header line
}

// The struct for the standard fields
type ConfigStandardFieldStruct struct {
	Value   string               // The value to use
	Prefix  string               // The prefix to add to each value
	Options []ConfigHeaderStruct // The different options possible (only for ALT and FILTER)
}

type SliceConfigInfoFormatStruct []ConfigInfoFormatStruct

// The struct for the info and format fields
type ConfigInfoFormatStruct struct {
	Name        string // The name of the current INFO or FORMAT field
	Value       string // The value to use
	Prefix      string // The prefix to add to each value
	Description string // The description of the field
	Number      string // The number of values that can be included in the INFO field (e.g. 1, 2, A, R)
	Type        string // The type of the header field (e.g. Integer, Float, Character, Flag)
}

//
// VCF
//

// The main VCF struct
type Vcf struct {
	Header   Header    // The header of the VCF
	Variants []Variant // The variants of the VCF
}

// The struct for the header
type Header struct {
	HeaderLines []HeaderLine // All conventional header lines
	Version     string       // The version of the VCF file
	Sample      string       // The sample name
}

// The struct for one header line
type HeaderLine struct {
	Category    string // The category of header line (e.g INFO, FORMAT, FILTER)
	Id          string // The ID of the header line (e.g. SVLEN, GT, END)
	Number      string // The number of values that can be included in the INFO field (e.g. 1, 2, A, R)
	Type        string // The type of the header field (e.g. Integer, Float, Character, Flag)
	Description string // The description of the header line
	Length      string // The length of the contig (only for contig header lines)
	Content     string // The content of the header line (only for non usual header lines)
}

// The struct for one variant
type Variant struct {
	Chrom  string                 // The chromosome
	Pos    string                 // The position
	Id     string                 // The ID
	Ref    string                 // The reference allele
	Alt    string                 // The alternative allele
	Qual   string                 // The quality
	Filter string                 // The filter
	Info   SliceVariantInfoFormat // The info fields
	Format SliceVariantInfoFormat // The format fields
}

// The map for the info and format fields
type SliceVariantInfoFormat []VariantInfoFormat

// The struct for one info or format field
type VariantInfoFormat struct {
	Name   string // The name of the current INFO or FORMAT field
	Number string // The number of values that can be included in the INFO field (e.g. 1, 2, A, R)
	Type   string // The type of the header field (e.g. Integer, Float, Character, Flag)
	Value  string // The value of the field
}
