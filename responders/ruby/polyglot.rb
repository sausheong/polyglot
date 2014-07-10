require 'bunny'
require 'json'
require 'haml'
require 'base64'

module Polyglot
  
  class Responder

    def initialize(id)
      @route_id = id
    end

    def run
      # A route ID uniquely identifies a route that this responder will respond to
      conn = Bunny.new
      conn.start
      ch = conn.create_channel
      
      # Set up the queue for the acceptor to add messages to
      # If the acceptor cannot find a queue with this route ID it will return a 404
      # Set the queue to auto delete ie if there are no more messages or consumers  
      # on the queue, it will remove itself
      q  = ch.queue(@route_id, durable: true, auto_delete: true)
      exch = ch.default_exchange
      ch.prefetch(1)
      puts "[Responder ready]."
      
      loop do
        begin
          q.subscribe(ack: true, block: true) do |delivery_info, properties, body|
            # Only respond to this route ID
            if @route_id == properties[:app_id] then                          
              response = self.respond(body)        
              exch.publish(response.to_json, routing_key: properties.reply_to, correlation_id: properties.correlation_id)
              ch.ack(delivery_info.delivery_tag)
            end
          end
        rescue Interrupt => int
          puts
          puts "[Responder shutdown]."          
          exit(0)
        end
      end      
      conn.close
    end

    
    # Responders must override this method; by default it will return a 200 OK with 
    # no message
    def respond(json)
      [200, {}, ""]
    end
    
    # Convenience method to return HTML
    def html(body)
      [200, {"Content-Type" => "text/html"}, body]
    end
    
    # Convenience method to redirect to the given url
    def redirect(url)
      [302, {"Location" => url}, ""]
    end
    
  end
  
end