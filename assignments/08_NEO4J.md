# 8. Zadanie: Neo4j

**Timotej Zaťko**

## Zadanie

### Do kedy?

Keďže všetci robíte radi s predstihom, deadline je 24.12.2020 o 19:30. A teda nech vám kapor moc nevychladne ...
Čo odovzdať?

Dotazy v separátnom textovom editore + dotaz a screenshot výsledku v pdf.

### Hodnotenie
Úlohy dokopy 10 bodov rozdelené nasledovne:
1,2,3,4,5 a 6 sú za jeden bod,
7, 8 za 2 body.
Bonusová úloha je za 3 body.

### Zadanie – Neo4j
1. Zoberte si nami vytvorený dataset z vašich tweetov: https://drive.google.com/file/d/1nyVwjhbQhT3kbMlU3mXCmJ-_uQ8Yx9B8/view?usp=sharing a importujte ho cez neo4j-admin
2. Vypíšte 10 Accountov s najvyšším množstvom followerov. Nezoraďujte Accounty podľa poľa followers_count, hodnota je prevzatá z Twitteru a nezodpovedá našim vzťahom v datasete. Zaujímajú nás followujúce Accounty v našom datasete.
3. Nájdite najkratšie cesty medzi Katy Perry {screen_name: ‘katyperry‘} a Donaldom Trumpom {screen_name: ‘realDonaldTrump‘} cez vzťah FOLLOWS. Všetky cesty, kde Donald Trump followuje niekoho, kto followuje niekoho, kto..., kto followuje Katy Perry.
4. Vyhľadajte zlyhania influencerov. Vyhľadajte 20 najmenej retweetovanych tweetov od Accountov, ktoré sú na prvých 10 miestach v celkovom počte retweetov.
5. Vytvorte volaním iba jednej query nový Account s Vašim menom, ktorý bude followovať Donalda Trumpa {screen_name:"realDonaldTrump"} a v tom istom volaní vytvorte tweet, ktorý bude retweetom Donaldovho najretweetovanejšieho tweetu.
6. Odporučte používateľovi {screen_name:"777stl"} followovanie ďalších Accountov, na základe followovania rovnakých Accountov: Vyhľadajte 10 Accountov, ktoré followujú najviac rovnakých Accountov, ale Donald ich ešte nefollowuje.
7. Odporučte používateľovi {screen_name:"DaynerWilson"} followovanie ďalších Accountov na základe zhody v retweetovaní rovnakých tweetov: Vyhľadajte 10 accountov, ktoré retweetli najviac tweetov rovnakych, ako náš používateľ. Počítajú sa aj retweety tweetov, ktoré retweetujú rovnaký tweet.
8. Vyhľadajte 10 tweetov ostatných Accountov, ktoré do hĺbky 5 followujú account, ktorý napísal tweet {id: "1289380305728503808"}, ktoré síce nie sú retweetom vybraného tweetu, ale napriek tomu majú čo najviac rovnakých slov v poli content zhodných s vybraným tweetom (odporúčam pozrieť si split() a procedúry v APOC pracujúce s collections ). Account,  ktorý followuje account, ktorý followuje nami vybraný Account rozumieme hĺbkou 2.

BONUS: Nájdite najkratšie cesty medzi Katy Perry (katyperry) a Donaldom Trumpom cez vzťah RETWEETS (a tým pádom aj POST). Všetky cesty, kde Katy Perry retweetla post Accountu, ktorý retweetol post Accountu, ktorý..., ktorý retweetol post Donalda Trumpa


## Vypracovanie

...