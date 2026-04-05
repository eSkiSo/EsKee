package main

import (
	"os"
	"fmt"
	"strings"
	"github.com/tobischo/gokeepasslib/v3"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/lipgloss"
	"github.com/atotto/clipboard"
	"log"
)

const (
	colorWhite  = "\033[39m"
	colorRed    = "\033[91m"
	colorGreen  = "\033[32m"
	colorYellow = "\033[33m"
	colorBlue   = "\033[36m" //was 34
	colorCyan   = "\033[34m" //was 36
	italicStart = "\033[3m"
	boldStart   = "\033[1m"
	styleReset  = "\033[0m"
)

var (
	titleStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("12")).Bold(true)
	passStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("10"))
	notesStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("5"))
	urlStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("202"))
	userStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("201"))
	groupStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("86"))
	helpStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("241"))
)

var version = "0.2 (2026/04/05)"
var databaseFile = getDBPath()

type item struct {
	title    string
	password string
	username string
	url string
	notes string
	group string
}

func (i item) Title() string { return i.group + " › " + i.title }

func (i item) Description() string { return i.username }

func (i item) Url() string { return i.url }

func (i item) Notes() string { return i.notes }

func (i item) Group() string { return i.group }

func (i item) FilterValue() string { return i.title + " " + i.username + " " + i.group }

type EntryWithGroup struct {
	Entry gokeepasslib.Entry
	Group string
}

type model struct {
	list        list.Model
	selected    *item
	showPass    bool
	copied      bool
	width       int
	height      int
}

func newModel(entries []EntryWithGroup) model {
	items := []list.Item{}

	for _, e := range entries {
		items = append(items, item{
			title:    e.Entry.GetTitle(),
			password: e.Entry.GetPassword(),
			username: e.Entry.GetContent("UserName"),
			url: e.Entry.GetContent("URL"),
			notes: e.Entry.GetContent("Notes"),
			group: e.Group,
		})
	}

	l := list.New(items, list.NewDefaultDelegate(), 0, 0)
	l.Title = fmt.Sprintf("[ Kee %s ]", databaseFile)
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(true)
	l.SetShowHelp(true)
	l.SetShowFilter(true)
	l.SetFilterState(list.Filtering)

	return model{list: l}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.list.SetSize(msg.Width, msg.Height-4)

	case tea.KeyMsg:
		switch msg.String() {

		case "q", "ctrl+c":
			return m, tea.Quit

		case "enter":
			if i, ok := m.list.SelectedItem().(item); ok {
				clipboard.WriteAll(i.password)
				m.selected = &i
				return m, tea.Quit
			}

		case "v":
			m.showPass = !m.showPass

		case "c":
			if i, ok := m.list.SelectedItem().(item); ok {
				clipboard.WriteAll(i.password)
				m.copied = true
			}
		}
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m model) View() string {
	if m.selected != nil {
		return fmt.Sprintf(
			"\n[ %s › %s ]\nUser: %s\nPassword: %s\nURL: %s\nNotes:\n%s\n\n\n",
			groupStyle.Render(m.selected.group),
			titleStyle.Render(m.selected.title),
			userStyle.Render(m.selected.username),
			passStyle.Render(m.selected.password),
			urlStyle.Render(m.selected.url),
			notesStyle.Render(m.selected.notes),
		)
	}

	var pass string
	if i, ok := m.list.SelectedItem().(item); ok {
		if m.showPass {
			pass = passStyle.Render(i.password)
		} else {
			pass = "••••••••"
		}
	}

	footer := helpStyle.Render(
		"\n↑↓ navigate • enter select • v toggle pass • c copy • q quit",
	)

	if m.copied {
		footer += helpStyle.Render(" • copied!")
	}

	return "\n" + m.list.View() +
		"\n\nPassword: " + pass +
		footer
}

func main() {
	if len(os.Args) < 2 {
		fmt.Printf("Usage: %s <password> <database file|optional>\n", os.Args[0])
		os.Exit(0)
	}

	password := os.Args[1]

	if password == "-v" || password == "-version" {
		fmt.Printf("[%s%sEsKee%s] Version %s%s%s\n", colorCyan, boldStart, styleReset, colorYellow ,version, styleReset)
		os.Exit(0)
	}

	if(len(os.Args) > 2) {
		databaseFile = os.Args[2] //password db
	}

	file, err := os.Open(databaseFile)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()


    db := gokeepasslib.NewDatabase()
    db.Credentials = gokeepasslib.NewPasswordCredentials(password)
    decoder := gokeepasslib.NewDecoder(file)

    if err = decoder.Decode(db); err != nil {
    	//log.Fatal("Failed to decode DB:", err)
    	fmt.Println(colorRed, "[ERROR] Failed to decode DB, wrong password?", styleReset)
    	os.Exit(0)
    }

    if err = db.UnlockProtectedEntries(); err != nil {
    	fmt.Println(colorRed, "[ERROR] Wrong password or corrupted database", styleReset)
		os.Exit(0)
	}

    entries := collectEntriesWithGroup(db.Content.Root.Groups, "") //collectAllEntries(db.Content.Root.Groups)

    p := tea.NewProgram(newModel(entries))
    finalModel, err := p.Run()
    if err != nil {
    	log.Fatal(err)
    }

    if m, ok := finalModel.(model); ok && m.selected != nil {
    	//fmt.Println(colorBlue, m.selected.title, colorCyan, m.selected.password, styleReset)
    }

}

func findEntriesByTitle(groups []gokeepasslib.Group, query string) []gokeepasslib.Entry {
	var results []gokeepasslib.Entry

	for _, g := range groups {
		// search entries in this group
		for _, e := range g.Entries {
			if strings.Contains(strings.ToLower(e.GetTitle()), strings.ToLower(query)) {
				results = append(results, e)
			}
		}

		// recurse into subgroups
		results = append(results, findEntriesByTitle(g.Groups, query)...)
	}

	return results
}

func collectAllEntries(groups []gokeepasslib.Group) []gokeepasslib.Entry {
	var results []gokeepasslib.Entry

	for _, g := range groups {
		results = append(results, g.Entries...)
		results = append(results, collectAllEntries(g.Groups)...)
	}

	return results
}

func collectEntriesWithGroup(groups []gokeepasslib.Group, parent string) []EntryWithGroup {
	var results []EntryWithGroup

	for _, g := range groups {
		currentPath := g.Name
		if parent != "" {
			currentPath = parent + " › " + g.Name
		}

		// attach group to entries
		for _, e := range g.Entries {
			results = append(results, EntryWithGroup{
				Entry: e,
				Group: currentPath,
			})
		}

		// recurse
		results = append(results,
			collectEntriesWithGroup(g.Groups, currentPath)...)
	}

	return results
}

func getDBPath() string {
	if path := os.Getenv("KEE_DB"); path != "" {
		return path
	}

	home, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}

	return home + "/Database.kdbx"
}