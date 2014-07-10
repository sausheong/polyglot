require "./polyglot"

class Hello < Polyglot::Responder

  def respond(json)
    puts json
    haml = Haml::Engine.new(File.read("hello.haml"))
    content = haml.render(Object.new, show_me: "Hello, world!")
    html content    
  end
end 

responder = Hello.new("GET/_/ruby/haml")
responder.run


