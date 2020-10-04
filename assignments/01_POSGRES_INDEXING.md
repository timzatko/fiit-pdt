# Zadanie: Postgres - Indexácia

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

1. Plánovač vybral __Parallel Seqential Scan__. Query: `EXPLAIN SELECT * FROM accounts a WHERE a.screen_name='realDonaldTrump' LIMIT 1;`
2. __2__. Čím viac workerov, tak tým rýchlejšie sa query vykonaná. Čas vykonávania podľa počtu nastavených workerov:
   * 2: 600ms - 900ms
   * 3: 400ms - 700ms
   * 4: 400ms - 500ms
   * 5 > plánovač vybral maximálne 4 workery (tj. strop použitých workerov)
3. Plánovač nepotrebuje viac workerov, na riadky nebol aplikovaný žiaden filter pretože sa používa index (ktorý používa vyhľadávací strom). V indexe sa vyhľadáva na základe `screen_name`, keď sa nájde v indexe (vyhľadávacom strome) `screen_name=realDonaldTrump`, daný node obsahuje smerník na daný riadok, preto nie je potrebný filter.
4. Plánovač vybral __Seqential Scan__, prečo? pretože vyhodnotil, že použitie sekvenčného skenu je výhodnejšie ako použitie paralelného sekvenčného skenu. Napríklad kvôli tomu, že pri paralelnom sekvenčnom skene je potrebná synchronizácia medzi workermi, v tomto prípade vyhodnotil, že by synchronizácia medzi nimi bola drahá. Čo sa týka správania oproti druhej úlohe, keďže bol použitý iba sekvenčný sken, nie je sa s čím porovnávať. 
5. Query:`SELECT * FROM accounts a WHERE a.followers_count>=100 and a.followers_count<=200;`. __Bitmap Heap Scan__ najskôr prejde index a vyberie tie bloky z indexu ktoré splňujú (existuje zánam/y z tých blokov, ktorý splňuje...) podmienku (followers_count>=100 and followers_count<=200) - to je recheck condition. Potom __Bitmap Index Scan__ zoberie list týchto blokov a pre každý z nich sa pozerá do tabuľky, a z nich vyberá záznamy, ktoré spĺňajú stanovenú podmienky. [[ref]](https://www.cybertec-postgresql.com/en/postgresql-indexing-index-scan-vs-bitmap-scan-vs-sequential-scan-basics/?gclid=Cj0KCQjw5eX7BRDQARIsAMhYLP9eCopBtmKbz8ed_s3LKcQpV-YLTArz0E5NKGYzgTkdDfhsvLHoszEaAuvNEALw_wcB) 
6. V tomto prípade plánovač uprednostnil __Seqential Scan__, pravdepodobne vyberáme veľkú časť riadkov tabuľy a __Bitmap Scan__ by už nebol asi taký efektívny - tj. vybral by taký počet blokov, ktorý by bol blízko k celkovému počtu blokov, preto v tomto prípade plánovač vybral __Seqential Scan__. 
7. Query: `INSERT INTO accounts ("id", "screen_name", "name", "description", "followers_count", "friends_count", "statuses_count") VALUES (1, 'timzatko', 'timzatko', 'Hello!', 1000, 10000, 99999)`. Vloženie nového riadku trvalo _39.00s_. Keď som dropol indexy, trvalo to kratšie pretože pri vkladaní sa nemusel aktualizovať index - až tri indexy. Aktualizácia indexu môže je náročná operácia (napr. v prípade stromových indexov je nutné ich niekedy vyvážiť).
8. Dĺžka vytvárania indexu pre _retweet_count_ bola 58.50s a pre _content_ sa index vytváral 59.75s. Rozdiel je tu len minimálny, očakával by som, že rozdiel bude väčší, kvôli tomu, že _content_ obsahuje viac rôznych hodnôt a strom bude v tom prípade hlbší. Avšak čas bol rovnaký, pretože pri vytváraní indexu museli byť aj tak v oboch prípadoch prejdené všetky riadky tabuľky. TODO:    
9. TODO: