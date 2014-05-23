# Ikura consists of eggs for Sushi

## API

### Put a file

Request:

	PUT / HTTP/1.1
	Host: www.example.com
	Authorization:XXX-FIEXAMPLE=
	Content-Length: 65534
	[file content]

Response:

	HTTP/1.1 200 OK
	[FID]

FID Representation:

	machine_id(4bytes) + md5(32bytes) + file_size(8bytes)
	00014dc8263cf28b66502d0d758sdfasdfasfasdqwer

### Delete a file

Request:

	DELETE /[FID] HTTP/1.1
	Host: www.example.com
	Authorization:XXX-FIEXAMPLE=

Response:

	HTTP/1.1 200 OK

### Get a file

Request:

	GET /[FID] HTTP/1.1
	Host: www.example.com

Response:

	Content-Type:image/jpeg
	HTTP/1.1 200 OK
	[file content]

### File existence

Request:

	HEAD /[FID] HTTP/1.1
	Host: www.example.com

Response:

	Content-Type:image/jpeg
	HTTP/1.1 200 OK