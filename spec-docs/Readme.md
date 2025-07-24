The files in this directory serve as the specifications that are implemented
by this library.

## Notes

1. In `npm-cli-server.md`, under the "Advanced Range Syntax" section, it is
stated that range comparator sets may be separated by ` ` (whitespace) or
by `||` (a double pipe). However, under the previous section titled "Ranges,"
it is stated that a range is composed by separating comparator sets with `||`.
This library chooses to follow the earliest definition, as trying to make a
determination by whitespace seems unreasonably difficult.
