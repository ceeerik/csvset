# csvset
Used to perform set operations on csv files.

## makefile

### build

`make build`

### install

Will also build prior to installing.

`make install`

### uninstalling

`make clean`

## flags

**--input** - List of csv.files (comma-seperated, locations based on current location) to be used in the program. They can then be used in the formula based on the order in which the were added, indexed from `0`.

**--output** - Specifies the filename/path of the output. **Will quietly overwrite an existing file if it has the same filename; Even with default filename. (`output.csv`)**

**--formula** - Specifies the calculation to be performed.

## operands

- `+` - union
  - Returns all members in any set.
  - Can be chained.
- `*` - intersection
  - Returns members that are in all sets only.
  - Can be chained.
- `-` - set difference
  - Returns members of the first set that are not in any of the other sets.
  - Can be chained.
- `/` - symmetric difference
  - Returns members that belong to one, and only one, set.
  - Can be chained.

## nesting

There's no operation priority so only the one kind of operation can be performed at a given level of nesting; operation order has to be set explicitly using parantheses.

Example:

`(1+2+3+(4/5/6))-7-0`

## full example

`csvset --input a.csv,b.csv,c.csv --output test.csv --formula "(0+1+2)/(2-0)"`