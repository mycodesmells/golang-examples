# Loading config from ENV in Go

Most of the applications we build need some kind of configuration that can alter its behavior. The level of logging, HTTP to run the server on, database credentials - you name it. While you can ship your app with YAML file, thanks to Docker (among others) it's more popular to use environmental variables. How exactly do we do it, though? Let's see what are our options.

### Manual approach

The simplest thing we can do is to call `os.Getenv(var_name)` every time we need to use some configurable value. A simple greeting app would look like this:

    ...
    func main() {
        username := os.Getenv("USERNAME")
        fmt.Printf("Hello, %s!\n", username)
    }

It's plain and simple, but there are a few things need to be aware of. First of all, we need to handle the situation when the variable is not defined. This uglifies our simple code with additional three-line `if` clause. To be honest, it wouldn't be that bad, we are using Go after all. The more important _problem_ is that we are not sure what our configuration consists of. Obviously, in such a tiny app you remember that there is only one variable, but what if you have a large web application that requires twenty values? Would you remember them all? 

### Parsing ENV into struct

Fortunately, there is a way to tackle both problems we faced with the manual approach - we can use `github.com/caarlos0/env` to parse current environment into a struct that both defines exactly what configuration is (all the attributes) and allows us to define default values. An updated version of the greeting app looks like this:

    ...
    import "github.com/caarlos0/env"

    type config struct {
        Username string `env:"USERNAME" envDefault:"Slomek"`
    }

    func main() {
        var cfg config
        if err := env.Parse(&cfg); err != nil {
            log.Fatalf("Failed to parse ENV")
        }
        fmt.Printf("Hello, %s!\n", cfg.Username)
    }

This is obviously a bit longer than in the previous snippet, but an overall gain is much bigger. It's important to handle an error that can occur during parsing, but this is a general rule in Go, so you should be used to it.

Is there any downside of that? It's not very easy to use this with our aforementioned 15-variable configurations while working in local development environment. Do you really want to set those vars forever in your OS? Or do you want to define them every time you call the application like:

    `USERNAME=Tom go run main.go`

I thought so.

### .env to the rescue

One last component I recommend to use while reading environment variables is using a special file that overwrites current state of your OS. To do that you can use `github.com/joho/godotenv` that reads a file (`.env` by default, but you can use any file you want) and inserts values defined there into the application's ENV:

    import (
        "log"

        "github.com/caarlos0/env"
        "github.com/joho/godotenv"
    )

    type config struct {
        Username string `env:"USERNAME" envDefault:"Slomek"`
    }

    func main() {
        if err := godotenv.Load(); err != nil {
            log.Println("File .env not found, reading configuration from ENV")
        }

        var cfg config
        if err := env.Parse(&cfg); err != nil {
            log.Fatalln("Failed to parse ENV")
        }
        log.Printf("Hello, %s!\n", cfg.Username)
    }

As you can see, this works perfectly alongside `caarlos0/env`, which is why I like using them both at the same time.

What is interesting, it's still possible to overwrite some values _ad hoc_:

    # .env:
    # USERNAME=George
    $ go run main.go 
    2017/07/28 21:28:41 Hello, George!

    $ USERNAME=Oscar go run main.go 
    2017/07/28 21:28:24 Hello, Oscar!

### Summary

If you are not using ENV to create the configuration for your app, I totally recommend it. This approach follows directives described as [The twelve-factor methodology](https://12factor.net/config). It works great with Docker containers - you can use the same image for development and production by changing their ENV.

The full source code of the examples is available [on Github](https://github.com/mycodesmells/golang-examples/tree/master/misc/env-config).
