require "./polyglot"

class Hello < Polyglot::Responder  
  def respond(json)
    work
    html "<h1>Hello Ruby!</h1>"
  end
  
  def work
    sleep 0.5
  end
end 

responder = Hello.new("GET/_/ruby/hello")
responder.run


