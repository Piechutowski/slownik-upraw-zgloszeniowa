# Instrukcja: System zgłoszeniowy Słownik Upraw

Ten dokument prowadzi Cię krok po kroku przez uruchomienie i obsługę systemu
zgłoszeniowego do Słownika Upraw, zbudowanego na GitHub Issues. Instrukcja jest
podzielona na trzy części: konfiguracja (raz, na starcie), zgłaszanie uwag oraz
podejmowanie decyzji.

## Co robi ten system

Zgłaszający wypełnia gotowy formularz (typ uwagi, opis, proponowana zmiana,
uzasadnienie, wpływ), a zespół podejmuje decyzję i zostawia ślad ustaleń. Numer
zgłoszenia, data, czas oraz konto autora zapisują się automatycznie.

---

## Część 1: Konfiguracja (wykonujesz raz)

### Krok 1. Utwórz prywatne repozytorium i wypchnij pliki

Ta paczka po rozpakowaniu jest gotowym repozytorium git z pierwszym commitem na
gałęzi `main`. Wystarczy podłączyć je do swojego repozytorium na GitHub.

1. Na GitHub utwórz nowe, **puste** repozytorium (bez README, bez `.gitignore`,
   bez licencji). Ustaw je jako **prywatne**. To jest mechanizm kontroli dostępu,
   tylko zaproszone osoby będą widzieć i tworzyć zgłoszenia.
2. Rozpakuj tę paczkę i wejdź w terminalu do rozpakowanego katalogu.
3. Podłącz swoje repozytorium i wypchnij całość (Fish):

   ```fish
   git remote add origin git@github.com:ORGANIZACJA/REPOZYTORIUM.git
   git push -u origin main
   ```

   Zamień adres na URL swojego repozytorium. Zamiast wersji ssh możesz użyć https,
   na przykład `https://github.com/ORGANIZACJA/REPOZYTORIUM.git`.

### Krok 2. Zaproś osoby

1. Wejdź w Settings, potem Collaborators (lub przypisz zespół, jeśli korzystasz z
   organizacji).
2. Dodaj konta GitHub osób, które mają zgłaszać i obsługiwać uwagi.

Uwaga: GitHub rozpoznaje ludzi po koncie GitHub, a nie po adresie e-mail. Każdy,
kto ma zgłaszać uwagi, musi mieć konto GitHub i zostać zaproszony.

### Krok 3. Wpisz swój słownik upraw

Lista upraw w formularzu pochodzi z pliku `data/slownik-upraw.csv`. Na start są
tam pozycje testowe.

1. Otwórz `data/slownik-upraw.csv`.
2. Zostaw pierwszy wiersz nagłówka (`id,nazwa`).
3. Wpisz swoje pozycje, po jednej w wierszu, w formacie `identyfikator,nazwa`.
   Na przykład:

   ```
   id,nazwa
   PSZ-001,Pszenica ozima
   ZYT-001,Żyto ozime
   ```

4. Zatwierdź zmianę (commit) na gałęzi `main`.

Po zatwierdzeniu uruchomi się automat, który przebuduje formularz tak, aby lista
wyboru odpowiadała Twojemu słownikowi. Automat działa na serwerach GitHuba, nie
na Twoim komputerze, i włącza się tylko przy zmianie słownika. Po chwili na liście
zgłoszeń pojawi się commit „Aktualizacja formularza ze słownika upraw".

Możesz też przebudować formularz na własnym komputerze (wymaga Go), z katalogu
głównego repozytorium:

```fish
go run tools/genform.go
```

### Krok 4. Utwórz etykiety

Wejdź w zakładkę Issues, potem Labels, i utwórz:

- `zgłoszenie`
- `faza:1`

Te dwie etykiety formularz nadaje automatycznie każdemu zgłoszeniu. Opcjonalnie
możesz dodać etykiety pomocnicze do filtrowania: `decyzja:akceptacja`,
`decyzja:odrzucenie`, `decyzja:odłożenie`.

### Krok 5. Skonfiguruj tablicę decyzji (GitHub Projects)

Decyzję podejmuje się po utworzeniu zgłoszenia, dlatego nie ma jej w formularzu.
Służy do tego tablica Projektu z dodatkowymi polami.

1. Utwórz nowy Projekt (widok Table) i powiąż go z repozytorium.
2. Dodaj cztery pola niestandardowe:
   - **Decyzja**: typ Single select, opcje: Akceptacja, Odrzucenie, Odłożenie.
   - **Kto podjął decyzję**: typ Text.
   - **Kiedy podjęto decyzję**: typ Date.
   - **Uzasadnienie decyzji**: typ Text.
3. (Opcjonalnie) Aby nowe zgłoszenia trafiały do Projektu automatycznie, otwórz
   `tools/genform.go`, odkomentuj linię `projects`, wpisz identyfikator w formacie
   `organizacja/numer` i przebuduj formularz (commit pliku CSV uruchomi automat).

Konfiguracja jest gotowa. Reszta instrukcji to codzienne używanie.

---

## Część 2: Jak zgłosić uwagę (dla zgłaszającego)

1. Wejdź w zakładkę **Issues** w repozytorium.
2. Kliknij **New issue**. Otworzy się formularz „Zgłoszenie - Słownik Upraw".
3. Wypełnij pola. Te oznaczone gwiazdką są wymagane:
   - **Obszar (słownik upraw)**: wybierz pozycję z listy. Wymagane.
   - **ID pozycji w słowniku** oraz **Nazwa pozycji**: jeśli dotyczy.
   - **Typ uwagi**: Dodanie, Usunięcie albo Modyfikacja. Wymagane.
   - **Opis problemu**, **Proponowana zmiana**, **Uzasadnienie**: wymagane.
   - **Wpływ**: zaznacz jeden lub więcej obszarów (Statystyka, Prawo i
     administracja, Rachunkowość).
   - **Adres e-mail**: pole pomocnicze.
4. Kliknij **Submit new issue**.

Numer zgłoszenia, data, czas i Twoje konto zapiszą się same. Faza jest na tym
etapie ustawiona na 1 i nie wymaga uzupełniania.

---

## Część 3: Jak podjąć decyzję (dla zespołu)

1. Otwórz zgłoszenie z listy Issues.
2. W tablicy Projektu uzupełnij pola: **Decyzja**, **Kto podjął decyzję**,
   **Kiedy podjęto decyzję** oraz **Uzasadnienie decyzji**.
3. Dla czytelnego śladu dodaj pod zgłoszeniem komentarz według wzoru (autor i czas
   komentarza zapiszą się same):

   ```
   Decyzja: Akceptacja
   Uzasadnienie: <treść uzasadnienia>
   ```

4. Po zakończeniu obsługi zamknij zgłoszenie (Close).

Pola w Projekcie nie wymuszają wypełnienia, więc uzupełnienie uzasadnienia
pilnujesz proceduralnie. Jeśli chcesz twardą blokadę (na przykład brak możliwości
zamknięcia bez uzasadnienia), można ją dołożyć osobno.

---

## Jak zmienić słownik później

Postępujesz tak samo jak w Kroku 3: edytujesz `data/slownik-upraw.csv`, zatwierdzasz
zmianę na `main`, a formularz przebuduje się automatycznie. Pliku
`.github/ISSUE_TEMPLATE/zgloszenie-slownik-upraw.yml` nie edytujesz ręcznie, bo
jest generowany.

Zasady dla pozycji listy: nazwy muszą być unikalne i żadna nie może brzmieć
dokładnie `None` (to słowo jest zarezerwowane przez GitHub).

## Co warto wiedzieć

- Dostęp jest po koncie GitHub, nie po adresie e-mail. Pole e-mail w formularzu to
  tylko deklaracja, a tożsamość pewną jest konto autora zgłoszenia.
- Faza jest na sztywno ustawiona na 1. Obsługę faz od 0 do 6 najwygodniej dodać
  później jako kolejne pole Single select w Projekcie.
- Automat przebudowujący formularz działa na infrastrukturze GitHuba i uruchamia
  się tylko przy zmianie słownika, więc nie obciąża Twoich maszyn.
