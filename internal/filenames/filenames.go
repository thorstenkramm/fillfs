// Package filenames provides random, human-readable names for files and directories.
package filenames

import (
	cryptorand "crypto/rand"
	"fmt"
	"math/big"
	"math/rand"
	"sync"
	"time"
)

// Names for PDF, Word, or RTF documents
var documentNames = []string{
	"Annual Report",
	"Management Review",
	"Jahresabschluss Entwurf",
	"Raport Roczny",
	"Yillik Degerlendirme",
	"Quarterly Report",
	"Quartalsbericht",
	"Raport kwartalny",
	"Ceyrek Ozeti",
	"Management Review",
	"Fuehrungskreis Protokoll",
	"Notatka dla kierownictwa",
	"Yonetim Kurulu Notu",
	"Meeting Summary",
	"Besprechungsnotizen",
	"Toplanti Ozeti",
	"Podsumowanie spotkania",
	"Project Kickoff Notes",
	"Projektstart Protokoll",
	"Proje Baslangic Notu",
	"Notatki startowe projektu",
	"Contract Draft",
	"Vertragsentwurf",
	"Sozlesme Taslagi",
	"Projekt umowy",
	"Supplier Agreement",
	"Lieferanten Vertrag",
	"Tedarikci Anlasmasi",
	"Umowa dostawcy",
	"HR Policy Update",
	"Personalrichtlinie Update",
	"IK Politika Guncellemesi",
	"Aktualizacja polityki HR",
	"Data Privacy Statement",
	"Datenschutzhinweis",
	"Gizlilik Bildirimi",
	"Oswiadczenie o prywatnosci",
	"Incident Postmortem",
	"Vorfall Nachbetrachtung",
	"Olay Sonrasi Rapor",
	"Raport powypadkowy",
	"Architecture Overview",
	"Architektur Ueberblick",
	"Mimari Genel Bakis",
	"Przeglad architektury",
	"Onboarding Handbook",
	"Einarbeitungsleitfaden",
	"Uyum Kilavuzu",
	"Podrecznik wdrozenia",
	"Support Playbook",
	"Service Leitfaden",
}

// Names for spreadsheets (XLSX, ODS).
var spreadsheetNames = []string{
	"Q1 Sales",
	"Q2 Sales",
	"Q3 Sales",
	"Q4 Sales",
	"Umsatz Q1",
	"Umsatz Q2",
	"Umsatz Q3",
	"Umsatz Q4",
	"Sprzedaz Q1",
	"Sprzedaz Q2",
	"Sprzedaz Q3",
	"Sprzedaz Q4",
	"Satis Q1",
	"Satis Q2",
	"Satis Q3",
	"Satis Q4",
	"Sales Forecast",
	"Umsatz Prognose",
	"Prognoza sprzedazy",
	"Satis Tahmini",
	"Expense Tracker",
	"Kosten Uebersicht",
	"Gider Listesi",
	"Lista wydatkow",
	"Budget Plan",
	"Budget Plan DE",
	"Butce Plani",
	"Plan budzetowy",
	"Payroll Sheet",
	"Lohnabrechnung",
	"Maas Bordrosu",
	"Lista plac",
	"Inventory List",
	"Lagerbestand",
	"Stok Listesi",
	"Lista magazynowa",
	"Marketing Spend",
	"Marketing Kosten",
	"Pazarlama Harcamasi",
	"Wydatki marketingowe",
	"Hiring Funnel",
	"Bewerber Pipeline",
	"Aday Havuzu",
	"Lejek rekrutacyjny",
	"OKR Tracking",
	"OKR Nachverfolgung",
	"OKR Takibi",
	"Monitorowanie OKR",
	"Risk Register",
	"Rejestr ryzyka",
}

// Names for images.
var imageNames = []string{
	"City Sunset",
	"Stadt Sonnenuntergang",
	"Sehir Gunes Batimi",
	"Zachod slonca miasto",
	"Mountain Trail",
	"Bergpfad",
	"Doga Yolu",
	"Gorski szlak",
	"Forest Path",
	"Waldweg",
	"Orman Patikasi",
	"Lesna sciezka",
	"Ocean Waves",
	"Ozean Wellen",
	"Okyanus Dalgasi",
	"Fale oceanu",
	"Desert Dunes",
	"Wueste Duene",
	"Col Tepesi",
	"Wydmy pustynne",
	"Winter Cabin",
	"Winterhuette",
	"Kisin Dag Evi",
	"Zimowa chata",
	"City Market",
	"Stadt Markt",
	"Sehir Pazari",
	"Targ miejski",
	"Street Food",
	"Strassen Essen",
	"Sokak Yemegi",
	"Jedzenie uliczne",
	"Night Skyline",
	"Nacht Skyline",
	"Gece Silueti",
	"Nocna panorama",
	"Flower Field",
	"Blumenwiese",
	"Cicek Tarlasi",
	"Lakowe kwiaty",
	"Drone View",
	"Drohnenblick",
	"Drone Gorunum",
	"Widok z drona",
	"Product Hero",
	"Produkt Foto",
	"Urun Kapagi",
	"Zdjecie produktu",
	"UI Mockup",
	"Arayuz Taslagi",
	"Makieta UI",
}

// Names for sound files.
var soundNames = []string{
	"Ambient Rain",
	"Stadtregen",
	"Sehir Yagmuru",
	"Deszcz w miescie",
	"Forest Birds",
	"Wald Voegeln",
	"Orman Kuslari",
	"Ptaki w lesie",
	"Ocean Surf",
	"Meeresrauschen",
	"Okyanus Dalga Sesi",
	"Szum oceanu",
	"Campfire Crackle",
	"Lagerfeuer",
	"Kamp Atasi",
	"Ognisko trzask",
	"Night Crickets",
	"Nacht Grillen",
	"Gece Bocekleri",
	"Swierzcze nocne",
	"Cafe Murmur",
	"Kaffeehaus Geraeusch",
	"Kafe Ugrultusu",
	"Szum kawiarni",
	"Office Floor",
	"Buero Hintergrund",
	"Ofis Zemin Sesi",
	"Biuro w tle",
	"Keyboard Typing",
	"Tastatur Klick",
	"Klavye Tiklama",
	"Stukanie klawiatury",
	"Train Ride",
	"Zugfahrt",
	"Tren Yolculugu",
	"Podroz pociagiem",
	"Airport Lounge",
	"Flughafen Lounge",
	"Havalimani Salonu",
	"Poczekalnia lotnisko",
	"Drum Groove",
	"Schlagzeug Groove",
	"Davul Ritim",
	"Rytm perkusji",
	"Guitar Riff",
	"Gitarren Riff",
	"Gitar Rifi",
	"Riff gitary",
	"Podcast Intro",
	"Podcast Wstep",
}

// Names for slide decks (PPT).
var powerpointNames = []string{
	"Quarterly Business Review",
	"Quartalsreview",
	"Ceyrek Degerlendirme",
	"Przeglad kwartalny",
	"Board Update",
	"Vorstand Update",
	"Yonetim Kurulu Sunumu",
	"Aktualizacja zarzadu",
	"Strategy Offsite",
	"Strategie Workshop",
	"Strateji Calistayi",
	"Warsztat strategii",
	"Product Vision",
	"Produkt Vision",
	"Urun Vizyonu",
	"Wizja produktu",
	"Release Plan",
	"Release Plan DE",
	"Surum Plani",
	"Plan wydania",
	"Sales Kickoff",
	"Vertriebs Kickoff",
	"Satis Baslangici",
	"Start sprzedazy",
	"Marketing Plan",
	"Marketing Plan DE",
	"Pazarlama Plani",
	"Plan marketingowy",
	"Brand Guidelines",
	"Markenrichtlinien",
	"Marka Rehberi",
	"Wytyczne marki",
	"Investor Pitch",
	"Investoren Pitch",
	"Yatirimci Sunumu",
	"Prezentacja inwestorska",
	"Onboarding Deck",
	"Einarbeitung Folien",
	"Uyum Sunumu",
	"Slajdy wdrozenia",
	"Security Awareness",
	"Sicherheits Schulung",
	"Guvenlik Egitimi",
	"Szkolenie bezpieczenstwa",
	"Architecture Deep Dive",
	"Architektur Deepdive",
	"Mimari Derin Inceleme",
	"Szczegoly architektury",
	"Town Hall Slides",
	"All Hands Folien",
}

// Names for directories (English only).
var directoryNames = []string{
	"alpha-team",
	"beta-team",
	"gamma-team",
	"delta-team",
	"project-atlas",
	"project-beacon",
	"project-comet",
	"project-delta",
	"project-ember",
	"project-falcon",
	"project-galaxy",
	"project-harbor",
	"project-ionic",
	"project-jade",
	"project-keystone",
	"project-lighthouse",
	"project-meridian",
	"project-northstar",
	"project-orbit",
	"project-pioneer",
	"project-quartz",
	"project-ridge",
	"project-summit",
	"project-trail",
	"project-umbra",
	"project-voyager",
	"project-willow",
	"project-xenon",
	"project-yonder",
	"project-zephyr",
	"archive-2021",
	"archive-2022",
	"archive-2023",
	"archive-2024",
	"assets-shared",
	"assets-brand",
	"assets-product",
	"assets-raw",
	"assets-processed",
	"backups-daily",
	"backups-weekly",
	"backups-monthly",
	"compliance-reports",
	"customer-success",
	"customer-feedback",
	"design-explorations",
	"design-system",
	"dev-tools",
	"docs-guides",
	"docs-reference",
	"docs-api",
	"engineering-handbook",
	"environment-dev",
	"environment-staging",
	"environment-prod",
	"experiments-a",
	"experiments-b",
	"feature-flags",
	"financials-q1",
	"financials-q2",
	"financials-q3",
	"financials-q4",
	"growth-ideas",
	"growth-tests",
	"handovers",
	"incident-reports",
	"infra-terraform",
	"infra-kubernetes",
	"infra-scripts",
	"legal-contracts",
	"legal-templates",
	"logs-application",
	"logs-audit",
	"logs-security",
	"marketing-campaigns",
	"marketing-creative",
	"metrics-dashboards",
	"metrics-exports",
	"ops-runbooks",
	"ops-schedules",
	"ops-checklists",
	"partners",
	"planning-q1",
	"planning-q2",
	"planning-q3",
	"planning-q4",
	"product-discovery",
	"product-research",
	"product-specs",
	"qa-cases",
	"qa-results",
	"releases",
	"roadmap-archive",
	"sales-collateral",
	"sales-enablement",
	"security-audits",
	"support-guides",
	"team-photos",
	"training-materials",
	"user-interviews",
	"workshops",
}

// Suffixes for directories (English only)
var directorySuffixes = []string{
	"approved",
	"rejected",
	"archived",
	"secret",
	"public",
	"top secret",
	"reviewed",
	"to be deleted",
	"duplicated",
	"staging",
	"production",
	"development",
	"testing",
	"legacy",
	"deprecated",
	"experimental",
	"hotfix",
	"backup",
	"cold-storage",
	"long-term",
	"short-term",
	"internal",
	"external",
	"shared",
	"private",
	"readonly",
	"archivable",
	"in-progress",
	"completed",
	"pending",
	"final",
	"draft",
	"wip",
	"candidate",
	"beta",
	"alpha",
	"release",
	"release-candidate",
	"stable",
	"unstable",
	"nightly",
	"daily",
	"weekly",
	"monthly",
	"quarterly",
	"yearly",
	"historical",
	"current",
	"future",
	"mirrored",
	"synced",
	"offline",
	"online",
	"locked",
	"unlocked",
	"verified",
	"unverified",
	"clean",
	"dirty",
}

var (
	rnd        *rand.Rand
	rndOnce    sync.Once
	separators = "._- +="
	nameCount  uint64
	nameMu     sync.Mutex
	startDate  = time.Date(1974, 4, 25, 0, 0, 0, 0, time.UTC)
)

func initRand() {
	rndOnce.Do(func() {
		seed, err := cryptorand.Int(cryptorand.Reader, big.NewInt(1<<62))
		if err != nil {
			rnd = rand.New(rand.NewSource(time.Now().UnixNano())) //nolint:gosec // fallback when crypto fails
			return
		}
		rnd = rand.New(rand.NewSource(seed.Int64())) //nolint:gosec // seeded from crypto/rand
	})
}

// RandomDirectoryName builds a random directory name composed of a padded number, separator,
// base name, separator, and suffix.
func RandomDirectoryName() string {
	initRand()

	number := fmt.Sprintf("%03d", rnd.Intn(1000))
	sep1 := string(separators[rnd.Intn(len(separators))])
	sep2 := string(separators[rnd.Intn(len(separators))])
	name := directoryNames[rnd.Intn(len(directoryNames))]
	suffix := directorySuffixes[rnd.Intn(len(directorySuffixes))]

	return number + sep1 + name + sep2 + suffix
}

// RandomDocumentFileName returns a random document-style name with optional date suffix.
func RandomDocumentFileName() string {
	return randomFileNameFrom(documentNames)
}

// RandomSpreadsheetFileName returns a random spreadsheet-style name with optional date suffix.
func RandomSpreadsheetFileName() string {
	return randomFileNameFrom(spreadsheetNames)
}

// RandomImageFileName returns a random image-style name with optional date suffix.
func RandomImageFileName() string {
	return randomFileNameFrom(imageNames)
}

// RandomSoundFileName returns a random sound-style name with optional date suffix.
func RandomSoundFileName() string {
	return randomFileNameFrom(soundNames)
}

// RandomPowerpointFileName returns a random slide-deck name with optional date suffix.
func RandomPowerpointFileName() string {
	return randomFileNameFrom(powerpointNames)
}

func randomFileNameFrom(items []string) string {
	initRand()

	base := items[rnd.Intn(len(items))]
	sep1 := string(separators[rnd.Intn(len(separators))])
	version := fmt.Sprintf("v%d", rnd.Intn(25)+1)

	if shouldAppendDate() {
		sep2 := string(separators[rnd.Intn(len(separators))])
		return base + sep1 + version + sep2 + wrapDate(randomDate())
	}

	return base + sep1 + version
}

func shouldAppendDate() bool {
	nameMu.Lock()
	defer nameMu.Unlock()
	nameCount++
	return nameCount%3 == 0
}

func randomDate() string {
	now := time.Now().UTC()
	span := now.Unix() - startDate.Unix()
	if span <= 0 {
		return startDate.Format("2006-01-02")
	}
	seconds := rnd.Int63n(span+1) + startDate.Unix()
	return time.Unix(seconds, 0).UTC().Format("2006-01-02")
}

func wrapDate(date string) string {
	switch rnd.Intn(4) { // 0: none, 1: (), 2: {}, 3: []
	case 1:
		return "(" + date + ")"
	case 2:
		return "{" + date + "}"
	case 3:
		return "[" + date + "]"
	default:
		return date
	}
}
