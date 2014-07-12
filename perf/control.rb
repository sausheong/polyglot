require 'sinatra'
configure { set :server, :thin; set :port, 8080 }

get "/_/ruby/hello" do  
  sleep 0.5
  "<h1>Hello Ruby!</h1>"
end