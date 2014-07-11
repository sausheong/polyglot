require 'httperf'
require 'json'

File.open('results.json', 'w') do |file|
  results = {}
  (10..100).step(5) do |num_connections|
    perf = HTTPerf.new "port" => 8080, "uri" => "/_/ruby/hello", "num-conns" => num_connections
    perf.parse = true
    results[num_connections] =  perf.run    
  end
  file.write results.to_json
end