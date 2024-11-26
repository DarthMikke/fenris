# Fenris

__Uttømmande klimadata for Noreg__

Fenris skal setje saman data frå Meteorologisk institutt for å gi detaljert
oversikt over klimaet ulike stadar i landet. Det finst allereie tenester
som [Se klima](https://seklima.met.no), men dei gir ikkje detaljerte data.
Fenris skal støtte førespørslar som til dømes:

- Gjennomsnittlege lågaste og høgaste temperatur pr månad i løpet av
  ein normalperiode.
- Antal vinterdagar/sumardagar pr månad.
- Korleis klimaet _hadde vore_ dersom temperaturane var lågare/høgare.
- Samanlikning av klimadata ved ulike målestasjonar.

## Implementerte features

- Hent inn og mellomlagre data om målestasjonar.

## Roadmap

- Gjennomsnittleg temperatur, minimums- og maksimumstemperatur ved ein målestasjon __(AKTIV)__
- Finn målestasjonar etter fylke/kommune

## Fenris API

API-et er tilgjengeleg på `https://localhost:8081/api` under utvikling.
Tilgjengelege endepunkt er:

- `s/{stationId}` hentar informasjon om ein målestasjon
- `s/{stationId}/from/{fromYear}/to/{toYear}` hentar ut klimadata for ein
  gitt periode

ID-ar til målestasjonar kan hentast ut frå [Se klima](https://seklima.met.no),
men skal som målsetning vera tilgjengeleg frå Fenris. Eksempelvis er
__SN76917__ Ekofiskfeltet i Nordsjøen, __SN18700__ Blindern, og __SN68173__
Gløshaugen.

## Oppset av utviklingsmiljø

Oppsettet har som føresetnad at du har installert [`ansible-playbook`](https://docs.ansible.com/ansible/latest/installation_guide/intro_installation.html).

```console
$ cd tools
$ ansible-playbook setup.play.yml
```

Dette set opp eit docker-compose-miljø i `tools/dev`. Miljøet består av:

- reverse proxy (Apache)
- Go med autoreload
- Redis

### Køyre API
Eit høveleg skript `tools/compose` kan brukast til å køyre kommandoar frå `tools`:

```console
$ cd tools
$ ./compose up
```

API-et vil vera tilgjengeleg på `https://localhost:8081/api`.

### Køyre frontend for utvikling

```console
$ cd tools
$ npm install
$ npm run dev
```

## Mappeoppsett
- `api` -- kjeldekode for Go-api-et.
- `frontend` -- kjeldekode for frontend.
- `public` -- blir fylt ut automatisk ut med kompilert HTML, Javascript og CSS. Ikkje røyr dette, det vil bli overskrive ved neste kompillering. Alle andre førespørslar enn `/api/*` går til denne mappa.
- `tools` -- verktøy for oppsett av utviklingsmiljø, og på sikt for deployment.
