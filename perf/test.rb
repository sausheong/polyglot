require 'httperf'
require 'json'

STDOUT.sync = true
results = {}
(10..1000).step(10) do |num_connections|
  perf = HTTPerf.new "port" => 8080, "uri" => "/perf", "num-conns" => num_connections, "rate" => num_connections
  perf.parse = true  
  results[num_connections] =  perf.run    
  p results[num_connections]
  puts "#{num_connections} - #{results[num_connections][:total_test_duration]}s"
end

hash = {}
results.each do |num_conn, res|
  hash[num_conn] = res
end

File.open("results_#{Time.now.strftime("%y%m%d%H%M%S")}.json", 'w') do |file|  
  file.write hash.to_json
end