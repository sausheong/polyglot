import pika
import json

route_id = "GET/_/py/hello"
connection = pika.BlockingConnection(pika.ConnectionParameters(host='localhost'))
channel = connection.channel()
channel.queue_declare(queue=route_id, durable=True, auto_delete=True)

print "[Responder ready]"

def callback(ch, method, props, body):
  if props.app_id == route_id:
    response = [200, {"Content-Type" : "text/html"}, "<h1>Hello Python!</h1>"]
    response_json = json.dumps(response)
    
    ch.basic_publish(exchange='',
                     routing_key=props.reply_to,
                     properties=pika.BasicProperties(correlation_id = props.correlation_id),
                     body=str(response_json))
    
    ch.basic_ack(delivery_tag = method.delivery_tag)


channel.basic_qos(prefetch_count=1)
channel.basic_consume(callback, queue=route_id)
try:
  channel.start_consuming()
except KeyboardInterrupt:
  print "[Responder shutdown]"