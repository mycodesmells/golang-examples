# Watching Metrics With Prometheus

Have you ever noticed your web application working a bit slower than you expected? Or maybe it looks OK, but some client complains that it's slow for them? If you don't monitor your app, it's virtually impossible to verify if this is true. That's why you should use Prometheus to find out.

# What is Prometheus

[Prometheus](https://prometheus.io/) is an application that stores and aggregates various _metrics_ registered in your application in order to visualize everything that may help you debug, detect and predict problems with your app.

Using Prometheus is really simple, as we just need to set up its instance, then register some metrics, send some data and browse them in the web interface. This might sound difficult, but in a few moments, you'll see how little effort is necessary to have your application's statistics. 

# Setting up the environment

Since we need to create an instance of an external application, my first thought was to look for a Docker image. Thanks to its awesomeness, we don't need to install anything and get lost in its internals, but just spin up something that is ready. Great, isn't it?

I decided to start with creating a `docker-compose.yml` configuration and create two services: one for metrics and another for our own application:

    # docker-compose.yml
    version: '2'
services:
  web:
    build: docker/web
    ports: 
      - "3000:3000"
  metrics:
    image: prom/prometheus
    ports: 
      - "9090:9090"
    volumes:
      - ./docker/metrics/prometheus.yml:/etc/prometheus/prometheus.yml

As you can see, we defined one volume for Prometheus image, because we'd like to inject a configuration file. Following official docs, we start by adding a Prometheus itself to the config, so that we are able to monitor what is happening there as well. Then, we add another _job_ pointing to our web application which will be exposed on `web` service, on port `3000`. 

    # prometheus.yml
    scrape_configs:
        - job_name: 'self'
            scrape_interval: 5s
            static_configs:
            - targets: ['localhost:9090']
        - job_name: 'web'
            scrape_interval: 5s
            static_configs:
            - targets: ['web:3000']
                labels:
                group: 'production'

# Sample web application

The most common thing that we might want to monitor in our application is how much time does it take to finish handling a single request. To make the example a bit more interesting, we'll introduce random sleeps in our handler to make those times differentiate at least a bit:

    func main() {
        r := rand.New(rand.NewSource(time.Now().UnixNano()))

        handler := http.NewServeMux()
        handler.HandleFunc("/", func(rw http.ResponseWriter, req *http.Request) {
            delay := int32(r.Float64() * 10000)
            time.Sleep(time.Millisecond * time.Duration(delay))

            fmt.Printf("Responed with 'Hello' in %dms\n", delay)
            rw.Write([]byte("Hello!"))
        })

        log.Fatal(http.ListenAndServe(":3000", withMetrics))
    }

In order to measure time and see it in the Prometheus dashboard, we first need to register a new metric. There are many metrics types available, but in the case of response time, we need to use a thing called _summary_, since we want to get time values for percentiles (0.5, 0.9 and 0.99 by default). Our metric will be named `http_response_time_seconds` (with `http` serving as a namespace prefix): 

    // metrics/middleware.go
    var (
        httpRequestsResponseTime prometheus.Summary
    )

    func init() {
        httpRequestsResponseTime = prometheus.NewSummary(prometheus.SummaryOpts{
            Namespace: "http",
            Name:      "response_time_seconds",
            Help:      "Request response times",
        })

        prometheus.MustRegister(httpRequestsResponseTime)
    }

The middleware itself is fairly straightforward: first record a time before the requests processing started, then do the actual work, and finally get the current time again, calculate the difference and send that value to the metrics server. To do that we use `Observe(..)` function on our metric object.

    func Middleware(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            start := time.Now()

            next.ServeHTTP(w, r)

            httpRequestsResponseTime.Observe(float64(time.Since(start).Seconds()))
        })
    }

Then add necessary middleware to send that data to Prometheus:

    // main.go
    handler := http.NewServeMux()
    handler.Handle("/metrics", promhttp.Handler())
    ...
    withMetrics := metrics.Middleware(handler)

Everything seems to be in place, let's run some load tests to see how does it work.

# Watching it work

It would be a bit tedious to send requests manually, especially since our dummy handler may take up to 10 seconds to finish its work. Thankfully, there are tools that can make our life a bit easier, like [hey](https://github.com/rakyll/hey) (previously known as boom), which hit a chosen endpoint many times and with some of the requests sent concurrently (both values are configurable). For example we'd like to send 10k requests, with 200 in parallel:

    hey -n 10000 -c 200 http://localhost:3000

Now we can log into the dashboard and see how the graph changes:

<img src="https://raw.githubusercontent.com/mycodesmells/golang-examples/master/prometheus/posts/img/response-time-graph.png"/>

As you can see, half of the requests take 4.32 seconds, but one percent are slower than 9.82! Might be useful to know one day.

The whole source code of this example is available [on Github](https://github.com/mycodesmells/golang-examples/tree/master/prometheus).
