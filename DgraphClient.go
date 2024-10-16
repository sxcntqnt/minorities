// Function to create a Dgraph client
func createDgraphClient() *dgo.Dgraph {
    conn, err := grpc.Dial(dgraphServer, grpc.WithInsecure())
    if err != nil {
        log.Fatal("Unable to connect to Dgraph server:", err)
    }

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

// Server mode: Start a simple HTTP server
func dgraphServerExample() {
    http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
        fmt.Fprintf(w, "Dgraph server is running...")
    })

    fmt.Println("Server is listening at http://localhost:8080")
    log.Fatal(http.ListenAndServe(":8080", nil))
}
