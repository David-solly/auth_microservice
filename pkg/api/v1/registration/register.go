package register

import (
	"math/rand"
	"os"
	"strconv"
	"time"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/sd"
	consulsd "github.com/go-kit/kit/sd/consul"
	"github.com/hashicorp/consul/api"
)

func RegisterService(consulAddress, consulPort, advertiseAddress, advertisePort, healthPort string) (registrar sd.Registrar) {

	// Logging domain.
	var logger log.Logger
	{
		logger = log.NewLogfmtLogger(os.Stderr)
		logger = log.With(logger, "ts", log.DefaultTimestampUTC)
		logger = log.With(logger, "caller", log.DefaultCaller)
	}

	rand.Seed(time.Now().UTC().UnixNano())

	// Service discovery domain. In this example we use Consul.
	var client consulsd.Client
	{
		consulConfig := api.DefaultConfig()
		consulConfig.Address = consulAddress + ":" + consulPort
		consulClient, err := api.NewClient(consulConfig)
		if err != nil {
			logger.Log("err", err)
			os.Exit(1)
		}
		client = consulsd.NewClient(consulClient)
	}

	check := api.AgentServiceCheck{
		GRPC:       "" + advertiseAddress + ":" + healthPort + "/" + "grpc.health.v1.Health",
		Interval:   "7s",
		Timeout:    "1s",
		Notes:      "Basic health checks",
		CheckID:    "health-check-util",
		Name:       "Service health status",
		GRPCUseTLS: false,
	}

	port, _ := strconv.Atoi(advertisePort)
	num := rand.Intn(100) // to make service ID unique
	asr := api.AgentServiceRegistration{
		ID:      "JWT" + strconv.Itoa(num), //unique service ID
		Name:    "JWT-Service",
		Address: advertiseAddress,
		Port:    port,
		Tags:    []string{"jwt generate", "yea-buddy!!!"},
		Check:   &check,
	}
	return consulsd.NewRegistrar(client, &asr, logger)

}
