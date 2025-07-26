# Order Book

An order book contains all the Limit orders of a trading pair, whether buy or sell side. So, a trading pair contains two order books.

The order book is initialized on the corresponding side, and it points to the
minimum or maximum price depending on the side. 

Run the benchmark tests with:

```
go test -bench=. ./engine/orderbook > ./engine/orderbook/benchmark.txt
```

> There are a few allocations on Delete benchmarks, because the Node pool
is growing without another process to clean it with pool.Get.

To dive deeper into memory allocations, run `pprof`:

```
go test -bench=^Benchmark_Delete$ \
  exchange/engine/orderbook -benchmem -cpuprofile=cpu.out -memprofile=mem.out

go tool pprof ./orderbook.test ./mem.out

(pprof) top
(pprof) web
```

To perform escape analysis, run:

```
go build -gcflags="-m=2" > escape.txt  2>&1  
```
