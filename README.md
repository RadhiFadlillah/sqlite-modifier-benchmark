Simple SQLite Modifier Benchmark
-----

While working on my job, I found an issue where SQLite took too long to load a small number of data. Turns out it's happened because I was using `localtime` modifier in `SELECT` clause like this :

```sql
SELECT COUNT(*) FROM purchase
WHERE DATE(input_time, "localtime") >= "2019-01-01"
AND DATE(input_time, "localtime") <= "2019-08-10";
```

With that said, I decided to replace `localtime` modifier into `NNN hours` modifier. The change is significant, with `NNN hours` is around 5x faster than `localtime` :

```
N Days    : 10
Rows      : 2649
Localtime : 3.522 s
Hours     : 0.641 s
Hours is 5.49x faster than Localtime

N Days    : 20
Rows      : 5059
Localtime : 3.533 s
Hours     : 0.634 s
Hours is 5.57x faster than Localtime

N Days    : 40
Rows      : 9831
Localtime : 3.503 s
Hours     : 0.616 s
Hours is 5.69x faster than Localtime

N Days    : 80
Rows      : 19481
Localtime : 3.532 s
Hours     : 0.633 s
Hours is 5.58x faster than Localtime

N Days    : 160
Rows      : 38674
Localtime : 3.501 s
Hours     : 0.636 s
Hours is 5.50x faster than Localtime

N Days    : 320
Rows      : 77206
Localtime : 3.590 s
Hours     : 0.658 s
Hours is 5.46x faster than Localtime

N Days    : 640
Rows      : 154173
Localtime : 3.598 s
Hours     : 0.658 s
Hours is 5.47x faster than Localtime

N Days    : 1280
Rows      : 308009
Localtime : 3.568 s
Hours     : 0.658 s
Hours is 5.42x faster than Localtime

N Days    : 2560
Rows      : 615760
Localtime : 3.657 s
Hours     : 0.665 s
Hours is 5.50x faster than Localtime
```
