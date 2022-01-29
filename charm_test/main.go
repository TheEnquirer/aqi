package main

import (
    "fmt"
    "os"
    "io"
    //"os/exec"
    "log"
    //"time"
    "net/http"
    //import "encoding/json"
    "github.com/erikgeiser/promptkit/selection"
    "strings"

    tea "github.com/charmbracelet/bubbletea"
)
import "encoding/json"

type Request struct {
    Operation string      `json:"operation"`
    Key string            `json:"key"`
    Value string          `json:"value"`
    Status string         `json:"'success'"`
    Data DataStruct       `json:"'data'"`
}

type DataStruct struct {
    City string      `json:"'city'"`
    State string      `json:"'state'"`
    Country string      `json:"'country'"`
    Location LocationStruct `json:"'location'"`
    Current CurrentStruct `json:"'current'"`
}

type LocationStruct struct {
    Type string      `json:"'type'"`
}

type CurrentStruct struct {
    Pollution PollutionStruct      `json:"'pollution'"`
}

type PollutionStruct struct {
    Aqius int      `json:"'aqius'"`
}

///////////////////////

type CityRequest struct {
    Status string `json:"'status'"`
    //Data CityDataStruct `json:"'data'"`
    Data []CityNameStruct `json:"'data'"`
}

type CityDataStruct struct {

}

type CityNameStruct struct {
    City string   `json:"'city'"`
}

type model struct {
    choices  []string           // items on the to-do list
    //cursor   int                // which to-do list item our cursor is pointing at
    //selected map[int]struct{}   // which to-do items are selected
    aqi string
}

func initialModel(aqii int) model {
    //fmt.Println(aqii, "the aqi!!")
    return model{
	aqi: string(aqii),
	// Our shopping list is a grocery list
	choices:  []string{"Buy carrots", "Buy celery", "Buy kohlrabi"},

	// A map which indicates which choices are selected. We're using
	// the  map like a mathematical set. The keys refer to the indexes
	// of the `choices` slice, above.
	//selected: make(map[int]struct{}),
    }
}


func getWeather(city *selection.Choice) int {
    req := fmt.Sprintf("http://api.airvisual.com/v2/city?city=%s&state=California&country=USA&key=5b78dd52-e51a-47f6-b63f-d580b4e33b83", strings.ReplaceAll(city.Value.(string), " ", "%20"))

    resp, err := http.Get(req) // get the resp
    if err != nil {
	log.Fatalln(err)
    }
    body, err := io.ReadAll(resp.Body) 
    defer resp.Body.Close()

    s := string(body)
    data := Request{}

    eerr := json.Unmarshal([]byte(s), &data)
    if eerr != nil {
	fmt.Println(eerr.Error()) 
    } // process it into a request

    aqi := data.Data.Current.Pollution.Aqius // navigate the tree and get out aqi 
    //fmt.Println(aqi)
    return aqi
}


func getCities() []string {
    resp, err := http.Get("http://api.airvisual.com/v2/cities?state=California&country=USA&key=5b78dd52-e51a-47f6-b63f-d580b4e33b83") // get the resp
    if err != nil {
	log.Fatalln(err)
    }
    body, err := io.ReadAll(resp.Body) 
    defer resp.Body.Close()

    s := string(body)
    data := CityRequest{}

    eerr := json.Unmarshal([]byte(s), &data)
    if eerr != nil {
	fmt.Println(eerr.Error()) 
    } // process it into a request

    cityData := data
    cityNames := make([]string, len(cityData.Data))

    for i := 0; i < len(cityData.Data); i++ {
	cityNames = append(cityNames, cityData.Data[i].City)
    }
    return cityNames
}


func (m model) Init() tea.Cmd {
    // Just return `nil`, which means "no I/O right now, please."
    return nil
}
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    //fmt.Println("updating")
    switch msg := msg.(type) {

    // Is it a key press?
    case tea.KeyMsg:

        // Cool, what was the actual key pressed?
        switch msg.String() {

        // These keys should exit the program.
        case "ctrl+c", "q":
            return m, tea.Quit

	//case "l":
	//    //aqi := getWeather()
	//    //fmt.Println(aqi)
	//    //return 
	//    fmt.Println("")
	}
    }


    // Return the updated model to the Bubble Tea runtime for processing.
    // Note that we're not returning a command.
    return m, nil
}

func (m model) View() string {
    // The header
    s := string(setAqi)
    s += "\nPress q to quit.\n"
    fmt.Println(m, "the aqii")

    // Send the UI for rendering
    return s
}

var setAqi int = -1

func main() {

    cities := getCities()
    sp := selection.New("What city do you want to search for?",
	selection.Choices(cities))
    sp.PageSize = 3

    choice, err := sp.RunPrompt()
    if err != nil {
	fmt.Printf("Error: %v\n", err)
	os.Exit(1)
    }

    //_ = choice
    aqi := getWeather(choice)
    setAqi = aqi
    //_ = aqi
    //fmt.Println(aqi)

    p := tea.NewProgram(initialModel(aqi))
    if err := p.Start(); err != nil {
	fmt.Printf("Alas, there's been an error: %v", err)
	os.Exit(1)
    }
    //os.Exit(0)

}






















