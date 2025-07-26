# Price Level

A Price Level contains all the orders queued for processing on a specific
price of a specific symbol.

They are processed FIFO, and it has to manage cancellations.

The list of orders is not "ordered" by arrival timestamp. Instead, they
are pushed back individually as they are received, because the architecture
of the system will handle concurrent messages at a higher level, not using
locks in the matching engine.

Using a double linked list we get O(1) additions to the back, and reads from the
front. We also get O(1) deletions if we have a pointer to the element in the
doubly linked list. For deletions, we keep a map of order ID to pointer in the
list, which even though increases the overall memory consumption, allows for
lightning fast deletions.

Without this order map, deleting an element from the list using an order ID
would take O(n), as we would have to iterate over the entire list to find
the element and delete it.

For data streaming, each time an order is added or deleted, the Volume of the
price level is updated accordingly, being available in O(1) when needed. This
Volume is designed to be read only when creating snapshots, not to broadcast
on every single change.

Matching a market order requires extracting the open orders in the list, FIFO,
until the required volume is extracted. This is an O(n) operation, and there
is likely no way around it.

Run the benchmark tests with:

```
go test -bench=. ./engine/orderbook/pricelevel > ./engine/orderbook/pricelevel/benchmark.txt
```

To dive deeper into memory allocations, run `pprof`:

```
go test -bench=^Benchmark_Remove$ \
  exchange/engine/orderbook/pricelevel -benchmem -cpuprofile=cpu.out -memprofile=mem.out

go tool pprof ./pricelevel.test ./mem.out

(pprof) top
(pprof) web
```