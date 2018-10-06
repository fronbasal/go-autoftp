# go-autoftp

go-autoftp is a utility used to automagically upload a directory to a given FTP server on change.

## Installation

To install this application, use the default go toolchain.

```bash
go get -u -v github.com/fronbasal/go-autoftp
```

## Usage

```
usage: autoftp --server=SERVER --username=USERNAME --password=PASSWORD --directory=DIRECTORY [<flags>]

Flags:
      --help                 Show context-sensitive help (also try --help-long and --help-man).
  -s, --server=SERVER        The FTP host to connect to
  -u, --username=USERNAME    The FTP username
  -p, --password=PASSWORD    The FTP password
  -d, --directory=DIRECTORY  The directory to watch
  -v, --verbose              Enable verbose output
  -f, --overwrite            Overwrite existing files
```

## License

This software is licensed under the AGPLv3 license.

## Maintainers

- [Daniel Malik (mail@fronbasal.de)](https://github.com/fronbasal)
