require "./polyglot"

class Hello < Polyglot::Responder
  
  def initialize
    super
    @method, @path = "GET", "ruby/hello"
  end
  
  def respond(json)
    html "<h1>Hello Ruby!</h1>"
  end
end 

responder = Hello.new
responder.run


