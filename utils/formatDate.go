package forum

import (
        "fmt"
        "time"
)


func FormatDate(date *time.Time) string {
        if date == nil {
                return "no message yet"
        }

        duration := time.Since(*date) 
        day := 24 * time.Hour
        yesterday := 48 * time.Hour
        month := 30 * day
        year := 365 * day

        switch {
        case duration < time.Minute:
                return "now"
        case duration < time.Hour:
                minutes := int(duration.Minutes())
                if minutes == 1 {
                        return "a minute ago"
                }

                return fmt.Sprintf("%d minutes ago", minutes)
        case duration < day:
                hours := int(duration.Hours())
                if hours == 1 {
                        return "an hour ago"
                }

                return fmt.Sprintf("%d hours ago", hours)
        case duration < yesterday:
                return "yesterday"
        case duration < month:
                days := int(duration.Hours() / 24)
                return fmt.Sprintf("%d days ago", days)
        case duration < year:
                months := int(duration.Hours()) / (30 * 24)
                if months == 1 {
                        return "a month ago"
                }

                return fmt.Sprintf("%d months ago", months)
        default:
                years := int(duration.Hours()) / (365 * 24)
                if years == 1 {
                        return "a year ago"
                }

                return fmt.Sprintf("%d years ago", years)
        }
}
