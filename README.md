pdfcombiner
===========

This is an HTTP endpoint that downloads a list of PDFs from Amazon S3, 
combines them using `cpdf`, uploads the combined file, and POSTs the job
status to a provided callback URL.  The format of the request should look
like:

    POST  http://localhost:8080?docs=1.pdf&docs=1.pdf?callback=http://myurl/update
    
Both parameters are required.

The endpoint will immediately either return `200 OK` if it understood the 
request, or `400 Bad Request` if not.  Work is done asynchronously in the
background, and when complete, the provided callback URL will recieve a 
json-encoded message indicating success or failure.
