# Counting stuff with Redis

Being one of the most popular key-value databases, Redis has been perceived my the majority of their users as a dictionary that is used for two operations: setting values and getting them. What is less popular, yet very useful is the way we can keep track of some count values, by simply incrementing them. Let's implement a small "pagehit" summary page then.

### Increment command

The clue behind our use-case lays in one of the built-in commands in Redis, [INCR](https://redis.io/commands/incr). It is very useful for managing numeric data since it's popular to set some value based on its previous state. Implementing a counter is very simple since we just take a key and increment the value by one:

    val, err := client.Get("counter").Int64()
    ...
    fmt.Printf("Old value: %d\n", val)
    // Old value: 2
    
    err = client.Incr("counter").Err()
    ...
    
    val, err = client.Get("counter").Int64()
    ...
    fmt.Printf("New value: %d\n", val)
    // New value: 3

Alternatively, we can use a variations of this command (errors omitted for readability):

    // Value: 3
    client.Decr("counter")
    // Value: 2
    client.IncrBy("counter", 3)
    // Value: 5
    client.DecrBy("counter", 4)
    // Value: 1

### Use case

In order to count page hits all you need to do is increment a per-site (per-URL) counter, right? Since we can access the URL in the middleware, we can increment a counter there as well, leaving each handler clean.

First, let's create an abstraction over Redis DB that has a `Hit(page string)` function:

    type Store interface {
        Hit(page string) error
    }

This will connect to the database and make `INCR`s left and right. 

    func (r redisStore) Hit(url string) error {
        _, err = r.client.Incr(url).Result()

        return err
    }

Now our middleware is a simple one:

    func Middleware(hs Store) echo.MiddlewareFunc {
        return func(next echo.HandlerFunc) echo.HandlerFunc {
            return func(ctx echo.Context) error {
                hs.Hit(ctx.Path())
                return next(ctx)
            }
        }
    }

And that's it! Sure, we should sanitize our page URLs somehow to create some user-friendly keys in the DB (eg. `pagehit.index` instead of `pagehit./`), but let's not bother with this at the moment.

### Downside

The problem with this solution is that in order to see counters for all pages we need to know all URLs and iterate over known keys to list the stats, right? To make it possible, let's create another key in the DB that will store a set of URLs. It requires us to alter the `Hit(page)` function, so that each time a page is accessed, it's added to the set. To do that, we need to use [SADD](https://redis.io/commands/sadd) command:

    func (r redisStore) Hit(url string) error {
        i, err := r.client.SAdd("pages", url).Result()
        if err != nil {
            return err
        }

        if i == 1 {
            fmt.Printf("First time hit to %s/n", url)
        }

        _, err = r.client.Incr(url).Result()

        return err
    }

Last but not least, we'd like to expose the stats to see which page was accessed how many times, so let's add another function to `Store` interface:

    type Stats map[string]int
    ...
    func (r redisStore) GetStats() (Stats, error) {
        var pages []string
        pages, err := r.client.SMembers("pages").Result()
        if err != nil {
            return Stats{}, err
        }

        stats := make(map[string]int)
        for _, p := range pages {
            count, err := r.client.Get(p).Int64()
            if err == nil {
                stats[p] = int(count)
            }
        }

        return stats, nil
    }

This reads all pages from the set, iterates over them and gets al value for respective values. Last but not least, let's create an endpoint to print the values:

    func Handler(hs Store) echo.HandlerFunc {
        return func(ctx echo.Context) error {
            stats, err := hs.GetStats()
            if err != nil {
                return echo.NewHTTPError(http.StatusInternalServerError, err)
            }

            return ctx.JSON(http.StatusOK, stats)
        }
    }

Let's see this in action:

    (redis) $ curl localhost:5000
    Hello World #0
    (redis) $ curl localhost:5000
    Hello World #0
    (redis) $ curl localhost:5000/stats
    {"/":2,"/stats":1}
    (redis) $ curl localhost:5000
    Hello World #0
    (redis) $ curl localhost:5000
    Hello World #0
    (redis) $ curl localhost:5000/stats
    {"/":4,"/stats":2}

Fortunately, there is a way to do it all in a slightly simpler fashion, stick around for the next post on Redis.

The whole source code of this example is available [on Github](https://github.com/mycodesmells/golang-examples/tree/master/redis).
