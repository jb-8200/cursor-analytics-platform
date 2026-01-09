package models

// CopilotReportPeriod represents the report period for Copilot usage data.
type CopilotReportPeriod string

const (
	CopilotPeriodD7   CopilotReportPeriod = "D7"
	CopilotPeriodD30  CopilotReportPeriod = "D30"
	CopilotPeriodD90  CopilotReportPeriod = "D90"
	CopilotPeriodD180 CopilotReportPeriod = "D180"
	CopilotPeriodAll  CopilotReportPeriod = "All"
)

// Days returns the number of days for this report period.
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
		return 180
	default:
		return 180 // Default to 180 days for unknown periods
	}
}

// CopilotApp represents supported Microsoft 365 Copilot applications.
type CopilotApp string

const (
	CopilotAppTeams      CopilotApp = "Teams"
	CopilotAppWord       CopilotApp = "Word"
	CopilotAppExcel      CopilotApp = "Excel"
	CopilotAppPowerPoint CopilotApp = "PowerPoint"
	CopilotAppOutlook    CopilotApp = "Outlook"
	CopilotAppOneNote    CopilotApp = "OneNote"
	CopilotAppLoop       CopilotApp = "Loop"
	CopilotAppCopilot    CopilotApp = "Copilot"
)

// AllCopilotApps returns all supported Copilot application constants.
func AllCopilotApps() []CopilotApp {
	return []CopilotApp{
		CopilotAppTeams,
		CopilotAppWord,
		CopilotAppExcel,
		CopilotAppPowerPoint,
		CopilotAppOutlook,
		CopilotAppOneNote,
		CopilotAppLoop,
		CopilotAppCopilot,
	}
}

// CopilotUsageUserDetail represents Microsoft 365 Copilot usage data for a single user.
// Field names match the Microsoft Graph API schema exactly.
// See: https://learn.microsoft.com/en-us/graph/api/reportroot-getmicrosoft365copilotusageuserdetail
type CopilotUsageUserDetail struct {
	ReportRefreshDate                     string  `json:"reportRefreshDate"`
	ReportPeriod                          int     `json:"reportPeriod"`
	UserPrincipalName                     string  `json:"userPrincipalName"`
	DisplayName                           string  `json:"displayName"`
	LastActivityDate                      *string `json:"lastActivityDate,omitempty"`
	MicrosoftTeamsCopilotLastActivityDate *string `json:"microsoftTeamsCopilotLastActivityDate,omitempty"`
	WordCopilotLastActivityDate           *string `json:"wordCopilotLastActivityDate,omitempty"`
	ExcelCopilotLastActivityDate          *string `json:"excelCopilotLastActivityDate,omitempty"`
	PowerPointCopilotLastActivityDate     *string `json:"powerPointCopilotLastActivityDate,omitempty"`
	OutlookCopilotLastActivityDate        *string `json:"outlookCopilotLastActivityDate,omitempty"`
	OneNoteCopilotLastActivityDate        *string `json:"oneNoteCopilotLastActivityDate,omitempty"`
	LoopCopilotLastActivityDate           *string `json:"loopCopilotLastActivityDate,omitempty"`
	CopilotChatLastActivityDate           *string `json:"copilotChatLastActivityDate,omitempty"`
}

// GetAppLastActivityDate returns the last activity date for the specified app.
func (d *CopilotUsageUserDetail) GetAppLastActivityDate(app CopilotApp) *string {
	switch app {
	case CopilotAppTeams:
		return d.MicrosoftTeamsCopilotLastActivityDate
	case CopilotAppWord:
		return d.WordCopilotLastActivityDate
	case CopilotAppExcel:
		return d.ExcelCopilotLastActivityDate
	case CopilotAppPowerPoint:
		return d.PowerPointCopilotLastActivityDate
	case CopilotAppOutlook:
		return d.OutlookCopilotLastActivityDate
	case CopilotAppOneNote:
		return d.OneNoteCopilotLastActivityDate
	case CopilotAppLoop:
		return d.LoopCopilotLastActivityDate
	case CopilotAppCopilot:
		return d.CopilotChatLastActivityDate
	default:
		return nil
	}
}

// CopilotUsageResponse is the OData response format from Microsoft Graph API.
type CopilotUsageResponse struct {
	Context  string                   `json:"@odata.context,omitempty"`
	NextLink string                   `json:"@odata.nextLink,omitempty"`
	Value    []CopilotUsageUserDetail `json:"value"`
}
