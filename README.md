# gping

a ping tool

## usage

```bash
gping -C 4 google.com
```

```
PING google.com 56 bytes of data.
64 bytes from google.com: icmp_seq=1 ttl=50 time=0ms
64 bytes from google.com: icmp_seq=2 ttl=50 time=0ms
64 bytes from google.com: icmp_seq=3 ttl=50 time=0ms
64 bytes from google.com: icmp_seq=4 ttl=50 time=1ms

--- google.com ping statistics ---
4 packets transmitted, 4 packets received, 0% packet loss, time 3085ms
rtt min/avg/max/mdev = 0.683/0.866/1.045/0.128 ms
```

## install
```bash
make install
```
or 

```bash
go install github.com/equationzhao/gping@latest
```
