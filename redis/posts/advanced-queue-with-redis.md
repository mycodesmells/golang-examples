# Advanced Queue in Redis

First, check out [the post](http://mycodesmells.com/post/implementing-queue-in-redis) on the queue using Redis lists.

We've already seen how you can implement a simple queue in Redis using a list key and smart pushing and popping the elements. But is there any better, cleaner and more suitable way to do such a thing? Does Redis provide any other possible approach? You bet! Enter Pub/Sub!

### What is Pub/Sub?

Apart from an artificial approach, which was the one with a list shared between message producer and consumer, Redis provides a solution that implements a well-known paradigm called [publish-subscribe pattern](https://en.wikipedia.org/wiki/Publish%E2%80%93subscribe_pattern). It major principle is that the producer does not know anything about the consumers, it doesn't even know if there are any. It builds some messages that are sent to some broker or event bus, from where they are available to all the consumers. 

The major downside of this solution is, that the messages are sent to the database, but if there is no client subscribed they are lost. This is an important thing to remember because it might be a bit surprising for someone used to eg. Kafka, which keeps all the messages no matter what.

### Implementing Publisher

As with our previous queue approach, this side of the equation is simpler, as message producer just needs to send some data to Redis and never worry about it anymore. Instead of pushing it to the list, this time we use [PUBLISH](https://redis.io/commands/publish) command, which sends some data to given key in the database:

    if err := client.Publish("pubsub-key", task).Err(); err != nil {
        log.Fatalf("Failed to put stuff into queue: %v", err)
    }
    log.Printf("'%v' task put into queue", task)

It just cannot get any simpler than that, can it?

### Implementing Subscriber

Now here comes a slightly more complicated part. In order to receive information about data added do our `pubsub-key` queue, we first need to subscribe to it using, surprise surprise, [SUBSCRIBE](https://redis.io/commands/subscribe) command (we can always take a step back with [UNSUBSCRIBE](https://redis.io/commands/unsubscribe)). We can subscribe to multiple keys at the same time, which is important later on when reading incoming data.

    pubsub, err := client.Subscribe("pubsub-key")
    if err != nil {
        log.Fatalf("Failed to get task from queue: %v\n", err)
    }

Once we are subscribed (and received `pubsub` object of type `*redis.PubSub`), we can proceed and handle whatever is thrown at us. We would like to check every five seconds if something was published, then process the data, and repeat the process:

    for {
        msgi, err := pubsub.ReceiveTimeout(5 * time.Second)
        if err != nil {
            break
        }
        ...

There are three things that can be returned in `msgi` variable: subscription information (we are notified that we are in fact subscribed to the channel), an actual message or an error. In case we get the first one, we would like to print the information to the user (as there is not much we can _do_ with this information). An error does not mean that there is no data to be processed, it generally means that there is something seriously wrong with our database, so in our one-task program we should panic. Finally, we arrive at the place we wanted to be - do something about the message.

As you can recall, previously we were working on strings, as we needed to decode database key and value from an incoming slice of strings. PubSub abstraction gives us more friendly API because now we have `msg.Channel` and `msg.Payload` work with. As our example application is extremely dumb, we print the value and are happy about it:

    for {
        ...
        switch msg := msgi.(type) {
        case *redis.Subscription:
            fmt.Println("subscribed to", msg.Channel)
        case *redis.Message:
            fmt.Println("received", msg.Payload, "from", msg.Channel)
        default:
            panic(fmt.Errorf("unknown message: %#v", msgi))
        }
    }

### Pub/Sub in action

First, we start the producer without any subscribers to see the first message getting lost. Then once the client is started, all data is delivered properly.

    // running publisher
    (redis) $ go run cmd/pserver/main.go -task super-task-1
    2017/05/07 21:24:54 'super-task-1' task put into queue
    (redis) $ go run cmd/pserver/main.go -task super-task-2
    2017/05/07 21:25:01 'super-task-2' task put into queue
    (redis) $ go run cmd/pserver/main.go -task super-task-3
    2017/05/07 21:53:24 'super-task-3' task put into queue
    (redis) $ go run cmd/pserver/main.go -task super-task-4
    2017/05/07 21:53:28 'super-task-4' task put into queue

    // running subscriber
    (redis) $ go run cmd/pclient/main.go 
    subscribed to pubsub-key
    // note that super-task-1 was never delivered
    received super-task-2 from pubsub-key
    received super-task-3 from pubsub-key
    received super-task-4 from pubsub-key

The whole source code of this example is available [on Github](https://github.com/mycodesmells/golang-examples/tree/master/redis).
