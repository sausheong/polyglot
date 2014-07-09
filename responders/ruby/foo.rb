require "./polyglot"
require 'base64'

class Hello < Polyglot::Responder
  
  def initialize
    super
    @method, @path = "GET", "foo/bar"
  end
  
  def respond(json)
    pic = File.read('monalisa.jpg')
    [200, {"Content-Type" => "image/jpeg"}, Base64.encode64(pic)]
  end
end 

responder = Hello.new
responder.run


