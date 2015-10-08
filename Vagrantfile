# -*- mode: ruby -*-
# vi: set ft=ruby :

$provision = <<SCRIPT
apt-get update -y
apt-get upgrade -y
apt-get install -y git build-essential \
  libncurses5-dev \
  stow \
  keychain \
  gnupg-agent \
  pass \
  openjdk-8-jre \
  postgresql-9.4
sudo -i -u vagrant /home/vagrant/src/github.com/calebdoxsey/devenv/bin/devenv
SCRIPT

Vagrant.configure(2) do |config|
  config.vm.box = "ubuntu/vivid64"
  config.vm.provision "shell", inline: $provision

  # config.vm.network "forwarded_port", guest: 80, host: 8080
  config.vm.network "private_network", ip: "192.168.33.10"
  # config.vm.network "public_network"

  config.vm.synced_folder "../../..", "/home/vagrant/src", type: "nfs"
  config.vm.synced_folder "../../../../storage/repositories", "/home/vagrant/repositories", type: "nfs"
  config.vm.synced_folder "../../../../storage/keys", "/home/vagrant/keys", type: "nfs"

  config.vm.provider "virtualbox" do |v|
    v.memory = 1024
    v.cpus = 4
  end
end
