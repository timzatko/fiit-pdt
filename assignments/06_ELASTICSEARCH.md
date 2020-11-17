# 6. Zadanie: Elastic Search

**Timotej Zaťko**

## Do kedy?

Keďže sa preklápame do času, kde toho začína byť veľa aj z pohľadu ostatným predmetov a diplomoviek, toto aj ďalšie zadanie (obe sú Elasticsearch) bude mať spoločný čas odovzdania a to je 4.12. 23:59.

## Zadanie č.1 – práca s indexom a dokumentami v Elasticsearch

1. Rozbehajte si 3 inštancie Elasticsearch-u
2. Vytvorte index pre Tweety, ktorý bude mať “optimálny“ počet shardov a replík pre 3 nody
(aby tam bola distribúcia dotazov vo vyhľadávaní, aj distribúcia uložených dát)
3. Vytvorte mapping pre normalizované dáta z Postgresu - Tweet musí obsahovať údaje rovnaké ako máte už uložené v PostgreSQL. Dbajte na to, aby ste vytvorili polia v správnom dátovom type (polia ktoré má zmysel analyzovať analyzujte správne, tie ktoré nemá, aby neboli zbytočne analyzované (keyword analyzer)) tak aby index nebol zbytočne veľký. Mapovanie musí byť striktné.
4. Vytvorte bulk import pre vaše normalizované Tweety.
5. Importujete dáta do Elasticsearchu
6. Experimentujte s nódami, a zistite koľko nódov musí bežať (a ktoré) aby vám Elasticsearch vedel pridávať dokumenty, mazať dokumenty, prezerať dokumenty a vyhľadávať nad nimi?
7. Upravujte počet retweetov pre vami vybraný tweet pomocou vašeho jednoduchého
scriptu (v rámci Elasticsearchu) a sledujte ako sa mení _seq_no a_primary_term pri tom ako zabíjate a spúšťate nódy.

### Hodnotenie

Zadanie č.1 - hodnotenie - dokopy 7,5 boda rozdelené percentuálne:

- 2 - prečo ste vytvorili index tak ako ste vytvorili? / 10%
- 3 - ako ste vytvorili mapovanie? / 25%
- 4 - naindexovali ste správny počet dokumentov (buď vaše číslo z PostgreSQL, alebo správny
počet)? / 25%
- 5 - napísali ste skript tak aby dával zmysel (je ten kód �� alebo ��)? / 15%
- 6 + 7 napíšte v skratke aké prípady ste si definovali a prečo. Popísali ste výstupy
zmysluplne? Tj. oťukali ste si ako to funguje? / 25%

## Odpovede
