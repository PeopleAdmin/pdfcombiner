[![Build Status](https://travis-ci.org/PeopleAdmin/pdfcombiner.png)](https://travis-ci.org/PeopleAdmin/pdfcombiner)
pdfcombiner
===========

This is an HTTP endpoint that downloads a list of PDFs from Amazon S3,
combines and uploads the combined file to S3, and POSTs the job
status to a provided callback URL.  The format of the request to `/` should
look like:

```json
{
  "bucket_name": "somebucket",
  "employer_id": 123,
  "doc_list": [
    "1.pdf",
    "2.pdf"
  ],
  "callback": "http://mycallbackurl.com/combination_result/12345"
}
```

The server will immediately respond either with:

    HTTP/1.1 200 OK
    {"response":"ok"}

or

    HTTP/1.1 400 Bad Request
    {"response":"invalid params"}

and begin processing the file.  When work is complete, the provided
callback URL will recieve a POST with a JSON body similar to:

```json
{
  "success": true,
  "combined_file": "path/to/combined/file.pdf",
  "job": {
    "bucket_name": "somebucket",
    "employer_id": 123,
    "doc_list": [
      "realfile.pdf",
      "nonexistent_file"
    ],
    "downloaded": [
      "realfile.pdf"
    ],
    "callback": "http://mycallbackurl.com/combination_result/12345",
    "errors": {
      "nonexistent_file": "The specified key does not exist."
    },
    "perf_stats": {
      "realfile.pdf": {
        "Filename": "realfile.pdf",
        "Size": 1292244,
        "DlTime": 1229882563
      }
    }
  }
}
```

`"success"` is true if at least one file downloaded successfully.
`"combined_file"` may be `null` if `success` is false.


You can also combine files in standalone mode from the command line.
Use `./pdfcombiner -help` to get a list of options.
