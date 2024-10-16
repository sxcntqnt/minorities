package main

import (
	"context"
	"fmt"
	"github.com/golang-jwt/jwt"
	"log"
	"math"
	"net/http"
	"time"

	"github.com/dgraph-io/dgo/v210"
	"github.com/dgraph-io/dgo/v210/protos/api"
	"github.com/dixonwille/wmenu/v5"
	"github.com/spf13/viper"
)

// Dgraph server address (modify if necessary)
const dgraphServer = "localhost:9080"

// Function to initialize Viper and read from .env file
func initConfig() {
    // Set the name of the .env file without the extension
    viper.SetConfigName(".env")
    viper.SetConfigType("env") // Set the config type to ENV
    viper.AddConfigPath(".")   // Look for the .env file in the current directory

    // Read the configuration from the .env file
    if err := viper.ReadInConfig(); err != nil {
        log.Println("No .env file found, using environment variables instead")
    }
}

// Function to create a Dgraph client for Dgraph Cloud
func createDgraphClient() *dgo.Dgraph {
    // Initialize Viper to read from .env file
    initConfig()

    endpoint := viper.GetString("DGRAPH_ENDPOINT") // Load the endpoint from the .env file
    apiToken := viper.GetString("DGRAPH_API_TOKEN") // Load the API token from the .env file

    // Create a connection to Dgraph Cloud
    conn, err := dgo.DialCloud(endpoint, apiToken)
    if err != nil {
        log.Fatal("Unable to connect to Dgraph Cloud:", err)
    }

    // Ensure that the connection is closed when done
    defer conn.Close()

    dgraphClient := dgo.NewDgraphClient(api.NewDgraphClient(conn))
    return dgraphClient
}

// Client mode: Perform a simple Dgraph query
func dgraphClientExample() {
	dgraphClient := createDgraphClient()
	ctx := context.Background()

	// Create a query to fetch data from Dgraph
	query := `{
        all(func: has(name)) {
            uid
            name
        }
    }`

	// Perform the query
	response, err := dgraphClient.NewTxn().Query(ctx, query)
	if err != nil {
		log.Fatal("Failed to query Dgraph:", err)
	}

	fmt.Printf("Query result: %s\n", response.Json)
}

// Server mode: Start a simple HTTP server to simulate a service (Dgraph is already running separately)
func dgraphServerExample() {
	// Create a simple HTTP server that could potentially serve Dgraph GraphQL queries
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Dgraph server is running...")
	})

	fmt.Println("Server is listening at http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

// Function to calculate distance between two points
func distance(lat1, lng1, lat2, lng2 float64, unit ...string) float64 {
	const PI = 3.141592653589793
	radlat1 := PI * lat1 / 180
	radlat2 := PI * lat2 / 180
	theta := lng1 - lng2
	radtheta := PI * theta / 180

	dist := math.Sin(radlat1)*math.Sin(radlat2) + math.Cos(radlat1)*math.Cos(radlat2)*math.Cos(radtheta)
	if dist > 1 {
		dist = 1
	}

	dist = math.Acos(dist)
	dist = dist * 180 / PI
	dist = dist * 60 * 1.1515 // Default: miles

	if len(unit) > 0 && unit[0] == "K" {
		dist = dist * 1.609344 // Kilometers
	}

	return dist
}

// Function to calculate speed based on distance and time
func speed(s, dist float64) float64 {
	if s == 0 {
		return 0 // Avoid division by zero
	}
	return dist / s
}

// Function to generate JWT token
func GenerateToken() (string, error) {
	secretKey := []byte("kadaddy")

	claims := jwt.StandardClaims{
		Issuer:    "my-app",
		ExpiresAt: time.Now().Add(24 * time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secretKey)
}

// Function to handle person management menu
func handlePersonFunc(opts []wmenu.Opt) {
    switch opts[0].Value {
    case 0:
        fmt.Println("Add a new Person selected")
        // Add logic to handle adding a person
    case 1:
        fmt.Println("Find a Person selected")
        // Add logic to handle finding a person
    case 2:
        fmt.Println("Update a Person's information selected")
        // Add logic to handle updating a person
    case 3:
        fmt.Println("Delete a Person by ID selected")
        // Add logic to handle deleting a person
    default:
        fmt.Println("Unknown option selected")
    }
}

// Function to show the person management menu
func showPersonMenu() {
    // Create a command-line menu for person management
    personMenu := wmenu.NewMenu("What would you like to do with Persons?")
    personMenu.Action(func(opts []wmenu.Opt) error {
        handlePersonFunc(opts) // Handle person management options
        return nil
    })

    // Add options to the person menu
    personMenu.Option("Add a new Person", 0, true, nil)
    personMenu.Option("Find a Person", 1, false, nil)
    personMenu.Option("Update a Person's information", 2, false, nil)
    personMenu.Option("Delete a Person by ID", 3, false, nil)

    // Run the person menu
    if err := personMenu.Run(); err != nil {
        log.Fatal(err)
    }
}

// Function to handle Dgraph mode selection
func handleDgraphMode(opts []wmenu.Opt) {
    switch opts[0].Value {
    case 0: // Client mode
        fmt.Println("Selected: Dgraph Client Mode")
        dgraphClientExample()
    case 1: // Server mode
        fmt.Println("Selected: Dgraph Server Mode")
        dgraphServerExample()
    default:
        fmt.Println("Unknown option selected")
    }
}

// Function to show the Dgraph mode selection menu
func showDgraphMenu() {
    // Create a command-line menu for Dgraph mode selection
    dgraphMenu := wmenu.NewMenu("Choose Dgraph Mode:")
    dgraphMenu.Action(func(opts []wmenu.Opt) error {
        handleDgraphMode(opts) // Handle the selected Dgraph option
        return nil
    })

    // Add options to the Dgraph menu
    dgraphMenu.Option("Dgraph Client Mode", 0, true, nil)
    dgraphMenu.Option("Dgraph Server Mode", 1, false, nil)

    // Run the Dgraph menu
    if err := dgraphMenu.Run(); err != nil {
        log.Fatal(err)
    }
}

// Main function
func main() {
	// Generate the JWT token
	tokenString, err := GenerateToken()
	if err != nil {
		log.Fatal("Error generating token:", err)
	}
	fmt.Printf("JWT Token: %s\n", tokenString)

	// Display Dgraph mode menu
	showPersonMenu()

        // Show the Dgraph mode selection menu afterward
        showDgraphMenu()
}
