package main

import (
	"encoding/json"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/knadh/listmonk/internal/subimporter"
	"github.com/knadh/listmonk/models"
	"github.com/labstack/echo/v4"
)

// ImportSubscribers handles the uploading and bulk importing of
// a ZIP file of one or more CSV files.
func (a *App) ImportSubscribers(c echo.Context) error {
	// Is an import already running?
	if a.importer.GetStats().Status == subimporter.StatusImporting {
		return echo.NewHTTPError(http.StatusBadRequest, a.i18n.T("import.alreadyRunning"))
	}

	// Unmarshal the JSON params.
	var opt subimporter.SessionOpt
	if err := json.Unmarshal([]byte(c.FormValue("params")), &opt); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest,
			a.i18n.Ts("import.invalidParams", "error", err.Error()))
	}

	// Validate mode.
	if opt.Mode != subimporter.ModeSubscribe && opt.Mode != subimporter.ModeBlocklist {
		return echo.NewHTTPError(http.StatusBadRequest, a.i18n.T("import.invalidMode"))
	}

	// If no status is specified, pick a default one.
	if opt.SubStatus == "" {
		switch opt.Mode {
		case subimporter.ModeSubscribe:
			opt.SubStatus = models.SubscriptionStatusUnconfirmed
		case subimporter.ModeBlocklist:
			opt.SubStatus = models.SubscriptionStatusUnsubscribed
		}
	}

	if opt.SubStatus != models.SubscriptionStatusUnconfirmed &&
		opt.SubStatus != models.SubscriptionStatusConfirmed &&
		opt.SubStatus != models.SubscriptionStatusUnsubscribed {
		return echo.NewHTTPError(http.StatusBadRequest, a.i18n.T("import.invalidSubStatus"))
	}

	if len(opt.Delim) != 1 {
		return echo.NewHTTPError(http.StatusBadRequest, a.i18n.T("import.invalidDelim"))
	}

	// Open the HTTP file.
	file, err := c.FormFile("file")
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest,
			a.i18n.Ts("import.invalidFile", "error", err.Error()))
	}

	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	// Copy it to a temp location.
	out, err := os.CreateTemp("", "listmonk")
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError,
			a.i18n.Ts("import.errorCopyingFile", "error", err.Error()))
	}
	defer out.Close()

	if _, err = io.Copy(out, src); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError,
			a.i18n.Ts("import.errorCopyingFile", "error", err.Error()))
	}

	// Start the importer session.
	opt.Filename = file.Filename
	sess, err := a.importer.NewSession(opt)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError,
			a.i18n.Ts("import.errorStarting", "error", err.Error()))
	}
	go sess.Start()

	if strings.HasSuffix(strings.ToLower(file.Filename), ".csv") {
		go sess.LoadCSV(out.Name(), rune(opt.Delim[0]))
	} else {
		// Only 1 CSV from the ZIP is considered. If multiple files have
		// to be processed, counting the net number of lines (to track progress),
		// keeping the global import state (failed / successful) etc. across
		// multiple files becomes complex. Instead, it's just easier for the
		// end user to concat multiple CSVs (if there are multiple in the first)
		// place and upload as one in the first place.
		dir, files, err := sess.ExtractZIP(out.Name(), 1)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError,
				a.i18n.Ts("import.errorProcessingZIP", "error", err.Error()))
		}

		go sess.LoadCSV(dir+"/"+files[0], rune(opt.Delim[0]))
	}

	return c.JSON(http.StatusOK, okResp{a.importer.GetStats()})
}

// GetImportSubscribers returns import statistics.
func (a *App) GetImportSubscribers(c echo.Context) error {
	s := a.importer.GetStats()
	return c.JSON(http.StatusOK, okResp{s})
}

// GetImportSubscriberStats returns import statistics.
func (a *App) GetImportSubscriberStats(c echo.Context) error {
	return c.JSON(http.StatusOK, okResp{string(a.importer.GetLogs())})
}

// StopImportSubscribers sends a stop signal to the importer.
// If there's an ongoing import, it'll be stopped, and if an import
// is finished, it's state is cleared.
func (a *App) StopImportSubscribers(c echo.Context) error {
	a.importer.Stop()
	return c.JSON(http.StatusOK, okResp{a.importer.GetStats()})
}
