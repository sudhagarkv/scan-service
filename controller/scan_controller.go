package controller

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"

	"scan-service/models"
	"scan-service/service"
)

type ScanController struct {
	scanService service.ScanService
}

func (s ScanController) ProcessRequest(ctx *gin.Context) {
	var scanRequest models.ScanRequest
	err := ctx.ShouldBindBodyWith(&scanRequest, binding.JSON)
	if err != nil {
		log.Printf("Unable to bind request body %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	err = s.scanService.ProcessScanRequest(ctx.Request.Context(), scanRequest)
	if err != nil {
		log.Printf("Unable to process scan request %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status": "queued",
	})
}

func NewScanController(scanService service.ScanService) ScanController {
	return ScanController{
		scanService: scanService,
	}
}
