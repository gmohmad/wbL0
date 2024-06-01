# Wildberries internship L0 task

## Quick demo

https://github.com/gmohmad/wbL0/assets/95181543/e3a0a043-fa52-4650-90aa-aa157414bc04

## Load tests:
```
===VEGETA===
Requests      [total, rate, throughput]         100000, 10000.17, 9999.92
Duration      [total, attack, wait]             10s, 10s, 242.825µs
Latencies     [min, mean, 50, 90, 95, 99, max]  87.559µs, 269.751µs, 180.864µs, 534µs, 792.036µs, 1.416ms, 7.211ms
Bytes In      [total, mean]                     86400000, 864.00
Bytes Out     [total, mean]                     0, 0.00
Success       [ratio]                           100.00%
Status Codes  [code:count]                      200:100000

-------------------------------------------------------------------------------------------------------------------

===WRK===
Thread Stats   Avg      Stdev     Max   +/- Stdev
    Latency    29.83ms   12.64ms 141.59ms   80.62%
    Req/Sec    17.10k     2.35k   21.38k    72.50%
  340374 requests in 10.09s, 315.84MB read
  Socket errors: connect 8981, read 0, write 0, timeout 0
Requests/sec:  33740.41
Transfer/sec:     31.31MB
```

## Project startup

### 1. Clone the repo
```
git clone https://github.com/gmohmad/wbL0.git
```

### 2. Add .env file into the root directory of the project and fill it up like .env.example
**you can just copy everything from .env.example and paste it into .env, or if you're using vim/neovim etc, 
just open .env file you created and type ":r .env.example" in normal mode
<br><br>
### 3. Start project containers
```
make up
```

### Stopping the project:
```
make down # run 'make downv' to also remove volumes
```
## User interface for order lookup - http://localhost:8000/
