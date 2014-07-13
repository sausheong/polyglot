require 'gruff'
require 'json'

NUM_POINTS = {100 => 1, 50 => 20, 25 => 40, 20 => 25, 10 => 100}

def process(chart, metric)
  d = {}
  chart.each do |k, v|
    d[k.to_i] = v[metric].to_f
  end
  d.delete_if do |k, v|
    k % NUM_POINTS[10] != 0
  end
  
  return d
end


# Sinatra + Puma with min 10 threads, max 100 threads
c1 = JSON.parse(File.read('results_puma.json')) 

# Go (Go spins up as many goroutines as necessary to handle the workload)
c2 = JSON.parse(File.read('results_go.json')) 

# Polyglot with 50 Ruby responders in the same server, and 50 Ruby responders in another server
c3 = JSON.parse(File.read('results_poly.json'))


metrics = ["connection_rate_per_sec", "connection_time_avg", "total_test_duration", "reply_status_2xx"]

metrics.each do |metric|
  d1 = process(c1, metric)
  d2 = process(c2, metric)
  d3 = process(c3, metric)

  labels = {}
  d1.keys.size.times do |i|
    labels[i] = d1.keys[i].to_s
  end

  chart = Gruff::Line.new(1024)
  chart.theme_pastel
  chart.labels = labels
  chart.title = metric
  chart.marker_font_size = 10
  chart.data "Puma", d1.values, "#00EEFF" # blue
  chart.data "Go", d2.values, "#AAEE00" # green
  chart.data "Polyglot", d3.values, "#009900" # dark green
  chart.write "#{metric}.png"
end