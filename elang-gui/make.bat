@echo off

set NAME=elang-gui

:entry
if "%1"=="" goto make
if "%1"=="test" goto test
if "%1"=="clean" goto clean
goto end

:make
call :clean
echo MAKE:
echo WINDRES -o %NAME%_res.syso %NAME%_res.rc
windres -o %NAME%_res.syso %NAME%_res.rc || goto nowindres
echo go build -ldflags="-H windowsgui"
go build -ldflags="-H windowsgui" || goto end
::echo mt -manifest "%NAME%.exe.manifest"
::mt /nologo -manifest "%NAME%.exe.manifest" -outputresource:"%NAME%.exe" || echo you need mt.exe to compile portable exefile
echo.
call :test
goto end

:clean
echo CLEAN:
if exist *.exe del *.exe
if exist *.syso del *.syso
echo.
goto end

:nowindres
echo windres is mingw's res pack toolkit
echo if you need icon. please download it.
call :test
echo.
goto end

:test
echo TEST:
if exist %NAME%.exe %NAME%
echo.
goto end

:end