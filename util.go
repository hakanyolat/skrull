package skrull

import (
	"github.com/gofiber/fiber"
	"github.com/valyala/fasthttp"
	"go.uber.org/zap"
	"os"
	"os/signal"
)

// statusIsMoved ...
func statusIsMoved(statusCode int) bool {
	return statusCode == fasthttp.StatusMovedPermanently ||
		statusCode == fasthttp.StatusFound ||
		statusCode == fasthttp.StatusSeeOther ||
		statusCode == fasthttp.StatusTemporaryRedirect ||
		statusCode == fasthttp.StatusPermanentRedirect
}

// configureGracefulShutdown ...
func configureGracefulShutdown(logger *zap.Logger, app *fiber.App) {
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit

	if err := app.Shutdown(); err != nil {
		logger.Fatal(err.Error())
	}
}

// createFullPath ...
func createFullPath(basePath, relativePath string) string {
	return basePath + relativePath
}

// createFullRelativePath
func createFullRelativePath(relativePath, queryString string) string {
	fp := relativePath
	if queryString != "" {
		fp += "?" + queryString
	}
	return fp
}