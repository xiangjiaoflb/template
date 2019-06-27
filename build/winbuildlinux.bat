:: windows下编译为linux可执行程序
set appname=template
set str=%0 
set filepath=
set filepath2=

:GOON
for /f "delims=\, tokens=1,*" %%i in (%str%) do (
    set filepath=%filepath%%filepath2%
    set str="%%j"
    set filepath2=%%i/
    goto GOON
)

cd %filepath%../main

for /f "delims=" %%t in ('go version') do set goversion=%%t

set hour=%time:~,2%
if "%time:~,1%"==" " set hour=0%time:~1,1%

for /F %%i in ('git rev-parse --short HEAD') do ( set gitversion=%%i)

set GOOS=linux
set GOARCH=amd64

go build -mod=vendor -i -v -o ../bin/%appname% -ldflags  "-s -w -X 'main.BUILDTIME=%date:~0,4%-%date:~5,2%-%date:~8,2% %hour%:%time:~3,2%:%time:~6,2%' -X 'main.GOVERSION=%goversion%' -X 'main.GITHASH=%gitversion%'"
