package main

// TODO(jeanhaley) - The following items need to be addressed
// Refactor Code, then figure out what to do with it.
import (
	"flag"
	"fmt"
	"net"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/common-nighthawk/go-figure"
)

const (
	corgi = "⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⣰⣧⣼⣧⠀⠀⠀⠀⠀⠀\n" +
		"⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⣿⣿⣭⣭⣤⣄⠀⠀⠀⠀\n" +
		"⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⣼⣿⣿⣿⣿⣿⣷⣤⣤⡄\n" +
		"⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⣿⣿⣿⣿⣿⣿⣿⣿⡿⠀\n" +
		"⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⣀⣼⣿⣮⣍⣉⣉⣀⣀⠀⠀⠀\n" +
		"⠀⠀⣠⣶⣶⣶⣶⣶⣶⣶⣶⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣧⠀⠀\n" +
		"⣴⣿⣿⣿⣿⣿⣯⡛⢿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⠀⠀\n" +
		"⠉⣿⣿⣿⣿⣿⣿⣷⠘⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⠏⠀⠀\n" +
		"⠀⣿⣿⣿⣿⣿⣿⡟⠸⠿⠿⠿⣿⣿⣿⣿⣿⣿⣿⠿⠋⠀⠀⠀\n" +
		"⠀⠘⢿⣿⣿⠿⠋⠀⠀⠀⠀⠀⠀⠉⠉⣿⣿⡏⠁⠀⠀⠀⠀⠀\n" +
		"⠀⠀⢸⣿⣿⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⣿⣿⡇⠀⠀⠀⠀⠀⠀\n"
	gnome = "⣿⣿⣿⣿⣿⠟⠉⠁⠄⠄⠄⠈⠙⠿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿\n" +
		"⣿⣿⣿⣿⠏⠄⠄⠄⠄⠄⠄⠄⠄⠄⠸⢿⣿⣿⣿⣿⣿⣿⣿⣿\n" +
		"⣿⣿⣿⣏⠄⡠⡤⡤⡤⡤⡤⡤⡠⡤⡤⣸⣿⣿⣿⣿⣿⣿⣿⣿\n" +
		"⣿⣿⣿⣗⢝⢮⢯⡺⣕⢡⡑⡕⡍⣘⢮⢿⣿⣿⣿⣿⣿⣿⣿⣿\n" +
		"⣿⣿⣿⣿⡧⣝⢮⡪⡪⡪⡎⡎⡮⡲⣱⣻⣿⣿⣿⣿⣿⣿⣿⣿\n" +
		"⣿⣿⣿⠟⠁⢸⡳⡽⣝⢝⢌⢣⢃⡯⣗⢿⣿⣿⣿⣿⣿⣿⣿⣿\n" +
		"⣿⠟⠁⠄⠄⠄⠹⡽⣺⢽⢽⢵⣻⢮⢯⠟⠿⠿⢿⣿⣿⣿⣿⣿\n" +
		"⡟⢀⠄⠄⠄⠄⠄⠙⠽⠽⡽⣽⣺⢽⠝⠄⠄⢰⢸⢝⠽⣙⢝⢿\n" +
		"⡄⢸⢹⢸⢱⢘⠄⠄⠄⠄⠄⠈⠄⠄⠄⣀⠄⠄⣵⣧⣫⣶⣜⣾\n" +
		"⣧⣬⣺⠸⡒⠬⡨⠄⠄⠄⠄⠄⠄⠄⣰⣿⣿⣿⣿⣿⣷⣽⣿⣿\n" +
		"⣿⣿⣿⣷⠡⠑⠂⠄⠄⠄⠄⠄⠄⠄⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿\n" +
		"⣿⣿⣿⣿⣄⠠⢀⢀⢀⡀⡀⠠⢀⢲⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿\n" +
		"⣿⣿⣿⣿⣿⢐⢀⠂⢄⠇⠠⠈⠄⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿\n" +
		"⣿⣿⣿⣿⣧⠄⠠⠈⢈⡄⠄⢁⢀⣾⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿\n" +
		"⣿⣿⣿⣿⣿⡀⠠⠐⣼⠇⠄⡀⠸⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿\n" +
		"⣿⣿⣿⣿⣯⠄⠄⡀⠈⠂⣀⠄⢀⠄⠈⣿⣿⣿⣿⣿⣿⣿⣿⣿\n" +
		"⣿⣿⣿⣿⣿⣶⣄⣀⠐⢀⣸⣷⣶⣶⣶⣿⣿⣿⣿⣿⣿⣿⣿⣿"
	clearScreen = "\033[H\033[2J"
)

// ___ Global Variables ___

var (
	ip, netp, banner, banmsgfp, SpecialMessage, socket                                      string
	buffersize, logerTime, clientchannelbuffer, logchannelbuffer, systemchannelbuffer, port int
	ClientMessageCount                                                                      int // Sets limit for messages show to client
	clientChan, logChan, sysChan                                                            ch  // Global Channels
	currentstate                                                                            state
	globalState                                                                             []msg
	ServerStartTime                                                                         time.Time
)

const (
	carRet uint8 = 13
)

var (
	// creating a blank global branding variable.
	// this needs to be done, because the type figure.figure is not exported.
	branding = figure.NewColorFigure("", "nancyj-fancy", "Blue", true)
	// Carraige return for Windows.
)

// uses init function to set set up global flag variables, and channels.
func init() {
	ServerStartTime = time.Now()
	// setting Global Flags
	flag.StringVar(&ip, "true", "0.0.0.0", "IP for server to listen on, default is 0.0.0.0")
	flag.StringVar(&netp, "netp", "tcp", "Network protocol to use")
	flag.IntVar(&port, "port", 6000, "Port for server to listen on; defaults to 6000")
	flag.IntVar(&buffersize, "bufferSize", 1024, "Message Buffer size.")
	flag.IntVar(&logerTime, "logerTime", 120, "time in between server status check, in seconds.")
	flag.StringVar(&banner, "banner", "TCPServ", "Banner to display on startup")
	flag.IntVar(&clientchannelbuffer, "clientchannelbuffer", 20, "size of client channel buffer")
	flag.IntVar(&logchannelbuffer, "logchannelbuffer", 20, "size of log channel buffer")
	flag.IntVar(&systemchannelbuffer, "systemchannelbuffer", 20, "size of system channel buffer")
	flag.IntVar(&ClientMessageCount, "ClientMessageCount", 20, "Number of messages to show to client")
	flag.StringVar(&banmsgfp, "banmsgfp", "msg.txt", "Banner Message File Path")
	flag.Parse()

	globalState = make([]msg, 0, ClientMessageCount)
	// instantiating global channels.
	clientChan = make(chan msg, clientchannelbuffer)
	logChan = make(chan msg, logchannelbuffer)
	sysChan = make(chan msg, systemchannelbuffer)
	branding = figure.NewColorFigure(banner, "nancyj-fancy", "Blue", true) // sets banner to value passed by terminal flags.
	socket = net.JoinHostPort(ip, strconv.Itoa(port))
}

func main() {
	for _, v := range branding.Slicify() {
		fmt.Println(colorWrap(Blue, v))
		time.Sleep(100 * time.Millisecond)
	}
	fmt.Print(corgi)
	var wg sync.WaitGroup
	wg.Add(3) // adding two goroutines
	// Handle panicked state. Logs and recovers from panic
	go func() {
		MessageBroker() // starting the Event Handler go routine
		wg.Done()       // decrementing the counter when done
	}()
	go func() {
		connListener(ip)
		wg.Done() // decrementing the counter when done
	}()
	go func() {
		TimeKeeper()
		wg.Done()
	}()
	wg.Wait() // waiting for all goroutines to finish
}

// returns payload based on action preceeding ":"
func parseAction(m msg) payload {
	switch strings.Split(m.GetPayload().String(), ":")[0] {
	case "ascii":
		byteMsg := m.payload
		// if the last character is a carraige return, remove that character.
		// This breaks the library we're using to convert text to ascii.
		if byteMsg[len(byteMsg)-1] == carRet {
			byteMsg = byteMsg[:len(byteMsg)-1]
		}
		ascii := figure.NewColorFigure(strings.Split(byteMsg.String(), ":")[1], "nancyj-fancy", "Green", true)
		asciiMsg := ascii.ColorString()
		return payload(asciiMsg)
	default:
		return payload(fmt.Sprintf("Invalid action: %v", strings.Split(m.GetPayload().String(), ":")[0]))
	}
}

// Returns Splash Screen elements.
func splashScreen() string {
	welcome := "Welcome to the Void!"
	activeconn := colorWrap(
		Green, fmt.Sprintf(
			"There are currently %v active connections.", currentstate.ActiveConnections()))
	directions := colorWrap(Purple, "Type 'ascii:' before your message to display ascii art")
	splashmessage := fmt.Sprintf("\t\t%v\n\t  %v\n  %v\v", welcome, activeconn, directions)
	return splashmessage
}

// HasString returns true if a string contains another string
func HasString(str, match string) bool {
	bool, err := regexp.MatchString(match, str)
	if err != nil {
		return false
	}
	return bool
}

// Prints text with border
func printWithBorder(text string) string {
	horizontalBorder := "+" + strings.Repeat("-", len(text)+2) + "+"
	return fmt.Sprintf("%v\n| %v |\n%v", horizontalBorder, text, horizontalBorder)
}
