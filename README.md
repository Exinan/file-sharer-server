# file-sharer-server
Create a server that allows you to download, upload and delete files 

Comands
:
http://localhost:8080/upload

http://localhost:8080/download

http://localhost:8080/delete

http://localhost:8080/shutdown

Examples:

curl -X POST -F "file=@FILE PATH/FILE NAME" http://localhost:8080/upload

curl -O http://localhost:8080/download/FILE NAME

curl -X DELETE http://localhost:8080/delete/FILE NAME

curl -X POST -k http://localhost:8080/shutdown
