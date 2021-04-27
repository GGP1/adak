package main

import (
	"context"
	"flag"
	"os"

	"github.com/GGP1/adak/internal/config"
	"github.com/GGP1/adak/pkg/auth"
	"github.com/GGP1/adak/pkg/http/rest"
	"github.com/GGP1/adak/pkg/postgres"
	"github.com/GGP1/adak/pkg/product"
	"github.com/GGP1/adak/pkg/review"
	"github.com/GGP1/adak/pkg/shop"
	"github.com/GGP1/adak/pkg/shopping/cart"
	"github.com/GGP1/adak/pkg/shopping/ordering"
	"github.com/GGP1/adak/pkg/user"
	"github.com/GGP1/adak/pkg/user/account"

	_ "github.com/lib/pq"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
)

type server interface {
	Run(port int) error
}

func main() {
	var (
		port         = flag.Int("port", 2727, "The service port")
		accountaddr  = flag.String("accountaddr", "account:2727", "Account service address")
		productaddr  = flag.String("productaddr", "product:2727", "Product server address")
		reviewaddr   = flag.String("reviewaddr", "review:2727", "Review server address")
		shopaddr     = flag.String("shopaddr", "shop:2727", "Shop service address")
		useraddr     = flag.String("useraddr", "user:2727", "User service address")
		orderingaddr = flag.String("orderingaddr", "ordering:2727", "Ordering service address")
		sessionaddr  = flag.String("sessionaddr", "session:2727", "Session service address")
		shoppingaddr = flag.String("shoppingaddr", "shopping:2727", "Shopping service address")
	)
	flag.Parse()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	zerolog.SetGlobalLevel(zerolog.DebugLevel)
	conf, err := config.New()
	if err != nil {
		log.Fatal().Err(err)
	}

	db, err := postgres.Connect(ctx, &conf.Postgres, os.Args[1])
	if err != nil {
		log.Fatal().Err(err)
	}
	defer db.Close()

	if len(os.Args) < 2 {
		log.Fatal().Msg("no service was specified")
	}

	var srv server
	switch os.Args[1] {
	case "account":
		srv = account.NewService(db, GRPCDial(*useraddr))
	case "product":
		srv = product.NewService(db, GRPCDial(*reviewaddr))
	case "review":
		srv = review.NewService(db)
	case "shop":
		srv = shop.NewService(
			db,
			GRPCDial(*productaddr),
			GRPCDial(*reviewaddr),
		)
	case "user":
		srv = user.NewService(
			db,
			GRPCDial(*orderingaddr),
			GRPCDial(*shoppingaddr),
		)
	case "ordering":
		srv = ordering.NewService(db, GRPCDial(*shoppingaddr))
	case "session":
		srv = auth.NewSession(db, GRPCDial(*useraddr))
	case "shopping":
		srv = cart.NewService(db)
	case "api":
		srv = rest.NewAPI(
			conf,
			GRPCDial(*accountaddr),
			GRPCDial(*productaddr),
			GRPCDial(*reviewaddr),
			GRPCDial(*shopaddr),
			GRPCDial(*useraddr),
			GRPCDial(*orderingaddr),
			GRPCDial(*sessionaddr),
			GRPCDial(*shoppingaddr),
		)
	default:
		log.Fatal().Msgf("unknown command %s", os.Args[1])
	}

	if err := srv.Run(*port); err != nil {
		log.Fatal().Err(errors.Wrapf(err, "failed running %s server", os.Args[1]))
	}
}

// GRPCDial initializes a connection on the address provided.
func GRPCDial(addr string) *grpc.ClientConn {
	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		log.Fatal().Err(errors.Wrapf(err, "dial on %s", addr))
	}
	return conn
}
