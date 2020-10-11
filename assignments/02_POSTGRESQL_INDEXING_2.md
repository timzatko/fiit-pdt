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

