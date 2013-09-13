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
  "doc_list": [
    {
      "key": "s3/path/to/file.pdf",
      "title": "name of pdf"
    },
    {
      "key": "s3/path/to/other/file.pdf",
      "name": "name of pdf"
    }
  ],
  "callback": "http://mycallbackurl.com/combination_result/12345",
  "combined_key": "path/to/upload/combined/file.pdf",
  "title": "Combined Doc for Some Applicant"
}
```

The server will immediately respond either with:

    HTTP/1.1 200 OK
    {"response":"ok"}

or

    HTTP/1.1 400 Bad Request
    {"response":"invalid params"}

When work is complete, the provided callback URL will recieve a POST
with a JSON body similar to:

```json
{
  "success": true,
  "errors": {},
  "callback": "http://mycallbackurl.com/combination_result/12345",
  "perf_stats": {
      "606/docs/1068.pdf": {
      "s3/path/to/file.pdf": {
          "Filename": "s3/path/to/file.pdf",
          "Size": 1234,
          "PageCount": 5,
          "DlTime": 622469262
      },
      "s3/path/to/other/file.pdf": {
          "Filename": "s3/path/to/other/file.pdf",
          "Size": 3456,
          "PageCount": 3,
          "DlTime": 622469262
      }
    }
  }
}
```

`"success"` is true if at least one file downloaded successfully.

This application requires authentication, please put a file in
~/.pdfcombiner.json with content like:

```json
{"remote_user":"u","remote_password":"pass"}
```

You can also combine files in standalone mode from the command line.
Use `./pdfcombiner -help` to get a list of options.
