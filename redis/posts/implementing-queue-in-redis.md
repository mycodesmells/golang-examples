# Implementing queue in Redis

As you probably already know, Redis goes well beyond getting and setting values. One of the less often used features in the database is the ability to implement a pretty simple task queue and share it between a producer and a consumer. Let's take a look how you can do that in a few steps.

# List as a queue

The secret of creating a queue in Redis lies in the fact that the database in single-threaded, which means that it provides consistent read and write operations to its clients. Whenever two clients attempt to write to the same key, one has to wait for another to finish before starting its work. That means each operation is transactional, which allows us to create a queue under given key.

The other implementation details you need to know is _a list_. This is the easiest way to store queue data, as you can have one party (_a producer_) add stuff to the list, while the other side (_a consumer_) reads from it. In order to implement both _FIFO_ (_First-In-First-Out_) and _LIFO_ (_Last-In-First-Out_), Redis provides us with operations to put elements to each end of the list([RPUSH](https://redis.io/commands/rpush), [LPUSH](https://redis.io/commands/lpush)), as well as pop them both from the beginning ([LPOP](https://redis.io/commands/lpop)) and the end ([RPOP](https://redis.io/commands/rpop)).

In our example we'd like to have a _FIFO_ queue, so we'll push elements to the end, and pop them from the beginning.

# Producer (server)

Adding element to the queue list is straightforward:

    // cmd/qserver/main.go
    client := redis.NewClient(...)
    ...

    if err := client.RPush("queue-key", task).Err(); err != nil {
        log.Fatalf("Failed to put stuff into queue: %v", err)
    }
    log.Printf("'%v' task put into queue", task)

This simple code is more than enough to fill our queue (kept under `queue-key` key) with some _tasks_:

    (redis) $ go run cmd/qserver/main.go -task super-task-1
    2017/05/05 22:05:06 'super-task-1' task put into queue
    (redis) $ go run cmd/qserver/main.go -task super-task-2
    2017/05/05 22:05:08 'super-task-2' task put into queue
    (redis) $ go run cmd/qserver/main.go -task super-task-3
    2017/05/05 22:05:10 'super-task-3' task put into queu
    (redis) $ go run cmd/qserver/main.go -task super-task-4
    2017/05/05 22:04:30 'super-task-4' task put into queue

# Consumer (client)

Our first approach would be to make it as simple as possible and just read data from the list from `queue-key`:

    // cmd/qclient/main.go
    for {
        task, err := client.LPop("queue-key").Result()
        if err != nil {
            log.Fatalf("Failed to get task from queue: %v\n", err)
        }

        log.Printf("Working on '%s' task...\n", task)
    }

The problem here is, that whenever we try to `LPop(..)` from an empty (nil) list, it throws an error:

    (redis) $ go run cmd/qclient/main.go 
    2017/05/05 22:26:21 Working on 'super-task-1' task...
    ...
    2017/05/05 22:26:21 Working on 'super-task-4' task...
    2017/05/05 22:26:21 Failed to get task from queue: redis: nil
    exit status 1

If we don't want to exit the customer service every time we run out of elements in the queue (and we don't in most cases) then we need to find a better solution. Fortunately, we get another gift from Redis itself, called _blocking list pop_ ([BRPOP](https://redis.io/commands/blpop), alternatively [BRPOP](https://redis.io/commands/brpop) for popping from the other side). What it does is wait some time (we can wait forever if we set the time to zero) before pulling value from the list. The catch here, however, is that we can provide multiple lists at the same time, and the response contains both a newly added list element and the key of that list. This is important if we want to do something with the value because `Result()` returns now an array:

    // cmd/qclient/main.go
    ...
    task, err := client.BLPop("queue-key").Result()
    if err != nil {
        log.Fatalf("Failed to get task from queue: %v\n", err)
    }
    
    log.Printf("Working on '%s' task...\n", task[1])
    ...

If you looked carefully on the example above (the one with `LPop`), we read from the database constantly inside neverending `for` loop. If you prefer to cut Redis some slack, you can always wait between read attempts: 

    // cmd/qclient/main.go
    ...
    time.Sleep(5 * time.Second) // wait 5s in each `for` loop iteration
    task, err := client.BLPop(0, "queue-key").Result()
    ...
    log.Printf("Working on '%s' task...\n", task[1])
    ...

The end result looks as follows, from the server side:

    (redis) $ go run cmd/qserver/main.go -task super-task-1
    2017/05/05 22:35:05 'super-task-1' task put into queue
    (redis) $ go run cmd/qserver/main.go -task super-task-2
    2017/05/05 22:35:07 'super-task-2' task put into queue
    (redis) $ go run cmd/qserver/main.go -task super-task-3
    2017/05/05 22:35:08 'super-task-3' task put into queue
    (redis) $ go run cmd/qserver/main.go -task super-task-4
    2017/05/05 22:35:10 'super-task-4' task put into queue
    (redis) $ go run cmd/qserver/main.go -task super-task-5
    2017/05/05 22:35:12 'super-task-5' task put into queue

While from the client side:

    (redis) $ go run cmd/qclient/main.go 
    2017/05/05 22:35:21 Working on 'super-task-1' task...
    2017/05/05 22:35:26 Working on 'super-task-2' task...
    2017/05/05 22:35:31 Working on 'super-task-3' task...
    2017/05/05 22:35:36 Working on 'super-task-4' task...
    2017/05/05 22:35:41 Working on 'super-task-5' task...
    // still waiting for new tasks

The whole source code of this example is available [on Github](https://github.com/mycodesmells/golang-examples/tree/master/redis).
