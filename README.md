# zfaas

zfaas is a CLI application to create zrunner pipeline plugin.

Steps to use zfaas:

First initialize a go module in a folder. 

Then use `go get github.com/Zettablock/zfaas` to install the zfaas cli application.

Now you can use `init` command to initialize a new pipeline in your module folder. 

```
zfaas init mypipeline
```

The init command will create a new folder `mypipeline` with several scaffold files. You can edit these files to add your business logic.

