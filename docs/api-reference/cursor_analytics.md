# Analytics API

The Analytics API provides comprehensive insights into your team's Cursor usage, including AI-assisted coding metrics, active users, model usage, and more.

- The Analytics API uses [Basic Authentication](/docs/api#basic-authentication). You can generate an API key from your [team settings page](https://cursor.com/settings).
- For details on authentication, rate limits, and best practices, see the [API Overview](/docs/api).
- **Availability**: Only for enterprise teams

### Available Endpoints

### Agent Edits

GET`/analytics/team/agent-edits`Get metrics on AI-suggested code edits accepted by your team in the Cursor IDE.

#### Parameters

`startDate` string

Start date for analytics period (default: 7 days ago). See [Date Formats](#date-formats)`endDate` string

End date for analytics period (default: today). See [Date Formats](#date-formats)`users` string

Filter data to specific users (comma-separated emails or user IDs, e.g., `alice@example.com,user_abc123`)```
curl -X GET "https://api.cursor.com/analytics/team/agent-edits" \
  -u YOUR_API_KEY:
```

**Response:**

```
{
  "data": [
    {
      "event_date": "2025-01-15",
      "total_suggested_diffs": 145,
      "total_accepted_diffs": 98,
      "total_rejected_diffs": 47,
      "total_green_lines_accepted": 820,
      "total_red_lines_accepted": 160,
      "total_green_lines_rejected": 210,
      "total_red_lines_rejected": 60,
      "total_green_lines_suggested": 1030,
      "total_red_lines_suggested": 220,
      "total_lines_suggested": 1250,
      "total_lines_accepted": 980
    },
    {
      "event_date": "2025-01-16",
      "total_suggested_diffs": 132,
      "total_accepted_diffs": 89,
      "total_rejected_diffs": 43,
      "total_green_lines_accepted": 740,
      "total_red_lines_accepted": 150,
      "total_green_lines_rejected": 185,
      "total_red_lines_rejected": 55,
      "total_green_lines_suggested": 925,
      "total_red_lines_suggested": 175,
      "total_lines_suggested": 1100,
      "total_lines_accepted": 890
    }
  ],
  "params": {
    "metric": "agent-edits",
    "teamId": 12345,
    "startDate": "2025-01-01",
    "endDate": "2025-01-31"
  }
}
```

### Tab Usage

GET`/analytics/team/tabs`Get metrics on Tab autocomplete usage across your team.

#### Parameters

`startDate` string

Start date for analytics period (default: 7 days ago). See [Date Formats](#date-formats)`endDate` string

End date for analytics period (default: today). See [Date Formats](#date-formats)`users` string

Filter data to specific users (comma-separated emails or user IDs, e.g., `alice@example.com,user_abc123`)```
curl -X GET "https://api.cursor.com/analytics/team/tabs" \
  -u YOUR_API_KEY:
```

**Response:**

```
{
  "data": [
    {
      "event_date": "2025-01-15",
      "total_suggestions": 5420,
      "total_accepts": 3210,
      "total_rejects": 2210,
      "total_green_lines_accepted": 4120,
      "total_red_lines_accepted": 2000,
      "total_green_lines_rejected": 1480,
      "total_red_lines_rejected": 730,
      "total_green_lines_suggested": 5600,
      "total_red_lines_suggested": 2740,
      "total_lines_suggested": 8340,
      "total_lines_accepted": 6120
    },
    {
      "event_date": "2025-01-16",
      "total_suggestions": 4980,
      "total_accepts": 3050,
      "total_rejects": 1930,
      "total_green_lines_accepted": 3890,
      "total_red_lines_accepted": 1890,
      "total_green_lines_rejected": 1350,
      "total_red_lines_rejected": 580,
      "total_green_lines_suggested": 5240,
      "total_red_lines_suggested": 2650,
      "total_lines_suggested": 7890,
      "total_lines_accepted": 5780
    }
  ],
  "params": {
    "metric": "tabs",
    "teamId": 12345,
    "startDate": "2025-01-01",
    "endDate": "2025-01-31"
  }
}
```

### Daily Active Users (DAU)

GET`/analytics/team/dau`Get daily active user counts for your team. DAU is the number of unique users who have used Cursor in a given day.
An active user is a user who has used at least one AI feature in the Cursor IDE.

Response includes DAU breakdown metrics for the Cursor CLI, Cloud Agents, and BugBot.

#### Parameters

`startDate` string

Start date for analytics period (default: 7 days ago). See [Date Formats](#date-formats)`endDate` string

End date for analytics period (default: today). See [Date Formats](#date-formats)`users` string

Filter data to specific users (comma-separated emails or user IDs, e.g., `alice@example.com,user_abc123`)```
curl -X GET "https://api.cursor.com/analytics/team/dau?startDate=14d&endDate=today" \
  -u YOUR_API_KEY:
```

**Response:**

```
{
  "data": [
    {
      "date": "2025-01-15",
      "dau": 42,
      "cli_dau": 5,
      "cloud_agent_dau": 37,
      "bugbot_dau": 10
    },
    {
      "date": "2025-01-16",
      "dau": 38,
      "cli_dau": 4,
      "cloud_agent_dau": 34,
      "bugbot_dau": 12
    }
  ],
  "params": {
    "metric": "dau",
    "teamId": 12345,
    "startDate": "2025-01-01",
    "endDate": "2025-01-31"
  }
}
```

### Client Versions

GET`/analytics/team/client-versions`Get distribution of Cursor IDE client versions used by your team (defaults to last 7 days). We report the latest version for each user per day (if a user has installed multiple versions, we report the latest).

#### Parameters

`startDate` string

Start date for analytics period (default: 7 days ago). See [Date Formats](#date-formats)`endDate` string

End date for analytics period (default: today). See [Date Formats](#date-formats)`users` string

Filter data to specific users (comma-separated emails or user IDs, e.g., `alice@example.com,user_abc123`)```
curl -X GET "https://api.cursor.com/analytics/team/client-versions" \
  -u YOUR_API_KEY:
```

**Response:**

```
{
  "data": [
    {
      "event_date": "2025-01-01",
      "client_version": "0.42.3",
      "user_count": 35,
      "percentage": 0.833
    },
    {
      "event_date": "2025-01-01",
      "client_version": "0.42.2",
      "user_count": 7,
      "percentage": 0.167
    }
  ],
  "params": {
    "metric": "client-versions",
    "teamId": 12345,
    "startDate": "2025-01-01",
    "endDate": "2025-01-31"
  }
}
```

### Model Usage

GET`/analytics/team/models`Get metrics on AI model usage across your team.

#### Parameters

`startDate` string

Start date for analytics period (default: 7 days ago). See [Date Formats](#date-formats)`endDate` string

End date for analytics period (default: today). See [Date Formats](#date-formats)`users` string

Filter data to specific users (comma-separated emails or user IDs, e.g., `alice@example.com,user_abc123`)```
curl -X GET "https://api.cursor.com/analytics/team/models" \
  -u YOUR_API_KEY:
```

**Response:**

```
{
  "data": [
    {
      "date": "2025-01-15",
      "model_breakdown": {
        "claude-sonnet-4.5": {
          "messages": 1250,
          "users": 28
        },
        "gpt-4o": {
          "messages": 450,
          "users": 15
        },
        "claude-opus-4": {
          "messages": 320,
          "users": 12
        }
      }
    },
    {
      "date": "2025-01-16",
      "model_breakdown": {
        "claude-sonnet-4.5": {
          "messages": 1180,
          "users": 26
        },
        "gpt-4o": {
          "messages": 420,
          "users": 14
        }
      }
    }
  ],
  "params": {
    "metric": "models",
    "teamId": 12345,
    "startDate": "2025-01-01",
    "endDate": "2025-01-31"
  }
}
```

### Top File Extensions

GET`/analytics/team/top-file-extensions`Get the most frequently edited files across your team in the Cursor IDE. Returns the top 5 file extensions per day by suggestion volume.

#### Parameters

`startDate` string

Start date for analytics period (default: 7 days ago). See [Date Formats](#date-formats)`endDate` string

End date for analytics period (default: today). See [Date Formats](#date-formats)`users` string

Filter data to specific users (comma-separated emails or user IDs, e.g., `alice@example.com,user_abc123`)```
curl -X GET "https://api.cursor.com/analytics/team/top-file-extensions?startDate=30d&endDate=today" \
  -u YOUR_API_KEY:
```

**Response:**

```
{
  "data": [
    {
      "event_date": "2025-01-15",
      "file_extension": "tsx",
      "total_files": 156,
      "total_accepts": 98,
      "total_rejects": 45,
      "total_lines_suggested": 3230,
      "total_lines_accepted": 2340,
      "total_lines_rejected": 890
    },
    {
      "event_date": "2025-01-15",
      "file_extension": "ts",
      "total_files": 142,
      "total_accepts": 89,
      "total_rejects": 38,
      "total_lines_suggested": 2850,
      "total_lines_accepted": 2100,
      "total_lines_rejected": 750
    }
  ],
  "params": {
    "metric": "top-file-extensions",
    "teamId": 12345,
    "startDate": "2025-01-01",
    "endDate": "2025-01-31"
  }
}
```

### MCP Adoption

GET`/analytics/team/mcp`Get metrics on MCP (Model Context Protocol) tool adoption across your team. Returns daily adoption counts broken down by tool name and MCP server name.

#### Parameters

`startDate` string

Start date for analytics period (default: 7 days ago). See [Date Formats](#date-formats)`endDate` string

End date for analytics period (default: today). See [Date Formats](#date-formats)`users` string

Filter data to specific users (comma-separated emails or user IDs, e.g., `alice@example.com,user_abc123`)```
curl -X GET "https://api.cursor.com/analytics/team/mcp" \
  -u YOUR_API_KEY:
```

**Response:**

```
{
  "data": [
    {
      "event_date": "2025-01-15",
      "tool_name": "read_file",
      "mcp_server_name": "filesystem",
      "usage": 245
    },
    {
      "event_date": "2025-01-15",
      "tool_name": "search_web",
      "mcp_server_name": "brave-search",
      "usage": 128
    },
    {
      "event_date": "2025-01-16",
      "tool_name": "read_file",
      "mcp_server_name": "filesystem",
      "usage": 231
    }
  ],
  "params": {
    "metric": "mcp",
    "teamId": 12345,
    "startDate": "2025-01-01",
    "endDate": "2025-01-31"
  }
}
```

### Commands Adoption

GET`/analytics/team/commands`Get metrics on Cursor command adoption across your team. Returns daily adoption counts broken down by command name.

#### Parameters

`startDate` string

Start date for analytics period (default: 7 days ago). See [Date Formats](#date-formats)`endDate` string

End date for analytics period (default: today). See [Date Formats](#date-formats)`users` string

Filter data to specific users (comma-separated emails or user IDs, e.g., `alice@example.com,user_abc123`)```
curl -X GET "https://api.cursor.com/analytics/team/commands" \
  -u YOUR_API_KEY:
```

**Response:**

```
{
  "data": [
    {
      "event_date": "2025-01-15",
      "command_name": "explain",
      "usage": 89
    },
    {
      "event_date": "2025-01-15",
      "command_name": "refactor",
      "usage": 45
    },
    {
      "event_date": "2025-01-16",
      "command_name": "explain",
      "usage": 92
    }
  ],
  "params": {
    "metric": "commands",
    "teamId": 12345,
    "startDate": "2025-01-01",
    "endDate": "2025-01-31"
  }
}
```

### Plans Adoption

GET`/analytics/team/plans`Get metrics on Plan mode adoption across your team. Returns daily adoption counts broken down by AI model used for plan generation.

#### Parameters

`startDate` string

Start date for analytics period (default: 7 days ago). See [Date Formats](#date-formats)`endDate` string

End date for analytics period (default: today). See [Date Formats](#date-formats)`users` string

Filter data to specific users (comma-separated emails or user IDs, e.g., `alice@example.com,user_abc123`)```
curl -X GET "https://api.cursor.com/analytics/team/plans" \
  -u YOUR_API_KEY:
```

**Response:**

```
{
  "data": [
    {
      "event_date": "2025-01-15",
      "model": "claude-sonnet-4.5",
      "usage": 156
    },
    {
      "event_date": "2025-01-15",
      "model": "gpt-4o",
      "usage": 42
    },
    {
      "event_date": "2025-01-16",
      "model": "claude-sonnet-4.5",
      "usage": 148
    }
  ],
  "params": {
    "metric": "plans",
    "teamId": 12345,
    "startDate": "2025-01-01",
    "endDate": "2025-01-31"
  }
}
```

### Ask Mode Adoption

GET`/analytics/team/ask-mode`Get metrics on Ask mode adoption across your team. Returns daily adoption counts broken down by AI model used for Ask mode queries.

#### Parameters

`startDate` string

Start date for analytics period (default: 7 days ago). See [Date Formats](#date-formats)`endDate` string

End date for analytics period (default: today). See [Date Formats](#date-formats)`users` string

Filter data to specific users (comma-separated emails or user IDs, e.g., `alice@example.com,user_abc123`)```
curl -X GET "https://api.cursor.com/analytics/team/ask-mode" \
  -u YOUR_API_KEY:
```

**Response:**

```
{
  "data": [
    {
      "event_date": "2025-01-15",
      "model": "claude-sonnet-4.5",
      "usage": 203
    },
    {
      "event_date": "2025-01-15",
      "model": "gpt-4o",
      "usage": 67
    },
    {
      "event_date": "2025-01-16",
      "model": "claude-sonnet-4.5",
      "usage": 198
    }
  ],
  "params": {
    "metric": "ask-mode",
    "teamId": 12345,
    "startDate": "2025-01-01",
    "endDate": "2025-01-31"
  }
}
```

### Leaderboard

GET`/analytics/team/leaderboard`Get a leaderboard of team members ranked by AI usage metrics.

**Behavior:**

- **Without user filtering**: Returns users ranked by the specified metric (default: combined lines accepted)
- **With user filtering**: Returns users that match the filter (with their actual team-wide rankings)
- Supports pagination for teams with many members

#### Parameters

`startDate` string

Start date for analytics period (default: 7 days ago). See [Date Formats](#date-formats)`endDate` string

End date for analytics period (default: today). See [Date Formats](#date-formats)`page` number

Page number for pagination (1-indexed). Default: `1``pageSize` number

Number of users per page (default: 10, max: 500)`users` string

Filter to specific users (comma-separated emails or user IDs, e.g., `alice@example.com,user_abc123`)Returns separate leaderboards for Tab autocomplete and Agent edits. When filtering by users, those users appear with their **actual team-wide rank**, not a filtered rank. For example, if you request a user who ranks #45 overall, they'll appear with `rank: 45`.

```
# Get first page of leaderboard (top 10 users)
curl -X GET "https://api.cursor.com/analytics/team/leaderboard" \
  -u YOUR_API_KEY:
```

```
# Get second page with custom page size
curl -X GET "https://api.cursor.com/analytics/team/leaderboard?page=2&pageSize=20" \
  -u YOUR_API_KEY:
```

```
# Filter by specific users
curl -X GET "https://api.cursor.com/analytics/team/leaderboard?users=alice@example.com,bob@example.com" \
  -u YOUR_API_KEY:
```

**Response:**

```
{
  "data": {
    "tab_leaderboard": {
      "data": [
        {
          "email": "alice@example.com",
          "user_id": "user_abc123",
          "total_accepts": 1334,
          "total_lines_accepted": 3455,
          "total_lines_suggested": 15307,
          "line_acceptance_ratio": 0.226,
          "accept_ratio": 0.233,
          "rank": 1
        },
        {
          "email": "bob@example.com",
          "user_id": "user_def789",
          "total_accepts": 796,
          "total_lines_accepted": 2090,
          "total_lines_suggested": 7689,
          "line_acceptance_ratio": 0.272,
          "accept_ratio": 0.273,
          "rank": 2
        }
      ],
      "total_users": 142
    },
    "agent_leaderboard": {
      "data": [
        {
          "email": "alice@example.com",
          "user_id": "user_abc123",
          "total_accepts": 914,
          "total_lines_accepted": 65947,
          "total_lines_suggested": 201467,
          "line_acceptance_ratio": 0.327,
          "favorite_model": "claude-sonnet-4.5",
          "rank": 1
        },
        {
          "email": "bob@example.com",
          "user_id": "user_def789",
          "total_accepts": 843,
          "total_lines_accepted": 61709,
          "total_lines_suggested": 51092,
          "line_acceptance_ratio": 1.208,
          "favorite_model": "claude-sonnet-4.5",
          "rank": 2
        }
      ],
      "total_users": 142
    }
  },
  "pagination": {
    "page": 1,
    "pageSize": 10,
    "totalUsers": 142,
    "totalPages": 15,
    "hasNextPage": true,
    "hasPreviousPage": false
  },
  "params": {
    "metric": "leaderboard",
    "teamId": 12345,
    "startDate": "2025-01-01",
    "endDate": "2025-01-31",
    "page": 1,
    "pageSize": 10
  }
}
```

---

## By-User Endpoints

By-user endpoints provide the same metrics as team-level endpoints, but organized by individual users with pagination support. These are ideal for generating per-user reports or processing large teams in batches.

### Common Query Parameters

ParameterTypeRequiredDescription`startDate`Date stringNoStart date for the analytics period (default: 7 days ago)`endDate`Date stringNoEnd date for the analytics period (default: today)`page`numberNoPage number (default: 1)`pageSize`numberNoNumber of users per page (default: 100, max: 500)`users`stringNoLimit pagination to specific users (comma-separated emails or IDs, e.g., `alice@example.com,user_abc123`)
**User Filtering:**
When you provide the `users` parameter to by-user endpoints:

- **Pagination is filtered**: Only the specified users are included in the result set and pagination counts
- **Useful for**: Getting detailed data for specific team members without paginating through all users
- Example: If you have 500 users but only want data for 3 specific users, filter by their emails to get all 3 in a single page

**Note:** By-user endpoints support the same date formats and shortcuts as team-level endpoints. See the [Date Formats](#date-formats) section above.

### Response Format

All by-user endpoints return data in this format:

```
{
  "data": {
    "user1@example.com": [ /* user's data */ ],
    "user2@example.com": [ /* user's data */ ]
  },
  "pagination": {
    "page": 1,
    "pageSize": 100,
    "totalUsers": 250,
    "totalPages": 3,
    "hasNextPage": true,
    "hasPreviousPage": false
  },
  "params": {
    "metric": "agent-edits",
    "teamId": 12345,
    "startDate": "2025-01-01",
    "endDate": "2025-01-31",
    "page": 1,
    "pageSize": 100,
    "userMappings": [
      { "id": "user_abc123", "email": "user1@example.com" },
      { "id": "user_def456", "email": "user2@example.com" }
    ]
  }
}
```

**Response Structure:**

- `data` - Object keyed by user email addresses, each containing an array of that user's metrics
- `pagination` - Pagination information
- `params` - Request parameters echoed back
- `userMappings` - Array mapping email addresses to public user IDs for this page. Useful for cross-referencing with other APIs or creating links to user profiles.

### Available Endpoints

All by-user endpoints follow the pattern: `/analytics/by-user/{metric}`

- `GET /analytics/by-user/agent-edits` - Agent edits by user
- `GET /analytics/by-user/tabs` - Tab usage by user
- `GET /analytics/by-user/models` - Model usage by user
- `GET /analytics/by-user/top-file-extensions` - Top files by user
- `GET /analytics/by-user/client-versions` - Client versions by user
- `GET /analytics/by-user/mcp` - MCP adoption by user
- `GET /analytics/by-user/commands` - Commands adoption by user
- `GET /analytics/by-user/plans` - Plans adoption by user
- `GET /analytics/by-user/ask-mode` - Ask mode adoption by user

### Agent Edits By User

GET`/analytics/by-user/agent-edits`Get agent edits metrics organized by individual users with pagination support.

#### Parameters

`startDate` string

Start date for analytics period (default: 7 days ago). See [Date Formats](#date-formats)`endDate` string

End date for analytics period (default: today). See [Date Formats](#date-formats)`page` number

Page number (1-indexed). Default: `1``pageSize` number

Number of users per page (default: 100, max: 500)`users` string

Limit pagination to specific users (comma-separated emails or user IDs, e.g., `alice@example.com,user_abc123`)```
curl -X GET "https://api.cursor.com/analytics/by-user/agent-edits?page=1&pageSize=50" \
  -u YOUR_API_KEY:
```

```
curl -X GET "https://api.cursor.com/analytics/by-user/agent-edits?users=alice@example.com,bob@example.com,carol@example.com" \
  -u YOUR_API_KEY:
```

**Response:**

```
{
  "data": {
    "alice@example.com": [
      {
        "event_date": "2025-01-15",
        "suggested_lines": 125,
        "accepted_lines": 98
      },
      {
        "event_date": "2025-01-16",
        "suggested_lines": 110,
        "accepted_lines": 89
      }
    ],
    "bob@example.com": [
      {
        "event_date": "2025-01-15",
        "suggested_lines": 95,
        "accepted_lines": 72
      },
      {
        "event_date": "2025-01-16",
        "suggested_lines": 88,
        "accepted_lines": 65
      }
    ]
  },
  "pagination": {
    "page": 1,
    "pageSize": 50,
    "totalUsers": 120,
    "totalPages": 3,
    "hasNextPage": true,
    "hasPreviousPage": false
  },
  "params": {
    "metric": "agent-edits",
    "teamId": 12345,
    "startDate": "2025-01-01",
    "endDate": "2025-01-31",
    "page": 1,
    "pageSize": 50,
    "userMappings": [
      { "id": "user_abc123", "email": "alice@example.com" },
      { "id": "user_def456", "email": "bob@example.com" }
    ]
  }
}
```

### MCP Adoption By User

GET`/analytics/by-user/mcp`Get MCP tool adoption metrics organized by individual users with pagination support.

#### Parameters

`startDate` string

Start date for analytics period (default: 7 days ago). See [Date Formats](#date-formats)`endDate` string

End date for analytics period (default: today). See [Date Formats](#date-formats)`page` number

Page number (1-indexed). Default: `1``pageSize` number

Number of users per page (default: 100, max: 500)`users` string

Limit pagination to specific users (comma-separated emails or user IDs, e.g., `alice@example.com,user_abc123`)```
curl -X GET "https://api.cursor.com/analytics/by-user/mcp?page=1&pageSize=50" \
  -u YOUR_API_KEY:
```

**Response:**

```
{
  "data": {
    "alice@example.com": [
      {
        "event_date": "2025-01-15",
        "tool_name": "read_file",
        "mcp_server_name": "filesystem",
        "usage": 45
      },
      {
        "event_date": "2025-01-16",
        "tool_name": "read_file",
        "mcp_server_name": "filesystem",
        "usage": 38
      }
    ],
    "bob@example.com": [
      {
        "event_date": "2025-01-15",
        "tool_name": "search_web",
        "mcp_server_name": "brave-search",
        "usage": 23
      }
    ]
  },
  "pagination": {
    "page": 1,
    "pageSize": 50,
    "totalUsers": 120,
    "totalPages": 3,
    "hasNextPage": true,
    "hasPreviousPage": false
  },
  "params": {
    "metric": "mcp",
    "teamId": 12345,
    "startDate": "2025-01-01",
    "endDate": "2025-01-31",
    "page": 1,
    "pageSize": 50,
    "userMappings": [
      { "id": "user_abc123", "email": "alice@example.com" },
      { "id": "user_def456", "email": "bob@example.com" }
    ]
  }
}
```

### Commands Adoption By User

GET`/analytics/by-user/commands`Get command adoption metrics organized by individual users with pagination support.

#### Parameters

`startDate` string

Start date for analytics period (default: 7 days ago). See [Date Formats](#date-formats)`endDate` string

End date for analytics period (default: today). See [Date Formats](#date-formats)`page` number

Page number (1-indexed). Default: `1``pageSize` number

Number of users per page (default: 100, max: 500)`users` string

Limit pagination to specific users (comma-separated emails or user IDs, e.g., `alice@example.com,user_abc123`)```
curl -X GET "https://api.cursor.com/analytics/by-user/commands?page=1&pageSize=50" \
  -u YOUR_API_KEY:
```

**Response:**

```
{
  "data": {
    "alice@example.com": [
      {
        "event_date": "2025-01-15",
        "command_name": "explain",
        "usage": 12
      },
      {
        "event_date": "2025-01-16",
        "command_name": "explain",
        "usage": 15
      }
    ],
    "bob@example.com": [
      {
        "event_date": "2025-01-15",
        "command_name": "refactor",
        "usage": 8
      }
    ]
  },
  "pagination": {
    "page": 1,
    "pageSize": 50,
    "totalUsers": 120,
    "totalPages": 3,
    "hasNextPage": true,
    "hasPreviousPage": false
  },
  "params": {
    "metric": "commands",
    "teamId": 12345,
    "startDate": "2025-01-01",
    "endDate": "2025-01-31",
    "page": 1,
    "pageSize": 50,
    "userMappings": [
      { "id": "user_abc123", "email": "alice@example.com" },
      { "id": "user_def456", "email": "bob@example.com" }
    ]
  }
}
```

### Plans Adoption By User

GET`/analytics/by-user/plans`Get Plan mode adoption metrics organized by individual users with pagination support.

#### Parameters

`startDate` string

Start date for analytics period (default: 7 days ago). See [Date Formats](#date-formats)`endDate` string

End date for analytics period (default: today). See [Date Formats](#date-formats)`page` number

Page number (1-indexed). Default: `1``pageSize` number

Number of users per page (default: 100, max: 500)`users` string

Limit pagination to specific users (comma-separated emails or user IDs, e.g., `alice@example.com,user_abc123`)```
curl -X GET "https://api.cursor.com/analytics/by-user/plans?page=1&pageSize=50" \
  -u YOUR_API_KEY:
```

**Response:**

```
{
  "data": {
    "alice@example.com": [
      {
        "event_date": "2025-01-15",
        "model": "claude-sonnet-4.5",
        "usage": 23
      },
      {
        "event_date": "2025-01-16",
        "model": "claude-sonnet-4.5",
        "usage": 19
      }
    ],
    "bob@example.com": [
      {
        "event_date": "2025-01-15",
        "model": "gpt-4o",
        "usage": 12
      }
    ]
  },
  "pagination": {
    "page": 1,
    "pageSize": 50,
    "totalUsers": 120,
    "totalPages": 3,
    "hasNextPage": true,
    "hasPreviousPage": false
  },
  "params": {
    "metric": "plans",
    "teamId": 12345,
    "startDate": "2025-01-01",
    "endDate": "2025-01-31",
    "page": 1,
    "pageSize": 50,
    "userMappings": [
      { "id": "user_abc123", "email": "alice@example.com" },
      { "id": "user_def456", "email": "bob@example.com" }
    ]
  }
}
```

### Ask Mode Adoption By User

GET`/analytics/by-user/ask-mode`Get Ask mode adoption metrics organized by individual users with pagination support.

#### Parameters

`startDate` string

Start date for analytics period (default: 7 days ago). See [Date Formats](#date-formats)`endDate` string

End date for analytics period (default: today). See [Date Formats](#date-formats)`page` number

Page number (1-indexed). Default: `1``pageSize` number

Number of users per page (default: 100, max: 500)`users` string

Limit pagination to specific users (comma-separated emails or user IDs, e.g., `alice@example.com,user_abc123`)```
curl -X GET "https://api.cursor.com/analytics/by-user/ask-mode?page=1&pageSize=50" \
  -u YOUR_API_KEY:
```

**Response:**

```
{
  "data": {
    "alice@example.com": [
      {
        "event_date": "2025-01-15",
        "model": "claude-sonnet-4.5",
        "usage": 34
      },
      {
        "event_date": "2025-01-16",
        "model": "claude-sonnet-4.5",
        "usage": 28
      }
    ],
    "bob@example.com": [
      {
        "event_date": "2025-01-15",
        "model": "gpt-4o",
        "usage": 15
      }
    ]
  },
  "pagination": {
    "page": 1,
    "pageSize": 50,
    "totalUsers": 120,
    "totalPages": 3,
    "hasNextPage": true,
    "hasPreviousPage": false
  },
  "params": {
    "metric": "ask-mode",
    "teamId": 12345,
    "startDate": "2025-01-01",
    "endDate": "2025-01-31",
    "page": 1,
    "pageSize": 50,
    "userMappings": [
      { "id": "user_abc123", "email": "alice@example.com" },
      { "id": "user_def456", "email": "bob@example.com" }
    ]
  }
}
```

---

## Team-Level Endpoints

Team-level endpoints provide aggregated metrics for your entire team or filtered subsets of users. All endpoints support date range filtering and optional user filtering.

### Common Query Parameters

ParameterTypeRequiredDescription`startDate`Date stringNoStart date for the analytics period (default: 7 days ago)`endDate`Date stringNoEnd date for the analytics period (default: today)`users`stringNoFilter data to specific users (comma-separated). Each value can be an email (e.g., `alice@example.com`) or public user ID (e.g., `user_abc123`). You can mix both formats.
**User Filtering:**
The `users` parameter accepts a comma-separated list of identifiers. Each identifier can be:

- **Email address** (e.g., `alice@example.com`) - Auto-detected by the presence of `@`
- **Public user ID** (e.g., `user_abc123`) - Auto-detected by the `user_` prefix
- **Mixed format** - You can combine emails and IDs in the same request

**Examples:**

```
# Filter by emails only
?users=alice@example.com,bob@example.com,carol@example.com

# Filter by public user IDs only
?users=user_abc123,user_def456,user_ghi789

# Mix emails and IDs
?users=alice@example.com,user_def456,bob@example.com
```

When you filter by users, the API returns data **only for those specific users**. This is useful for:

- Analyzing specific team members or groups (e.g., engineering leads, specific project teams)
- Generating reports for a subset of users
- Comparing metrics across selected individuals

### Date Formats

**Default Behavior:**
If you omit both `startDate` and `endDate`, the API defaults to the **last 7 days** (from 7 days ago until today). This is perfect for quick queries without specifying dates.

**Standard Formats:**

- `YYYY-MM-DD` - Simple date format (e.g., `2025-01-15`) **← Recommended**
- ISO 8601 timestamps (e.g., `2025-01-15T00:00:00Z`)

**Shortcuts:**

- `now` or `today` - Current date (at 00:00:00)
- `yesterday` - Yesterday's date (at 00:00:00)
- `<number>d` - Days ago (e.g., `7d` = 7 days ago, `30d` = 30 days ago)

**Important Notes:**

- **Time is ignored**: All dates are resolved to the day level (00:00:00 UTC). Sending `2025-01-15T14:30:00Z` is the same as `2025-01-15`.
- **Use recommended formats**: Use `YYYY-MM-DD` or shortcuts for better HTTP caching support. Different time values (like `T14:30:00Z` vs `T08:00:00Z`) prevent cache hits even though they resolve to the same day.
- **Date ranges**: Limited to a maximum of 30 days.

**Examples:**

```
# Omit dates for last 7 days (simplest and best for caching)
curl "https://api.cursor.com/analytics/team/agent-edits"

# Using YYYY-MM-DD format for specific date range (recommended)
?startDate=2025-01-01&endDate=2025-01-31

# Using shortcuts for last 30 days
?startDate=30d&endDate=today

# Using shortcuts for last 14 days
?startDate=14d&endDate=now

# ❌ Don't use timestamps - prevents caching and time is ignored anyway
?startDate=2025-01-15T14:30:00Z&endDate=2025-01-31T23:59:59Z
```

## Rate Limits

Rate limits are enforced per team and reset every minute:

- **Team-level endpoints**: 100 requests per minute per team
- **By-user endpoints**: 50 requests per minute per team

**What happens when you exceed the rate limit?**

When you exceed the rate limit, you'll receive a `429 Too Many Requests` response:

```
{
  "error": "Too Many Requests",
  "message": "Rate limit exceeded. Please try again later."
}
```

## Best Practices

For general API best practices including exponential backoff, caching strategies, and error handling, see the [API Overview Best Practices](/docs/api#best-practices).

1. **Use pagination for large teams**: If your team has more than 100 users, use the by-user endpoints with pagination to avoid timeouts.
2. **Leverage caching**: Both Team and User level endpoints support ETags. Store the ETag and use `If-None-Match` headers to reduce unnecessary data transfer.
3. **Filter by users when possible**: If you only need data for specific users, use the `users` parameter to reduce query time.
4. **Date ranges**: Keep date ranges reasonable (e.g., 1-3 months) for optimal performance.