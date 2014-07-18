using RabbitMQ.Client;
using RabbitMQ.Client.Events;

var routeId = "GET/_/csharp/hello";
var factory = new ConnectionFactory() { HostName = "localhost", UserName = "guest", Password = "guest" };
using (var connection = factory.CreateConnection())
{
  using (var channel = connection.CreateModel())
  {
    channel.QueueDeclare(routeId, true, false, true, null);
    var consumer = new QueueingBasicConsumer(channel);
    channel.BasicConsume(routeId, false, consumer);

    Console.Write("[Responder ready]");

    // while (true) makes mono freak out, so I'm applying some hackery
    var i = 0;
    while (i < 1)
    {
      var ea = (BasicDeliverEventArgs)consumer.Queue.Dequeue();
      var props = ea.BasicProperties;

      Console.WriteLine(" [x] Received {0}", Encoding.UTF8.GetString(ea.Body));

      channel.BasicAck(ea.DeliveryTag, false);

      var json = @"[200, { ""Content-Type"": ""text/html"" }, ""<h1>Hello C#!</h1>""]";

      channel.BasicPublish("", props.ReplyTo, null, Encoding.UTF8.GetBytes(json));
      Console.WriteLine(" [x] Sent {0}", json);
    }
  }
}
