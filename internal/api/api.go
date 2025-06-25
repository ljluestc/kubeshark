// GetConnections returns all tracked connections
func (a *API) GetConnections(c *gin.Context) {
	// Check if half-connections should be included
	includeHalf := false
	includeHalfParam := c.Query("includeHalf")
	if includeHalfParam == "true" {
		includeHalf = true
	}

	connections, err := a.connectionService.GetConnections(includeHalf)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, connections)
}
