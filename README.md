# Bedgovcf

BedGoVcf is a simple tool to dynamically convert BED files to VCF files using a YAML configuration file. 

## Usage
```bash
bedgovcf --bed <input.bed> --config <config.yaml> --fai <reference.fai>
```

### Required Arguments
| Argument | Description |
| --- | --- |
| `--bed <path>` | Path to the BED file to convert |
| `--config <path>` | Path to the YAML configuration file |
| `--fai <path>` | Path to the FASTA index file of the reference genome |

### Optional Arguments
| Argument | Description |
| --- | --- |
| `--output <path>` | Path to the output VCF file (default: stdout) |
| `--skip <integer>` | Skip the first N lines of the BED file (default: 0) |
| `--header` | The BED file has a header (default: false) |
| `--sample <string>` | Sample name to use in the VCF file (default: prefix of the BED file) |

## The configuration file
The configuration file can be used to tell `bedgovcf` how to handle the BED file. It is a YAML file with the following structure:

:warning: All names should be lowercase, otherwise the tool won't recognize them. :warning:

```yaml
# Optional headers to add to the VCF file
header:
  - name: header_name # The name of the header
    content: this header does something # The content of the header
  - name: other_header_name
    content: this header does something else

# Optional chromosome field (will default to the first column of the BED file)
chrom:
  value: $0 # The value to use for the chromosome field
  prefix: chr # A prefix to add to the chromosome field

# Optional position field (will default to the second column of the BED file)
pos:
  value: $1 # The value to use for the position field

# Optional ID field (will default to an incrementing number)
id:
  value: $5 # The value to use for the ID field
  prefix: toolname_ # A prefix to add to the ID field

# Optional reference field (will default to N)
ref:
  value: "N" # The value to use for the reference field

# Optional alternate field (will default to <CNV>)
alt:
  value: ~if $4 < 0 <DEL> <DUP> # The value to use for the alternate field
  options: # All possibilities for the alternate field
    - name: DEL
      description: Deletion
    - name: DUP
      description: Duplication

# Optional quality field (will default to '.')
qual:
  value: $6 # The value to use for the quality field

# Optional filter field (will default to PASS)
filter:
  value: ~if $6 >= 0 PASS LOWQUAL # The value to use for the filter field
  options: # All possibilities for the filter field
    - name: PASS
      description: Passed all filters
    - name: LOWQUAL
      description: Low quality

# Optional info fields (will default to no info fields)
info:
  - name: SVTYPE # The name of the info field
    value: ~if $4 < 0 DEL DUP # The value to use for the info field
    number: 1 # The number of values for the info field
    type: String # The type of the info field
    description: Type of structural variant # The description of the info field
  - name: SVLEN
    value: ~min $2 $1
    number: 1
    type: Integer
    description: Length of structural variant
  - name: END
    value: $2
    number: 1
    type: Integer
    description: End position of structural variant

# Optional format fields (will default to no format fields)
format:
  - name: GT # The name of the format field
    value: ~if $4 < 0 0/1 1/1 # The value to use for the format field
    number: 1 # The number of values for the format field
    type: String # The type of the format field
    description: Genotype of the sample # The description of the format field
  - name: CN
    value: ~round $4
    number: 1
    type: Integer
    description: Copy number of the sample
```

### Dynamically fetching fields from the BED file
All `value` fields in the config can be resolved by column names. For example, if you have a BED file with the following columns:

```tsv
start	end	chromosome
1	2	chr1
```

You can use the following config to resolve the `chrom` and `pos` fields (when using the `--header` option):

```yaml
chrom:
  value: $chromosome
pos:
  value: $start
```

When no header is present you can also use the 0-based index of the column like this:

```yaml
chrom:
  value: $2
pos:
  value: $0
```

### Functions
The `value` fields in the config can also be resolved by using functions (words starting with `~`). The following functions are available:

#### `~round`
Pattern: `~round <value>`

Rounds the given value to the nearest integer.

:warning: This function will panic if the given value is not an integer or a float. :warning:

#### `~sum`
Pattern: `~sum <value1> <value2> ...`

Adds all values together.

:warning: This function will panic if the given values are not integers or floats. :warning:

#### `~min`
Pattern: `~min <value1> <value2> ...`

Substracts all values from the first value.

:warning: This function will panic if the given values are not integers or floats. :warning:

#### `~if`
Pattern: `~if <value1> <operator> <value2> <value_if_true> <value_if_false>`

Checks if the given condition (`<value1> <operator> <value2>`) is true for the given values. If it is true, the `value_if_true` will be returned, otherwise the `value_if_false` will be returned. `value_if_false` can also be a new function (this way you can create nested if statements).

Supported operators: `<`, `<=`, `>`, `>=`, `==`, `!=`

## Installation
### Mamba/Conda
This is the preffered way of installing BedGoVcf.

```bash
mamba install -c bioconda bedgovcf
```

or with conda:
  
```bash 
conda install -c bioconda bedgovcf
```

### Precompiled binaries
Precompiled binaries are available for Linux and macOS on the [releases page](https://github.com/nvnieuwk/bedgovcf/releases).


### Installation from source
Make sure you have go installed on your machine (or [install](https://go.dev/doc/install) it if you don't currently have it)

Then run these commands to install bedgovcf:

```bash
go get .
go build .
sudo mv bedgovcf /usr/local/bin/
```

Next run this command to check if it was correctly installed:

```bash
bedgovcf --help
```
