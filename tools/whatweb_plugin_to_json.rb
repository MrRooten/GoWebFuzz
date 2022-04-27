require 'json'
$target = ""
$mode = ""
$cur_file = ""
def randstr
  rand(36**8).to_s(36)
end
class Plugin
  class << self
    attr_reader :registered_plugins, :attributes
    private :new
  end

  @registered_plugins = {}
  @attributes = %i(
    aggressive
    authors
    description
    dorks
    matches
    name
    passive
    version
    website
  )

  @attributes.each do |symbol|
    define_method(symbol) do |*value, &block|
      name = "@#{symbol}"
      if block
        instance_variable_set(name, block)
      elsif !value.empty?
        instance_variable_set(name, *value)
      else
        instance_variable_get(name)
      end
    end
  end

  def initialize
    @matches = []
    @dorks = []
    @passive = nil
    @aggressive = nil
    @variables = {}
    @website = nil
  end

  def self.define(&block)
    # TODO: plugins should isolated
    p = new
    p.instance_eval(&block)
    p.startup
    # TODO: make sure required attributes are set

    # Freeze the plugin attributes so they cannot be self-modified by a plugin
    Plugin.attributes.each { |symbol| p.instance_variable_get("@#{symbol}").freeze }

      tar = {}
      tar[:name] = p.name
      tar[:matches] = p.matches
      tar[:dorks] = p.dorks
      tar[:website] = p.website
      s = JSON.pretty_generate tar
      if $mode == "dir"
        target_file = $target+$cur_file.chop.chop.chop+'.json'
        File.write(target_file,s)
      end
  end

  def self.shutdown_all
    Plugin.registered_plugins.each { |_, plugin| plugin.shutdown }
  end

  def version_detection?
    return false unless @matches
    !@matches.map { |m| m[:version] }.compact.empty?
  end

  # individual plugins can override this
  def startup; end

  # individual plugins can override this
  def shutdown; end

end
help = "./whatweb_plugin_to_json.rb [dir|file] $file $target_file"
if ARGV[0] == "dir"
  if ARGV.length() < 3
    puts help
  end
  source_dir = ARGV[1]
  target_dir = ARGV[2]
  $target = target_dir
  $mode = "dir"
  Dir.foreach(source_dir) do |item|
    next if item=='.' or item=='..' or !item.end_with?(".rb")
    $cur_file = item
    require source_dir + item
  end
elsif ARGV[0] == "file"

end

