require 'bundler'
Bundler.require

module Polyglot
  
  class Responder

    def run
      conn = Bunny.new
      conn.start
      ch = conn.create_channel
      q  = ch.queue("polyglot", durable: true)
      exch = ch.default_exchange
      ch.prefetch(1)
      puts "[Responder ready]."

      while true
        begin
          q.subscribe(manual_ack: true, block: true) do |delivery_info, properties, body|
            response = self.respond(body)            
            exch.publish(response.to_json, routing_key: properties.reply_to, correlation_id: properties.correlation_id)
            ch.ack(delivery_info.delivery_tag)
          end
        rescue Interrupt => int
          p int
          exit(0)
        end
      end      
      conn.close
    end
    
    def respond(json)
      # parse json and return a value
      [200]
    end
    
    def html(body)
      [200, {"Content-Type" => "text/html"}, body]
    end
    
    def redirect(url)
      [302, {"Location" => url}, ""]
    end
    
  end
  
end