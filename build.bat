@echo off

SET GOROOT=C:\Program Files\Go
SET GOPATH=C:\Users\gmaga\go

mkdir build 2> NUL

:: *** Build Windows ***
echo * Building Windows executable
mkdir build\windows 2> NUL
SET GOOS=windows
SET GOARCH=amd64
del build\windows\datetimerunner.exe
"%GOROOT%\bin\go.exe" build -o build\windows\datetimerunner.exe src\datetimerunner.go

:: *** Build Linux ***
echo * Building Linux executable
mkdir build\linux 2> NUL
SET GOOS=linux
SET GOARCH=amd64
del build\linux\datetimerunner
"%GOROOT%\bin\go.exe" build -o build\linux\datetimerunner src\datetimerunner.go
