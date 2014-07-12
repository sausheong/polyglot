require 'gruff'
require 'json'

# Sinatra + Puma with min 50 threads, max 50 threads
c1 = JSON.parse(File.read('puma.json')) 
# Polyglot with 50 Ruby responders in the same server, and 50 Ruby responders in another server
c2 = JSON.parse(File.read('polyglot.json'))

chart = Gruff::Line.new(4096)
chart.theme_pastel
chart.data "Sinatra Puma", c2.values.map{|item| item.to_i}, "#00EEFF" # blue
chart.data "Polyglot", c3.values.map{|item| item.to_i}, "#AAEE00" # green
chart.write "results.png"