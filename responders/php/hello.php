<?php

require_once __DIR__ . '/vendor/autoload.php';
use PhpAmqpLib\Connection\AMQPConnection;
use PhpAmqpLib\Message\AMQPMessage;

$route_id = 'GET/_/php/hello';
$connection = new AMQPConnection('localhost', 5672, 'guest', 'guest');
$channel = $connection->channel();

$channel->queue_declare($route_id);

echo "[Responder ready]\n";

$callback = function($req){
  global $route_id;
  if ($req->get('app_id') == $route_id) {
    $payload = [
      "200",
      ['Content-Type' => 'text/html'],
      "<h1>Hello World!</h1>"
    ];
    $json = json_encode($payload);

    $msg = new AMQPMessage(
        $json,
        array('correlation_id' => $req->get('correlation_id'), 'content_type' => 'text/html')
        );

    $req->delivery_info['channel']->basic_publish($msg, '', $req->get('reply_to'));    
    $req->delivery_info['channel']->basic_ack($req->delivery_info['delivery_tag']);    
  }
};

$channel->basic_qos(null, 1, null);
$channel->basic_consume($route_id, '', false, false, false, false, $callback);

while(count($channel->callbacks)) {
    $channel->wait();
}

$channel->close();
$connection->close();
echo "[Responder shutdown]", "\n";
?>