gomon
=====

gomon is a package that simplifies the process of sending metrics to AWS
CloudWatch.  It provides a few things:

  - A registry and event loop that collects monitoring data as frequently as
    you like, then batches them up to be periodically delivered to AWS,
    minimizing API calls and cost, while preserving granularity.
  - Helper functions to aggregate and transform metric data in useful ways.
  - A few prepackaged monitors to collect basic system data like disk and
    memory statistics, so that you don't need to set up redundant monitoring
    scripts to get OS metrics in addition to application-specific ones.

## How To Use

#### Prerequisites
You must have a set of AWS credentials that allow the PutMetricData API call on
the resources you want to create.  The following
[IAM](http://docs.aws.amazon.com/IAM/latest/UserGuide/PoliciesOverview.html)
policy should work:

```json
{
  "Statement" : [
    {
      "Effect" : "Allow",
      "Action" : "cloudwatch:PutMetricData",
      "Resource" : "*"
    }
  ]
}
```

Make sure these credentials are exported as environment variables at
`AWS_ACCESS_KEY_ID` and `AWS_SECRET_ACCESS_KEY`.

#### Register Emitters
Choose some metrics you want to watch, and call one of the
[`Register*`](https://github.com/PeopleAdmin/gomon/blob/master/registry.go)
functions, passing it a method value or anonymous function that when called,
returns the current value of the metric.  Given a function in `myapp` called
`ConnectionCount()`, to track its value 60 times per second, you would register
as:

```go
gomon.RegisterInt("CurrentConnections", 60, "Count", myapp.ConnectionCount)
```

or, to report on something more complicated:
```go
currentRegisteredUsers := func() int { return currentUsers() - currentGuests() }
gomon.RegisterInt("CurrentRegisteredUsers", 60, "Count", currentRegisteredUsers)
```

If you have a metric which represents some sort of accumulating total, you can
use `DeltaSinceLastCall()` to track its rate of change.  For example:

```go
processedSinceLastCall = DeltaSinceLastCall(myapp.TotalJobsProcessed)
gomon.Register("JobsProcessedPerMinute", 1, "Count", ProcessedSinceLastCall)
```

There's an equivalent shorthand for the above common pattern:
```go
gomon.RegisterDelta("JobsProcessedPerMinute", 1, "Count", myapp.TotalJobsProcessed)
```

#### Enable Dimensions
Cloudwatch datapoints frequently have
[Dimensions](http://docs.aws.amazon.com/AmazonCloudWatch/latest/DeveloperGuide/cloudwatch_concepts.html#Dimension)
attached to them to characterize the data, e.g. marking which server it came
from or what application produced it. You can add dimensions that will be
present in all transmitted data:

```go
gomon.AddDimension("InstanceId", "i-31a74258")
gomon.AddDimension("ApplicationName", "MyApp")
```

#### Start Sending
After all emitters are registered, start the loop with `gomon.Start()`.
