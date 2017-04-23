# Counting stuff in hashes

Having implemented counting in Redis using one counter per value seems the way to go, as long as we know about all the things we want to keep track of. Once we want to monitor some more dynamic structures, this quickly becomes an issue. There is, however, a way to do this smart. We need to know what are hashes.

### Simple solution

In the [previous post](http://mycodesmells.com/post/counting-stuff-with-redis) we tried to implement a simple page-hit counter for all paths in the application using separate counters for every page. This quickly hit us back, when we wanted to list all pages with their respective page hit count. In order to do that, we needed to add another key to the store, where we kept a set of all pages, which we later iterated through to ask Redis for each counter.

### Redis Hashes

Apart from strings and sets (that were used in the previous solution), Redis allows us to keep data in _hashes_, which represent maps (string key to string value), which are perfect for our case. Not only can we keep our counter there (alongside with the pages their represent which really saves us a lot of time nad database queries), but it comes with a utility functions like [HINCRBY](https://redis.io/commands/hincrby) to change numeric value within the map, and, obviously, [HGETALL](https://redis.io/commands/hgetall) to fetch all key-value pairs.

### Source code

In order to prove that we make our existing solution better, we need to stick to an existing interface:

    type Stats map[string]int

    type Store interface {
        GetStats() (Stats, error)
        Hit(page string) error
    }

    func NewRedisStore() Store {
        client := redis.NewClient(&redis.Options{
            Addr:     "localhost:6379",
            Password: "", // no password set
            DB:       0,  // use default DB
        })

        _, err := client.Ping().Result()
        if err != nil {
            log.Fatalf("Failed to ping Redis: %v", err)
        }

        return &redisStore{
            client: client,
        }
    }

Now, adding another _hit_ for a given page is much simpler. Just call an existing function on Redis client, and return any error if it happens to occur:

    func (r redisStore) Hit(url string) error {
        return r.client.HIncrBy("pagehits", url, 1).Err()
    }

Here comes a true test of our solution, which is a handler for listing all page-hit counter values. We are going to make just a single query to Redis (possibly saving us a lot of I/O time), then iterate over the map to convert string values to ints (after all, we want the end-user to see numbers for counters, right?) and return the map:

    func (r redisStore) GetStats() (Stats, error) {
        pagehits, err := r.client.HGetAll("pagehits").Result()
        if err != nil {
            return Stats{}, err
        }

        stats := make(map[string]int)
        for key, val := range pagehits {
            count, err := strconv.Atoi(val)
            if err != nil {
                count = 0
            }

            stats[key] = count
        }

        return stats, nil
    }

One last thing we need to do is test if this still works:

(redis) $ curl http://localhost:5000
Hello World #0
(redis) $ curl http://localhost:5000
Hello World #0
(redis) $ curl http://localhost:5000/stats
{"/":2,"/stats":1}
(redis) $ curl http://localhost:5000/stats
{"/":2,"/stats":2}
(redis) $ curl http://localhost:5000/stats
{"/":2,"/stats":3}
(redis) $ curl http://localhost:5000
Hello World #0
(redis) $ curl http://localhost:5000
Hello World #0
(redis) $ curl http://localhost:5000/stats
{"/":4,"/stats":4}

Of course it does, otherwise, I wouldn't publish this post (duh!). 

The whole source code of this example is available [on Github](https://github.com/mycodesmells/golang-examples/tree/master/redis).
