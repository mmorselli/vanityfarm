
# Vanityfarm

Vanityfarm is Algorand vanity address generator written in Go. Unlike other similar programs, Vanityfarm will look for interesting words in a dictionary for you, you only need to choose the minimum and maximum length of these words

# Usage

`vanityfarm <minchar> <maxchar>`

# Examples

`vanityfarm 5 12` - program will try to find words between 5 and 12 characters

It is infinity process, to stop you may send CTRL+C to the program.

# Build

To build the vanityfarm, you need to download and install the Go compiler

Then just run the following command in the source code directory:

```
go build
```

After successful build, it will produce `vanityfarm` (or `vanityfarm.exe` on Windows OS) named binary file


# License

vanityfarm is distributed under the terms GNU General Public License (Version 3).

See [LICENSE](./LICENSE) for details.