package main

import (
    "context"
    "time"
    "encoding/json"
    "github.com/dgraph-io/dgo/v200"
    "github.com/uber/h3-go/v4"
)

// Extended Route structure with additional fields
type Route struct {
    Uid            string      `json:"uid,omitempty"`
    RouteNumber    string      `json:"route_number"`
    PickupPoint    string      `json:"pickup_point"`
    Destinations   []string    `json:"destinations"`
    PickupH3Index  string      `json:"pickup_h3_index"`
    DestH3Index    string      `json:"dest_h3_index"`
    PickupLat      float64     `json:"pickup_lat"`
    PickupLng      float64     `json:"pickup_lng"`
    DestLat        float64     `json:"dest_lat"`
    DestLng        float64     `json:"dest_lng"`
    Schedule       []Schedule  `json:"schedule"`
    Fare           FareInfo    `json:"fare"`
    ActiveDays     []string    `json:"active_days"`
    LastUpdated    time.Time   `json:"last_updated"`
}

type Schedule struct {
    StartTime string `json:"start_time"`
    EndTime   string `json:"end_time"`
    Frequency int    `json:"frequency_minutes"`
}

type FareInfo struct {
    Regular     float64 `json:"regular_fare"`
    PeakHours   float64 `json:"peak_fare"`
    OffPeakHours float64 `json:"off_peak_fare"`
}

// Extended schema for Dgraph
const extendedSchema = `
    route_number: string @index(exact) .
    pickup_point: string @index(term) .
    destinations: [string] @index(term) .
    pickup_h3_index: string @index(exact) .
    dest_h3_index: string @index(exact) .
    pickup_lat: float .
    pickup_lng: float .
    dest_lat: float .
    dest_lng: float .
    schedule: [uid] @reverse .
    fare: uid @reverse .
    active_days: [string] @index(term) .
    last_updated: datetime @index(hour) .

    type Route {
        route_number
        pickup_point
        destinations
        pickup_h3_index
        dest_h3_index
        pickup_lat
        pickup_lng
        dest_lat
        dest_lng
        schedule
        fare
        active_days
        last_updated
    }

    type Schedule {
        start_time
        end_time
        frequency_minutes
    }

    type FareInfo {
        regular_fare
        peak_fare
        off_peak_fare
    }
`

// Advanced search function with multiple criteria
type SearchCriteria struct {
    NearLocation    *Location  `json:"near_location,omitempty"`
    Destination     string     `json:"destination,omitempty"`
    MaxDistance     float64    `json:"max_distance,omitempty"`
    DayOfWeek      string     `json:"day_of_week,omitempty"`
    MaxFare         float64    `json:"max_fare,omitempty"`
    TimeOfDay      string     `json:"time_of_day,omitempty"`
}

type Location struct {
    Lat float64 `json:"lat"`
    Lng float64 `json:"lng"`
}

func SearchRoutes(ctx context.Context, dgraphClient *dgo.Dgraph, criteria SearchCriteria) ([]Route, error) {
    var queryBuilder strings.Builder
    queryBuilder.WriteString(`
        query SearchRoutes($location: string, $maxDistance: float, $day: string, $maxFare: float) {
            routes(func: type(Route)) @filter(`)

    // Build dynamic filter conditions
    var conditions []string
    variables := make(map[string]interface{})

    if criteria.NearLocation != nil {
        h3Index := h3.LatLngToCell(h3.LatLng{
            Lat: criteria.NearLocation.Lat,
            Lng: criteria.NearLocation.Lng,
        }, 9)
        neighbors := h3.GridDisk(h3Index, 1)
        h3Indexes := make([]string, len(neighbors))
        for i, n := range neighbors {
            h3Indexes[i] = n.String()
        }
        conditions = append(conditions, "eq(pickup_h3_index, val(h3Indexes))")
        variables["h3Indexes"] = h3Indexes
    }

    if criteria.Destination != "" {
        conditions = append(conditions, "anyofterms(destinations, $destination)")
        variables["destination"] = criteria.Destination
    }

    if criteria.DayOfWeek != "" {
        conditions = append(conditions, "anyofterms(active_days, $day)")
        variables["day"] = criteria.DayOfWeek
    }

    // Combine conditions and complete the query
    queryBuilder.WriteString(strings.Join(conditions, " AND "))
    queryBuilder.WriteString(`) {
        uid
        route_number
        pickup_point
        destinations
        pickup_lat
        pickup_lng
        dest_lat
        dest_lng
        schedule {
            start_time
            end_time
            frequency_minutes
        }
        fare {
            regular_fare
            peak_fare
            off_peak_fare
        }
        active_days
        last_updated
    }
}`)

    resp, err := dgraphClient.NewTxn().QueryWithVars(ctx, queryBuilder.String(), variables)
    if err != nil {
        return nil, err
    }

    var result struct {
        Routes []Route `json:"routes"`
    }
    if err := json.Unmarshal(resp.Json, &result); err != nil {
        return nil, err
    }

    // Post-process results for additional filtering
    filteredRoutes := filterRoutes(result.Routes, criteria)
    return filteredRoutes, nil
}

// Example usage of advanced search
func searchExample() {
    criteria := SearchCriteria{
        NearLocation: &Location{
            Lat: -1.2865,
            Lng: 36.815,
        },
        Destination: "Ngong Road",
        MaxDistance: 1000, // meters
        DayOfWeek: "Monday",
        MaxFare: 100,
        TimeOfDay: "08:00",
    }

    routes, err := SearchRoutes(context.Background(), dgraphClient, criteria)
    if err != nil {
        log.Fatal(err)
    }

    // Process and display results
    for _, route := range routes {
        fmt.Printf("Route %s: %s -> %s\n", 
            route.RouteNumber, 
            route.PickupPoint, 
            strings.Join(route.Destinations, " -> "))
    }
}

// Helper function to update route information
func UpdateRoute(ctx context.Context, dgraphClient *dgo.Dgraph, route Route) error {
    // Update H3 indexes
    resolution := 9
    route.PickupH3Index = h3.LatLngToCell(h3.LatLng{
        Lat: route.PickupLat,
        Lng: route.PickupLng,
    }, resolution).String()

    route.DestH3Index = h3.LatLngToCell(h3.LatLng{
        Lat: route.DestLat,
        Lng: route.DestLng,
    }, resolution).String()

    route.LastUpdated = time.Now()

    // Prepare mutation
    mutation := &api.Mutation{
        SetJson: getMutationJSON(route),
        CommitNow: true,
    }

    _, err := dgraphClient.NewTxn().Mutate(ctx, mutation)
    return err
}

// Cache implementation for frequently accessed routes
type RouteCache struct {
    cache    map[string][]Route
    mu       sync.RWMutex
    maxAge   time.Duration
}

func NewRouteCache(maxAge time.Duration) *RouteCache {
    return &RouteCache{
        cache:  make(map[string][]Route),
        maxAge: maxAge,
    }
}

func (rc *RouteCache) Get(key string) ([]Route, bool) {
    rc.mu.RLock()
    defer rc.mu.RUnlock()
    routes, exists := rc.cache[key]
    return routes, exists
}

func (rc *RouteCache) Set(key string, routes []Route) {
    rc.mu.Lock()
    defer rc.mu.Unlock()
    rc.cache[key] = routes
}

// RouteAnalytics tracks usage and performance metrics
type RouteAnalytics struct {
    RouteID        string    `json:"route_id"`
    UsageCount     int       `json:"usage_count"`
    AverageDelay   float64   `json:"avg_delay_minutes"`
    PeakHours      []string  `json:"peak_hours"`
    Reliability    float64   `json:"reliability_score"` // 0-1
    LastAnalyzed   time.Time `json:"last_analyzed"`
}

// RouteSet manages sets of routes with common characteristics
type RouteSet struct {
    ID          string    `json:"id"`
    Name        string    `json:"name"`
    Routes      []string  `json:"route_ids"` // Route UIDs
    Coverage    []string  `json:"h3_coverage"` // H3 indexes covered
    Properties  map[string]interface{} `json:"properties"`
    Created     time.Time `json:"created_at"`
    Updated     time.Time `json:"updated_at"`
}

// RoutePlanner handles route optimization and planning
type RoutePlanner struct {
    cache       *RouteCache
    analytics   map[string]*RouteAnalytics
    routeSets   map[string]*RouteSet
    mu          sync.RWMutex
    dgraph      *dgo.Dgraph
}

func NewRoutePlanner(dgraph *dgo.Dgraph, cache *RouteCache) *RoutePlanner {
    return &RoutePlanner{
        cache:     cache,
        analytics: make(map[string]*RouteAnalytics),
        routeSets: make(map[string]*RouteSet),
        dgraph:    dgraph,
    }
}

// CreateRouteSet creates a new set of routes based on criteria
func (rp *RoutePlanner) CreateRouteSet(ctx context.Context, name string, criteria SearchCriteria) (*RouteSet, error) {
    routes, err := SearchRoutes(ctx, rp.dgraph, criteria)
    if err != nil {
        return nil, err
    }

    // Create new RouteSet
    routeSet := &RouteSet{
        ID:         generateUUID(),
        Name:       name,
        Routes:     make([]string, len(routes)),
        Coverage:   make([]string, 0),
        Properties: make(map[string]interface{}),
        Created:    time.Now(),
        Updated:    time.Now(),
    }

    // Calculate H3 coverage for the route set
    coverageMap := make(map[string]bool)
    for i, route := range routes {
        routeSet.Routes[i] = route.Uid

        // Add pickup and destination H3 indexes to coverage
        coverageMap[route.PickupH3Index] = true
        coverageMap[route.DestH3Index] = true

        // Add intermediate H3 indexes along the route
        intermediateIndexes := calculateIntermediateH3Indexes(
            route.PickupLat, route.PickupLng,
            route.DestLat, route.DestLng,
            9, // resolution
        )
        for _, idx := range intermediateIndexes {
            coverageMap[idx] = true
        }
    }

    // Convert coverage map to slice
    for idx := range coverageMap {
        routeSet.Coverage = append(routeSet.Coverage, idx)
    }

    // Store route set
    rp.mu.Lock()
    rp.routeSets[routeSet.ID] = routeSet
    rp.mu.Unlock()

    return routeSet, nil
}

// OptimizeRouteSet optimizes routes within a set based on various metrics
func (rp *RoutePlanner) OptimizeRouteSet(ctx context.Context, setID string, optimizationCriteria map[string]float64) error {
    rp.mu.RLock()
    routeSet, exists := rp.routeSets[setID]
    rp.mu.RUnlock()

    if !exists {
        return fmt.Errorf("route set not found: %s", setID)
    }

    // Gather analytics for all routes in the set
    routeAnalytics := make(map[string]*RouteAnalytics)
    for _, routeID := range routeSet.Routes {
        analytics, err := rp.getRouteAnalytics(ctx, routeID)
        if err != nil {
            continue
        }
        routeAnalytics[routeID] = analytics
    }

    // Optimize based on criteria
    optimizedRoutes := optimizeRoutes(routeSet.Routes, routeAnalytics, optimizationCriteria)

    // Update route set with optimized routes
    rp.mu.Lock()
    routeSet.Routes = optimizedRoutes
    routeSet.Updated = time.Now()
    rp.mu.Unlock()

    return nil
}

// AnalyzeRouteSet performs analysis on a set of routes
func (rp *RoutePlanner) AnalyzeRouteSet(ctx context.Context, setID string) (*RouteSetAnalysis, error) {
    rp.mu.RLock()
    routeSet, exists := rp.routeSets[setID]
    rp.mu.RUnlock()

    if !exists {
        return nil, fmt.Errorf("route set not found: %s", setID)
    }

    analysis := &RouteSetAnalysis{
        SetID:           setID,
        RouteCount:      len(routeSet.Routes),
        CoverageArea:    calculateCoverageArea(routeSet.Coverage),
        AverageMetrics:  make(map[string]float64),
        PeakHours:       identifySetPeakHours(routeSet),
        ReliabilityScore: calculateSetReliability(routeSet),
        AnalyzedAt:      time.Now(),
    }

    return analysis, nil
}

// Helper function to calculate intermediate H3 indexes between two points
func calculateIntermediateH3Indexes(startLat, startLng, endLat, endLng float64, resolution int) []string {
    var indexes []string
    
    // Calculate number of points based on distance
    distance := calculateDistance(startLat, startLng, endLat, endLng)
    pointCount := int(distance / 100) // One point every 100 meters

    for i := 0; i <= pointCount; i++ {
        fraction := float64(i) / float64(pointCount)
        lat := startLat + (endLat-startLat)*fraction
        lng := startLng + (endLng-startLng)*fraction

        h3Index := h3.LatLngToCell(h3.LatLng{
            Lat: lat,
            Lng: lng,
        }, resolution)
        indexes = append(indexes, h3Index.String())
    }

    return indexes
}

// RouteSetAnalysis contains analytical data for a route set
type RouteSetAnalysis struct {
    SetID            string             `json:"set_id"`
    RouteCount       int                `json:"route_count"`
    CoverageArea     float64            `json:"coverage_area_km2"`
    AverageMetrics   map[string]float64 `json:"average_metrics"`
    PeakHours        []string           `json:"peak_hours"`
    ReliabilityScore float64            `json:"reliability_score"`
    AnalyzedAt       time.Time          `json:"analyzed_at"`
}

// Calculate coverage area from H3 indexes
func calculateCoverageArea(h3Indexes []string) float64 {
    uniqueCells := make(map[string]bool)
    for _, idx := range h3Indexes {
        uniqueCells[idx] = true
    }

    // H3 resolution 9 has an average hexagon area of about 0.1 km²
    return float64(len(uniqueCells)) * 0.1
}

// Identify peak hours for a route set
func identifySetPeakHours(routeSet *RouteSet) []string {
    hourCounts := make(map[string]int)
    peakThreshold := len(routeSet.Routes) / 2

    // Count usage per hour
    for _, routeID := range routeSet.Routes {
        // Aggregate schedule information
        // This is a simplified example - you'd typically get this from your analytics data
    }

    // Find hours exceeding threshold
    var peakHours []string
    for hour, count := range hourCounts {
        if count >= peakThreshold {
            peakHours = append(peakHours, hour)
        }
    }

    return peakHours
}

// Calculate reliability score for a route set
func calculateSetReliability(routeSet *RouteSet) float64 {
    var totalScore float64
    var count int

    for _, routeID := range routeSet.Routes {
        // Aggregate reliability scores
        // This is a simplified example - you'd typically get this from your analytics data
        totalScore += 0.9 // placeholder
        count++
    }

    if count == 0 {
        return 0
    }
    return totalScore / float64(count)
}

// Example usage of the route planner
func routePlannerExample() {
    ctx := context.Background()
    cache := NewRouteCache(15 * time.Minute)
    planner := NewRoutePlanner(dgraphClient, cache)

    // Create a new route set
    criteria := SearchCriteria{
        NearLocation: &Location{
            Lat: -1.2865,
            Lng: 36.815,
        },
        MaxDistance: 5000, // 5km
    }

    routeSet, err := planner.CreateRouteSet(ctx, "CBD Routes", criteria)
    if err != nil {
        log.Fatal(err)
    }

    // Optimize the route set
    optimizationCriteria := map[string]float64{
        "reliability": 0.7,
        "coverage": 0.3,
    }
    err = planner.OptimizeRouteSet(ctx, routeSet.ID, optimizationCriteria)
    if err != nil {
        log.Fatal(err)
    }

    // Analyze the route set
    analysis, err := planner.AnalyzeRouteSet(ctx, routeSet.ID)
    if err != nil {
        log.Fatal(err)
    }

    fmt.Printf("Route Set Analysis:\n")
    fmt.Printf("Coverage Area: %.2f km²\n", analysis.CoverageArea)
    fmt.Printf("Reliability Score: %.2f\n", analysis.ReliabilityScore)
    fmt.Printf("Peak Hours: %v\n", analysis.PeakHours)
}

package main

import (
    "context"
    "time"
    "sync"
    "math"
    "github.com/dgraph-io/dgo/v200"
    "github.com/uber/h3-go/v4"
)

// TrafficData represents real-time traffic information
type TrafficData struct {
    H3Index      string    `json:"h3_index"`
    Speed        float64   `json:"speed_kmh"`
    Congestion   float64   `json:"congestion_level"` // 0-1
    LastUpdated  time.Time `json:"last_updated"`
    Source       string    `json:"source"`
}

// PredictiveModel represents route demand and performance predictions
type PredictiveModel struct {
    RouteID      string                `json:"route_id"`
    Predictions  map[string]Prediction `json:"predictions"` // hour -> prediction
    Accuracy     float64               `json:"accuracy_score"`
    LastTrained  time.Time             `json:"last_trained"`
}

type Prediction struct {
    Demand       float64 `json:"expected_demand"`
    TravelTime   float64 `json:"expected_travel_time"`
    Reliability  float64 `json:"reliability_score"`
}

// RealTimeManager handles real-time updates and predictions
type RealTimeManager struct {
    trafficData    map[string]*TrafficData
    predictions    map[string]*PredictiveModel
    updates        chan RouteUpdate
    mu            sync.RWMutex
    dgraph        *dgo.Dgraph
    planner       *RoutePlanner
}

type RouteUpdate struct {
    RouteID     string
    UpdateType  string
    Data        interface{}
    Timestamp   time.Time
}

func NewRealTimeManager(dgraph *dgo.Dgraph, planner *RoutePlanner) *RealTimeManager {
    rtm := &RealTimeManager{
        trafficData: make(map[string]*TrafficData),
        predictions: make(map[string]*PredictiveModel),
        updates:     make(chan RouteUpdate, 1000),
        dgraph:      dgraph,
        planner:     planner,
    }

    // Start update processing
    go rtm.processUpdates()
    return rtm
}

// ProcessUpdates handles real-time updates
func (rtm *RealTimeManager) processUpdates() {
    for update := range rtm.updates {
        switch update.UpdateType {
        case "traffic":
            rtm.handleTrafficUpdate(update)
        case "demand":
            rtm.handleDemandUpdate(update)
        case "incident":
            rtm.handleIncidentUpdate(update)
        }
    }
}

// UpdateTraffic updates traffic data for a specific area
func (rtm *RealTimeManager) UpdateTraffic(h3Index string, speed float64, congestion float64) {
    rtm.mu.Lock()
    defer rtm.mu.Unlock()

    rtm.trafficData[h3Index] = &TrafficData{
        H3Index:     h3Index,
        Speed:       speed,
        Congestion:  congestion,
        LastUpdated: time.Now(),
        Source:      "real-time-sensors",
    }

    // Send update to processing channel
    rtm.updates <- RouteUpdate{
        UpdateType: "traffic",
        Data:      rtm.trafficData[h3Index],
        Timestamp: time.Now(),
    }
}

// PredictDemand predicts route demand for the next 24 hours
func (rtm *RealTimeManager) PredictDemand(ctx context.Context, routeID string) (map[string]Prediction, error) {
    rtm.mu.RLock()
    model, exists := rtm.predictions[routeID]
    rtm.mu.RUnlock()

    if !exists || time.Since(model.LastTrained) > 24*time.Hour {
        // Train new model if doesn't exist or is outdated
        if err := rtm.trainModel(ctx, routeID); err != nil {
            return nil, err
        }
        model = rtm.predictions[routeID]
    }

    // Get predictions for next 24 hours
    predictions := make(map[string]Prediction)
    now := time.Now()
    for i := 0; i < 24; i++ {
        hour := now.Add(time.Duration(i) * time.Hour)
        hourKey := hour.Format("15:00")
        predictions[hourKey] = model.Predictions[hourKey]
    }

    return predictions, nil
}

// TrainModel trains the predictive model for a route
func (rtm *RealTimeManager) trainModel(ctx context.Context, routeID string) error {
    // Fetch historical data
    historicalData, err := rtm.getHistoricalData(ctx, routeID)
    if err != nil {
        return err
    }

    // Train model using historical data
    model := &PredictiveModel{
        RouteID:     routeID,
        Predictions: make(map[string]Prediction),
        LastTrained: time.Now(),
    }

    // Process historical data and generate predictions
    for hour := 0; hour < 24; hour++ {
        hourStr := fmt.Sprintf("%02d:00", hour)
        prediction := rtm.calculatePrediction(historicalData, hour)
        model.Predictions[hourStr] = prediction
    }

    // Calculate model accuracy
    model.Accuracy = rtm.calculateModelAccuracy(historicalData, model.Predictions)

    // Store the model
    rtm.mu.Lock()
    rtm.predictions[routeID] = model
    rtm.mu.Unlock()

    return nil
}

// OptimizeRealTime performs real-time route optimization
func (rtm *RealTimeManager) OptimizeRealTime(ctx context.Context, routeSetID string) error {
    // Get current traffic conditions
    trafficConditions := rtm.getCurrentTrafficConditions()
    
    // Get demand predictions
    predictions, err := rtm.getPredictionsForRouteSet(ctx, routeSetID)
    if err != nil {
        return err
    }

    // Optimize based on current conditions and predictions
    optimizationCriteria := map[string]float64{
        "current_traffic": 0.4,
        "predicted_demand": 0.3,
        "reliability": 0.3,
    }

    return rtm.planner.OptimizeRouteSet(ctx, routeSetID, optimizationCriteria)
}

// MonitorRouteHealth continuously monitors route health
func (rtm *RealTimeManager) MonitorRouteHealth(ctx context.Context, routeID string) {
    ticker := time.NewTicker(5 * time.Minute)
    defer ticker.Stop()

    for {
        select {
        case <-ctx.Done():
            return
        case <-ticker.C:
            health := rtm.calculateRouteHealth(routeID)
            if health.Score < 0.7 {
                rtm.handleUnhealthyRoute(routeID, health)
            }
        }
    }
}

// Example usage
func realTimeExample() {
    ctx := context.Background()
    cache := NewRouteCache(15 * time.Minute)
    planner := NewRoutePlanner(dgraphClient, cache)
    rtm := NewRealTimeManager(dgraphClient, planner)

    // Start route monitoring
    go func() {
        for _, routeID := range planner.GetAllRouteIDs() {
            go rtm.MonitorRouteHealth(ctx, routeID)
        }
    }()

    // Update traffic data
    rtm.UpdateTraffic("8928308281fffff", 45.5, 0.3)

    // Get demand predictions
    predictions, err := rtm.PredictDemand(ctx, "route-123")
    if err != nil {
        log.Fatal(err)
    }

    // Optimize routes in real-time
    err = rtm.OptimizeRealTime(ctx, "set-456")
    if err != nil {
        log.Fatal(err)
    }

    // Example of handling real-time updates
    update := RouteUpdate{
        RouteID:    "route-123",
        UpdateType: "incident",
        Data: map[string]interface{}{
            "type": "accident",
            "location": map[string]float64{
                "lat": -1.2865,
                "lng": 36.815,
            },
            "severity": "high",
        },
        Timestamp: time.Now(),
    }
    rtm.updates <- update
}

// HealthMetrics represents route health information
type HealthMetrics struct {
    Score        float64   `json:"health_score"`
    Issues       []string  `json:"issues"`
    LastChecked  time.Time `json:"last_checked"`
}

func (rtm *RealTimeManager) calculateRouteHealth(routeID string) HealthMetrics {
    rtm.mu.RLock()
    defer rtm.mu.RUnlock()

    health := HealthMetrics{
        Score:       1.0,
        LastChecked: time.Now(),
    }

    // Check traffic conditions
    for _, traffic := range rtm.trafficData {
        if traffic.Congestion > 0.8 {
            health.Score *= 0.8
            health.Issues = append(health.Issues, "Heavy traffic detected")
        }
    }

    // Check reliability
    if model, exists := rtm.predictions[routeID]; exists {
        if model.Accuracy < 0.7 {
            health.Score *= 0.9
            health.Issues = append(health.Issues, "Low prediction accuracy")
        }
    }

    return health
}

func (rtm *RealTimeManager) handleUnhealthyRoute(routeID string, health HealthMetrics) {
    // Log the issue
    log.Printf("Unhealthy route detected: %s (Score: %.2f)", routeID, health.Score)

    // Notify relevant systems
    for _, issue := range health.Issues {
        rtm.updates <- RouteUpdate{
            RouteID:    routeID,
            UpdateType: "health_issue",
            Data:       issue,
            Timestamp:  time.Now(),
        }
    }

    // Trigger re-optimization if necessary
    if health.Score < 0.5 {
        ctx := context.Background()
        rtm.OptimizeRealTime(ctx, routeID)
    }
}
