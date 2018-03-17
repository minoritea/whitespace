#!/usr/bin/env ruby
testfile_lines = open(ARGV[0]).readlines rescue begin
  $stderr.puts "No valid test file was given."
  exit 1
end

expect = testfile_lines.first&.chomp || ""
if expect.empty? or expect =~ /[ \t]/
  $stderr.puts "No valid expectation was given."
  exit 1
end
source = testfile_lines[1..-1]
  .join
  .gsub(/[ \t\n]/, "")
  .gsub(/S/, " ")
  .gsub(/T/, "\t")
  .gsub(/L/, "\n")

Dir.chdir File.expand_path("../..", __FILE__)

result = IO.popen(%w[go run main.go], 'r+') do |io|
  io.write source
  io.close_write
  io.read
end

unless result == expect
  $stderr.write """
The expected result is `#{expect}`,
but the acutal result is `#{result}`.
  """
  exit 1
end
