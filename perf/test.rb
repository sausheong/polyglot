require 'httperf'
require 'json'

STDOUT.sync = true
results = {}
(10..200).step(10) do |num_connections|
  perf = HTTPerf.new "port" => 8080, "uri" => "/_/ruby/hello", "num-conns" => num_connections, "rate" => num_connections
  perf.parse = true  
  results[num_connections] =  perf.run    
  puts "#{num_connections} - #{results[num_connections][:total_test_duration]}s"
end

hash = {}
results.each do |num_conn, res|
  hash[num_conn.to_s.rjust(3, "0")] = res[:connection_time_avg]
end

File.open('results.json', 'w') do |file|  
  file.write hash.to_json
end