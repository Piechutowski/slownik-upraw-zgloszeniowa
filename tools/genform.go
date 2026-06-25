//usr/bin/env go run "$0" "$@"; exit "$?"

// genform.go buduje formularz zgłoszeniowy GitHub (Issue Form) na podstawie
// słownika upraw w pliku data/slownik-upraw.csv. Słownik jest jedynym źródłem
// prawdy, a wygenerowany plik formularza jest artefaktem.
//
// Uruchomienie lokalne (z katalogu głównego repozytorium):
//   ./tools/genform.go
// albo:
//   go run tools/genform.go

package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"sort"
	"strings"
)

// formTemplate zawiera całą strukturę formularza. Jedyny element generowany
// dynamicznie to lista opcji pola "Obszar", wstawiana w miejsce %s.
const formTemplate = `name: Zgłoszenie - Słownik Upraw
description: Internetowy System zgłoszeniowy dla Słownika Upraw
title: "[ZGŁOSZENIE] "
labels: ["zgłoszenie", "faza:1"]
# Aby automatycznie dodawać zgłoszenia do tablicy Projektu, odkomentuj poniższą
# linię i ustaw właściwy projekt, na przykład: ["moja-organizacja/3"].
# projects: ["ORGANIZACJA/NUMER"]
body:
  - type: markdown
    attributes:
      value: |
        Dziękujemy za zgłoszenie do Słownika Upraw.
        Pola oznaczone gwiazdką są wymagane.
        Numer zgłoszenia, data i czas oraz konto zgłaszającego zapisują się automatycznie.
        Faza zgłoszenia: 1 (ustawiana automatycznie, na tym etapie bez możliwości zmiany).
  - type: dropdown
    id: obszar
    attributes:
      label: Obszar (słownik upraw)
      description: Wybierz pozycję ze słownika upraw, której dotyczy zgłoszenie.
      options:
%s
    validations:
      required: true
  - type: input
    id: id-pozycji
    attributes:
      label: ID pozycji w słowniku
      description: Identyfikator konkretnej pozycji, jeśli dotyczy.
    validations:
      required: false
  - type: input
    id: nazwa-pozycji
    attributes:
      label: Nazwa pozycji
      description: Nazwa konkretnej pozycji, jeśli dotyczy.
    validations:
      required: false
  - type: dropdown
    id: typ-uwagi
    attributes:
      label: Typ uwagi
      options:
        - Dodanie
        - Usunięcie
        - Modyfikacja
    validations:
      required: true
  - type: textarea
    id: opis-problemu
    attributes:
      label: Opis problemu
      description: Opisz, czego dotyczy problem.
    validations:
      required: true
  - type: textarea
    id: proponowana-zmiana
    attributes:
      label: Proponowana zmiana
      description: Opisz proponowaną zmianę.
    validations:
      required: true
  - type: textarea
    id: uzasadnienie
    attributes:
      label: Uzasadnienie
      description: Uzasadnij proponowaną zmianę.
    validations:
      required: true
  - type: dropdown
    id: wplyw
    attributes:
      label: Wpływ
      description: Można wskazać więcej niż jeden obszar.
      multiple: true
      options:
        - Statystyka
        - Prawo i administracja
        - Rachunkowość
    validations:
      required: false
  - type: input
    id: email
    attributes:
      label: Adres e-mail zgłaszającego
      description: Pole pomocnicze. Zweryfikowaną tożsamością jest konto GitHub, z którego wysłano zgłoszenie.
      placeholder: jan.kowalski@example.com
    validations:
      required: false
`

const (
	csvPath = "data/slownik-upraw.csv"
	outPath = ".github/ISSUE_TEMPLATE/zgloszenie-slownik-upraw.yml"
)

func main() {
	options, err := wczytajSlownik(csvPath)
	if err != nil {
		fmt.Fprintln(os.Stderr, "błąd:", err)
		os.Exit(1)
	}

	var b strings.Builder
	for _, opt := range options {
		// Cudzysłowy zabezpieczają znaki specjalne YAML, w tym nawiasy.
		safe := strings.ReplaceAll(opt, `"`, `\"`)
		fmt.Fprintf(&b, "        - \"%s\"\n", safe)
	}

	out := fmt.Sprintf(formTemplate, strings.TrimRight(b.String(), "\n"))

	if err := os.WriteFile(outPath, []byte(out), 0o644); err != nil {
		fmt.Fprintf(os.Stderr, "nie można zapisać %s: %v\n", outPath, err)
		os.Exit(1)
	}
	fmt.Printf("Zapisano %s (%d pozycji słownika)\n", outPath, len(options))
}

// wczytajSlownik czyta CSV, waliduje go i zwraca posortowaną listę etykiet
// w formacie "Nazwa (ID)".
func wczytajSlownik(path string) ([]string, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("nie można otworzyć %s: %w", path, err)
	}
	defer f.Close()

	rows, err := csv.NewReader(f).ReadAll()
	if err != nil {
		return nil, fmt.Errorf("odczyt CSV: %w", err)
	}
	if len(rows) < 2 {
		return nil, fmt.Errorf("słownik jest pusty (oczekiwano nagłówka i co najmniej jednej pozycji)")
	}

	seen := map[string]bool{}
	var options []string
	for i, row := range rows {
		if i == 0 {
			continue // nagłówek
		}
		if len(row) < 2 {
			continue
		}
		id := strings.TrimSpace(row[0])
		nazwa := strings.TrimSpace(row[1])
		if id == "" || nazwa == "" {
			continue
		}
		label := fmt.Sprintf("%s (%s)", nazwa, id)
		if label == "None" {
			return nil, fmt.Errorf("pozycja o nazwie None jest zarezerwowana przez GitHub")
		}
		if seen[label] {
			return nil, fmt.Errorf("zduplikowana pozycja: %s", label)
		}
		seen[label] = true
		options = append(options, label)
	}
	if len(options) == 0 {
		return nil, fmt.Errorf("brak prawidłowych pozycji w słowniku")
	}
	sort.Strings(options)
	return options, nil
}
