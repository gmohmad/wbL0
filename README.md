# wb_l0 task (devd in neovim btw)

## Load tests:

### Vegeta
```
Requests      [total, rate, throughput]         100000, 10000.17, 10000.03
Duration      [total, attack, wait]             10s, 10s, 139.759µs
Latencies     [min, mean, 50, 90, 95, 99, max]  81.461µs, 268.403µs, 182.298µs, 557.429µs, 776.225µs, 1.215ms, 7.443ms
Bytes In      [total, mean]                     86400000, 864.00
Bytes Out     [total, mean]                     0, 0.00
Success       [ratio]                           100.00%
Status Codes  [code:count]                      200:100000  
```

### Wrk
```
Thread Stats   Avg      Stdev     Max   +/- Stdev
    Latency    29.13ms   11.96ms 129.47ms   80.19%
    Req/Sec    17.52k     2.73k   22.59k    57.50%
  348677 requests in 10.10s, 323.55MB read
  Socket errors: connect 8981, read 0, write 0, timeout 0
Requests/sec:  34518.55
Transfer/sec:     32.03MB
```
