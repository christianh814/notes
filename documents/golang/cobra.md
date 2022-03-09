# Get the binary

Install it using go

```shell
go install github.com/spf13/cobra/cobra@latest
```

Put in in your path

```shell
sudo mv $GOPATH/bin/cobra /usr/local/bin/
```

Bash completion

```shell
source <(cobra completion bash)
```

# Init

First, init go modules (replace with the repo where it's going to be stored)

```shell
go mod init github.com/christian/mycli
```

Then, use `cobra` to scaffold (use `--help` if you want to know more)

```shell
cobra init --author="Christian Hernandez christian@email.com" --license=apache --viper=true
```

You now have a program

```shell
go run main.go
```

# Commands

Add commands. Like `foobar`

```shell
cobra add foobar --author="Christian Hernandez christian@email.com" --license=apache --viper=true
```

You can now pass `foobar` to your command

```shell
go run main.go foobar
```

You can add as many as you need

# Subcommands

If you want to add a command to `foobar`, you pass the `--parent` command. For example, if you want to add `bazz` to the `foobar` command.

```shell
cobra add bazz --parent="foobarCmd" --author="Christian Hernandez christian@email.com" --license=apache --viper=true
```

Now you can run

```shell
go run main.go foobar bazz
```


See it with...

```shell
go run main.go foobar -h
```

It's good practice to rename the subcommand (in this case `bazz`) with the parent command (in this case `foobar`) as *parent_subcommand.go*

```shell
mv cmd/bazz.go cmd/foobar_bazz.go
```

It should look like this...

```
$ tree  .
.
├── cmd
│   ├── foobar_bazz.go
│   ├── foobar.go
│   └── root.go
├── go.mod
├── go.sum
├── LICENSE
└── main.go
```
