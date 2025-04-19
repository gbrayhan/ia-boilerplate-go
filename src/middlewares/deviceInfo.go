package middlewares

import "github.com/gin-gonic/gin"

import "github.com/mssola/user_agent"

type DeviceInfo struct {
	IPAddress      string `json:"ip_address"`
	UserAgent      string `json:"user_agent"`
	DeviceType     string `json:"device_type"`
	Browser        string `json:"browser"`
	BrowserVersion string `json:"browser_version"`
	OS             string `json:"os"`
	Language       string `json:"language"`
}

func DeviceInfoInterceptor() gin.HandlerFunc {
	return func(c *gin.Context) {
		userAgentStr := c.GetHeader("User-Agent")
		language := c.GetHeader("Accept-Language")
		ipAddress := c.ClientIP()

		ua := user_agent.New(userAgentStr)
		browserName, browserVersion := ua.Browser()
		deviceType := "desktop"
		if ua.Mobile() {
			deviceType = "mobile"
		}
		osInfo := ua.OS()

		deviceInfo := DeviceInfo{
			IPAddress:      ipAddress,
			UserAgent:      userAgentStr,
			DeviceType:     deviceType,
			Browser:        browserName,
			BrowserVersion: browserVersion,
			OS:             osInfo,
			Language:       language,
		}

		c.Set("deviceInfo", deviceInfo)

		c.Next()
	}
}
