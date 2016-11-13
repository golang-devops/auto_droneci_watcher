@echo off
cls
SET ERRORLEVEL=0

::——————————————————— This will make it quit immediately upon pressing Ctrl+C instead of first asking yes/no
IF "%~1"=="–FIX_CTRL_C" (
SHIFT
) ELSE (
CALL <NUL %0 –FIX_CTRL_C %*
GOTO EOF
)
::———————————————————

CALL .vars_win.bat & if errorlevel 1 goto ERROR

REM call .git_version_win.bat & if errorlevel 1 goto ERROR
call .git_version_win.bat

echo install^
  & go install %RELATIVE_GO_SRC_PATH%/...^
  & if errorlevel 1 goto ERROR
echo build^
  & go build -o "%OUT_BIN%" -ldflags "-X main.GitSha1=%GIT_SHA%" ^
  & if errorlevel 1 goto ERROR





goto SUCCESS

:SUCCESS
echo Success!!
goto EOF

:ERROR
echo ERROR!!! See the last ran command
pause
goto EOF

:EOF
