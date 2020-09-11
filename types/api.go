package types

import (
	"github.com/gin-gonic/gin"
)

type HandlerFunc func(c *gin.Context) error
