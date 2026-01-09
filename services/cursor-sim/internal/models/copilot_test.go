package models

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCopilotReportPeriod_Days(t *testing.T) {
	tests := []struct {
		period CopilotReportPeriod
		days   int
	}{
		{CopilotPeriodD7, 7},
		{CopilotPeriodD30, 30},
		{CopilotPeriodD90, 90},
		{CopilotPeriodD180, 180},
		{CopilotPeriodAll, 180}, // "All" defaults to 180 days
	}

	for _, tt := range tests {
		t.Run(string(tt.period), func(t *testing.T) {
			assert.Equal(t, tt.days, tt.period.Days())
		})
	}
}

func TestCopilotReportPeriod_InvalidPeriod(t *testing.T) {
	invalidPeriod := CopilotReportPeriod("D365")
	// Should default to 180 for unknown period
	assert.Equal(t, 180, invalidPeriod.Days())
}

func TestAllCopilotApps(t *testing.T) {
	apps := AllCopilotApps()

	// Should return all 8 apps
	assert.Len(t, apps, 8)

	// Verify all expected apps are present
	expectedApps := []CopilotApp{
		CopilotAppTeams,
		CopilotAppWord,
		CopilotAppExcel,
		CopilotAppPowerPoint,
		CopilotAppOutlook,
		CopilotAppOneNote,
		CopilotAppLoop,
		CopilotAppCopilot,
	}

	for _, expected := range expectedApps {
		assert.Contains(t, apps, expected)
	}
}

func TestCopilotUsageUserDetail_GetAppLastActivityDate(t *testing.T) {
	teamsDate := "2026-01-08"
	wordDate := "2026-01-05"
	excelDate := "2026-01-03"

	detail := CopilotUsageUserDetail{
		MicrosoftTeamsCopilotLastActivityDate: &teamsDate,
		WordCopilotLastActivityDate:           &wordDate,
		ExcelCopilotLastActivityDate:          &excelDate,
	}

	// Test existing activity dates
	assert.Equal(t, &teamsDate, detail.GetAppLastActivityDate(CopilotAppTeams))
	assert.Equal(t, &wordDate, detail.GetAppLastActivityDate(CopilotAppWord))
	assert.Equal(t, &excelDate, detail.GetAppLastActivityDate(CopilotAppExcel))

	// Test nil activity dates
	assert.Nil(t, detail.GetAppLastActivityDate(CopilotAppPowerPoint))
	assert.Nil(t, detail.GetAppLastActivityDate(CopilotAppOutlook))
	assert.Nil(t, detail.GetAppLastActivityDate(CopilotAppOneNote))
	assert.Nil(t, detail.GetAppLastActivityDate(CopilotAppLoop))
	assert.Nil(t, detail.GetAppLastActivityDate(CopilotAppCopilot))
}

func TestCopilotUsageUserDetail_GetAppLastActivityDate_AllApps(t *testing.T) {
	teamsDate := "2026-01-08"
	wordDate := "2026-01-07"
	excelDate := "2026-01-06"
	powerPointDate := "2026-01-05"
	outlookDate := "2026-01-04"
	oneNoteDate := "2026-01-03"
	loopDate := "2026-01-02"
	copilotDate := "2026-01-01"

	detail := CopilotUsageUserDetail{
		MicrosoftTeamsCopilotLastActivityDate: &teamsDate,
		WordCopilotLastActivityDate:           &wordDate,
		ExcelCopilotLastActivityDate:          &excelDate,
		PowerPointCopilotLastActivityDate:     &powerPointDate,
		OutlookCopilotLastActivityDate:        &outlookDate,
		OneNoteCopilotLastActivityDate:        &oneNoteDate,
		LoopCopilotLastActivityDate:           &loopDate,
		CopilotChatLastActivityDate:           &copilotDate,
	}

	// Verify all apps return correct dates
	assert.Equal(t, &teamsDate, detail.GetAppLastActivityDate(CopilotAppTeams))
	assert.Equal(t, &wordDate, detail.GetAppLastActivityDate(CopilotAppWord))
	assert.Equal(t, &excelDate, detail.GetAppLastActivityDate(CopilotAppExcel))
	assert.Equal(t, &powerPointDate, detail.GetAppLastActivityDate(CopilotAppPowerPoint))
	assert.Equal(t, &outlookDate, detail.GetAppLastActivityDate(CopilotAppOutlook))
	assert.Equal(t, &oneNoteDate, detail.GetAppLastActivityDate(CopilotAppOneNote))
	assert.Equal(t, &loopDate, detail.GetAppLastActivityDate(CopilotAppLoop))
	assert.Equal(t, &copilotDate, detail.GetAppLastActivityDate(CopilotAppCopilot))
}

func TestCopilotUsageUserDetail_GetAppLastActivityDate_InvalidApp(t *testing.T) {
	detail := CopilotUsageUserDetail{
		ReportRefreshDate: "2026-01-09",
	}

	// Test invalid/unknown app returns nil
	invalidApp := CopilotApp("InvalidApp")
	assert.Nil(t, detail.GetAppLastActivityDate(invalidApp))
}

func TestCopilotUsageUserDetail_JSONMarshal(t *testing.T) {
	date := "2026-01-08"
	detail := CopilotUsageUserDetail{
		ReportRefreshDate:                     "2026-01-09",
		ReportPeriod:                          30,
		UserPrincipalName:                     "user@company.com",
		DisplayName:                           "Jane Dev",
		LastActivityDate:                      &date,
		MicrosoftTeamsCopilotLastActivityDate: &date,
	}

	data, err := json.Marshal(detail)
	require.NoError(t, err)

	// Verify exact field names match Microsoft API
	assert.Contains(t, string(data), `"reportRefreshDate":"2026-01-09"`)
	assert.Contains(t, string(data), `"reportPeriod":30`)
	assert.Contains(t, string(data), `"userPrincipalName":"user@company.com"`)
	assert.Contains(t, string(data), `"displayName":"Jane Dev"`)
	assert.Contains(t, string(data), `"lastActivityDate":"2026-01-08"`)
	assert.Contains(t, string(data), `"microsoftTeamsCopilotLastActivityDate":"2026-01-08"`)
}

func TestCopilotUsageUserDetail_JSONMarshal_NullableDates(t *testing.T) {
	detail := CopilotUsageUserDetail{
		ReportRefreshDate: "2026-01-09",
		ReportPeriod:      30,
		UserPrincipalName: "user@company.com",
		DisplayName:       "Jane Dev",
		// All activity dates are nil
	}

	data, err := json.Marshal(detail)
	require.NoError(t, err)

	jsonStr := string(data)

	// Verify required fields are present
	assert.Contains(t, jsonStr, `"reportRefreshDate":"2026-01-09"`)
	assert.Contains(t, jsonStr, `"reportPeriod":30`)
	assert.Contains(t, jsonStr, `"userPrincipalName":"user@company.com"`)
	assert.Contains(t, jsonStr, `"displayName":"Jane Dev"`)

	// Verify omitempty works - null fields should not be present
	// or should be explicitly null if included
	assert.NotContains(t, jsonStr, `"lastActivityDate":""`)
	assert.NotContains(t, jsonStr, `"microsoftTeamsCopilotLastActivityDate":""`)
}

func TestCopilotUsageUserDetail_JSONUnmarshal(t *testing.T) {
	jsonData := `{
		"reportRefreshDate": "2026-01-09",
		"reportPeriod": 30,
		"userPrincipalName": "user@company.com",
		"displayName": "Jane Dev",
		"lastActivityDate": "2026-01-08",
		"microsoftTeamsCopilotLastActivityDate": "2026-01-07"
	}`

	var detail CopilotUsageUserDetail
	err := json.Unmarshal([]byte(jsonData), &detail)
	require.NoError(t, err)

	assert.Equal(t, "2026-01-09", detail.ReportRefreshDate)
	assert.Equal(t, 30, detail.ReportPeriod)
	assert.Equal(t, "user@company.com", detail.UserPrincipalName)
	assert.Equal(t, "Jane Dev", detail.DisplayName)
	require.NotNil(t, detail.LastActivityDate)
	assert.Equal(t, "2026-01-08", *detail.LastActivityDate)
	require.NotNil(t, detail.MicrosoftTeamsCopilotLastActivityDate)
	assert.Equal(t, "2026-01-07", *detail.MicrosoftTeamsCopilotLastActivityDate)
}

func TestCopilotUsageResponse_JSONMarshal(t *testing.T) {
	date := "2026-01-08"
	response := CopilotUsageResponse{
		Context:  "https://graph.microsoft.com/v1.0/$metadata#reports",
		NextLink: "https://graph.microsoft.com/v1.0/reports?$skiptoken=abc123",
		Value: []CopilotUsageUserDetail{
			{
				ReportRefreshDate:                     "2026-01-09",
				ReportPeriod:                          30,
				UserPrincipalName:                     "user@company.com",
				DisplayName:                           "Jane Dev",
				LastActivityDate:                      &date,
				MicrosoftTeamsCopilotLastActivityDate: &date,
			},
		},
	}

	data, err := json.Marshal(response)
	require.NoError(t, err)

	jsonStr := string(data)

	// Verify OData field names
	assert.Contains(t, jsonStr, `"@odata.context"`)
	assert.Contains(t, jsonStr, `"@odata.nextLink"`)
	assert.Contains(t, jsonStr, `"value"`)
	assert.Contains(t, jsonStr, `"userPrincipalName":"user@company.com"`)
}

func TestCopilotUsageResponse_JSONMarshal_NoNextLink(t *testing.T) {
	date := "2026-01-08"
	response := CopilotUsageResponse{
		Context: "https://graph.microsoft.com/v1.0/$metadata#reports",
		// NextLink is empty (no pagination)
		Value: []CopilotUsageUserDetail{
			{
				ReportRefreshDate: "2026-01-09",
				ReportPeriod:      30,
				UserPrincipalName: "user@company.com",
				DisplayName:       "Jane Dev",
				LastActivityDate:  &date,
			},
		},
	}

	data, err := json.Marshal(response)
	require.NoError(t, err)

	jsonStr := string(data)

	// Verify context is present
	assert.Contains(t, jsonStr, `"@odata.context"`)

	// NextLink should be omitted when empty
	assert.NotContains(t, jsonStr, `"@odata.nextLink":""`)
}

func TestCopilotApp_Constants(t *testing.T) {
	assert.Equal(t, CopilotApp("Teams"), CopilotAppTeams)
	assert.Equal(t, CopilotApp("Word"), CopilotAppWord)
	assert.Equal(t, CopilotApp("Excel"), CopilotAppExcel)
	assert.Equal(t, CopilotApp("PowerPoint"), CopilotAppPowerPoint)
	assert.Equal(t, CopilotApp("Outlook"), CopilotAppOutlook)
	assert.Equal(t, CopilotApp("OneNote"), CopilotAppOneNote)
	assert.Equal(t, CopilotApp("Loop"), CopilotAppLoop)
	assert.Equal(t, CopilotApp("Copilot"), CopilotAppCopilot)
}

func TestCopilotReportPeriod_Constants(t *testing.T) {
	assert.Equal(t, CopilotReportPeriod("D7"), CopilotPeriodD7)
	assert.Equal(t, CopilotReportPeriod("D30"), CopilotPeriodD30)
	assert.Equal(t, CopilotReportPeriod("D90"), CopilotPeriodD90)
	assert.Equal(t, CopilotReportPeriod("D180"), CopilotPeriodD180)
	assert.Equal(t, CopilotReportPeriod("All"), CopilotPeriodAll)
}
