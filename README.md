# zetta-go-cli

zetta-go-cli is a CLI application to create zrunner pipeline plugin.

Steps to use zetta-go-cli:

First initialize a go module in a folder. 

Then use `go get github.com/Zettablock/zetta-go-client` to install the zetta cli application.

Now you can use `init` command to initialize a new pipeline in your module folder. 

```
zetta-go init
```

The init command will create a new folder `example-pipeline` with several scaffold files. You can edit these files to add your business logic.

