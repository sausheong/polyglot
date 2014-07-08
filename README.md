# Polyglot

Programmers use web frameworks to simplify the effort to develop web applications. Web frameworks reduce the overhead of work that needs to be done in most web applications, like session management, view templates and so on. This reduces time and effort needed and often improve the stability and create a consistent and maintainable system. Of course, as with any framework, web application frameworks are written with a single, major programming language, like Ruby with Ruby on Rails, Python with Django, Java with Spring or Javascript with Angular.js or Ember.js. 

So what good is a web framework that increases the complexity and effort of developing web applications? The answer is trade-offs. As a programmer you trade-off complexity and effort for something you think is more important for the web application you're creating.

In Polyglot, we are trading complexity and effort for:

1. Performance scalability -- Polyglot responders are distributed and independent processes that can reside anywhere on a connected network
2. Chained processing -- Polyglot responders can be chained, each doing an individual piece of processing, encouraging reusability of code
3. Multi-lingual development -- Polyglot responders can be developed in multiple programming languages, **at the same time**


Polyglot is not for all types of web applications. You should only use Polyglot for web applications that need to be scale in a highly performant way and/or need to be incrementally developed in multiple programming languages. For example, if your web application never needs to scale beyond a single server, you're probably better off using some other single language framework. And if once you create your web application and you or anyone else never need to add new features, Polyglot is probably not for you either.

The first and second are understandable, but the third is quite strange, why would you want to develop a web application in multiple programming languages? As it turns out, there are good reasons, very often for practical purposes:

1. Web applications you write are systems and they change over time and can be written or maintained by different groups of people. If you're not restricted to a particular platform or language, then the chances of getting an incrementally better piece of software is higher. 
2. Also, by forcing the deliberate use of different programming languages, you are forced to separate the layers and make each component more independent and robust, being able to switch out the poor-performing responders and replacing them with higher-performing ones
3. Different responders can have different criteria for performance, ease-of-development, ease-of-maintenance or quick turnaround in development. With a single programming language you are often forced to accept a compromise. With multiple programming languages, you can choose the platform and language as what you need for that particular responder
4. Different responders can be written for specific performance gains or maintenability.


## Architecture

Polyglot consists of 3 main components:

1. **Acceptor** - the acceptor is a HTTP interface that takes in a HTTP request and provides a HTTP response. The acceptor takes in a HTTP request converts it into a generic message and drops it into the message queue. Then depending on what is asked, it will return the appropriate HTTP response. The implementation of the acceptor is in Go.
2. **Messsage queue** - a queue that receives the messages that represent the HTTP request. the acceptor accepts HTTP requests and converts the requests into messages that goes into the message queue. The messages then gets picked up by the next component, the responder. The implementation of the message queue is a RabbitMQ server.
3. **Responder** - the responder is a standalone process written in any language that picks up messages from the message queue and responds to it accordingly. In most web frameworks this is usually called the controller. However unlike most web framework controllers, the responders are actual independent processes that can potentially sit in any connected remote server. Responders contain the "business logic" of your code. 

The normal flow goes like this:

1. Client sents a HTTP request to the acceptor
2. The acceptor converts the request into JSON and adds the JSON message into the message queue, and waits for a response
3. A responder detects the message and starts processing it
4. The responder completes processing the message and adds a response message back to the message queue, with the correlation ID set to the same ID that was sent as part of the request message. The response message is an array that has 3 elements -- the HTTP status code, a map of headers, and a body.
5. The acceptor detects a message with the same correlation ID on the queue and picks it up
6. The acceptor uses the HTTP status code, response headers and the body, creates a HTTP response and sends it to the client

The chained flow goes like this:

1. Client sents a HTTP request to the acceptor
2. The acceptor converts the request into JSON and adds the JSON message into the message queue, and waits for a response
3. A responder detects the message and starts processing it
4. Once the responder completes processing, it will send a create another message on the queue for another responder to process, then waits for a response
5. This results in a chain of responders
6. Once the final responder completes processing, the results are gathered and rolled back to the first responder
7. The first responder adds a response message back to the message queue, with the correlation ID set to the same ID that was sent as part of the request message
8. The acceptor detects a message with the same correlation ID on the queue and picks it up
9. The acceptor uses the HTTP status code, response headers and the body, creates a HTTP response and sends it to the client

