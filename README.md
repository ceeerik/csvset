# csvset
Used to perform set operations on csv files.

## flags

**--input** - List of csv.files (comma-seperated, locations based on current location) to be used in the program. They can then be used in the formula based on the order in which the were added, indexed from `0`.

**--output** - Specifies the filename/path of the output.

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

Example:

`(1+2+3+(4/5/6))-7-8`

## examples