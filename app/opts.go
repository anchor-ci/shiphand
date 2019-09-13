package app

// TODO: Decouple this from the application, investigate Traefik's CLI args

import (
  "github.com/urfave/cli"
)

func GetOpts() []cli.Flag {
	return []cli.Flag{
		cli.StringFlag{
			Name:  "run, r",
			Usage: "For running a local CI definition",
            EnvVar: "DEBUG",
		},
        cli.StringFlag{
            Name: "redis-host, rh",
            Value: "localhost",
            Usage: "Defines a redis host to connect to",
            EnvVar: "REDIS_URL",
        },
        cli.StringFlag{
          Name: "redis-port, rp",
          Value: "6379",
          Usage: "Defines a redis port to connect to",
          EnvVar: "REDIS_PORT",
        },
        cli.StringFlag{
          Name: "key, k",
          Value: "job:v1:*",
          Usage: "Defines the Redis key to look for, for new jobs",
        },
	}
}
