package cmd

import (
	"net"
	"strings"

	"github.com/go-acme/lego/v4/challenge"
	"github.com/go-acme/lego/v4/challenge/http01"
	"github.com/go-acme/lego/v4/lego"
	"github.com/go-acme/lego/v4/log"
	"github.com/go-acme/lego/v4/providers/http/memcached"
	"github.com/go-acme/lego/v4/providers/http/webroot"
	"github.com/urfave/cli/v2"
)

func setupChallenges(ctx *cli.Context, client *lego.Client) {
	if !ctx.Bool("http") && !ctx.Bool("tls") && !ctx.IsSet("dns") {
		log.Fatal("No challenge selected. You must specify at least one challenge: `--http`.")
	}

	if ctx.Bool("http") {
		err := client.Challenge.SetHTTP01Provider(setupHTTPProvider(ctx))
		if err != nil {
			log.Fatal(err)
		}
	}
}

//nolint:gocyclo // the complexity is expected.
func setupHTTPProvider(ctx *cli.Context) challenge.Provider {
	switch {
	case ctx.IsSet("http.webroot"):
		ps, err := webroot.NewHTTPProvider(ctx.String("http.webroot"))
		if err != nil {
			log.Fatal(err)
		}
		return ps
	case ctx.IsSet("http.memcached-host"):
		ps, err := memcached.NewMemcachedProvider(ctx.StringSlice("http.memcached-host"))
		if err != nil {
			log.Fatal(err)
		}
		return ps
	case ctx.IsSet("http.port"):
		iface := ctx.String("http.port")
		if !strings.Contains(iface, ":") {
			log.Fatalf("The --http switch only accepts interface:port or :port for its argument.")
		}

		host, port, err := net.SplitHostPort(iface)
		if err != nil {
			log.Fatal(err)
		}

		srv := http01.NewProviderServer(host, port)
		if header := ctx.String("http.proxy-header"); header != "" {
			srv.SetProxyHeader(header)
		}
		return srv
	case ctx.Bool("http"):
		srv := http01.NewProviderServer("", "")
		if header := ctx.String("http.proxy-header"); header != "" {
			srv.SetProxyHeader(header)
		}
		return srv
	default:
		log.Fatal("Invalid HTTP challenge options.")
		return nil
	}
}
