@echo off
cls
SET ERRORLEVEL=0

SET PROJ_NAME=auto_droneci_watcher
SET RELATIVE_GO_SRC_PATH=github.com/golang-devops/auto_droneci_watcher
SET OUT_BIN=%GOPATH%\bin\%PROJ_NAME%.exe

exit /b 0
