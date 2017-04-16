# Using Redis as sessions store

Since its initial release in 2009, Redis has become one of the most popular NoSQL solutions and almost a synonym for a key-value database. Thanks to being around for some time, there are a few use cases that it is more than suitable for. One of them is a persistent store for users' sessions in a web application. While it has already been implemented and available using various open source libraries, let's see if we can do it by ourselves.

# Redis basics

As you may already know, Redis is a single-threaded, key-value store that is very easy to use, therefore is very popular. It has just a few value data types, with _string_ as a most often-used one. Even when storing numeric data, a string value is created underneath. Then there are more complex data types like lists, sets,  hashes (map-like structures) and more.

When connecting to Redis you need to define which _database_ you want to select, but the choice here is quite limited. There are 16 databases available, represented by their consecutive number (zero-based). By default, all connections are heading to 0, but it can obviously be modified. Those databases can be considered as something similar to _schemas_ in relational databases - they provide easy data separation.

Basic usage of Redis is as easy as working with a map (dictionary) - you can set and get values, which will be enough for our first example use case.

# Keeping sessions in memory

Imagine having a web application, which allows users to do multiple things during a single session and we'd like to monitor what is going on: which pages are being visited and how many times something happens. In our case, we'll limit ourselves to a simple counter - whenever a user connects to the site, we'll update it in Redis and return a value (also we'll log it to see how it changes on the server side of things).

To do this we'll start with defining an interface for sessions store and a struct for a session:

    // sessions/sessions.go
    type Session struct {
        VisitCount int `json:"visitCount"`
    }

    type Store interface {
        Get(string) (Session, error)
        Set(string, Session) error
    }

Then let's create an in-memory solution to see how it compares to Redis later on:

    // sessions/memory.go
    type memoryStore struct {
        sessions map[string]Session
    }

    func NewMemoryStore() Store {
        return &memoryStore{
            sessions: make(map[string]Session),
        }
    }

    func (m memoryStore) Get(id string) (Session, error) {
        session, ok := m.sessions[id]
        if !ok {
            return Session{}, errors.New("session not found")
        }

        return session, nil
    }

    func (m *memoryStore) Set(id string, session Session) error {
        m.sessions[id] = session
        return nil
    }

The most of the stuff happens in a middleware, which binds a request with a session using a cookie value. Then, the counter value is being incremented (or set it doesn't exist) and saved back to the store:

    //sessions/middleware.go
    func Middleware(store Store) echo.MiddlewareFunc {
        return func(hf echo.HandlerFunc) echo.HandlerFunc {
            return func(ectx echo.Context) error {
                cookie, err := ectx.Cookie("sessionID")
                if err != nil {
                    sessionID := uuid.NewV4().String()
                    ectx.SetCookie(&http.Cookie{
                        Name:  "sessionID",
                        Value: sessionID,
                    })
                    ectx.Set("sessionID", sessionID)
                    store.Set(sessionID, Session{})
                    return hf(ectx)
                }

                sessionID := cookie.Value
                ectx.Set("sessionID", sessionID)
                return hf(ectx)
            }
        }
    }

To make logging easier, we'll also use a custom middleware that will create a logger instance. The logger, using `logrus` Fields, will print a session ID alongside whatever we log at given time:

    //logger/logger.go
    ...
    func FromContext(ectx echo.Context) *logrus.Entry {
        sessionID := ectx.Get("sessionID").(string)
        return logrus.WithField("sessionID", sessionID)
    }

Finally, we can create a handler in our simple server (built using [echo](https://echo.labstack.com) :

//cmd/server/main.go
...
func main() {
    sessionsStore := sessions.NewMemoryStore()

    e := echo.New()
    e.Use(sessions.Middleware(sessionsStore))

    e.GET("/", func(ectx echo.Context) error {
        log := logger.FromContext(ectx)
        log.Info("Hello world")

        sessionID := ectx.Get("sessionID").(string)
        s, err := sessionsStore.Get(sessionID)
        if err != nil {
            log.Errorf("err: %v", err)
        }

        log.Infof("Visits: %d", s.VisitCount)
        response := fmt.Sprintf("Hello World #%d\n", s.VisitCount)

        s.VisitCount = s.VisitCount + 1
        err = sessionsStore.Set(sessionID, s)
        if err != nil {
            log.Errorf("err: %v", err)
        }

        return ectx.String(http.StatusOK, response)
    })

    e.Start(":5000")
}

Now when we run a server and a client, the value updates beautifully:

    // server
    $ go run cmd/server/main.go 
    ⇛ http server started on [::]:5000
    INFO[0031] Hello world                                   sessionID=ebe956d6-0f60-4fde-831a-691dbba736ca
    INFO[0031] Visits: 0                                     sessionID=ebe956d6-0f60-4fde-831a-691dbba736ca
    INFO[0033] Hello world                                   sessionID=ebe956d6-0f60-4fde-831a-691dbba736ca
    INFO[0033] Visits: 1                                     sessionID=ebe956d6-0f60-4fde-831a-691dbba736ca
    INFO[0033] Hello world                                   sessionID=ebe956d6-0f60-4fde-831a-691dbba736ca
    INFO[0033] Visits: 2                                     sessionID=ebe956d6-0f60-4fde-831a-691dbba736ca

    // client
    (redis) $ curl http://localhost:5000 --cookie-jar /tmp/cookie-jar -b /tmp/cookie-jar
    Hello World #0
    (redis) $ curl http://localhost:5000 --cookie-jar /tmp/cookie-jar -b /tmp/cookie-jar
    Hello World #1
    (redis) $ curl http://localhost:5000 --cookie-jar /tmp/cookie-jar -b /tmp/cookie-jar
    Hello World #2

The prolems start when we want to restart a server and keep the counter alive:

    // server
    (redis) $ go run cmd/server/main.go 
    ⇛ http server started on [::]:5000
    INFO[0001] Hello world                                   sessionID=ebe956d6-0f60-4fde-831a-691dbba736ca
    INFO[0001] Visits: 0                                     sessionID=ebe956d6-0f60-4fde-831a-691dbba736ca
    INFO[0002] Hello world                                   sessionID=ebe956d6-0f60-4fde-831a-691dbba736ca
    INFO[0002] Visits: 1                                     sessionID=ebe956d6-0f60-4fde-831a-691dbba736ca
    ^Csignal: interrupt
    (redis) $ go run cmd/server/main.go 
    ⇛ http server started on [::]:5000
    INFO[0001] Hello world                                   sessionID=ebe956d6-0f60-4fde-831a-691dbba736ca
    INFO[0001] Visits: 0                                     sessionID=ebe956d6-0f60-4fde-831a-691dbba736ca

    // client
    redis) $ curl http://localhost:5000 --cookie-jar /tmp/cookie-jar -b /tmp/cookie-jar
    Hello World #0
    (redis) $ curl http://localhost:5000 --cookie-jar /tmp/cookie-jar -b /tmp/cookie-jar
    Hello World #1
    // now server restarts
    (redis) $ curl http://localhost:5000 --cookie-jar /tmp/cookie-jar -b /tmp/cookie-jar
    Hello World #0

Let's see how Redis helps with this one.

# Using Redis

Since we already have an interface, all we need to do is implement a store using Redis Go client:

    // sessions/redis.go
    type redisStore struct {
        client *redis.Client
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

Setting a value to the database will require us to Marshall a `Session` to bytes and put that into Redis:

    func (r redisStore) Set(id string, session Session) error {
        bs, err := json.Marshal(session)
        if err != nil {
            return errors.Wrap(err, "failed to save session to redis")
        }

        if err := r.client.Set(id, bs, 0).Err(); err != nil {
            return errors.Wrap(err, "failed to save session to redis")
        }

        return nil
    }

Three arguments that are needed to put something into the key-value store are key, value (obviously!) and expiration, which means that data in Redis can expire (excellent for caching). In our case we want to keep the data forever and ever, so we set the value to zero. The result of the command provides some utility functions to see how the database responded to our query. In this case, we get a `StatusCmd` response that returns _OK_ as value if the command was completed successfully, and an error otherwise.

Getting value from the store is analogous: we get raw bytes data, unmarshal it into a struct and return a session:

    func (r redisStore) Get(id string) (Session, error) {
        var session Session

        bs, err := r.client.Get(id).Bytes()
        if err != nil {
            return session, errors.Wrap(err, "failed to get session from redis")
        }

        if err := json.Unmarshal(bs, &session); err != nil {
            return session, errors.Wrap(err, "failed to unmarshall session data")
        }

        return session, nil
    }

Now when we change an implementation of `sessions.Store` being used, the data is persisted even after restarting our server:

    //main.go
    func main() {
        sessionsStore = sessions.NewRedisStore()
        ...
    }

    // server
    (redis) $ go run cmd/server/main.go 
    ⇛ http server started on [::]:5000
    INFO[0003] Hello world                                   sessionID=ebe956d6-0f60-4fde-831a-691dbba736ca
    INFO[0003] Visits: 0                                     sessionID=ebe956d6-0f60-4fde-831a-691dbba736ca
    INFO[0003] Hello world                                   sessionID=ebe956d6-0f60-4fde-831a-691dbba736ca
    INFO[0003] Visits: 1                                     sessionID=ebe956d6-0f60-4fde-831a-691dbba736ca
    ^Csignal: interrupt
    (redis) $ go run cmd/server/main.go 
    ⇛ http server started on [::]:5000
    INFO[0001] Hello world                                   sessionID=ebe956d6-0f60-4fde-831a-691dbba736ca
    INFO[0001] Visits: 2                                     sessionID=ebe956d6-0f60-4fde-831a-691dbba736ca

    // client
    (redis) $ curl http://localhost:5000 --cookie-jar /tmp/cookie-jar -b /tmp/cookie-jar
    Hello World #0
    (redis) $ curl http://localhost:5000 --cookie-jar /tmp/cookie-jar -b /tmp/cookie-jar
    Hello World #1
    // now server restarts
    (redis) $ curl http://localhost:5000 --cookie-jar /tmp/cookie-jar -b /tmp/cookie-jar
    Hello World #2

As you can see, we created a persistent sessions store with a little help from Redis and a minimal effort from our side. A complete source code of this solution is available [on Github](https://github.com/mycodesmells/golang-examples/tree/master/redis).
