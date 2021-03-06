package main

import (
	"log"

	"github.com/docopt/docopt-go"
	"github.com/dotcloud/docker/daemon/graphdriver"
	_ "github.com/dotcloud/docker/daemon/graphdriver/aufs"
	_ "github.com/dotcloud/docker/daemon/graphdriver/btrfs"
	_ "github.com/dotcloud/docker/daemon/graphdriver/devmapper"
	_ "github.com/dotcloud/docker/daemon/graphdriver/vfs"
	"github.com/flynn/pinkerton/store"
)

func init() {
	log.SetFlags(log.Lshortfile)
}

func main() {
	usage := `Pinkerton manages Docker images.

Usage:
  pinkerton pull [options] <image-url>
  pinkerton checkout [options] <id> <image-id>
  pinkerton cleanup [options] <id>
  pinkerton -h | --help

Commands:
  pull      Download a Docker image
  checkout  Create a working copy of an image
  cleanup   Destroy a working copy of an image

Examples:
  pinkerton pull https://registry.hub.docker.com/redis
  pinkerton pull https://registry.hub.docker.com/ubuntu?tag=trusty
  pinkerton pull https://registry.hub.docker.com/flynn/slugrunner?id=1443bd6a675b959693a1a4021d660bebbdbff688d00c65ff057c46702e4b8933
  pinkerton checkout slugrunner-test 1443bd6a675b959693a1a4021d660bebbdbff688d00c65ff057c46702e4b8933
  pinkerton cleanup slugrunner-test

Options:
  -h, --help       show this message and exit
  --driver=<name>  storage driver [default: aufs]
  --root=<path>    storage root [default: /var/lib/docker]
`

	args, _ := docopt.Parse(usage, nil, true, "", false)

	root := args["--root"].(string)
	driver, err := graphdriver.GetDriver(args["--driver"].(string), root)
	if err != nil {
		log.Fatal(err)
	}

	s, err := store.New(root, driver)
	if err != nil {
		log.Fatal(err)
	}
	ctx := &Context{Store: s, driver: driver}

	switch {
	case args["pull"].(bool):
		ctx.Pull(args["<image-url>"].(string))
	case args["checkout"].(bool):
		ctx.Checkout(args["<id>"].(string), args["<image-id>"].(string))
	case args["cleanup"].(bool):
		ctx.Cleanup(args["<id>"].(string))
	}
}
