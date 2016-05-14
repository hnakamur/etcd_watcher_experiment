package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"golang.org/x/net/context"

	"github.com/coreos/etcd/client"
)

var usage = `Usage etcd_watcher_example [Globals] <Command> [Options]
Commands:
  set       set the value to the key
  get       get and print the value of the key
  delete    delete the key
  watch     start a process and watch the specified key
Globals Options:
`

var subcommandOptionsUsageFormat = "\nOptions for subcommand \"%s\":\n"

func subcommandUsageFunc(subcommand string, fs *flag.FlagSet) func() {
	return func() {
		flag.Usage()
		fmt.Printf(subcommandOptionsUsageFormat, subcommand)
		fs.PrintDefaults()
	}
}

func main() {
	var help bool
	flag.BoolVar(&help, "h", false, "show help")

	flag.Usage = func() {
		fmt.Print(usage)
		flag.PrintDefaults()
	}
	flag.Parse()
	args := flag.Args()
	if len(args) == 0 {
		flag.Usage()
		os.Exit(1)
	}
	switch args[0] {
	case "get":
		getCommand(args[1:])
	case "set":
		setCommand(args[1:])
	case "delete":
		deleteCommand(args[1:])
	case "watch":
		watchCommand(args[1:])
	default:
		flag.Usage()
		os.Exit(1)
	}
}

func getCommand(args []string) {
	fs := flag.NewFlagSet("get", flag.ExitOnError)
	var key string
	fs.StringVar(&key, "key", "key1", "key to watch")
	var endpointsStr string
	fs.StringVar(&endpointsStr, "endpoints", "http://127.0.0.1:2379", "comma separated etcd endpoints")
	fs.Usage = subcommandUsageFunc("get", fs)
	fs.Parse(args)

	endpoints := strings.Split(endpointsStr, ",")
	c, err := newEtcdClient(endpoints)
	if err != nil {
		log.Fatal(err)
	}
	kAPI := client.NewKeysAPI(c)
	ctx := context.Background()
	resp, err := kAPI.Get(ctx, key, nil)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("got value=%s for key=%s\n", resp.Node.Value, key)
}

func setCommand(args []string) {
	fs := flag.NewFlagSet("set", flag.ExitOnError)
	var key string
	fs.StringVar(&key, "key", "key1", "key to watch")
	var value string
	fs.StringVar(&value, "value", "value1", "value to set")
	var endpointsStr string
	fs.StringVar(&endpointsStr, "endpoints", "http://127.0.0.1:2379", "comma separated etcd endpoints")
	fs.Usage = subcommandUsageFunc("set", fs)
	fs.Parse(args)

	endpoints := strings.Split(endpointsStr, ",")
	c, err := newEtcdClient(endpoints)
	if err != nil {
		log.Fatal(err)
	}
	kAPI := client.NewKeysAPI(c)
	ctx := context.Background()
	_, err = kAPI.Set(ctx, key, value, nil)
	if err != nil {
		log.Fatal(err)
	}
}

func deleteCommand(args []string) {
	fs := flag.NewFlagSet("delete", flag.ExitOnError)
	var key string
	fs.StringVar(&key, "key", "key1", "key to watch")
	var endpointsStr string
	fs.StringVar(&endpointsStr, "endpoints", "http://127.0.0.1:2379", "comma separated etcd endpoints")
	fs.Usage = subcommandUsageFunc("delete", fs)
	fs.Parse(args)

	endpoints := strings.Split(endpointsStr, ",")
	c, err := newEtcdClient(endpoints)
	if err != nil {
		log.Fatal(err)
	}
	kAPI := client.NewKeysAPI(c)
	ctx := context.Background()
	_, err = kAPI.Delete(ctx, key, nil)
	if err != nil {
		log.Fatal(err)
	}
}

func watchCommand(args []string) {
	fs := flag.NewFlagSet("watch", flag.ExitOnError)
	var key string
	fs.StringVar(&key, "key", "key1", "key to watch")
	var endpointsStr string
	fs.StringVar(&endpointsStr, "endpoints", "http://127.0.0.1:2379", "comma separated etcd endpoints")
	fs.Usage = subcommandUsageFunc("watch", fs)
	fs.Parse(args)

	endpoints := strings.Split(endpointsStr, ",")
	c, err := newEtcdClient(endpoints)
	if err != nil {
		log.Fatal(err)
	}
	kAPI := client.NewKeysAPI(c)
	w := kAPI.Watcher(key, nil)
	ctx := context.Background()
	for {
		resp, err := w.Next(ctx)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("got value=%s while watching key=%s\n", resp.Node.Value, key)
	}
}

func newEtcdClient(endpoints []string) (client.Client, error) {
	cfg := client.Config{
		Endpoints: endpoints,
		Transport: client.DefaultTransport,
	}
	return client.New(cfg)
}
