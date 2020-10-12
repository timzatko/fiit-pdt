# Zadanie: Postgres - Indexácia 2

**Timotej Zaťko**

## Zadanie

V celom tomto zadaní nepoužívajte trigramy ani inverzný index.
1. Vyhľadajte v tweets.content meno „Gates“ na ľubovoľnom mieste a porovnajte
výsledok po tom, ako content naindexujete. V čom je rozdiel a prečo?
2. Vyhľadajte tweet, ktorý začína “DANGER: WARNING:”. Použil sa index?
3. Teraz naindexujte content tak, aby sa použil index a zhodnoťte prečo sa pred tým
nad “DANGER: WARNING:” nepoužil. Použije sa teraz na „Gates“ na ľubovoľnom
mieste?
4. Vytvorte nový index, tak aby ste vedeli vyhľadať tweet, ktorý konči reťazcom
„LUCIFERASE“ a nezáleží na tom ako to napíšete.
5. Nájdite účty, ktoré majú follower_count menší ako 10 a friends_count väčší ako 1000
a výsledok zoraďte podľa statuses_count. Následne spravte jednoduché indexy tak,
aby to malo zmysel a popíšte výsledok.
6. Na predošlú query spravte zložený index a porovnajte výsledok s tým, keď je sú
indexy separátne.
7. Upravte query tak, aby bol follower_count menší ako 1000 a friends_count vačší ako
1000. V čom je rozdiel a prečo?
8. Vyhľadajte všetky tweety (content), ktoré spomenul autor, ktorý obsahuje v popise
(description) reťazec „comedian” (case insensitive), tweety musia obsahovať reťazec
„conspiracy“ (case insensitive), tweety nesmú mať priradený hashtag a počet
retweetov tweetu (retweet_count) je buď menší rovný 10, alebo väčší ako 50.
Zobrazte len rozdielne záznamy a zoraďte ich podľa počtu followerov DESC a pobavte
sa. Mimo to si nad tým spravte analýzu a tú popíšte do protokolu - čo všetko sa tam
deje (explain analyse).

## Odpovede

### 1. Vyhľadajte v tweets.content meno „Gates“ na ľubovoľnom mieste a porovnajte výsledok po tom, ako content naindexujete. V čom je rozdiel a prečo?

Query:
```postgresql
EXPLAIN ANALYZE SELECT * FROM tweets t WHERE t.content LIKE '%Gates%';
```

Result:
```postgresql
                                                           QUERY PLAN
---------------------------------------------------------------------------------------------------------------------------------
 Gather  (cost=1000.00..252509.96 rows=131067 width=251) (actual time=6.170..7534.621 rows=21021 loops=1)
   Workers Planned: 4
   Workers Launched: 1
   ->  Parallel Seq Scan on tweets  (cost=0.00..238403.26 rows=32767 width=251) (actual time=3.612..7214.827 rows=10510 loops=2)
         Filter: (content ~~ '%Gates%'::text)
         Rows Removed by Filter: 3479450
 Planning Time: 1.421 ms
 JIT:
   Functions: 4
   Options: Inlining false, Optimization false, Expressions true, Deforming true
   Timing: Generation 0.953 ms, Inlining 0.000 ms, Optimization 0.392 ms, Emission 5.791 ms, Total 7.136 ms
 Execution Time: 7886.510 ms
(12 rows)
```

Nepoužil sa _Index Scan_ ale _Parallel Sequential Scan_.

Pridáme index pre `content`.
```
CREATE INDEX tweets_content_index ON tweets(content);
```

Query:
```postgresql
EXPLAIN ANALYZE SELECT * FROM tweets t WHERE t.content LIKE '%Gates%';
```

Result:
```postgresql
                                                           QUERY PLAN
---------------------------------------------------------------------------------------------------------------------------------
 Gather  (cost=1000.00..252509.96 rows=131067 width=251) (actual time=5.557..2225.156 rows=21021 loops=1)
   Workers Planned: 4
   Workers Launched: 1
   ->  Parallel Seq Scan on tweets  (cost=0.00..238403.26 rows=32767 width=251) (actual time=4.279..2080.790 rows=10510 loops=2)
         Filter: (content ~~ '%Gates%'::text)
         Rows Removed by Filter: 3479450
 Planning Time: 0.433 ms
 JIT:
   Functions: 4
   Options: Inlining false, Optimization false, Expressions true, Deforming true
   Timing: Generation 0.797 ms, Inlining 0.000 ms, Optimization 0.515 ms, Emission 7.316 ms, Total 8.627 ms
 Execution Time: 2395.448 ms
(12 rows)
```

Po vytvorení obyčajného indexu sa index nepoužil, použil sa `Parallel Sequential Scan`, pretože obyčajný index nevie pracovať s `LIKE` queries, ktoré obsahujú `%%`, na to treba iný index - _trigram index_.

### 2. Vyhľadajte tweet, ktorý začína “DANGER: WARNING:”. Použil sa index?

Query:
```postgresql
EXPLAIN ANALYZE SELECT * FROM tweets t WHERE t.content LIKE 'DANGER: WARNING:%';
```

Result:
```postgresql
                                                         QUERY PLAN
----------------------------------------------------------------------------------------------------------------------------
 Gather  (cost=1000.00..239468.16 rows=649 width=251) (actual time=70.761..926.153 rows=1 loops=1)
   Workers Planned: 4
   Workers Launched: 1
   ->  Parallel Seq Scan on tweets  (cost=0.00..238403.26 rows=162 width=251) (actual time=464.350..890.460 rows=0 loops=2)
         Filter: (content ~~ 'DANGER: WARNING:%'::text)
         Rows Removed by Filter: 3489960
 Planning Time: 0.244 ms
 JIT:
   Functions: 4
   Options: Inlining false, Optimization false, Expressions true, Deforming true
   Timing: Generation 0.809 ms, Inlining 0.000 ms, Optimization 0.520 ms, Emission 5.871 ms, Total 7.200 ms
 Execution Time: 926.648 ms
(12 rows)

```

Index sa nepoužil.

### 3. Teraz naindexujte content tak, aby sa použil index a zhodnoťte prečo sa pred tým nad “DANGER: WARNING:” nepoužil. Použije sa teraz na „Gates“ na ľubovoľnom mieste?

Podľa [dokumentácie](https://www.postgresql.org/docs/9.5/indexes-opclass.html), vytvoríme index tak, aby sa používal index aj pri `LIKE search%` queries. Nepoužijeme pri tom trigram index.

> The operator classes text_pattern_ops, varchar_pattern_ops, and bpchar_pattern_ops support __B-tree__ indexes on the types text, varchar, and char respectively. 
> The difference from the default operator classes is that the values are compared strictly character by character rather than according to the locale-specific collation rules. 
> __This makes these operator classes suitable__ for use by queries involving __pattern matching expressions__ (__LIKE__ or POSIX regular expressions) when the database does not use the standard "C" locale.


Vytvoríme INDEX tak aby sa pri query s `LIKE search%` používal index. 
```postgresql
CREATE INDEX index_content ON tweets (content varchar_pattern_ops);
```

Result:
```postgresql
                                                        QUERY PLAN
--------------------------------------------------------------------------------------------------------------------------
 Index Scan using index_content on tweets  (cost=0.68..8.70 rows=649 width=251) (actual time=5.908..5.928 rows=1 loops=1)
   Index Cond: ((content ~>=~ 'DANGER: WARNING:'::text) AND (content ~<~ 'DANGER: WARNING;'::text))
   Filter: (content ~~ 'DANGER: WARNING:%'::text)
 Planning Time: 46.968 ms
 Execution Time: 6.106 ms
(5 rows)
```

Index sa použil.

Teraz skúsime query s `LIKE %Gates%`.

Query:
```postgresql
EXPLAIN ANALYZE SELECT * FROM tweets t WHERE t.content LIKE '%Gates%';
```

Result:
```postgresql
                                                            QUERY PLAN
----------------------------------------------------------------------------------------------------------------------------------
 Gather  (cost=1000.00..252509.96 rows=131067 width=251) (actual time=40.084..3714.861 rows=21021 loops=1)
   Workers Planned: 4
   Workers Launched: 1
   ->  Parallel Seq Scan on tweets  (cost=0.00..238403.26 rows=32767 width=251) (actual time=13.077..3504.296 rows=10510 loops=2)
         Filter: (content ~~ '%Gates%'::text)
         Rows Removed by Filter: 3479450
 Planning Time: 6.622 ms
 JIT:
   Functions: 4
   Options: Inlining false, Optimization false, Expressions true, Deforming true
   Timing: Generation 6.524 ms, Inlining 0.000 ms, Optimization 0.488 ms, Emission 11.590 ms, Total 18.602 ms
 Execution Time: 4012.166 ms
(12 rows)
```

Použil sa `Parallel Seqential Scan`.
Index sa stále nepoužil, pretože _LIKE_ query s `%search%` alebo s `%search` potrebuje špeciálny - trigram/inverzný index.

### 4. Vytvorte nový index, tak aby ste vedeli vyhľadať tweet, ktorý konči reťazcom „LUCIFERASE“ a nezáleží na tom ako to napíšete.

Najskôr vložíme takú tweet do tabuľky, keďže sa to tam nenachádza. Vložíme jeden string, ktorý sa končí na _LUCIFERASE_ aj na _luciferase_ (tj. malým aj veľkým).

```postgresql
INSERT INTO "tweets"("id","content","location","retweet_count","favorite_count","happened_at","author_id","country_id","parent_id") VALUES (E'1293030235974381573',E'#Christians #CovidHoax #COVID19 #Patriots #Trump #MAGA #WeThePeople Not a hoax - can\'t allow vaccines produced by any of the 9 companies associated with Gates to be used. Twitter dans it Titled - Enzyme that will make Bill Gates’ vaccine microchip implant work known as LUCIFERASE',NULL,0,0,E'2020-08-11 05:43:26+02',898564909,NULL,NULL)
INSERT INTO "tweets" ("id","content","location","retweet_count","favorite_count","happened_at","author_id","country_id","parent_id") VALUES (E'12930302359743815733',E'#Christians #CovidHoax #COVID19 #Patriots #Trump #MAGA #WeThePeople Not a hoax - can\'t allow vaccines produced by any of the 9 companies associated with Gates to be used. Twitter dans it Titled - Enzyme that will make Bill Gates’ vaccine microchip implant work known as luciferase',NULL,0,0,E'2020-08-11 05:43:26+02',898564909,NULL,NULL);
```

Vytvoríme index.

```postgresql
CREATE INDEX index_content_luciferase
ON tweets(content) 
WHERE (tweets.content ILIKE '%LUCIFERASE');
```

Query 1:
```postgresql
EXPLAIN ANALYSE SELECT * FROM tweets t WHERE t.content ILIKE UPPER('%LUCIFERASE');
```

Result:
```postgresql
                                                            QUERY PLAN
-----------------------------------------------------------------------------------------------------------------------------------
 Bitmap Heap Scan on tweets t  (cost=4.29..2501.83 rows=649 width=251) (actual time=0.035..0.072 rows=2 loops=1)
   Recheck Cond: (content ~~* '%LUCIFERASE'::text)
   Heap Blocks: exact=2
   ->  Bitmap Index Scan on index_content_luciferase  (cost=0.00..4.13 rows=649 width=0) (actual time=0.016..0.023 rows=2 loops=1)
 Planning Time: 9.087 ms
 Execution Time: 0.144 ms
(6 rows)
```

Query 2:
```postgresql
EXPLAIN ANALYSE SELECT * FROM tweets t WHERE t.content ILIKE UPPER('%luciferase');
```

Result:

```postgresql
                                                            QUERY PLAN
-----------------------------------------------------------------------------------------------------------------------------------
 Bitmap Heap Scan on tweets t  (cost=4.29..2501.83 rows=649 width=251) (actual time=0.045..0.088 rows=2 loops=1)
   Recheck Cond: (content ~~* '%LUCIFERASE'::text)
   Heap Blocks: exact=2
   ->  Bitmap Index Scan on index_content_luciferase  (cost=0.00..4.13 rows=649 width=0) (actual time=0.021..0.029 rows=2 loops=1)
 Planning Time: 12.682 ms
 Execution Time: 0.539 ms
(6 rows)
```

Query 3:
```postgresql
EXPLAIN ANALYSE SELECT * FROM tweets t WHERE t.content ILIKE UPPER('%lucifERase');
```

Result:
```postgresql
                                                            QUERY PLAN
-----------------------------------------------------------------------------------------------------------------------------------
 Bitmap Heap Scan on tweets t  (cost=4.29..2501.83 rows=649 width=251) (actual time=0.034..0.070 rows=2 loops=1)
   Recheck Cond: (content ~~* '%LUCIFERASE'::text)
   Heap Blocks: exact=2
   ->  Bitmap Index Scan on index_content_luciferase  (cost=0.00..4.13 rows=649 width=0) (actual time=0.015..0.022 rows=2 loops=1)
 Planning Time: 7.007 ms
 Execution Time: 0.141 ms
(6 rows)
```

Nezáleží na tom, ako je _LUCIFERASE_ napísané v query, či malým alebo veľkým. Index sa použije v každom prípade a nájde záznamy, ktoré končia na _LUCIFERASE_ (case insensitive).

### 5. Nájdite účty, ktoré majú follower_count menší ako 10 a friends_count väčší ako 1000 a výsledok zoraďte podľa statuses_count. Následne spravte jednoduché indexy tak, aby to malo zmysel a popíšte výsledok.

Query:
```postgresql
EXPLAIN ANALYZE SELECT * FROM accounts a WHERE a.followers_count < 10 AND a.friends_count > 1000 ORDER BY a.statuses_count;
```

Result:
```postgresql
                                                             QUERY PLAN
------------------------------------------------------------------------------------------------------------------------------------
 Gather Merge  (cost=81774.14..96930.63 rows=126584 width=112) (actual time=396.561..402.175 rows=170 loops=1)
   Workers Planned: 4
   Workers Launched: 1
   ->  Sort  (cost=80774.08..80853.20 rows=31646 width=112) (actual time=335.728..336.553 rows=85 loops=2)
         Sort Key: statuses_count
         Sort Method: quicksort  Memory: 43kB
         Worker 0:  Sort Method: quicksort  Memory: 35kB
         ->  Parallel Seq Scan on accounts  (cost=0.00..76567.59 rows=31646 width=112) (actual time=9.510..334.649 rows=85 loops=2)
               Filter: ((followers_count < 10) AND (friends_count > 1000))
               Rows Removed by Filter: 1734126
 Planning Time: 0.311 ms
 Execution Time: 403.850 ms
(12 rows)
```

Najskôr sa použil parallel _Parallel Sequential Scan_.

Vytvoríme indexy.

```postgresql
CREATE INDEX index_friends_count
    ON accounts (friends_count);

CREATE INDEX index_followers_count
    ON accounts (followers_count);
```

A znovu vykonáme dopyt.

Query:
```postgresql
EXPLAIN ANALYZE SELECT * FROM accounts a WHERE a.followers_count < 10 AND a.friends_count > 1000 ORDER BY a.statuses_count;
```

```postgresql
                                                                            QUERY PLAN
-------------------------------------------------------------------------------------------------------------------------------------------------------------------
 Gather Merge  (cost=81469.26..96684.18 rows=127072 width=112) (actual time=331.715..337.975 rows=170 loops=1)
   Workers Planned: 4
   Workers Launched: 1
   ->  Sort  (cost=80469.20..80548.62 rows=31768 width=112) (actual time=277.299..278.025 rows=85 loops=2)
         Sort Key: statuses_count
         Sort Method: quicksort  Memory: 39kB
         Worker 0:  Sort Method: quicksort  Memory: 39kB
         ->  Parallel Bitmap Heap Scan on accounts  (cost=10551.20..76245.71 rows=31768 width=112) (actual time=55.533..275.537 rows=85 loops=2)
               Recheck Cond: (followers_count < 10)
               Filter: (friends_count > 1000)
               Rows Removed by Filter: 283170
               Heap Blocks: exact=32924
               ->  Bitmap Index Scan on accounts_followers_count_index  (cost=0.00..10519.44 rows=568934 width=0) (actual time=81.599..81.606 rows=566510 loops=1)
                     Index Cond: (followers_count < 10)
 Planning Time: 2.731 ms
 Execution Time: 339.593 ms
(16 rows)
```

Použil sa index `index_followers_count` na ktorom sa vykonal _Bitmap Index Scan_.


### 6. Na predošlú query spravte zložený index a porovnajte výsledok s tým, keď je sú indexy separátne.

Vytvoríme index.

```postgresql
CREATE INDEX index_followers_count_friends_count ON accounts(followers_count, friends_count);
```

Query:
```postgresql
EXPLAIN ANALYZE SELECT * FROM accounts a WHERE a.followers_count < 10 AND a.friends_count > 1000 ORDER BY a.statuses_count;
```

Result:
```postgresql
                                                                          QUERY PLAN
---------------------------------------------------------------------------------------------------------------------------------------------------------------
 Sort  (cost=95600.54..95918.23 rows=127074 width=112) (actual time=72.172..74.037 rows=170 loops=1)
   Sort Key: statuses_count
   Sort Method: quicksort  Memory: 53kB
   ->  Bitmap Heap Scan on accounts a  (cost=11969.54..77442.65 rows=127074 width=112) (actual time=48.019..65.081 rows=170 loops=1)
         Recheck Cond: ((followers_count < 10) AND (friends_count > 1000))
         Heap Blocks: exact=169
         ->  Bitmap Index Scan on index_followers_count_friends_count  (cost=0.00..11937.77 rows=127074 width=0) (actual time=47.948..47.956 rows=170 loops=1)
               Index Cond: ((followers_count < 10) AND (friends_count > 1000))
 Planning Time: 0.388 ms
 Execution Time: 75.891 ms
(10 rows)
```

Query sa vykonala rýchlejšie - __75.891 ms__, v prechádzajúcej úlohe sa vykonávala __339.593 ms__. Nemusel sa použiť filter, pretože v _Recheck Cond_ sa použila podmienka `((followers_count < 10) AND (friends_count > 1000))`. Kompozitný (zložený) index query teda zrýchlil. 

### 7. Upravte query tak, aby bol follower_count menší ako 1000 a friends_count vačší ako 1000. V čom je rozdiel a prečo?

### 8. Vyhľadajte všetky tweety (content), ktoré spomenul autor, ktorý obsahuje v popise (description) reťazec „comedian” (case insensitive), tweety musia obsahovať reťazec „conspiracy“ (case insensitive), tweety nesmú mať priradený hashtag a počet retweetov tweetu (retweet_count) je buď menší rovný 10, alebo väčší ako 50. Zobrazte len rozdielne záznamy a zoraďte ich podľa počtu followerov DESC a pobavte sa. Mimo to si nad tým spravte analýzu a tú popíšte do protokolu - čo všetko sa tam deje (explain analyse).