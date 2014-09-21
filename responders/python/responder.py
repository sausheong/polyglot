import zmq
import uuid

identity = str(uuid.uuid4())
routeid = "GET/_/hello/python"

context = zmq.Context(1)
responder = context.socket(zmq.REQ)
responder.setsockopt(zmq.IDENTITY, identity)
responder.connect("tcp://localhost:4321")

print "%s - %s responder ready" % (routeid, identity)

responder.send(routeid)

while True:
  request = responder.recv()
  if not request:
    break
  response = [routeid, "200", "{\"Content-Type\": \"text/html\"}", "Hello World"]
  responder.send_multipart(response)
