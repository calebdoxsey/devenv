package main

import (
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/badgerodon/stack/archive"
	"github.com/badgerodon/stack/storage"
)

var (
	home   string
	gopath string
	goroot string
)

func init() {
	home = os.Getenv("HOME")
	gopath = home
	goroot = filepath.Join(home, "dist", "go")
}

func download(src, name, dst string) error {
	url, err := storage.ParseLocation(src)
	if err != nil {
		return err
	}
	rc, err := storage.Get(url)
	if err != nil {
		return err
	}
	defer rc.Close()

	err = archive.ExtractReader(dst, rc, name)
	if err != nil {
		return err
	}

	return nil
}

func installGo() error {
	cmd := exec.Command(filepath.Join(goroot, "bin", "go"), "version")
	cmd.Env = []string{"GOROOT=" + goroot}
	goVersion, _ := cmd.CombinedOutput()
	if strings.TrimSpace(string(goVersion)) == "go version go1.5.1 linux/amd64" {
		return nil
	}

	log.Println("installing go")

	loc, err := storage.ParseLocation("https://storage.googleapis.com/golang/go1.5.1.linux-amd64.tar.gz")
	if err != nil {
		return err
	}

	rc, err := storage.Get(loc)
	if err != nil {
		return err
	}
	defer rc.Close()

	err = archive.ExtractReader("/tmp", rc, "go.tar.gz")
	if err != nil {
		return err
	}

	dst := goroot
	os.RemoveAll(dst)
	os.Rename("/tmp/go", dst)

	return nil
}

func installFish() error {
	fishVersion, _ := exec.Command(filepath.Join(home, "bin", "fish"), "--version").CombinedOutput()
	if strings.TrimSpace(string(fishVersion)) == "fish, version 2.2.0" {
		return nil
	}

	log.Println("installing fish")

	url, err := storage.ParseLocation("http://fishshell.com/files/2.2.0/fish-2.2.0.tar.gz")
	if err != nil {
		return err
	}

	rc, err := storage.Get(url)
	if err != nil {
		return err
	}
	defer rc.Close()

	err = archive.ExtractReader("/tmp", rc, "fish.tar.gz")
	if err != nil {
		return err
	}

	dst := filepath.Join(home, "dist", "fish")
	os.RemoveAll(dst)
	os.Rename("/tmp/fish-2.2.0", dst)

	cmd := exec.Command("bash", "-c", "./configure --prefix=$HOME && make && make install")
	cmd.Dir = dst
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		return err
	}

	log.Println("setting fish as default shell")
	cmd = exec.Command("bash", "-c", "sudo chsh -s "+filepath.Join(home, "bin", "fish")+" vagrant")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		return err
	}

	return nil
}

func installKeys() error {
	if _, err := os.Stat(filepath.Join(home, ".password-store")); err != nil {
		log.Println("installing passwords")
		cmd := exec.Command("bash", "-c", "git clone file:///home/vagrant/repositories/password-store .password-store")
		cmd.Dir = home
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err = cmd.Run()
		if err != nil {
			return err
		}
	}

	for _, c := range []string{
		"mkdir -p /home/vagrant/.ssh",
		"cp /home/vagrant/keys/id_rsa /home/vagrant/.ssh/",
		"cp /home/vagrant/keys/id_rsa.pub /home/vagrant/.ssh/",
		"chmod 600 /home/vagrant/.ssh/id_rsa",
		"chmod 600 /home/vagrant/.ssh/id_rsa.pub",
		"cp -R -f /home/vagrant/keys/gpg/ /home/vagrant/.gnupg/",
	} {
		cmd := exec.Command("bash", "-c", c)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err := cmd.Run()
		if err != nil {
			return err
		}
	}

	return nil
}

func installDotfiles() error {
	log.Println("installing dotfiles")
	if _, err := os.Stat(filepath.Join(home, "dotfiles")); err != nil {
		cmd := exec.Command("bash", "-c", "git clone file:///home/vagrant/repositories/dotfiles")
		cmd.Dir = home
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err = cmd.Run()
		if err != nil {
			return err
		}
	}

	cmd := exec.Command("bash", "-c", "stow fish git bash")
	cmd.Dir = filepath.Join(home, "dotfiles")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		return err
	}

	return nil
}

func installAppEngine() error {
	appEngineVersion, _ := exec.Command(filepath.Join(home, "dist", "go_appengine", "goapp"), "version").CombinedOutput()
	if strings.TrimSpace(string(appEngineVersion)) == "go version go1.4.2 (appengine-1.9.27) linux/amd64" {
		return nil
	}

	log.Println("installing app engine")

	url, err := storage.ParseLocation("https://storage.googleapis.com/appengine-sdks/featured/go_appengine_sdk_linux_amd64-1.9.27.zip")
	if err != nil {
		return err
	}
	rc, err := storage.Get(url)
	if err != nil {
		return err
	}
	defer rc.Close()

	err = archive.ExtractReader("/tmp", rc, "appengine.zip")
	if err != nil {
		return err
	}

	dst := filepath.Join(home, "dist", "go_appengine")
	os.RemoveAll(dst)
	os.Rename("/tmp/go_appengine", dst)

	return nil
}

func installRedis() error {
	redisVersion, _ := exec.Command(filepath.Join(home, "bin", "redis-server"), "--version").CombinedOutput()
	if strings.TrimSpace(string(redisVersion)) == "Redis server v=3.0.4 sha=00000000:0 malloc=jemalloc-3.6.0 bits=64 build=bbb5a106838ae534" {
		return nil
	}
	log.Println("installing redis")

	url, err := storage.ParseLocation("http://download.redis.io/releases/redis-3.0.4.tar.gz")
	if err != nil {
		return err
	}
	rc, err := storage.Get(url)
	if err != nil {
		return err
	}
	defer rc.Close()

	err = archive.ExtractReader("/tmp", rc, "redis.tar.gz")
	if err != nil {
		return err
	}

	dst := filepath.Join(home, "dist", "redis")
	os.RemoveAll(dst)
	os.Rename("/tmp/redis-3.0.4", dst)

	cmd := exec.Command("bash", "-c", "make && env PREFIX=$HOME make install")
	cmd.Dir = dst
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		return err
	}

	return nil
}

func installZookeeper() error {
	dst := filepath.Join(home, "dist", "zookeeper")

	if _, err := os.Stat(filepath.Join(dst, "zookeeper-3.4.6.jar")); err == nil {
		return nil
	}

	log.Println("installing zookeeper")

	err := download(
		"http://mirrors.sonic.net/apache/zookeeper/zookeeper-3.4.6/zookeeper-3.4.6.tar.gz",
		"zookeeper.tar.gz",
		"/tmp",
	)
	if err != nil {
		return err
	}

	os.RemoveAll(dst)
	os.Rename("/tmp/zookeeper-3.4.6", dst)

	return nil
}

func installKafka() error {
	dst := filepath.Join(home, "dist", "kafka")

	if _, err := os.Stat(filepath.Join(dst, "libs", "kafka_2.10-0.8.2.1.jar")); err == nil {
		return nil
	}

	log.Println("installing kafka")

	err := download(
		"http://mirrors.sonic.net/apache/kafka/0.8.2.1/kafka_2.10-0.8.2.1.tgz",
		"kafka.tar.gz",
		"/tmp",
	)
	if err != nil {
		return err
	}

	os.RemoveAll(dst)
	os.Rename("/tmp/kafka_2.10-0.8.2.1", dst)

	return nil
}

func installCassandra() error {
	dst := filepath.Join(home, "dist", "cassandra")

	if _, err := os.Stat(filepath.Join(dst, "lib", "apache-cassandra-2.2.1.jar")); err == nil {
		return nil
	}

	log.Println("installing cassandra")

	err := download(
		"http://mirrors.sonic.net/apache/cassandra/2.2.1/apache-cassandra-2.2.1-bin.tar.gz",
		"cassandra.tar.gz",
		"/tmp",
	)
	if err != nil {
		return err
	}

	os.RemoveAll(dst)
	os.Rename("/tmp/apache-cassandra-2.2.1", dst)

	return nil
}

func installElasticSearch() error {
	dst := filepath.Join(home, "dist", "elasticsearch")

	if _, err := os.Stat(filepath.Join(dst, "lib", "elasticsearch-1.7.2.jar")); err == nil {
		return nil
	}

	log.Println("installing elasticsearch")

	err := download(
		"https://download.elastic.co/elasticsearch/elasticsearch/elasticsearch-1.7.2.tar.gz",
		"elasticsearch.tar.gz",
		"/tmp",
	)
	if err != nil {
		return err
	}

	os.RemoveAll(dst)
	os.Rename("/tmp/elasticsearch-1.7.2", dst)

	return nil
}

func installGoTools() error {
	log.Println("installing go tools")
	cmd := exec.Command("go", "get", "-v",
		"github.com/mattn/goreman",
		"github.com/robfig/glock",
	)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		return err
	}
	return nil
}

func installHaproxy() error {
	dst := filepath.Join(home, "dist", "haproxy")

	haproxyVersion, _ := exec.Command(filepath.Join(dst, "haproxy"), "-v").CombinedOutput()
	if strings.HasPrefix(strings.TrimSpace(string(haproxyVersion)), "HA-Proxy version 1.6.1") {
		return nil
	}

	log.Println("installing haproxy")

	err := download(
		"http://www.haproxy.org/download/1.6/src/haproxy-1.6.1.tar.gz",
		"haproxy.tar.gz",
		"/tmp",
	)
	if err != nil {
		return err
	}

	os.RemoveAll(dst)
	os.Rename("/tmp/haproxy-1.6.1", dst)

	cmd := exec.Command("bash", "-c", "make TARGET=linux2628")
	cmd.Dir = dst
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		return err
	}

	cmd = exec.Command("bash", "-c", "sudo setcap 'cap_net_bind_service=+ep' haproxy")
	cmd.Dir = dst
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		return err
	}

	return nil
}

func main() {
	log.SetFlags(0)

	os.MkdirAll(filepath.Join(home, "dist"), 0777)

	todo := []func() error{
		installGo,
		installFish,
		installKeys,
		installDotfiles,
		installAppEngine,
		installRedis,
		installZookeeper,
		installKafka,
		installCassandra,
		installElasticSearch,
		installGoTools,
		installHaproxy,
	}
	for _, f := range todo {
		err := f()
		if err != nil {
			log.Fatalln(err)
		}
	}

}
