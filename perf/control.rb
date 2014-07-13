require 'sinatra'
configure { set :port, 8080 }

get "/perf" do  
  sleep 0.5
  "<h1>Hello Perf</h1>"
end