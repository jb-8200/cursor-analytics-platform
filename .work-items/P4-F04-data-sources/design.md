# Design Document: External Data Source Simulators

**Feature ID**: P4-F04-data-sources
**Created**: January 9, 2026
**Status**: Draft
**Architecture**: Extension of cursor-sim model/generator/handler pattern

---

## Overview

This design document specifies the implementation of three external data source simulators within cursor-sim:

1. **Harvey API** - AI legal assistant usage tracking
2. **Microsoft 365 Copilot Usage API** - Graph API for Copilot adoption metrics
3. **Qualtrics Survey Export API** - 3-step async state machine

Each simulator follows the established cursor-sim architecture patterns while introducing new capabilities for async workflows and enterprise AI tool correlation.

---

## Architecture Overview

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                              cursor-sim v2                                   │
├─────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│  ┌──────────────────────────────────────────────────────────────────────┐  │
│  │                         Seed Loader (Extended)                        │  │
│  │  ┌─────────────┐ ┌─────────────┐ ┌─────────────┐ ┌───────────────┐  │  │
│  │  │ Developers  │ │ Harvey      │ │ M365        │ │ Qualtrics     │  │  │
│  │  │ (existing)  │ │ Users       │ │ Tenants     │ │ Surveys       │  │  │
│  │  └─────────────┘ └─────────────┘ └─────────────┘ └───────────────┘  │  │
│  └──────────────────────────────────────────────────────────────────────┘  │
│                                        │                                     │
│                                        ▼                                     │
│  ┌──────────────────────────────────────────────────────────────────────┐  │
│  │                           Generators                                  │  │
│  │  ┌─────────────┐ ┌─────────────┐ ┌─────────────┐ ┌───────────────┐  │  │
│  │  │ Commit      │ │ Harvey      │ │ Copilot     │ │ Survey        │  │  │
│  │  │ Generator   │ │ Generator   │ │ Generator   │ │ Generator     │  │  │
│  │  │ (existing)  │ │ (NEW)       │ │ (NEW)       │ │ (NEW)         │  │  │
│  │  └─────────────┘ └─────────────┘ └─────────────┘ └───────────────┘  │  │
│  └──────────────────────────────────────────────────────────────────────┘  │
│                                        │                                     │
│                                        ▼                                     │
│  ┌──────────────────────────────────────────────────────────────────────┐  │
│  │                         In-Memory Storage                             │  │
│  │  ┌─────────────┐ ┌─────────────┐ ┌─────────────┐ ┌───────────────┐  │  │
│  │  │ Commits     │ │ Harvey      │ │ Copilot     │ │ Export Jobs   │  │  │
│  │  │ PRs         │ │ Events      │ │ Usage       │ │ Survey Data   │  │  │
│  │  └─────────────┘ └─────────────┘ └─────────────┘ └───────────────┘  │  │
│  └──────────────────────────────────────────────────────────────────────┘  │
│                                        │                                     │
│  ┌──────────────────────────────────────────────────────────────────────┐  │
│  │                           HTTP Router                                 │  │
│  │  ┌─────────────┐ ┌─────────────┐ ┌─────────────┐ ┌───────────────┐  │  │
│  │  │ /cursor/*   │ │ /harvey/*   │ │ /reports/*  │ │ /API/v3/*     │  │  │
│  │  │ (existing)  │ │ (NEW)       │ │ (NEW)       │ │ (NEW)         │  │  │
│  │  └─────────────┘ └─────────────┘ └─────────────┘ └───────────────┘  │  │
│  └──────────────────────────────────────────────────────────────────────┘  │
│                                                                              │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## 1. Harvey API Simulator

### 1.1 API Contract

#### Endpoint: GET /harvey/api/v1/history/usage

**Request:**
```http
GET /harvey/api/v1/history/usage HTTP/1.1
Host: localhost:8080
Authorization: Bearer <token>

Query Parameters:
  from: string (optional) - Start date YYYY-MM-DD
  to: string (optional) - End date YYYY-MM-DD
  user: string (optional) - Filter by user email
  task: string (optional) - Filter by task type (Assist|Draft|Review|Research)
  page: int (optional) - Page number (default 1)
  page_size: int (optional) - Items per page (default 100, max 500)
```

**Response:**
```json
{
  "data": [
    {
      "event_id": 103230489,
      "message_ID": "ab12a1ab-abcd-1a12-1234-1234ab123456",
      "Time": "2026-01-15T09:30:00.000Z",
      "User": "attorney@lawfirm.com",
      "Task": "Review",
      "Client Matter #": 2024.789,
      "Source": "Files",
      "Number of documents": 3,
      "Feedback Comments": "",
      "Feedback Sentiment": "positive"
    }
  ],
  "pagination": {
    "page": 1,
    "pageSize": 100,
    "totalPages": 3,
    "totalItems": 250,
    "hasNextPage": true,
    "hasPreviousPage": false
  },
  "params": {
    "from": "2026-01-01",
    "to": "2026-01-31"
  }
}
```

### 1.2 Data Model

```go
// internal/models/harvey.go

package models

import "time"

// HarveyTask represents the type of AI assistant task
type HarveyTask string

const (
    HarveyTaskAssist   HarveyTask = "Assist"   // General questions
    HarveyTaskDraft    HarveyTask = "Draft"    // Document drafting
    HarveyTaskReview   HarveyTask = "Review"   // Contract review
    HarveyTaskResearch HarveyTask = "Research" // Legal research
)

// HarveySource represents the data source for the task
type HarveySource string

const (
    HarveySourceFiles     HarveySource = "Files"     // Uploaded documents
    HarveySourceWeb       HarveySource = "Web"       // Web search
    HarveySourceKnowledge HarveySource = "Knowledge" // Knowledge base
)

// HarveySentiment represents feedback sentiment
type HarveySentiment string

const (
    HarveySentimentPositive HarveySentiment = "positive"
    HarveySentimentNegative HarveySentiment = "negative"
    HarveySentimentNeutral  HarveySentiment = "neutral"
)

// HarveyUsageEvent represents a single AI assistant usage event
type HarveyUsageEvent struct {
    EventID           int64           `json:"event_id"`
    MessageID         string          `json:"message_ID"`
    Time              time.Time       `json:"Time"`
    User              string          `json:"User"`
    Task              HarveyTask      `json:"Task"`
    ClientMatter      float64         `json:"Client Matter #"`
    Source            HarveySource    `json:"Source"`
    NumberOfDocuments int             `json:"Number of documents"`
    FeedbackComments  string          `json:"Feedback Comments"`
    FeedbackSentiment HarveySentiment `json:"Feedback Sentiment"`
}

// Validate ensures the event has required fields
func (e *HarveyUsageEvent) Validate() error {
    if e.EventID == 0 {
        return errors.New("event_id is required")
    }
    if e.User == "" {
        return errors.New("User is required")
    }
    if e.Task == "" {
        return errors.New("Task is required")
    }
    return nil
}
```

### 1.3 Generator Design

```go
// internal/generator/harvey_generator.go

package generator

import (
    "math/rand"
    "time"

    "github.com/cursor-analytics-platform/services/cursor-sim/internal/models"
    "github.com/cursor-analytics-platform/services/cursor-sim/internal/seed"
)

// HarveyGenerator generates Harvey AI usage events
type HarveyGenerator struct {
    seed   *seed.SeedData
    rng    *rand.Rand
    nextID int64
}

// HarveyConfig contains generation parameters
type HarveyConfig struct {
    EventsPerUserPerDay float64            // Average events per user per day
    TaskDistribution    map[models.HarveyTask]float64
    SentimentRates      map[models.HarveySentiment]float64
}

// DefaultHarveyConfig returns default generation parameters
func DefaultHarveyConfig() HarveyConfig {
    return HarveyConfig{
        EventsPerUserPerDay: 5.0,
        TaskDistribution: map[models.HarveyTask]float64{
            models.HarveyTaskAssist:   0.35,
            models.HarveyTaskDraft:    0.30,
            models.HarveyTaskReview:   0.25,
            models.HarveyTaskResearch: 0.10,
        },
        SentimentRates: map[models.HarveySentiment]float64{
            models.HarveySentimentPositive: 0.70,
            models.HarveySentimentNeutral:  0.20,
            models.HarveySentimentNegative: 0.10,
        },
    }
}

// NewHarveyGenerator creates a new generator
func NewHarveyGenerator(seedData *seed.SeedData) *HarveyGenerator {
    return NewHarveyGeneratorWithSeed(seedData, time.Now().UnixNano())
}

// NewHarveyGeneratorWithSeed creates a new generator with specific random seed
func NewHarveyGeneratorWithSeed(seedData *seed.SeedData, randSeed int64) *HarveyGenerator {
    return &HarveyGenerator{
        seed:   seedData,
        rng:    rand.New(rand.NewSource(randSeed)),
        nextID: 100000000 + rand.Int63n(1000000),
    }
}

// GenerateEvents generates Harvey usage events for the given time range
func (g *HarveyGenerator) GenerateEvents(from, to time.Time, config HarveyConfig) []models.HarveyUsageEvent {
    var events []models.HarveyUsageEvent

    // Get Harvey users from seed
    harveyUsers := g.getHarveyUsers()
    if len(harveyUsers) == 0 {
        return events
    }

    // Generate events for each day in range
    for day := from; !day.After(to); day = day.AddDate(0, 0, 1) {
        for _, user := range harveyUsers {
            userEvents := g.generateUserDayEvents(user, day, config)
            events = append(events, userEvents...)
        }
    }

    // Sort by time
    sort.Slice(events, func(i, j int) bool {
        return events[i].Time.Before(events[j].Time)
    })

    return events
}

// generateUserDayEvents generates events for a single user on a single day
func (g *HarveyGenerator) generateUserDayEvents(user seed.HarveyUser, day time.Time, config HarveyConfig) []models.HarveyUsageEvent {
    // Poisson-distributed event count
    count := g.poissonSample(config.EventsPerUserPerDay * user.ActivityMultiplier)

    var events []models.HarveyUsageEvent
    for i := 0; i < count; i++ {
        event := g.generateSingleEvent(user, day, config)
        events = append(events, event)
    }

    return events
}

// generateSingleEvent creates a single usage event
func (g *HarveyGenerator) generateSingleEvent(user seed.HarveyUser, day time.Time, config HarveyConfig) models.HarveyUsageEvent {
    g.nextID++

    // Random time during working hours
    hour := 8 + g.rng.Intn(10)  // 8 AM to 6 PM
    minute := g.rng.Intn(60)
    eventTime := time.Date(day.Year(), day.Month(), day.Day(), hour, minute, 0, 0, day.Location())

    task := g.selectTask(config.TaskDistribution)
    sentiment := g.selectSentiment(config.SentimentRates)

    return models.HarveyUsageEvent{
        EventID:           g.nextID,
        MessageID:         g.generateUUID(),
        Time:              eventTime,
        User:              user.Email,
        Task:              task,
        ClientMatter:      g.generateClientMatter(user),
        Source:            g.selectSource(),
        NumberOfDocuments: g.generateDocCount(task),
        FeedbackComments:  g.generateFeedbackComment(sentiment),
        FeedbackSentiment: sentiment,
    }
}
```

### 1.4 Handler Design

```go
// internal/api/harvey/handlers.go

package harvey

import (
    "net/http"
    "time"

    "github.com/cursor-analytics-platform/services/cursor-sim/internal/api"
    "github.com/cursor-analytics-platform/services/cursor-sim/internal/generator"
    "github.com/cursor-analytics-platform/services/cursor-sim/internal/storage"
)

// UsageHandler returns an HTTP handler for GET /harvey/api/v1/history/usage
func UsageHandler(store *storage.Store, gen *generator.HarveyGenerator) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // Parse query parameters
        from, to, err := parseDateRange(r)
        if err != nil {
            api.RespondError(w, http.StatusBadRequest, err.Error())
            return
        }

        // Get events from storage
        events := store.GetHarveyEvents(from, to)

        // Apply filters
        if user := r.URL.Query().Get("user"); user != "" {
            events = filterByUser(events, user)
        }
        if task := r.URL.Query().Get("task"); task != "" {
            events = filterByTask(events, task)
        }

        // Paginate
        page, pageSize := api.ParsePagination(r)
        paginatedEvents, pagination := api.Paginate(events, page, pageSize)

        // Respond
        response := map[string]interface{}{
            "data":       paginatedEvents,
            "pagination": pagination,
            "params": map[string]string{
                "from": from.Format("2006-01-02"),
                "to":   to.Format("2006-01-02"),
            },
        }
        api.RespondJSON(w, http.StatusOK, response)
    })
}
```

---

## 2. Microsoft 365 Copilot Usage API Simulator

### 2.1 API Contract (Microsoft Graph Beta)

#### Endpoint: GET /reports/getMicrosoft365CopilotUsageUserDetail(period='{period}')

**Request:**
```http
GET /reports/getMicrosoft365CopilotUsageUserDetail(period='D30')?$format=application/json HTTP/1.1
Host: localhost:8080
Authorization: Bearer <token>

Query Parameters:
  period: string (required) - D7, D30, D90, D180, or ALL
  $format: string (optional) - application/json (default) or text/csv
```

**JSON Response (200 OK):**
```json
{
  "@odata.nextLink": "https://localhost:8080/reports/getMicrosoft365CopilotUsageUserDetail(period='D30')?$format=application/json&$skiptoken=MDoyOg",
  "value": [
    {
      "reportRefreshDate": "2026-01-09",
      "reportPeriod": 30,
      "userPrincipalName": "user@company.com",
      "displayName": "Jane Developer",
      "lastActivityDate": "2026-01-08",
      "microsoftTeamsCopilotLastActivityDate": "2026-01-08",
      "wordCopilotLastActivityDate": "2026-01-05",
      "excelCopilotLastActivityDate": null,
      "powerPointCopilotLastActivityDate": "2025-12-20",
      "outlookCopilotLastActivityDate": "2026-01-07",
      "oneNoteCopilotLastActivityDate": null,
      "loopCopilotLastActivityDate": null,
      "copilotChatLastActivityDate": "2026-01-09",
      "copilotActivityUserDetailsByPeriod": [
        {
          "reportPeriod": 30
        }
      ]
    }
  ]
}
```

**CSV Response (302 Found):**
```http
HTTP/1.1 302 Found
Location: http://localhost:8080/reports/download/copilot-usage-abc123.csv
```

**CSV Content:**
```csv
Report Refresh Date,Report Period,User Principal Name,Display Name,Last Activity Date,Microsoft Teams Copilot Last Activity Date,Word Copilot Last Activity Date,Excel Copilot Last Activity Date,PowerPoint Copilot Last Activity Date,Outlook Copilot Last Activity Date,OneNote Copilot Last Activity Date,Loop Copilot Last Activity Date,Copilot Chat Last Activity Date
2026-01-09,30,user@company.com,Jane Developer,2026-01-08,2026-01-08,2026-01-05,,2025-12-20,2026-01-07,,,2026-01-09
```

### 2.2 Data Model

```go
// internal/models/copilot.go

package models

import "time"

// CopilotApp represents a Copilot-enabled application
type CopilotApp string

const (
    CopilotAppTeams      CopilotApp = "microsoftTeams"
    CopilotAppWord       CopilotApp = "word"
    CopilotAppExcel      CopilotApp = "excel"
    CopilotAppPowerPoint CopilotApp = "powerPoint"
    CopilotAppOutlook    CopilotApp = "outlook"
    CopilotAppOneNote    CopilotApp = "oneNote"
    CopilotAppLoop       CopilotApp = "loop"
    CopilotAppChat       CopilotApp = "copilotChat"
)

// AllCopilotApps returns all Copilot-enabled applications
func AllCopilotApps() []CopilotApp {
    return []CopilotApp{
        CopilotAppTeams,
        CopilotAppWord,
        CopilotAppExcel,
        CopilotAppPowerPoint,
        CopilotAppOutlook,
        CopilotAppOneNote,
        CopilotAppLoop,
        CopilotAppChat,
    }
}

// CopilotReportPeriod represents valid report periods
type CopilotReportPeriod string

const (
    CopilotPeriodD7   CopilotReportPeriod = "D7"
    CopilotPeriodD30  CopilotReportPeriod = "D30"
    CopilotPeriodD90  CopilotReportPeriod = "D90"
    CopilotPeriodD180 CopilotReportPeriod = "D180"
    CopilotPeriodAll  CopilotReportPeriod = "ALL"
)

// Days returns the number of days for this period
func (p CopilotReportPeriod) Days() int {
    switch p {
    case CopilotPeriodD7:
        return 7
    case CopilotPeriodD30:
        return 30
    case CopilotPeriodD90:
        return 90
    case CopilotPeriodD180:
        return 180
    case CopilotPeriodAll:
        return 180 // ALL includes all available data
    default:
        return 30
    }
}

// CopilotUsageUserDetail represents user-level Copilot usage data
// Matches Microsoft Graph API beta response schema
type CopilotUsageUserDetail struct {
    ReportRefreshDate                     string  `json:"reportRefreshDate"`
    ReportPeriod                          int     `json:"reportPeriod"`
    UserPrincipalName                     string  `json:"userPrincipalName"`
    DisplayName                           string  `json:"displayName"`
    LastActivityDate                      *string `json:"lastActivityDate"`
    MicrosoftTeamsCopilotLastActivityDate *string `json:"microsoftTeamsCopilotLastActivityDate"`
    WordCopilotLastActivityDate           *string `json:"wordCopilotLastActivityDate"`
    ExcelCopilotLastActivityDate          *string `json:"excelCopilotLastActivityDate"`
    PowerPointCopilotLastActivityDate     *string `json:"powerPointCopilotLastActivityDate"`
    OutlookCopilotLastActivityDate        *string `json:"outlookCopilotLastActivityDate"`
    OneNoteCopilotLastActivityDate        *string `json:"oneNoteCopilotLastActivityDate"`
    LoopCopilotLastActivityDate           *string `json:"loopCopilotLastActivityDate"`
    CopilotChatLastActivityDate           *string `json:"copilotChatLastActivityDate"`
    CopilotActivityUserDetailsByPeriod    []struct {
        ReportPeriod int `json:"reportPeriod"`
    } `json:"copilotActivityUserDetailsByPeriod"`
}

// CopilotUsageResponse represents the full API response
type CopilotUsageResponse struct {
    ODataNextLink *string                   `json:"@odata.nextLink,omitempty"`
    Value         []CopilotUsageUserDetail  `json:"value"`
}

// HasLastActivityDate returns true if user has any activity
func (u *CopilotUsageUserDetail) HasLastActivityDate() bool {
    return u.LastActivityDate != nil
}

// GetAppLastActivityDate returns the last activity date for a specific app
func (u *CopilotUsageUserDetail) GetAppLastActivityDate(app CopilotApp) *string {
    switch app {
    case CopilotAppTeams:
        return u.MicrosoftTeamsCopilotLastActivityDate
    case CopilotAppWord:
        return u.WordCopilotLastActivityDate
    case CopilotAppExcel:
        return u.ExcelCopilotLastActivityDate
    case CopilotAppPowerPoint:
        return u.PowerPointCopilotLastActivityDate
    case CopilotAppOutlook:
        return u.OutlookCopilotLastActivityDate
    case CopilotAppOneNote:
        return u.OneNoteCopilotLastActivityDate
    case CopilotAppLoop:
        return u.LoopCopilotLastActivityDate
    case CopilotAppChat:
        return u.CopilotChatLastActivityDate
    default:
        return nil
    }
}
```

### 2.3 Generator Design

```go
// internal/generator/copilot_generator.go

package generator

import (
    "math/rand"
    "time"

    "github.com/cursor-analytics-platform/services/cursor-sim/internal/models"
    "github.com/cursor-analytics-platform/services/cursor-sim/internal/seed"
)

// CopilotGenerator generates Microsoft 365 Copilot usage data
type CopilotGenerator struct {
    seed *seed.SeedData
    rng  *rand.Rand
}

// CopilotConfig contains generation parameters
type CopilotConfig struct {
    // App adoption rates (probability a user uses each app)
    AppAdoptionRates map[models.CopilotApp]float64
    // Activity frequency (average days between uses per app)
    ActivityFrequency map[models.CopilotApp]float64
}

// DefaultCopilotConfig returns default generation parameters
func DefaultCopilotConfig() CopilotConfig {
    return CopilotConfig{
        AppAdoptionRates: map[models.CopilotApp]float64{
            models.CopilotAppTeams:      0.85, // High - used for meetings/chat
            models.CopilotAppOutlook:    0.80, // High - email summarization
            models.CopilotAppWord:       0.60, // Medium - document drafting
            models.CopilotAppChat:       0.55, // Medium - general AI assistant
            models.CopilotAppPowerPoint: 0.35, // Lower - less frequent
            models.CopilotAppExcel:      0.30, // Lower - specialized use
            models.CopilotAppLoop:       0.15, // Low - newer product
            models.CopilotAppOneNote:    0.10, // Low - niche use
        },
        ActivityFrequency: map[models.CopilotApp]float64{
            models.CopilotAppTeams:      1.5,  // Every 1-2 days
            models.CopilotAppOutlook:    2.0,  // Every 2 days
            models.CopilotAppWord:       5.0,  // Weekly
            models.CopilotAppChat:       3.0,  // Every 3 days
            models.CopilotAppPowerPoint: 14.0, // Biweekly
            models.CopilotAppExcel:      10.0, // Every 10 days
            models.CopilotAppLoop:       21.0, // Less frequent
            models.CopilotAppOneNote:    14.0, // Biweekly
        },
    }
}

// NewCopilotGenerator creates a new generator
func NewCopilotGenerator(seedData *seed.SeedData) *CopilotGenerator {
    return NewCopilotGeneratorWithSeed(seedData, time.Now().UnixNano())
}

// NewCopilotGeneratorWithSeed creates a new generator with specific random seed
func NewCopilotGeneratorWithSeed(seedData *seed.SeedData, randSeed int64) *CopilotGenerator {
    return &CopilotGenerator{
        seed: seedData,
        rng:  rand.New(rand.NewSource(randSeed)),
    }
}

// GenerateUsageReport generates Copilot usage data for the specified period
func (g *CopilotGenerator) GenerateUsageReport(period models.CopilotReportPeriod, config CopilotConfig) []models.CopilotUsageUserDetail {
    var users []models.CopilotUsageUserDetail

    // Get M365 users from seed
    m365Users := g.getM365Users()
    periodDays := period.Days()
    reportDate := time.Now().Format("2006-01-02")

    for _, user := range m365Users {
        detail := g.generateUserDetail(user, periodDays, reportDate, config)
        users = append(users, detail)
    }

    return users
}

// generateUserDetail generates usage data for a single user
func (g *CopilotGenerator) generateUserDetail(user seed.M365User, periodDays int, reportDate string, config CopilotConfig) models.CopilotUsageUserDetail {
    detail := models.CopilotUsageUserDetail{
        ReportRefreshDate: reportDate,
        ReportPeriod:      periodDays,
        UserPrincipalName: user.Email,
        DisplayName:       user.DisplayName,
        CopilotActivityUserDetailsByPeriod: []struct {
            ReportPeriod int `json:"reportPeriod"`
        }{
            {ReportPeriod: periodDays},
        },
    }

    // Generate last activity dates for each app
    var latestActivity *time.Time

    for _, app := range models.AllCopilotApps() {
        if g.userAdoptsApp(user, app, config) {
            lastDate := g.generateLastActivityDate(periodDays, config.ActivityFrequency[app])
            if lastDate != nil {
                dateStr := lastDate.Format("2006-01-02")
                g.setAppActivityDate(&detail, app, &dateStr)

                if latestActivity == nil || lastDate.After(*latestActivity) {
                    latestActivity = lastDate
                }
            }
        }
    }

    // Set overall last activity date
    if latestActivity != nil {
        dateStr := latestActivity.Format("2006-01-02")
        detail.LastActivityDate = &dateStr
    }

    return detail
}

// setAppActivityDate sets the activity date for a specific app
func (g *CopilotGenerator) setAppActivityDate(detail *models.CopilotUsageUserDetail, app models.CopilotApp, date *string) {
    switch app {
    case models.CopilotAppTeams:
        detail.MicrosoftTeamsCopilotLastActivityDate = date
    case models.CopilotAppWord:
        detail.WordCopilotLastActivityDate = date
    case models.CopilotAppExcel:
        detail.ExcelCopilotLastActivityDate = date
    case models.CopilotAppPowerPoint:
        detail.PowerPointCopilotLastActivityDate = date
    case models.CopilotAppOutlook:
        detail.OutlookCopilotLastActivityDate = date
    case models.CopilotAppOneNote:
        detail.OneNoteCopilotLastActivityDate = date
    case models.CopilotAppLoop:
        detail.LoopCopilotLastActivityDate = date
    case models.CopilotAppChat:
        detail.CopilotChatLastActivityDate = date
    }
}
```

### 2.4 Handler Design

```go
// internal/api/microsoft/copilot_handlers.go

package microsoft

import (
    "encoding/csv"
    "net/http"
    "regexp"
    "strings"

    "github.com/cursor-analytics-platform/services/cursor-sim/internal/api"
    "github.com/cursor-analytics-platform/services/cursor-sim/internal/generator"
    "github.com/cursor-analytics-platform/services/cursor-sim/internal/models"
    "github.com/cursor-analytics-platform/services/cursor-sim/internal/storage"
)

// periodPattern matches the OData function call format
var periodPattern = regexp.MustCompile(`getMicrosoft365CopilotUsageUserDetail\(period='([^']+)'\)`)

// UsageUserDetailHandler returns handler for the Graph API endpoint
func UsageUserDetailHandler(store *storage.Store, gen *generator.CopilotGenerator) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // Extract period from path
        matches := periodPattern.FindStringSubmatch(r.URL.Path)
        if len(matches) < 2 {
            api.RespondError(w, http.StatusBadRequest, "Invalid endpoint format")
            return
        }
        periodStr := matches[1]
        period := models.CopilotReportPeriod(periodStr)

        // Validate period
        if !isValidPeriod(period) {
            api.RespondError(w, http.StatusBadRequest, "Invalid period value")
            return
        }

        // Get format preference
        format := r.URL.Query().Get("$format")
        if format == "" {
            format = "application/json"
        }

        // Get usage data from storage
        usageData := store.GetCopilotUsage(period)

        // Handle CSV format with redirect
        if format == "text/csv" {
            // Generate download token
            downloadID := generateDownloadToken()
            store.StoreCopilotCSVData(downloadID, usageData)

            // Return 302 redirect
            redirectURL := "/reports/download/" + downloadID + ".csv"
            http.Redirect(w, r, redirectURL, http.StatusFound)
            return
        }

        // Handle JSON format with pagination
        page, pageSize := api.ParsePagination(r)
        skipToken := r.URL.Query().Get("$skiptoken")

        paginatedData, nextLink := paginateWithSkipToken(usageData, page, pageSize, skipToken, r.URL.Path)

        response := models.CopilotUsageResponse{
            Value: paginatedData,
        }
        if nextLink != "" {
            response.ODataNextLink = &nextLink
        }

        api.RespondJSON(w, http.StatusOK, response)
    })
}

// CSVDownloadHandler serves pre-generated CSV files
func CSVDownloadHandler(store *storage.Store) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // Extract download ID from path
        downloadID := extractDownloadID(r.URL.Path)

        // Get cached CSV data
        usageData, exists := store.GetCopilotCSVData(downloadID)
        if !exists {
            api.RespondError(w, http.StatusNotFound, "Download not found or expired")
            return
        }

        // Write CSV response
        w.Header().Set("Content-Type", "text/csv")
        w.Header().Set("Content-Disposition", "attachment; filename=copilot-usage.csv")

        writer := csv.NewWriter(w)

        // Write header
        header := []string{
            "Report Refresh Date",
            "Report Period",
            "User Principal Name",
            "Display Name",
            "Last Activity Date",
            "Microsoft Teams Copilot Last Activity Date",
            "Word Copilot Last Activity Date",
            "Excel Copilot Last Activity Date",
            "PowerPoint Copilot Last Activity Date",
            "Outlook Copilot Last Activity Date",
            "OneNote Copilot Last Activity Date",
            "Loop Copilot Last Activity Date",
            "Copilot Chat Last Activity Date",
        }
        writer.Write(header)

        // Write data rows
        for _, user := range usageData {
            row := []string{
                user.ReportRefreshDate,
                fmt.Sprintf("%d", user.ReportPeriod),
                user.UserPrincipalName,
                user.DisplayName,
                nullableString(user.LastActivityDate),
                nullableString(user.MicrosoftTeamsCopilotLastActivityDate),
                nullableString(user.WordCopilotLastActivityDate),
                nullableString(user.ExcelCopilotLastActivityDate),
                nullableString(user.PowerPointCopilotLastActivityDate),
                nullableString(user.OutlookCopilotLastActivityDate),
                nullableString(user.OneNoteCopilotLastActivityDate),
                nullableString(user.LoopCopilotLastActivityDate),
                nullableString(user.CopilotChatLastActivityDate),
            }
            writer.Write(row)
        }

        writer.Flush()
    })
}
```

---

## 3. Qualtrics Survey Export API Simulator

### 3.1 API Contract (3-Step Async State Machine)

#### Step 1: POST /API/v3/surveys/{surveyId}/export-responses

**Request:**
```http
POST /API/v3/surveys/SV_abc123/export-responses HTTP/1.1
Host: localhost:8080
Authorization: Bearer <token>
Content-Type: application/json

{
  "format": "csv"
}
```

**Response (200 OK):**
```json
{
  "result": {
    "progressId": "ES_xyz789",
    "percentComplete": 0,
    "status": "inProgress"
  },
  "meta": {
    "requestId": "req_abc123",
    "httpStatus": "200 - OK"
  }
}
```

#### Step 2: GET /API/v3/surveys/{surveyId}/export-responses/{progressId}

**Request:**
```http
GET /API/v3/surveys/SV_abc123/export-responses/ES_xyz789 HTTP/1.1
Host: localhost:8080
Authorization: Bearer <token>
```

**Response (In Progress):**
```json
{
  "result": {
    "progressId": "ES_xyz789",
    "percentComplete": 45,
    "status": "inProgress"
  },
  "meta": {
    "requestId": "req_def456",
    "httpStatus": "200 - OK"
  }
}
```

**Response (Complete):**
```json
{
  "result": {
    "fileId": "FILE_abc123",
    "percentComplete": 100,
    "status": "complete"
  },
  "meta": {
    "requestId": "req_ghi789",
    "httpStatus": "200 - OK"
  }
}
```

#### Step 3: GET /API/v3/surveys/{surveyId}/export-responses/{fileId}/file

**Request:**
```http
GET /API/v3/surveys/SV_abc123/export-responses/FILE_abc123/file HTTP/1.1
Host: localhost:8080
Authorization: Bearer <token>
```

**Response:**
```
HTTP/1.1 200 OK
Content-Type: application/zip
Content-Disposition: attachment; filename=responses.zip

[ZIP file contents with survey_responses.csv]
```

### 3.2 Data Model

```go
// internal/models/qualtrics.go

package models

import (
    "archive/zip"
    "bytes"
    "encoding/csv"
    "time"
)

// ExportJobStatus represents the status of an export job
type ExportJobStatus string

const (
    ExportStatusInProgress ExportJobStatus = "inProgress"
    ExportStatusComplete   ExportJobStatus = "complete"
    ExportStatusFailed     ExportJobStatus = "failed"
)

// ExportJob represents an active survey export job
type ExportJob struct {
    ProgressID      string          `json:"progressId"`
    SurveyID        string          `json:"surveyId"`
    Status          ExportJobStatus `json:"status"`
    PercentComplete int             `json:"percentComplete"`
    FileID          string          `json:"fileId,omitempty"`
    StartedAt       time.Time       `json:"startedAt"`
    CompletedAt     *time.Time      `json:"completedAt,omitempty"`
    Error           string          `json:"error,omitempty"`
}

// ExportStartResponse is the response for starting an export
type ExportStartResponse struct {
    Result ExportJobResult    `json:"result"`
    Meta   ExportResponseMeta `json:"meta"`
}

// ExportProgressResponse is the response for checking export progress
type ExportProgressResponse struct {
    Result ExportJobResult    `json:"result"`
    Meta   ExportResponseMeta `json:"meta"`
}

// ExportJobResult contains the export job details
type ExportJobResult struct {
    ProgressID      string `json:"progressId,omitempty"`
    FileID          string `json:"fileId,omitempty"`
    PercentComplete int    `json:"percentComplete"`
    Status          string `json:"status"`
}

// ExportResponseMeta contains response metadata
type ExportResponseMeta struct {
    RequestID  string `json:"requestId"`
    HTTPStatus string `json:"httpStatus"`
}

// SurveyResponse represents a single survey response
type SurveyResponse struct {
    ResponseID     string    `json:"responseId"`
    StartDate      time.Time `json:"startDate"`
    EndDate        time.Time `json:"endDate"`
    Status         string    `json:"status"`
    Progress       int       `json:"progress"`
    Duration       int       `json:"duration"` // seconds
    Finished       bool      `json:"finished"`
    RecordedDate   time.Time `json:"recordedDate"`

    // Demographics
    RespondentEmail string `json:"respondentEmail"`
    Department      string `json:"department"`
    Role            string `json:"role"`
    Tenure          string `json:"tenure"`

    // AI Tool Satisfaction (1-5 scale)
    OverallAISatisfaction int    `json:"q1_overall_satisfaction"`
    CursorSatisfaction    int    `json:"q2_cursor_satisfaction"`
    CopilotSatisfaction   int    `json:"q3_copilot_satisfaction"`
    HarveySatisfaction    int    `json:"q4_harvey_satisfaction"`

    // Usage Frequency
    AIToolsUsageFrequency string `json:"q5_usage_frequency"` // Daily, Weekly, Monthly, Rarely
    MostUsedTool          string `json:"q6_most_used_tool"`  // Cursor, Copilot, Harvey, Other

    // Free Text
    PositiveFeedback  string `json:"q7_positive_feedback"`
    ImprovementAreas  string `json:"q8_improvement_areas"`
    AdditionalComment string `json:"q9_additional_comments"`
}

// SurveyConfig defines the structure of a simulated survey
type SurveyConfig struct {
    SurveyID       string   `json:"surveyId"`
    SurveyName     string   `json:"surveyName"`
    ResponseCount  int      `json:"responseCount"`
    Questions      []string `json:"questions"`
}

// GenerateZIPFile creates a ZIP containing the survey responses CSV
func GenerateZIPFile(responses []SurveyResponse) ([]byte, error) {
    buf := new(bytes.Buffer)
    zipWriter := zip.NewWriter(buf)

    // Create CSV file inside ZIP
    csvFile, err := zipWriter.Create("survey_responses.csv")
    if err != nil {
        return nil, err
    }

    csvWriter := csv.NewWriter(csvFile)

    // Write header
    header := []string{
        "ResponseID", "StartDate", "EndDate", "Duration", "Status",
        "Email", "Department", "Role", "Tenure",
        "Q1_OverallSatisfaction", "Q2_CursorSatisfaction",
        "Q3_CopilotSatisfaction", "Q4_HarveySatisfaction",
        "Q5_UsageFrequency", "Q6_MostUsedTool",
        "Q7_PositiveFeedback", "Q8_ImprovementAreas", "Q9_AdditionalComments",
    }
    csvWriter.Write(header)

    // Write data
    for _, r := range responses {
        row := []string{
            r.ResponseID,
            r.StartDate.Format(time.RFC3339),
            r.EndDate.Format(time.RFC3339),
            fmt.Sprintf("%d", r.Duration),
            r.Status,
            r.RespondentEmail,
            r.Department,
            r.Role,
            r.Tenure,
            fmt.Sprintf("%d", r.OverallAISatisfaction),
            fmt.Sprintf("%d", r.CursorSatisfaction),
            fmt.Sprintf("%d", r.CopilotSatisfaction),
            fmt.Sprintf("%d", r.HarveySatisfaction),
            r.AIToolsUsageFrequency,
            r.MostUsedTool,
            r.PositiveFeedback,
            r.ImprovementAreas,
            r.AdditionalComment,
        }
        csvWriter.Write(row)
    }

    csvWriter.Flush()
    zipWriter.Close()

    return buf.Bytes(), nil
}
```

### 3.3 State Machine Design

```go
// internal/services/qualtrics_export.go

package services

import (
    "sync"
    "time"

    "github.com/cursor-analytics-platform/services/cursor-sim/internal/models"
)

// ExportJobManager manages the lifecycle of export jobs
type ExportJobManager struct {
    jobs        map[string]*models.ExportJob
    files       map[string][]byte  // fileId -> ZIP content
    mu          sync.RWMutex
    generator   *SurveyGenerator

    // Configuration
    ProgressIncrementMin int   // Minimum progress increment per poll
    ProgressIncrementMax int   // Maximum progress increment per poll
    SimulatedDelayMs     int64 // Simulated processing delay between polls
}

// NewExportJobManager creates a new job manager
func NewExportJobManager(gen *SurveyGenerator) *ExportJobManager {
    return &ExportJobManager{
        jobs:                 make(map[string]*models.ExportJob),
        files:                make(map[string][]byte),
        generator:            gen,
        ProgressIncrementMin: 15,
        ProgressIncrementMax: 35,
        SimulatedDelayMs:     500, // 500ms simulated processing
    }
}

// StartExport initiates a new export job
func (m *ExportJobManager) StartExport(surveyID string) (*models.ExportJob, error) {
    m.mu.Lock()
    defer m.mu.Unlock()

    progressID := "ES_" + generateRandomID(8)

    job := &models.ExportJob{
        ProgressID:      progressID,
        SurveyID:        surveyID,
        Status:          models.ExportStatusInProgress,
        PercentComplete: 0,
        StartedAt:       time.Now(),
    }

    m.jobs[progressID] = job

    // Start background processing (simulated)
    go m.processExport(job)

    return job, nil
}

// GetProgress returns the current progress of an export job
func (m *ExportJobManager) GetProgress(progressID string) (*models.ExportJob, error) {
    m.mu.RLock()
    job, exists := m.jobs[progressID]
    m.mu.RUnlock()

    if !exists {
        return nil, errors.New("export job not found")
    }

    // Simulate progress advancement
    m.advanceProgress(job)

    return job, nil
}

// GetFile returns the generated file for a completed export
func (m *ExportJobManager) GetFile(fileID string) ([]byte, error) {
    m.mu.RLock()
    defer m.mu.RUnlock()

    data, exists := m.files[fileID]
    if !exists {
        return nil, errors.New("file not found or expired")
    }

    return data, nil
}

// processExport simulates the export generation process
func (m *ExportJobManager) processExport(job *models.ExportJob) {
    // This runs in background but actual progress is
    // advanced on-demand when polled (for predictable testing)
}

// advanceProgress advances the job progress on each poll
func (m *ExportJobManager) advanceProgress(job *models.ExportJob) {
    m.mu.Lock()
    defer m.mu.Unlock()

    if job.Status != models.ExportStatusInProgress {
        return
    }

    // Calculate progress increment
    increment := m.ProgressIncrementMin +
        rand.Intn(m.ProgressIncrementMax - m.ProgressIncrementMin)

    job.PercentComplete += increment

    if job.PercentComplete >= 100 {
        job.PercentComplete = 100
        job.Status = models.ExportStatusComplete
        now := time.Now()
        job.CompletedAt = &now

        // Generate the file
        fileID := "FILE_" + generateRandomID(8)
        job.FileID = fileID

        responses := m.generator.GenerateSurveyResponses(job.SurveyID)
        zipData, _ := models.GenerateZIPFile(responses)
        m.files[fileID] = zipData
    }
}
```

### 3.4 Handler Design

```go
// internal/api/qualtrics/handlers.go

package qualtrics

import (
    "net/http"
    "strings"

    "github.com/cursor-analytics-platform/services/cursor-sim/internal/api"
    "github.com/cursor-analytics-platform/services/cursor-sim/internal/services"
)

// ExportHandlers provides HTTP handlers for Qualtrics export API
type ExportHandlers struct {
    jobManager *services.ExportJobManager
}

// NewExportHandlers creates new handlers
func NewExportHandlers(manager *services.ExportJobManager) *ExportHandlers {
    return &ExportHandlers{jobManager: manager}
}

// StartExportHandler handles POST /API/v3/surveys/{surveyId}/export-responses
func (h *ExportHandlers) StartExportHandler() http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        if r.Method != http.MethodPost {
            api.RespondError(w, http.StatusMethodNotAllowed, "Method not allowed")
            return
        }

        surveyID := extractSurveyID(r.URL.Path)
        if surveyID == "" {
            api.RespondError(w, http.StatusBadRequest, "Invalid survey ID")
            return
        }

        job, err := h.jobManager.StartExport(surveyID)
        if err != nil {
            api.RespondError(w, http.StatusInternalServerError, err.Error())
            return
        }

        response := models.ExportStartResponse{
            Result: models.ExportJobResult{
                ProgressID:      job.ProgressID,
                PercentComplete: job.PercentComplete,
                Status:          string(job.Status),
            },
            Meta: models.ExportResponseMeta{
                RequestID:  generateRequestID(),
                HTTPStatus: "200 - OK",
            },
        }

        api.RespondJSON(w, http.StatusOK, response)
    })
}

// ProgressHandler handles GET /API/v3/surveys/{surveyId}/export-responses/{progressId}
func (h *ExportHandlers) ProgressHandler() http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        if r.Method != http.MethodGet {
            api.RespondError(w, http.StatusMethodNotAllowed, "Method not allowed")
            return
        }

        progressID := extractProgressID(r.URL.Path)
        if progressID == "" {
            api.RespondError(w, http.StatusBadRequest, "Invalid progress ID")
            return
        }

        job, err := h.jobManager.GetProgress(progressID)
        if err != nil {
            api.RespondError(w, http.StatusNotFound, err.Error())
            return
        }

        result := models.ExportJobResult{
            PercentComplete: job.PercentComplete,
            Status:          string(job.Status),
        }

        if job.Status == models.ExportStatusComplete {
            result.FileID = job.FileID
        } else {
            result.ProgressID = job.ProgressID
        }

        response := models.ExportProgressResponse{
            Result: result,
            Meta: models.ExportResponseMeta{
                RequestID:  generateRequestID(),
                HTTPStatus: "200 - OK",
            },
        }

        api.RespondJSON(w, http.StatusOK, response)
    })
}

// FileDownloadHandler handles GET /API/v3/surveys/{surveyId}/export-responses/{fileId}/file
func (h *ExportHandlers) FileDownloadHandler() http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        if r.Method != http.MethodGet {
            api.RespondError(w, http.StatusMethodNotAllowed, "Method not allowed")
            return
        }

        fileID := extractFileID(r.URL.Path)
        if fileID == "" {
            api.RespondError(w, http.StatusBadRequest, "Invalid file ID")
            return
        }

        data, err := h.jobManager.GetFile(fileID)
        if err != nil {
            api.RespondError(w, http.StatusNotFound, err.Error())
            return
        }

        w.Header().Set("Content-Type", "application/zip")
        w.Header().Set("Content-Disposition", "attachment; filename=responses.zip")
        w.Header().Set("Content-Length", fmt.Sprintf("%d", len(data)))
        w.Write(data)
    })
}
```

---

## 4. Seed Schema Extension

### 4.1 Extended Seed Schema

```yaml
# testdata/enterprise_seed.yaml
version: "1.0.0"

developers:
  - user_id: "dev_001"
    email: "jane.dev@company.com"
    name: "Jane Developer"
    # ... existing fields ...

# NEW: Harvey AI Users
harvey_users:
  - user_id: "atty_001"
    email: "john.attorney@lawfirm.com"
    name: "John Attorney"
    role: "partner"
    practice_area: "corporate"  # corporate, litigation, ip, regulatory
    activity_multiplier: 1.2    # Relative activity level
    client_matters:
      - 2024.001
      - 2024.045
      - 2024.089

  - user_id: "atty_002"
    email: "sarah.associate@lawfirm.com"
    name: "Sarah Associate"
    role: "associate"
    practice_area: "litigation"
    activity_multiplier: 0.8
    client_matters:
      - 2024.012
      - 2024.023

# NEW: Microsoft 365 Tenant Configuration
m365_tenant:
  tenant_id: "tenant_abc123"
  display_name: "Acme Corporation"

  users:
    - user_id: "m365_001"
      email: "jane.dev@company.com"
      display_name: "Jane Developer"
      department: "Engineering"
      # Link to existing developer
      linked_developer_id: "dev_001"
      copilot_enabled: true
      copilot_apps:  # Override default adoption
        teams: true
        word: true
        excel: false
        powerpoint: false
        outlook: true
        onenote: false
        loop: false
        chat: true

    - user_id: "m365_002"
      email: "bob.manager@company.com"
      display_name: "Bob Manager"
      department: "Product"
      copilot_enabled: true
      # Uses default app adoption rates

# NEW: Qualtrics Survey Configuration
qualtrics:
  surveys:
    - survey_id: "SV_aitools_q1_2026"
      name: "AI Tools Satisfaction Survey Q1 2026"
      response_count: 150  # Number of responses to generate

      # Response distribution
      satisfaction_distribution:
        1: 0.05   # 5% very dissatisfied
        2: 0.10   # 10% dissatisfied
        3: 0.25   # 25% neutral
        4: 0.40   # 40% satisfied
        5: 0.20   # 20% very satisfied

      # Link responses to users
      respondent_pools:
        - pool: "developers"
          weight: 0.60  # 60% of responses from developers
        - pool: "harvey_users"
          weight: 0.20  # 20% from legal
        - pool: "m365_users"
          weight: 0.20  # 20% from M365 users
```

### 4.2 Seed Types Extension

```go
// internal/seed/types.go (additions)

// HarveyUser represents a Harvey AI platform user
type HarveyUser struct {
    UserID             string    `json:"user_id" yaml:"user_id"`
    Email              string    `json:"email" yaml:"email"`
    Name               string    `json:"name" yaml:"name"`
    Role               string    `json:"role" yaml:"role"`
    PracticeArea       string    `json:"practice_area" yaml:"practice_area"`
    ActivityMultiplier float64   `json:"activity_multiplier" yaml:"activity_multiplier"`
    ClientMatters      []float64 `json:"client_matters" yaml:"client_matters"`
}

// M365Tenant represents a Microsoft 365 tenant configuration
type M365Tenant struct {
    TenantID    string     `json:"tenant_id" yaml:"tenant_id"`
    DisplayName string     `json:"display_name" yaml:"display_name"`
    Users       []M365User `json:"users" yaml:"users"`
}

// M365User represents a Microsoft 365 user with Copilot access
type M365User struct {
    UserID             string          `json:"user_id" yaml:"user_id"`
    Email              string          `json:"email" yaml:"email"`
    DisplayName        string          `json:"display_name" yaml:"display_name"`
    Department         string          `json:"department" yaml:"department"`
    LinkedDeveloperID  string          `json:"linked_developer_id,omitempty" yaml:"linked_developer_id,omitempty"`
    CopilotEnabled     bool            `json:"copilot_enabled" yaml:"copilot_enabled"`
    CopilotApps        *M365CopilotApps `json:"copilot_apps,omitempty" yaml:"copilot_apps,omitempty"`
}

// M365CopilotApps specifies which Copilot apps are enabled
type M365CopilotApps struct {
    Teams      bool `json:"teams" yaml:"teams"`
    Word       bool `json:"word" yaml:"word"`
    Excel      bool `json:"excel" yaml:"excel"`
    PowerPoint bool `json:"powerpoint" yaml:"powerpoint"`
    Outlook    bool `json:"outlook" yaml:"outlook"`
    OneNote    bool `json:"onenote" yaml:"onenote"`
    Loop       bool `json:"loop" yaml:"loop"`
    Chat       bool `json:"chat" yaml:"chat"`
}

// QualtricsConfig contains survey configurations
type QualtricsConfig struct {
    Surveys []SurveyConfig `json:"surveys" yaml:"surveys"`
}

// SurveyConfig defines a survey for simulation
type SurveyConfig struct {
    SurveyID                 string                 `json:"survey_id" yaml:"survey_id"`
    Name                     string                 `json:"name" yaml:"name"`
    ResponseCount            int                    `json:"response_count" yaml:"response_count"`
    SatisfactionDistribution map[int]float64        `json:"satisfaction_distribution" yaml:"satisfaction_distribution"`
    RespondentPools          []RespondentPoolConfig `json:"respondent_pools" yaml:"respondent_pools"`
}

// RespondentPoolConfig defines a pool of respondents
type RespondentPoolConfig struct {
    Pool   string  `json:"pool" yaml:"pool"`
    Weight float64 `json:"weight" yaml:"weight"`
}

// SeedData extended
type SeedData struct {
    Version      string          `json:"version" yaml:"version"`
    Developers   []Developer     `json:"developers" yaml:"developers"`
    Repositories []Repository    `json:"repositories" yaml:"repositories"`
    // ... existing fields ...

    // NEW: External data source configurations
    HarveyUsers  []HarveyUser    `json:"harvey_users,omitempty" yaml:"harvey_users,omitempty"`
    M365Tenant   *M365Tenant     `json:"m365_tenant,omitempty" yaml:"m365_tenant,omitempty"`
    Qualtrics    *QualtricsConfig `json:"qualtrics,omitempty" yaml:"qualtrics,omitempty"`
}
```

---

## 5. File Structure

### 5.1 New Files to Create

```
services/cursor-sim/
├── internal/
│   ├── models/
│   │   ├── harvey.go              # Harvey data types
│   │   ├── harvey_test.go         # Harvey model tests
│   │   ├── copilot.go             # Copilot data types
│   │   ├── copilot_test.go        # Copilot model tests
│   │   ├── qualtrics.go           # Qualtrics data types
│   │   └── qualtrics_test.go      # Qualtrics model tests
│   │
│   ├── generator/
│   │   ├── harvey_generator.go         # Harvey event generator
│   │   ├── harvey_generator_test.go    # Harvey generator tests
│   │   ├── copilot_generator.go        # Copilot usage generator
│   │   ├── copilot_generator_test.go   # Copilot generator tests
│   │   ├── survey_generator.go         # Survey response generator
│   │   └── survey_generator_test.go    # Survey generator tests
│   │
│   ├── services/
│   │   ├── qualtrics_export.go         # Export job state machine
│   │   └── qualtrics_export_test.go    # State machine tests
│   │
│   ├── api/
│   │   ├── harvey/
│   │   │   ├── handlers.go        # Harvey API handlers
│   │   │   └── handlers_test.go   # Handler tests
│   │   │
│   │   ├── microsoft/
│   │   │   ├── copilot_handlers.go      # Copilot API handlers
│   │   │   └── copilot_handlers_test.go # Handler tests
│   │   │
│   │   └── qualtrics/
│   │       ├── handlers.go        # Qualtrics API handlers
│   │       └── handlers_test.go   # Handler tests
│   │
│   ├── seed/
│   │   └── types.go               # MODIFY: Add new seed types
│   │
│   └── storage/
│       └── store.go               # MODIFY: Add storage for new data
│
├── test/e2e/
│   ├── harvey_test.go             # Harvey E2E tests
│   ├── copilot_test.go            # Copilot E2E tests
│   └── qualtrics_test.go          # Qualtrics E2E tests
│
└── testdata/
    └── enterprise_seed.yaml       # Example seed with all data sources
```

### 5.2 Files to Modify

| File | Changes |
|------|---------|
| `internal/seed/types.go` | Add HarveyUser, M365Tenant, QualtricsConfig types |
| `internal/seed/loader.go` | Parse new seed sections |
| `internal/storage/store.go` | Add storage methods for new data types |
| `internal/server/router.go` | Register new API routes |
| `cmd/simulator/main.go` | Initialize new generators |
| `SPEC.md` | Document new endpoints |

---

## 6. Router Integration

```go
// internal/server/router.go (additions)

func NewRouter(store *storage.Store, seedData *seed.SeedData) *http.ServeMux {
    mux := http.NewServeMux()

    // Existing routes...

    // Harvey API routes
    if len(seedData.HarveyUsers) > 0 {
        harveyGen := generator.NewHarveyGenerator(seedData)
        mux.Handle("/harvey/api/v1/history/usage",
            authMiddleware(harvey.UsageHandler(store, harveyGen)))
    }

    // Microsoft Graph API routes (Copilot)
    if seedData.M365Tenant != nil {
        copilotGen := generator.NewCopilotGenerator(seedData)

        // Pattern for Graph API function call
        mux.Handle("/reports/getMicrosoft365CopilotUsageUserDetail",
            authMiddleware(microsoft.UsageUserDetailHandler(store, copilotGen)))
        mux.Handle("/reports/download/",
            microsoft.CSVDownloadHandler(store))
    }

    // Qualtrics API routes
    if seedData.Qualtrics != nil && len(seedData.Qualtrics.Surveys) > 0 {
        surveyGen := generator.NewSurveyGenerator(seedData)
        jobManager := services.NewExportJobManager(surveyGen)
        handlers := qualtrics.NewExportHandlers(jobManager)

        mux.Handle("/API/v3/surveys/", authMiddleware(
            qualtricsRouter(handlers)))
    }

    return mux
}

// qualtricsRouter routes Qualtrics endpoints
func qualtricsRouter(h *qualtrics.ExportHandlers) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        path := r.URL.Path

        // POST /API/v3/surveys/{surveyId}/export-responses
        if r.Method == "POST" && strings.HasSuffix(path, "/export-responses") {
            h.StartExportHandler().ServeHTTP(w, r)
            return
        }

        // GET /API/v3/surveys/{surveyId}/export-responses/{id}/file
        if strings.HasSuffix(path, "/file") {
            h.FileDownloadHandler().ServeHTTP(w, r)
            return
        }

        // GET /API/v3/surveys/{surveyId}/export-responses/{progressId}
        if strings.Contains(path, "/export-responses/") {
            h.ProgressHandler().ServeHTTP(w, r)
            return
        }

        http.NotFound(w, r)
    })
}
```

---

## 7. Decision Log

| Date | Decision | Rationale |
|------|----------|-----------|
| Jan 9, 2026 | Use exact Microsoft Graph beta endpoint format | Ensures compatibility with clients expecting real API |
| Jan 9, 2026 | Implement CSV redirect for Copilot API | Matches real Microsoft Graph behavior |
| Jan 9, 2026 | Use on-demand progress advancement for Qualtrics | Simpler testing, predictable behavior |
| Jan 9, 2026 | Store ZIP files in memory | Acceptable for simulator; no persistence needed |
| Jan 9, 2026 | Extend seed schema rather than separate configs | Maintains correlation between data sources |
| Jan 9, 2026 | Default app adoption rates based on typical enterprise | Realistic simulation without complex configuration |

---

**Next Step**: Review task.md for implementation breakdown with subagent assignments.
