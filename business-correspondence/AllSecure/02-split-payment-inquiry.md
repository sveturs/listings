# Запрос о split платежах и ответ AllSecure

## Наш запрос

**От:** Dmitrij Vorošilov  
**Кому:** Andjela A  
**Дата:** Июль 2025  

```
Vidimo da imate Preauthorization i Partial Capture funkcije. Da li možete da ih koristite za:
1. Escrow - zadržavanje sredstava do potvrde isporuke
2. Split plaćanja - automatska podela između prodavca/platforme/logistike

Ako je moguće, molimo detaljnu ponudu za marketplace funkcionalnosti.
```

## Ответ AllSecure

**От:** Andjela A  
**Кому:** Dmitrij Vorošilov  
**Дата:** Июль 2025  

```
Poštovani Dmitrij,

Da biste omogućili split transakcije, potrebno je da uvedemo finansijsku instituciju Payspot u proces, koji mogu da Vam omoguće split transakcije. Ukoliko želite, možemo u Vaše ime da prosledimo zahtev Payspotu.

Što se tiče preautorizacije, to je tip transakcije gde se korisniku rezervišu sredstva (koja ne može da koristi više), 

1- sve dok trgovac ne kompletira transakciju nakon čega će sredstva biti preneta na njegov račun (period prenosa zavisi od banke), 

2- ili trgovac otkaže Preautorizaciju iz nekog razloga i korisniku se vrate sva sredstva.

Mi smo payment gateway, nismo finansijska institucija, tako da ne zadržavamo mi novac.

Srdačno,
Andjela A
```

## Ключевые выводы:

1. **AllSecure = Payment Gateway** (техническая обработка)
2. **PaySpot = Financial Institution** (может делать split)
3. **Preauthorization работает как escrow** - блокировка средств до подтверждения
4. AllSecure готов координировать с PaySpot от нашего имени