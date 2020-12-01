# 6. Zadanie: Elastic Search

**Timotej Zaťko**

## Zadanie:

### Časť 1. – analyzéry a mapovanie

Pre index tweets vytvorte 3 vlastné analyzéry (v settings) nasledovne:
1. Analyzér "englando". Tento analyzér bude obsahovať nasledovné: 
    1. fitre: english_possessive_stemmer, lowercase, english_stop, english_stemmer,
    1. char_filter: html_strip
    1. tokenizer: štandardný, ukážku nájdete na stránke elastic.co pre anglický analyzér
1. Analyzér custom_ngram:
    1. Filtre: lowercase, asciifolding, filter_ngrams (definujte si ho sami na rozmedzie 1-10)
    1. char_filter: html_strip
    1. tokenizer: štandardný
1. Analyzér custom_shingles:
    1. Filtre: lowercase, asciifolding, filter_shingles (definujte si ho sami a dajte token_separator: “”)
    1. char_filter: html_strip
    1. tokenizer: štandardný

Do mapovania pridajte:
1. Každý anglický text (rátajme že každý tweet a description u autora je primárne v angličtine) nech je analyzovaný novým analyzérom "englando"
1. Priraďte analýzery
    1. author.name nech má aj mapovania pre custom_ngram, a custom_shingles, 
    1. author. screen_name nech má aj custom_ngram, 
    1. author.description nech má aj custom_shingles. Toto platí aj pre mentions, ak tam tie záznamy máte.

Následne preindexujte dokumenty (ideálne cez reindex API).

### Časť 2. – vyhľadávanie

Vyhľadajte vo vašich tweetoch spojenie "gates s0ros vaccine micr0chip". V query použite function_score, kde jednotlivé medzikroky sú nasledovné:

Query:
1. Must - vyhľadajte vo viacerých poliach (konkrétne: author.name (pomocou shingle), content (cez analyzovaný anglický text), author.description (pomocou shingles), author.screen_name (pomocou ngram)) spojenie "gates s0ros vaccine micr0chip", zapojte podporu pre preklepy, operátor je OR.
1. Tieto polia vo vyhľadávaní boost-nite nasledovne - author.name * 6, content * 8, author.description * 6, author.screen_name * 10.
1. Filter - vyfiltrujte len tie, ktoré majú author.statuses_count > 1000
1. Should – boost-nite 10 krat tie, ktoré obsahujú v mentions.name (tento objekt je typu nested) cez ngram string "real".
1. Nastavte podmienené váhy cez functions nasledovne:
    1. retweet_count, ktorý je väčší rovný ako 100 a menší rovný ako 500 na 6,
    1. author.followers_count väčší ako 100 na 3

Zobrazte agregácie pre výsledky na konci. Vytvorte bucket hashtags podľa hashtagov a spočítajte hodnoty výskytov (na webe by to mohli byť facety). Následne “vyberte” ten, ktorý má najväčší počet a pridajte ho do filtra.

### Časť 3. – agregácie

Trumpovci vs Bidenovci. Pracujte zo všetkými tweetami, ktoré máte a sú vytvorené v auguste (použite tu dátumovú matematiku). Následne pre všetky augustové dni zistite pomocou vnorených agregácii, koľko retweet_count sumárne majú tweety Bidenových podporovateľov a koľko majú tweety Trumpových podporovateľov pre každý jeden deň. Podporovateľ je účet, ktorý má v account_description meno ako frázu jedného z kandidátov (“Joe Biden” a “Donald Trump”). Vo výsledku zhodnoťte kto je aktívnejší a víťaz na tejto sociálnej sieti☺.


## Vypracovanie

### Časť 1. – analyzéry a mapovanie

Pri vytváraní analyzérov musím najskôr zatvoriť index pomocou: **POST** - `http://localhost:9200/tweets-2/_close`. Potom ho môžem otvoriť pomocou: **POST** - `http://localhost:9200/tweets-2/_open` 

1. Vytvoríme analyzér *englando*.

**PUT:** `http://localhost:9200/tweets/_settings`

Body:
```json
{
    "analysis": {
        "analyzer": {
            **"englando": {
                "type": "custom",
                "filter": [
                    "english_possessive_stemmer",
                    "lowercase",
                    "english_stop",
                    "english_stemmer"
                ],
                "char_filter": [
                    "html_strip"
                ],
                "tokenizer": "standard"
            }
        },
        "filter": {
            "english_possessive_stemmer": {
                "type": "stemmer",
                "language": "possessive_english"
            },
            "english_stop": {
                "type": "stop",
                "stopwords": "_english_"
            },
            "english_stemmer": {
                "type": "stemmer",
                "language": "english"
            }
        }
    }
}
```

Response:
```json
{
    "acknowledged": true
}
```

2. Vytovríme analyzér *custom_ngram*.

**PUT:** `http://localhost:9200/tweets/_settings`

Body:
```json
{
    "index": {
        "max_ngram_diff": 9
    },
    "analysis": {
        "analyzer": {
            "custom_ngram": {
                "type": "custom",
                "filter": [
                    "lowercase",
                    "asciifolding",
                    "filter_ngram"
                ],
                "char_filter": [
                    "html_strip"
                ],
                "tokenizer": "standard"
            }
        },
        "filter": {
            "filter_ngram": {
                "type": "ngram",
                "min_gram": 1,
                "max_gram": 10
            }
        }
    }
}
```

Response:
```json
{
    "acknowledged": true
}
```

3. Vytovríme analyzér *custom_shingles*.

Body:
```json
{
    "analysis": {
        "analyzer": {
            "custom_shingles": {
                "type": "custom",
                "filter": [
                    "lowercase",
                    "asciifolding",
                    "filter_shingles"
                ],
                "char_filter": [
                    "html_strip"
                ],
                "tokenizer": "standard"
            }
        },
        "filter": {
            "filter_shingles": {
                "type": "shingle",
                "token_separator": ""
            }
        }
    }
}
```

Response:
```json
{
    "acknowledged": true
}
```

Vytvorené analyzéry použijem v mappingu, nový mapping bude vyzerať nasledovne.

```json
{
  "dynamic": "strict",
  "properties": {
    "content": {
      "type": "text",
      "analyzer": "englando"
    },
    "location": {
      "type": "geo_point"
    },
    "retweet_count": {
      "type": "integer"
    },
    "favorite_count": {
      "type": "integer"
    },
    "happened_at": {
      "type": "date"
    },
    "author": {
      "properties": {
        "id": {
          "type": "long"
        },
        "screen_name": {
          "type": "text",
          "fields": {
            "ngram": {
              "type": "text",
              "analyzer": "custom_ngram"
            }
          }
        },
        "name": {
          "type": "text",
          "fields": {
            "shingles": {
              "type": "text",
              "analyzer": "custom_shingles"
            },
            "ngram": {
              "type": "text",
              "analyzer": "custom_ngram"
            }
          }
        },
        "description": {
          "type": "text",
          "analyzer": "englando",
          "fields": {
            "shingles": {
              "type": "text",
              "analyzer": "custom_shingles"
            }
          }
        },
        "followers_count": {
          "type": "integer"
        },
        "friends_count": {
          "type": "integer"
        },
        "statuses_count": {
          "type": "integer"
        }
      }
    },
    "country": {
      "properties": {
        "code": {
          "type": "keyword"
        },
        "name": {
          "type": "text"
        }
      }
    },
    "hashtags": {
      "type": "text"
    },
    "mentions": {
      "properties": {
        "id": {
          "type": "long"
        },
        "screen_name": {
          "type": "text",
          "fields": {
            "ngram": {
              "type": "text",
              "analyzer": "custom_ngram"
            }
          }
        },
        "name": {
          "type": "text",
          "fields": {
            "shingles": {
              "type": "text",
              "analyzer": "custom_shingles"
            },
            "ngram": {
              "type": "text",
              "analyzer": "custom_ngram"
            }
          }
        }
      }
    },
    "parent_id": {
      "type": "keyword"
    }
  }
}
```

Následne som svoj index reindexoval.

### Časť 2. – vyhľadávanie


### Časť 3. – agregácie
