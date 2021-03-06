# jkl
jkl is a library for programmatically interacting with a JIRA installation.  It
comes with a command line program (also called `jkl`) which allows you to
interact with JIRA via the command line.

## Installation
To use the library, simply import it into your application:
`import "github.com/otremblay/jkl"`

To install the command line application:
First, make sure you have a working go environment:
https://golang.org/doc/install

Then, execute the following command from your shell:

`$ go get github.com/otremblay/jkl/cmd/jkl`

## Usage

Make sure you create a `~/.jklrc` file in your home directory, it should contain
at a minimum:

```
JIRA_ROOT="https://jira.example.com/"
JIRA_USER="myusername"
JIRA_PASSWORD="mypassword"
JIRA_PROJECT="DPK"
```
Those values are for example only, your setup will be different.

TODO: Finish writing usage instructions

## Contributing

1. Fork it!
2. Create your feature branch: `git checkout -b my-new-feature`
3. Commit your changes: `git commit -am 'Add some feature'`
4. Push to the branch: `git push origin my-new-feature`
5. Submit a pull request :D
