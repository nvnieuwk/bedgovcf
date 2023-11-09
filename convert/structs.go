package bedgovcf

type Config struct {
	Header map[string]ConfigHeaderStruct
	Chrom  ConfigStandardFieldStruct
	Pos    ConfigStandardFieldStruct
	Id     ConfigStandardFieldStruct
	Ref    ConfigStandardFieldStruct
	Alt    ConfigStandardFieldStruct
	Qual   ConfigStandardFieldStruct
	Filter ConfigStandardFieldStruct
	Info   MapConfigInfoFormatStruct
	Format MapConfigInfoFormatStruct
}

type ConfigHeaderStruct struct {
	Name        string
	Value       string
	Description string
}

type ConfigStandardFieldStruct struct {
	Field   string
	Value   string
	Prefix  string
	Options []ConfigHeaderStruct
}

type MapConfigInfoFormatStruct map[string]ConfigInfoFormatStruct

type ConfigInfoFormatStruct struct {
	Field       string
	Value       string
	Prefix      string
	Description string
	Number      string
	Type        string
}

type HeaderLine struct {
	Category    string // The category of header line (e.g INFO, FORMAT, FILTER)
	Id          string // The ID of the header line (e.g. SVLEN, GT, END)
	Number      string // The number of values that can be included in the INFO field (e.g. 1, 2, A, R)
	Type        string // The type of the header field (e.g. Integer, Float, Character, Flag)
	Description string // The description of the header line
	Length      string // The length of the contig (only for contig header lines)
	Content     string // The content of the header line (only for non usual header lines)
}

type Header struct {
	HeaderLines []HeaderLine // All conventional header lines
	Version     string       // The version of the VCF file
	Sample      string       // The sample name
}

type Vcf struct {
	Header   Header
	Variants []Variant
}

type Variant struct {
	Chrom  string
	Pos    string
	Id     string
	Ref    string
	Alt    string
	Qual   string
	Filter string
	Info   MapVariantInfoFormat
	Format MapVariantInfoFormat
}

type MapVariantInfoFormat map[string]VariantInfoFormat

type VariantInfoFormat struct {
	Number string
	Type   string
	Value  string
}
