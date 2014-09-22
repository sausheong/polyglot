require 'securerandom'
require 'bundler'
Bundler.require

broker = "tcp://localhost:4321"
routeid = "GET/_/hello/ruby"
identity = SecureRandom.uuid

puts "#{routeid} - #{identity} responder ready."

ctx = ZMQ::Context.new
client = ctx.socket ZMQ::REQ
client.identity = identity
client.connect broker

client.send_string routeid
loop do
  request = String.new
  client.recv_string request
  response = [routeid, "200", "{\"Content-Type\": \"text/html\"}", "Hello World"]
  client.send_strings response
end