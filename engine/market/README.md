# Market

One market holds a buy side and a sell side of a trading pair.

It fires events based on order additions and deletions, volume changes, and 
matches made.

Currently it supports two type of orders:

- Limit Order
- Market Order

If a limit order has a price that crosses the market boundary it becomes a
market order. In these cases, if the market does not have enough liquidity in
the boundary, price slippage may cause the order to match at a price very 
different from what the user intended.

Perform benchmark tests with:

```
go test -bench=. ./engine/market > ./engine/market/benchmark.txt
```

To dive deeper into memory allocations, run `pprof`:

```
go test -bench=^Benchmark_MatchTakerOrder$ \
  exchange/engine/market -benchmem -cpuprofile=cpu.out -memprofile=mem.out

go tool pprof ./market.test ./mem.out

(pprof) top
(pprof) web
```

To perform escape analysis, run:

```
go build -gcflags="-m=2" > escape.txt  2>&1  
```
