# Red Black Tree

Red Black Trees are a type of self-balancing binary trees. Other types include
AVL trees.

AVL trees are heavily balanced and provide faster lookups, so for a look-up
intensive task, use AVL trees.

For insert-intensive tasks, prefer Red Black Trees.

On an exchange, the load will be in insertions, when the price changes as new
prices are added and many prices are deleted. Each time an order is added to an
existing price, that's a lookup operation, and each time we insert or delete
a price that includes a lookup operation too. Still, the tradeoff of having less
rotations for a Red Black tree when there are insertions and deletions make them
more suitable for the order book than AVL trees.

Run logic tests with:

```
go test exchange/engine/orderbook/rbtree
```

Run fuzzy tests with:

```
go test -fuzz=Fuzz_Insert_Delete --fuzztime=30s exchange/engine/orderbook/rbtree
```

Update benchmark tests with:

```
go test -bench=. ./engine/orderbook/rbtree > ./engine/orderbook/rbtree/benchmark.txt
```

To dive deeper into memory allocations, run `pprof`:

```
go test -bench=^Benchmark_Insert_Delete$ \
  exchange/engine/orderbook/rbtree -benchmem -cpuprofile=cpu.out -memprofile=mem.out

go tool pprof ./rbtree.test ./mem.out

(pprof) top
(pprof) web
```