# Zadanie: Postgres - Indexácia

**Timotej Zaťko**

## Otázky

1. Vyhľadajte v accounts screen_name s presnou hodnotou ‘realDonaldTrump’ a analyzujte daný select. Akú metódu vám vybral plánovač?
2. Koľko workerov pracovalo na danom selecte? Zdvihnite počet workerov apovedzte ako to ovplyvňuje čas. Je tam nejaký strop?
3. Vytvorte index nad screen_name a pozrite ako sa zmenil a porovnajte výstup oproti požiadavke bez indexu. Potrebuje plánovač vtejto požiadavke viac workerov? Bol tu aplikovaný nejaký filter na riadky? Prečo?
4. Vyberte používateľov, ktorý majú followers_count väčší, rovný ako 100 a zároveň menší, rovný 200. Je správanie rovnaké v prvej úlohe? Je správanie rovnaké ako v druhej úlohe? Prečo?
5. Vytvorte index nad 4 úlohou a popíšte prácu s indexom. Čo je to Bitmap Index Scan a prečo je tam Bitmap Heap Scan? Prečo je tam recheck condition?
6. Vyberte používateľov, ktorí majú followers_count väčší, rovný ako 100 azároveň menší, rovný 1000? V čom je rozdiel, prečo?
7. Vytvorte daľšie 3 indexy na name, friends_count, a description a insertnite si svojho používateľa (to je jedno aké dáta) do accounts. Koľko to trvalo? Dropnite indexy a spravte to ešte raz. Prečo je tu rozdiel?
8. Vytvorte index nad tweetami pre retweet_count a pre content. Porovnajte ich dĺžku vytvárania. Prečo je tu taký rozdiel?
9. Porovnajte indexy pre retweet_count, content, followers_count, screen_name,...
včom sa líšia aprečo (stačí stručne)?
   1. `create extension pageinspect;`
   2. `select * from bt_metap('idx_content');`
   3. `select type, live_items, dead_items, avg_item_size, page_size, free_size from bt_page_stats('idx_content', 1000);`
   4. `select itemoffset, itemlen, data from bt_page_items('idx_content', 1) limit 1000;`

## Odpovede

### 1. Vyhľadajte v accounts screen_name s presnou hodnotou ‘realDonaldTrump’ a analyzujte daný select. Akú metódu vám vybral plánovač?

Plánovač vybral __Parallel Seqential Scan__. 

Query: `EXPLAIN ANALYZE SELECT * FROM accounts a WHERE a.screen_name='realDonaldTrump' LIMIT 1;`

Explain Analyze:
```
                                                            QUERY PLAN
-----------------------------------------------------------------------------------------------------------------------------------
 Limit  (cost=1000.00..82625.79 rows=1 width=112) (actual time=0.717..442.980 rows=1 loops=1)
   ->  Gather  (cost=1000.00..82625.79 rows=1 width=112) (actual time=0.696..442.932 rows=1 loops=1)
         Workers Planned: 2
         Workers Launched: 1
         ->  Parallel Seq Scan on accounts a  (cost=0.00..81625.69 rows=1 width=112) (actual time=217.876..217.884 rows=0 loops=2)
               Filter: ((screen_name)::text = 'realDonaldTrump'::text)
               Rows Removed by Filter: 1734207
 Planning Time: 1.241 ms
 Execution Time: 443.083 ms
(9 rows)
```

### 2. Koľko workerov pracovalo na danom selecte? Zdvihnite počet workerov apovedzte ako to ovplyvňuje čas. Je tam nejaký strop?

Query: `EXPLAIN ANALYZE SELECT * FROM accounts a WHERE a.screen_name='realDonaldTrump' LIMIT 1;`

Explain Analyze:
```
                                                            QUERY PLAN
-----------------------------------------------------------------------------------------------------------------------------------
 Limit  (cost=1000.00..82625.79 rows=1 width=112) (actual time=0.607..658.049 rows=1 loops=1)
   ->  Gather  (cost=1000.00..82625.79 rows=1 width=112) (actual time=0.590..658.011 rows=1 loops=1)
         Workers Planned: 2
         Workers Launched: 2
         ->  Parallel Seq Scan on accounts a  (cost=0.00..81625.69 rows=1 width=112) (actual time=324.833..324.841 rows=0 loops=2)
               Filter: ((screen_name)::text = 'realDonaldTrump'::text)
               Rows Removed by Filter: 1734207
 Planning Time: 0.085 ms
 Execution Time: 658.123 ms
(9 rows)
```


__2__ workery pracovali na danom selecte. Čím viac workerov, tak tým rýchlejšie sa query vykonaná. Čas vykonávania podľa počtu nastavených workerov (query som vykonal viac krát, rýchlosť vykonávania sa vždy viac-menej líšila):
* 2: 658.123ms
* 3: 448.429ms
* 4: 241.287ms
* 5 > plánovač vybral maximálne 4 workery (tj. strop použitých workerov)

### 3. Vytvorte index nad screen_name a pozrite ako sa zmenil a porovnajte výstup oproti požiadavke bez indexu. Potrebuje plánovač vtejto požiadavke viac workerov? Bol tu aplikovaný nejaký filter na riadky? Prečo? 

Vytvorenie indexu - `CREATE INDEX index_screen_name ON accounts (screen_name);`

Query: `EXPLAIN ANALYZE SELECT * FROM accounts a WHERE a.screen_name='realDonaldTrump' LIMIT 1;`

Explain Analyze:
```
                                                              QUERY PLAN
--------------------------------------------------------------------------------------------------------------------------------------
 Limit  (cost=0.43..8.45 rows=1 width=112) (actual time=0.190..0.576 rows=1 loops=1)
   ->  Index Scan using index_screen_name on accounts a  (cost=0.43..8.45 rows=1 width=112) (actual time=0.146..0.167 rows=1 loops=1)
         Index Cond: ((screen_name)::text = 'realDonaldTrump'::text)
 Planning Time: 0.249 ms
 Execution Time: 0.834 ms
(5 rows)
```

Plánovač nepotrebuje viac workerov, na riadky nebol aplikovaný žiaden filter pretože sa používa index (ktorý používa vyhľadávací strom). 
V indexe sa vyhľadáva na základe `screen_name`, keď sa nájde v indexe (vyhľadávacom strome) `screen_name=realDonaldTrump`, daný node obsahuje smerník na daný riadok, preto nie je potrebný filter.

### 4. Vyberte používateľov, ktorý majú followers_count väčší, rovný ako 100 a zároveň menší, rovný 200. Je správanie rovnaké v prvej úlohe? Je správanie rovnaké ako v druhej úlohe? Prečo?

Query: `EXPLAIN ANALYZE SELECT * FROM accounts a WHERE a.followers_count>=100 and a.followers_count <=200;`

Explain Analyze:
```
                                                        QUERY PLAN
--------------------------------------------------------------------------------------------------------------------------
 Seq Scan on accounts a  (cost=0.00..115587.32 rows=411857 width=112) (actual time=198.184..5035.599 rows=417112 loops=1)
   Filter: ((followers_count >= 100) AND (followers_count <= 200))
   Rows Removed by Filter: 3051309
 Planning Time: 0.092 ms
 JIT:
   Functions: 2
   Options: Inlining false, Optimization false, Expressions true, Deforming true
   Timing: Generation 9.618 ms, Inlining 0.000 ms, Optimization 32.026 ms, Emission 163.680 ms, Total 205.324 ms
 Execution Time: 9272.819 ms
(9 rows)
```

Plánovač vybral __Seqential Scan__, prečo? pretože vyhodnotil, že použitie sekvenčného skenu je výhodnejšie ako použitie paralelného sekvenčného skenu. 
Napríklad kvôli tomu, že pri paralelnom sekvenčnom skene je potrebná synchronizácia medzi workermi, v tomto prípade vyhodnotil, že by synchronizácia medzi nimi bola drahá. Čo sa týka správania oproti druhej úlohe, keďže bol použitý iba sekvenčný sken, nie je sa s čím porovnávať. 

### 5. Vytvorte index nad 4 úlohou a popíšte prácu s indexom. Čo je to Bitmap Index Scan a prečo je tam Bitmap Heap Scan? Prečo je tam recheck condition?

Vytvorenie indexu - `CREATE INDEX index_followers_count ON accounts (followers_count);`

Query: `EXPLAIN ANALYZE SELECT * FROM accounts a WHERE a.followers_count>=100 and a.followers_count<=200;`. 

Explain Analyze:
```
                                                                 QUERY PLAN
---------------------------------------------------------------------------------------------------------------------------------------------
 Bitmap Heap Scan on accounts a  (cost=8749.96..78488.82 rows=411857 width=112) (actual time=95.676..4236.551 rows=417112 loops=1)
   Recheck Cond: ((followers_count >= 100) AND (followers_count <= 200))
   Heap Blocks: exact=63067
   ->  Bitmap Index Scan on index_followers_count  (cost=0.00..8647.00 rows=411857 width=0) (actual time=83.067..83.075 rows=417112 loops=1)
         Index Cond: ((followers_count >= 100) AND (followers_count <= 200))
 Planning Time: 0.108 ms
 Execution Time: 7609.266 ms
(7 rows)
```

Bitmap scan je akoby kompromis medzi __Seqential Scan__ a __Index Scan__, použije sa najmä vtedy keď berieme taký počet riadkov, pri ktorých sa plánovač
rozhodne, že ich je stále málo na použitie __Seqential Scan__ ale príliš veľa na použitie __Index Scan__. 

__Bitmap Index Scan__ najskôr prejde index a vyberie tie bloky z indexu, ktoré splňujú podmienku, v ktorých sa nachádzajú také riadky spĺňajúce podmienku `followers_count>=100 and followers_count<=200` - to je _Index Cond_. 
Potom __Bitmap Heap Scan__ zoberie list týchto blokov, a pre všetky záznamy z týchto blokob sa následne pozerá do tabuľky, a vyberá iba také záznamy, ktoré spĺňajú podmienku `followers_count>=100 and followers_count<=200` - to je _Recheck Cond_.
Robí to preto, lebo v týchto blokoch sa môžu nachádzať aj záznamy, ktoré túto podmienku nespĺňajú. 
[[ref]](https://www.cybertec-postgresql.com/en/postgresql-indexing-index-scan-vs-bitmap-scan-vs-sequential-scan-basics/?gclid=Cj0KCQjw5eX7BRDQARIsAMhYLP9eCopBtmKbz8ed_s3LKcQpV-YLTArz0E5NKGYzgTkdDfhsvLHoszEaAuvNEALw_wcB) 

### 6. Vyberte používateľov, ktorí majú followers_count väčší, rovný ako 100 azároveň menší, rovný 1000? V čom je rozdiel, prečo?

Query: `EXPLAIN ANALYZE SELECT * FROM accounts a WHERE a.followers_count>=100 and a.followers_count<=1000;`. 

Explain Analyze:
```
                                                        QUERY PLAN
---------------------------------------------------------------------------------------------------------------------------
 Seq Scan on accounts a  (cost=0.00..115587.32 rows=1473569 width=112) (actual time=4.627..12603.684 rows=1480450 loops=1)
   Filter: ((followers_count >= 100) AND (followers_count <= 1000))
   Rows Removed by Filter: 1987971
 Planning Time: 0.269 ms
 JIT:
   Functions: 2
   Options: Inlining false, Optimization false, Expressions true, Deforming true
   Timing: Generation 3.289 ms, Inlining 0.000 ms, Optimization 0.597 ms, Emission 3.678 ms, Total 7.564 ms
 Execution Time: 23822.072 ms
(9 rows)
```

V tomto prípade plánovač uprednostnil __Seqential Scan__, pravdepodobne vyberáme veľkú časť riadkov tabuľy a __Bitmap Scan__ by už nebol asi taký efektívny - tj. vybral by taký počet blokov, ktorý by bol blízko k celkovému počtu blokov, preto v tomto prípade plánovač vybral __Seqential Scan__.
Ak si plánovač myslí, že by musel prejsť takmer všetky bloky, uprednostní radšej __Seqential Scan__.  

###  7. Vytvorte daľšie 3 indexy na name, friends_count, a description a insertnite si svojho používateľa (to je jedno aké dáta) do accounts. Koľko to trvalo? Dropnite indexy a spravte to ešte raz. Prečo je tu rozdiel?

Query: `INSERT INTO accounts ("id", "screen_name", "name", "description", "followers_count", "friends_count", "statuses_count") VALUES (1, 'timzatko', 'timzatko', 'Hello!', 1000, 10000, 99999)`. 

Explain Analyze - __without indexes__:
```
EXPLAIN ANALYZE INSERT INTO accounts VALUES (2, 'timqo', 'Timqo', 'To som ja!', 1, 10, 1000);
                                            QUERY PLAN
--------------------------------------------------------------------------------------------------
 Insert on accounts  (cost=0.00..0.01 rows=1 width=888) (actual time=2.887..2.910 rows=0 loops=1)
   ->  Result  (cost=0.00..0.01 rows=1 width=888) (actual time=0.008..0.024 rows=1 loops=1)
 Planning Time: 0.036 ms
 Execution Time: 2.964 ms
(4 rows)
```

Explain Analyze - __with indexes__:
```
                                            QUERY PLAN
--------------------------------------------------------------------------------------------------
 Insert on accounts  (cost=0.00..0.01 rows=1 width=888) (actual time=5.801..5.827 rows=0 loops=1)
   ->  Result  (cost=0.00..0.01 rows=1 width=888) (actual time=0.010..0.029 rows=1 loops=1)
 Planning Time: 0.088 ms
 Execution Time: 5.925 ms
(4 rows)
```

Vloženie nového riadku trvalo _5.925 ms_. Keď som dropol indexy, trvalo to kratšie (_2.964 ms_) pretože pri vkladaní sa nemusel aktualizovať index - až tri indexy. Aktualizácia indexu môže je náročná operácia (napr. v prípade stromových indexov je nutné ich niekedy vyvážiť).

### 8. Vytvorte index nad tweetami pre retweet_count a pre content. Porovnajte ich dĺžku vytvárania. Prečo je tu taký rozdiel?

Query 1: `CREATE INDEX index_retweet_count ON tweets (retweet_count);` (24 708 ms)

```
pdt=# CREATE INDEX index_retweet_count
pdt-#     ON tweets (retweet_count);
CREATE INDEX
Time: 24708.773 ms (00:24.709)
```

Query 1: `CREATE INDEX index_content ON tweets (content);` (71 853 ms)

```
pdt=# CREATE INDEX index_content
pdt-#     ON tweets (content);
CREATE INDEX
Time: 71853.356 ms (01:11.853)
```

Vytvorenie indexu pre _content_ trvá dlhšie, pretože porovnanie medzi reťazcami je výpočtovo náročnejšie, ako porovnávanie medzi číslami v _retweet_count_. Keď sa vytvára index a vkladá do stromu, porovnávajú sa jednotlivé prvky v strome aby sa vkladaný prvok vložil na správne miesto.  

### 9. Porovnajte indexy pre retweet_count, content, followers_count, screen_name,... včom sa líšia aprečo (stačí stručne)?

1. `create extension pageinspect;`
2. `select * from bt_metap('index_content');`
3. `select type, live_items, dead_items, avg_item_size, page_size, free_size from bt_page_stats('index_content', 1000);`
4. `select itemoffset, itemlen, data from bt_page_items('index_content', 1) limit 1000;`
  
    
#### select * from bt_metap('index_content');

```
pdt=# select * from bt_metap('index_content');
 magic  | version | root  | level | fastroot | fastlevel | oldest_xact | last_cleanup_num_tuples
--------+---------+-------+-------+----------+-----------+-------------+-------------------------
 340322 |       4 | 10285 |     4 |    10285 |         4 |           0 |                      -1
(1 row)
```

```
pdt=# select * from bt_metap('index_retweet_count');
 magic  | version | root | level | fastroot | fastlevel | oldest_xact | last_cleanup_num_tuples
--------+---------+------+-------+----------+-----------+-------------+-------------------------
 340322 |       4 |  209 |     2 |      209 |         2 |           0 |                      -1
(1 row)

Time: 9.688 ms
```

```
pdt=# select * from bt_metap('index_followers_count');
 magic  | version | root | level | fastroot | fastlevel | oldest_xact | last_cleanup_num_tuples
--------+---------+------+-------+----------+-----------+-------------+-------------------------
 340322 |       4 |  209 |     2 |      209 |         2 |           0 |                      -1
(1 row)

Time: 27.176 ms
```

```
pdt=# select * from bt_metap('index_screen_name');
 magic  | version | root | level | fastroot | fastlevel | oldest_xact | last_cleanup_num_tuples
--------+---------+------+-------+----------+-----------+-------------+-------------------------
 340322 |       4 |  220 |     2 |      220 |         2 |           0 |                      -1
(1 row)

Time: 4.677 ms
```

```
pdt=# select * from bt_metap('index_name');
 magic  | version | root | level | fastroot | fastlevel | oldest_xact | last_cleanup_num_tuples
--------+---------+------+-------+----------+-----------+-------------+-------------------------
 340322 |       4 |  177 |     2 |      177 |         2 |           0 |                      -1
(1 row)

Time: 3.757 ms
```

```
select * from bt_metap('index_friends_count');
 magic  | version | root | level | fastroot | fastlevel | oldest_xact | last_cleanup_num_tuples
--------+---------+------+-------+----------+-----------+-------------+-------------------------
 340322 |       4 |  209 |     2 |      209 |         2 |           0 |                      -1
(1 row)

Time: 2.352 ms
```

```
pdt=# select * from bt_metap('index_description');
 magic  | version | root  | level | fastroot | fastlevel | oldest_xact | last_cleanup_num_tuples
--------+---------+-------+-------+----------+-----------+-------------+-------------------------
 340322 |       4 | 40009 |     4 |    40009 |         4 |           0 |                      -1
(1 row)

Time: 2.984 ms
```

Indexy sa líšia v hodnotách pre _root_, _level_, _fastroot_, _fastlevel_.

- _root_ označuje číslo bloku pre koreň stromu
- _level_ označuje úroveň v strome na ktorej sa nachádza jeho koreň (root)
- _fastroot_ označuje číslo bloku pre koreň (fastroot) stromu
- _fastlevel_ označuje úroveň v strome na ktorej sa nachádza jeho koreň (fastroot)

_fastroot_, _fastlevel_ a fast level sa používajú pre optimalizáciu vyhľadávania. Namiesto začatia vyhľadávania z _root_, začne sa vyhľadávať cez _fast root_, tak sa vyhľadávač vyhne úrovniam stromu iba s jedným dieťaťom.
Tieto hodnoty môžu byť nastavené po vymazaní veľkého počtu záznamov. [[ref](http://www.louisemeta.com/blog/indexes-btree-algorithms/#deleting-from-a-btree)]

Každý index je tvorený vlastným stromom, preto má každý index odlišný koreň.

#### `select type, live_items, dead_items, avg_item_size, page_size, free_size from bt_page_stats('index_content', 1000);`

```
pdt=# select type, live_items, dead_items, avg_item_size, page_size, free_size from bt_page_stats('index_content', 1000);
 type | live_items | dead_items | avg_item_size | page_size | free_size
------+------------+------------+---------------+-----------+-----------
 l    |         36 |          0 |           200 |      8192 |       804
(1 row)

Time: 2.655 ms
```

```
pdt=# select type, live_items, dead_items, avg_item_size, page_size, free_size from bt_page_stats('index_retweet_count', 1000);
 type | live_items | dead_items | avg_item_size | page_size | free_size
------+------------+------------+---------------+-----------+-----------
 l    |        367 |          0 |            16 |      8192 |       800
(1 row)

Time: 3.830 ms
```

```
pdt=# select type, live_items, dead_items, avg_item_size, page_size, free_size from bt_page_stats('index_followers_count', 1000);
 type | live_items | dead_items | avg_item_size | page_size | free_size
------+------------+------------+---------------+-----------+-----------
 l    |        367 |          0 |            16 |      8192 |       800
(1 row)

Time: 2.490 ms
```

```
pdt=# select type, live_items, dead_items, avg_item_size, page_size, free_size from bt_page_stats('index_screen_name', 1000);
 type | live_items | dead_items | avg_item_size | page_size | free_size
------+------------+------------+---------------+-----------+-----------
 l    |        282 |          0 |            22 |      8192 |       804
(1 row)

Time: 4.591 ms
```

```
pdt=# select type, live_items, dead_items, avg_item_size, page_size, free_size from bt_page_stats('index_name', 1000);
 type | live_items | dead_items | avg_item_size | page_size | free_size
------+------------+------------+---------------+-----------+-----------
 l    |        220 |          0 |            29 |      8192 |       804
(1 row)

Time: 1.665 ms
```

```
pdt=# select type, live_items, dead_items, avg_item_size, page_size, free_size from bt_page_stats('index_friends_count', 1000);
 type | live_items | dead_items | avg_item_size | page_size | free_size
------+------------+------------+---------------+-----------+-----------
 l    |        367 |          0 |            16 |      8192 |       800
(1 row)

Time: 2.434 ms
```

```
pdt=# select type, live_items, dead_items, avg_item_size, page_size, free_size from bt_page_stats('index_description', 1000);
 type | live_items | dead_items | avg_item_size | page_size | free_size
------+------------+------------+---------------+-----------+-----------
 l    |        367 |          0 |            16 |      8192 |       800
(1 row)

Time: 1.928 ms
```

Indexy sa odlišujú v:
- live_items - pretože v každej stránke daného indexu je rôzny počet položiek
- avg_item_size - pretože v každý index indexuje niečo iné - čísla, text rôznych dĺžok
- free_size - pretože v indexoch/stránkach indexu je rôzny počet itemov o rôznych veľkostí, tak ostane rôzne voľného miesta

Indexy sa __neodlišujú__ v:
- type - pretože sú indexy rovnakého typu
- dead_items - pretože sme nič nemazali z tabuľky a nebolo treba nič mazať z indexu (dead_items vznikajú pri vymazávaní)
- page_size - pretože máme jednotnú page_size vo pre indexy v db

#### `select itemoffset, itemlen, data from bt_page_items('index_content', 1) limit 1000;`

Výpis je veľmi dlhý a nebudem ho sem dávať. V skratke, tento príkaz pre daný index vráti z prvej stánky max 1000 položiek.
Stĺpce _itemoffset_ a _itemlen_ slúžia na prečítanie hodnoty danej položky - kde sa nachádza v stránke a aká je dlhá. _data_ sú prečitané dáta.
Preto sa v týchto troch hodnotách budú líšiť všetky indexy.