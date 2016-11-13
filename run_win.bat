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

CALL build_win.bat & if errorlevel 1 goto ERROR

REM SET LOG_LEVEL=debug
SET LOG_LEVEL=info
SET CONFIG_PATH=%UserProfile%\.config\%PROJ_NAME%\config.yml
echo run^
  & "%OUT_BIN%" -config "%CONFIG_PATH%" -loglevel "%LOG_LEVEL%"^
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
