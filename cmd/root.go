package cmd

import (
	"errors"
	"fmt"
	"net"
	"os"
	"os/signal"
	"strconv"
	"time"

	"github.com/equationzhao/gping/tool"

	"github.com/go-ping/ping"

	"github.com/urfave/cli/v2"
)

var (
	ttl      int
	byteSize int
	interval int
)

// Execute is the entry point of the program
func Execute() error {
	app := cli.NewApp()
	app.Name = "gping"
	app.UsageText = "gping [options] hosts..."
	app.Usage = "a ping tool written in golang"
	app.Version = "0.0.2"
	count := 65535
	app.Flags = []cli.Flag{
		&cli.IntFlag{
			Name:    "c",
			Aliases: []string{"count", "C"},
			Usage:   "Ping a host only a specific number of times",
			Action: func(context *cli.Context, u int) error {
				if u > 65535 {
					return errors.New("count should be less than 65535")
				}
				count = u
				return nil
			},
		},
		&cli.IntFlag{
			Name:        "byte-size",
			Aliases:     []string{"bs"},
			Usage:       "Set byte size of ping packet",
			Value:       56,
			Destination: &byteSize,
		},
		&cli.IntFlag{
			Name:        "ttl",
			Usage:       "Set ttl of ping packet",
			Aliases:     []string{"TTL"},
			Value:       50,
			Destination: &ttl,
		},
		&cli.IntFlag{
			Name:        "i",
			Aliases:     []string{"I", "interval"},
			Usage:       "specifying the interval in seconds between requests",
			Value:       1,
			Destination: &interval,
		},
		&cli.BoolFlag{
			Name:    "6",
			Aliases: []string{"ipv6"},
			Usage:   "Use IPv6",
			Value:   false,
		},
		&cli.BoolFlag{
			Name:    "4",
			Aliases: []string{"ipv4"},
			Usage:   "Use IPv4",
			Value:   false,
		},
		&cli.BoolFlag{
			Name:    "a",
			Aliases: []string{"bell"},
			Usage:   "ring the bell when a packet is received (if your terminal supports it)",
			Value:   false,
		},
	}
	app.Action = func(c *cli.Context) error {
		if c.NArg() == 0 {
			tool.RedPrintln("usage error: Destination address required")
			return nil
		}

		v4 := c.Bool("4")
		v6 := c.Bool("6")

		if v4 && v6 {
			return errors.New("cannot use both -4 and -6")
		}

		hosts := c.Args().Slice()

		sig := make(chan os.Signal, 1)
		signal.Notify(sig, os.Interrupt) // catch ctrl+c
		l := len(hosts)
		for i, v := range hosts {

			var ip *net.IPAddr
			var err error
			if v4 {
				ip, err = net.ResolveIPAddr("ip4", v)
			} else if v6 {
				ip, err = net.ResolveIPAddr("ip6", v)
			} else {
				ip, err = net.ResolveIPAddr("ip", v)
			}
			if err != nil {
				tool.RedPrintln(err)
				continue
			}

			pinger := ping.New(ip.String())
			pinger.Size = byteSize
			pinger.TTL = ttl
			pinger.SetPrivileged(true)
			pinger.Count = count
			pinger.Interval = time.Duration(interval) * time.Second
			pinger.Timeout = time.Duration(count) * time.Second * 10

			Prompt := "PING " + v + " (" + ip.String() + ") " + strconv.Itoa(pinger.Size) + " bytes of data."

			tool.GreenPrintln(Prompt)
			start := time.Now()
			pinger.OnRecv = func(pkt *ping.Packet) {
				if c.Bool("a") {
					fmt.Print("\x07")
				}
				tool.GreenPrintln(fmt.Sprintf("%d bytes from %s: icmp_seq=%d ttl=%d time=%.3gms", pkt.Nbytes, pkt.Addr, pkt.Seq+1, pinger.TTL, float64(pkt.Rtt.Microseconds())/1000.0))
			}

			go func() {
				<-sig
				pinger.Stop()
			}()

			err = pinger.Run()

			if err != nil {
				tool.RedPrintln(err)
				continue
			}

			res := statisticString(pinger.Statistics(), time.Since(start))
			fmt.Println()
			tool.GreenPrintln(res)
			if i != l-1 {
				fmt.Println()
			}

		}
		return nil
	}
	app.HideHelpCommand = true

	return app.Run(os.Args)
}

func statisticString(s *ping.Statistics, total time.Duration) string {
	return fmt.Sprintf("--- %s ping statistics ---\n%d packets transmitted, %d packets received, %.f%% packet loss, time %dms\n"+
		"rtt min/avg/max/mdev = %.3f/%.3f/%.3f/%.3f ms", s.Addr,
		s.PacketsSent, s.PacketsRecv, s.PacketLoss, total.Milliseconds(),
		float64(s.MinRtt.Microseconds())/1000,
		float64(s.AvgRtt.Microseconds())/1000,
		float64(s.MaxRtt.Microseconds())/1000,
		float64(s.StdDevRtt.Microseconds())/1000)
}
