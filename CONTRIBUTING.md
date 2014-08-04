# Contributing

## Issues

When opening issues, try to be as clear as possible when describing the bug or feature request.
Tag the issue accordingly.

## Pull Requests

To hack on combine:

1. Install as usual (`go get github.com/walle/combine`)
2. Create your feature branch (`git checkout -b my-new-feature`)
3. Ensure everything works and the tests pass (see below)
4. Commit your changes (`git commit -am 'Add some feature'`)

Contribute upstream:

1. Fork combine on GitHub
2. Add your remote (`git remote add fork git@github.com:myuser/repo.git`)
3. Push to the branch (`git push fork my-new-feature`)
4. Create a new Pull Request on GitHub

For other team members:

1. Install as usual (`go get github.com/walle/combine`)
2. Add your remote (`git remote add fork git@github.com:myuser/repo.git`)
3. Pull your revisions (`git fetch fork; git checkout -b my-new-feature fork/my-new-feature`)

Notice: Always use the original import path by installing with `go get`.

## Testing

The integration tests use a bash script and the /tmp directory. They are only tested on OSX. Should work on UNIX based OSes but probably not on Windows.

To run the whole test suite, use the command

    $ go test -cover; ./run_integration_tests.sh

And you should see something like this

    [18:58:57] (master) ~/go_workspace/src/github.com/walle/combine $ go test -cover; ./run_integration_tests.sh
    PASS
    coverage: 40.0% of statements
    ok    github.com/walle/combine  0.011s
    ............
    All tests passed

You can also run the tests separately.

The integration tests will write a file structure to a temporary directory on your disk and try to clean up after it's done.
See [run_integration_tests.sh](run_integration_tests.sh) for more info.