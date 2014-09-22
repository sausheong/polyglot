var zmq = require('zmq')
  , sock = zmq.socket('req');
var uuid = require('node-uuid');

routeid = "GET/_/hello";
identity = uuid.v4();
sock.identity = identity;
sock.connect('tcp://localhost:4321');
sock.send(routeid);

console.log('%s - %s responder ready', routeid, identity);

sock.on('message', function(msg){  
  sock.send(routeid, zmq.ZMQ_SNDMORE);
  sock.send("200", zmq.ZMQ_SNDMORE);
  sock.send("{\"Content-Type\": \"text/html\"}", zmq.ZMQ_SNDMORE);
  sock.send(msg);  
});