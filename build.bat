@echo on

set BINARY_NAME="bot"
set PATH_MAIN="./main/"

go build -o %BINARY_NAME% %PATH_MAIN%
echo Build Complete!
pause 1
