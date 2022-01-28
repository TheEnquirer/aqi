package main

import (
    "fmt"
    "os"
    "io"
    //"os/exec"
    "log"
    "time"
    "net/http"
    //import "encoding/json"

    tea "github.com/charmbracelet/bubbletea"
)
import "encoding/json"

type model struct {
    choices  []string           // items on the to-do list
    cursor   int                // which to-do list item our cursor is pointing at
    selected map[int]struct{}   // which to-do items are selected
}


func initialModel() model {
	return model{
		// Our shopping list is a grocery list
		choices:  []string{"Buy carrots", "Buy celery", "Buy kohlrabi"},

		// A map which indicates which choices are selected. We're using
		// the  map like a mathematical set. The keys refer to the indexes
		// of the `choices` slice, above.
		selected: make(map[int]struct{}),
	}
}


func (m model) Init() tea.Cmd {
    // Just return `nil`, which means "no I/O right now, please."
    return nil
}

var myClient = &http.Client{Timeout: 10 * time.Second}

type Foo struct {
    Bar string
}

func getJson(url string, target interface{}) error {
    r, err := myClient.Get(url)
    if err != nil {
        return err
    }
    defer r.Body.Close()

    return json.NewDecoder(r.Body).Decode(target)
}

type DataStruct struct {
    City string      `json:"'city'"`
    State string      `json:"'state'"`
    Country string      `json:"'country'"`
    //Pollution PollutionStruct `json:"'pollution'"`
    //Pollution string `json:"'pollution'"`
    Location LocationStruct `json:"'location'"`
    Current CurrentStruct `json:"'current'"`
}

type CurrentStruct struct {
    Pollution PollutionStruct      `json:"'pollution'"`
}

type PollutionStruct struct {
    Aqius int      `json:"'aqius'"`
    //Mainus string      `json:"'mainus'"`
}

type LocationStruct struct {
    Type string      `json:"'type'"`
}

type Request struct {
    Operation string      `json:"operation"`
    Key string            `json:"key"`
    Value string          `json:"value"`
    Status string         `json:"'success'"`
    Data DataStruct       `json:"'data'"`
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {

    // Is it a key press?
    case tea.KeyMsg:

        // Cool, what was the actual key pressed?
        switch msg.String() {

        // These keys should exit the program.
        case "ctrl+c", "q":
            return m, tea.Quit

	case "l":
	    //temp := ""
	    resp, err := http.Get("http://api.airvisual.com/v2/city?city=San%20Mateo&state=California&country=USA&key=5b78dd52-e51a-47f6-b63f-d580b4e33b83")
	    if err != nil {
		log.Fatalln(err)
		// handle err
	    }
	    //fmt.Println(resp.Header)
	    body, err := io.ReadAll(resp.Body)
	    //fmt.Println(string(body))
	    defer resp.Body.Close()
	    //JsonString := string(body)
	    //json.Unmarshal([]byte(JsonString), temp)
	    //fmt.Println(temp)

	    //s := string("{'operation': 'get', 'key': 'example'}")
	    s := string(body)
	    data := Request{}

	    eerr := json.Unmarshal([]byte(s), &data)
	    if eerr != nil {
		fmt.Println(eerr.Error()) 
		//json: Unmarshal(non-pointer main.Request)
	    }

	    //fmt.Println(data)
	    //fmt.Printf("Operation: %s", data.Status)
	    //fmt.Printf("Operation: %s", data.Status)

	    //fmt.Printf("Operation: %s", data.Data.Pollution.Aqius)
	    //fmt.Printf("Operation: %s", data.Data.Location.Type)
	    aqi := data.Data.Current.Pollution.Aqius
	    fmt.Println(aqi)

	    //fmt.Printf(data.Data)
	    //fmt.Println(
	    //err := json.Unmarshal([]byte(s), &data)
	    //if err != nil {
	    //    fmt.Println(err.Error()) 
	    //    //invalid character '\'' looking for beginning of object key string
	    //}

        // The "up" and "k" keys move the cursor up
        case "up", "k":
            if m.cursor > 0 {
                m.cursor--
            }

        // The "down" and "j" keys move the cursor down
        case "down", "j":
            if m.cursor < len(m.choices)-1 {
                m.cursor++
            }

        // The "enter" key and the spacebar (a literal space) toggle
        // the selected state for the item that the cursor is pointing at.
        case "enter", " ":
            _, ok := m.selected[m.cursor]
            if ok {
                delete(m.selected, m.cursor)
            } else {
                m.selected[m.cursor] = struct{}{}
            }
        }
    }


    // Return the updated model to the Bubble Tea runtime for processing.
    // Note that we're not returning a command.
    return m, nil
}

func (m model) View() string {
    // The header
    s := "What should we buy at the market?\n\n"

    // Iterate over our choices
    for i, choice := range m.choices {

        // Is the cursor pointing at this choice?
        cursor := " " // no cursor
        if m.cursor == i {
            cursor = ">" // cursor!
        }

        // Is this choice selected?
        checked := " " // not selected
        if _, ok := m.selected[i]; ok {
            checked = "x" // selected!
        }

        // Render the row
        s += fmt.Sprintf("%s [%s] %s\n", cursor, checked, choice)
    }

    // The footer
    s += "\nPress q to quit.\n"

    // Send the UI for rendering
    return s
}


func main() {
    foo1 := new(Foo) // or &Foo{}
    getJson("http://api.airvisual.com/v2/city?city=San%20Mateo&state=California&country=USA&key=5b78dd52-e51a-47f6-b63f-d580b4e33b83", foo1)
    println(foo1.Bar)
    fmt.Print(foo1.Bar, "whee")
    // alternately:

    //foo2 := Foo{}
    //getJson("http://example.com", &foo2)

    p := tea.NewProgram(initialModel())
    if err := p.Start(); err != nil {
        fmt.Printf("Alas, there's been an error: %v", err)
        os.Exit(1)
    }
}

