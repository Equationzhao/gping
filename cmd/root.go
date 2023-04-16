package cmd

import (
	"errors"
	"fmt"
	"net"
	"os"
	"os/signal"
	"strconv"
	"time"

	"gping/tool"

	"github.com/go-ping/ping"

	"github.com/urfave/cli/v2"
)

const (
	ttl      = 50
	byteSize = 56
)

// Execute is the entry point of the program
func Execute() {
	app := cli.NewApp()
	app.Name = "gping"
	app.UsageText = "gping [options] hosts..."
	app.Usage = "a ping tool written in golang"
	app.Version = "0.0.1"
	var count int = 5
	app.Flags = []cli.Flag{
		&cli.BoolFlag{
			Name:  "O",
			Usage: "Also display a message if no response was received",
		},
		&cli.IntFlag{
			Name:  "C",
			Usage: "Ping a host only a specific number of times",
			Action: func(context *cli.Context, u int) error {
				if u > 65535 {
					return errors.New("count should be less than 65535")
				}
				count = u
				return nil
			},
		},
	}
	app.Action = func(c *cli.Context) error {
		if c.NArg() == 0 {
			tool.RedPrintln("usage error: Destination address required")
			return nil
		}

		hosts := c.Args().Slice()
		validIP := make([]string, 0, len(hosts))
		var err error
		for _, host := range hosts {
			// parse host into ip
			ip, errTemp := net.ResolveIPAddr("ip", host)
			if ip == nil {
				tool.RedPrintln("ping: unknown host " + host)
				err = errors.Join(err, errTemp)
				continue
			} else {
				validIP = append(validIP, host)
			}
		}
		l := len(validIP)
		if l == 0 {
			return cli.Exit(err, 1)
		}

		sig := make(chan os.Signal, 1)
		signal.Notify(sig, os.Interrupt) // catch ctrl+c

		for i, v := range validIP {
			pinger, err := ping.NewPinger(v)
			pinger.Size = byteSize
			pinger.TTL = ttl
			tool.GreenPrintln("PING " + pinger.Addr() + " " + strconv.Itoa(pinger.Size) + " bytes of data.")
			pinger.OnRecv = func(pkt *ping.Packet) {
				tool.GreenPrintln(fmt.Sprintf("%d bytes from %s: icmp_seq=%d ttl=%d time=%dms", pkt.Nbytes, pkt.Addr, pkt.Seq+1, pinger.TTL, pkt.Rtt.Milliseconds()))
			}
			pinger.SetPrivileged(true)
			if err != nil {
				panic(err)
			}
			pinger.Count = count
			pinger.Timeout = time.Second * 10 * time.Duration(count)
			start := time.Now()

			go func() {
				<-sig
				pinger.Stop()
			}()

			err = pinger.Run()
			if err != nil {
				panic(err)
			}

			statistic := pinger.Statistics()
			res := statisticString(statistic, time.Since(start))
			fmt.Println()
			tool.GreenPrintln(res)
			if i != l-1 {
				fmt.Println()
			}
		}
		return nil
	}
	app.HideHelpCommand = true

	app.Run(os.Args)
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
