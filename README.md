# file-sharer-server
Create a server that allows you to download, upload and delete files 

Commands:

http://localhost:8080/ping

http://localhost:8080/upload

http://localhost:8080/download

http://localhost:8080/delete

http://localhost:8080/shutdown

Examples:

curl http://localhost:8080/ping

curl -X POST -F "file=@FILE PATH/FILE_NAME" http://localhost:8080/upload

curl -O http://localhost:8080/download/FILE_NAME

curl -X DELETE http://localhost:8080/delete/FILE_NAME

curl -X POST -k http://localhost:8080/shutdown
