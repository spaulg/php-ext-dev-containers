def local_cache_path(basebox_name)
  File.expand_path(File.join('~/.vagrant.d/cache/apt/', basebox_name))
end

def total_cpus
  require 'etc'
  Etc.nprocessors
end

Vagrant.configure("2") do |config|
  config.vagrant.plugins = ["vagrant-hostmanager", "vagrant-mutagen-project"]

  # Detect correct platform architecture to set the
  # correct vagrant image
  config.vm.box = "bento/ubuntu-22.04-arm64" if `sysctl -n machdep.cpu.brand_string` =~ /M1/
  config.vm.box = "bento/ubuntu-22.04" if `sysctl -n machdep.cpu.brand_string` !~ /M1/

  # Configure the virtualbox provider
  config.vm.provider "virtualbox" do |p, override|
    p.cpus = total_cpus - 2
    p.memory = 4096
  end

  # Configure the parallels provider
  config.vm.provider "parallels" do |p, override|
    p.cpus = total_cpus - 2
    p.memory = 4096
  end

  # Use hostmanager to set the host name of the VM on both
  # the host and the guest automatically
  config.hostmanager.enabled = true
  config.hostmanager.manage_host = true
  config.hostmanager.manage_guest = true

  # Additional host aliases to assign with
  # hostmanager
  config.hostmanager.aliases = %w(node1.test)

  # Automatically provision mutagen to synchronise files
  # in to the VM
  config.mutagen.orchestrate = true
  config.mutagen.project_file = "mutagen.yml"

  config.vm.hostname = "php-ext-dev-containers"
  config.vm.network "private_network", ip: "192.168.56.30"

  # Synchronised folders
  # config.vm.synced_folder ".", "/home/vagrant/php-ext-dev-containers"
  config.vm.synced_folder ".", "/vagrant", disabled: true
  config.vm.synced_folder local_cache_path(config.vm.box), "/var/cache/apt/archives", create: true

  # Provisioner
  config.vm.provision "shell", path: "provisioner.sh"
end
