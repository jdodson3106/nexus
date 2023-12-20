package nexus

import "strings"

func ToSnakeLower(appName string) string {
    lowerName := strings.ToLower(appName)
    return strings.Join(strings.Split(lowerName, " "), "_")
}


func ToSnakeUpper(appName string) string {
    upperName := strings.ToUpper(appName)
    return strings.Join(strings.Split(upperName, " "), "_")
}
