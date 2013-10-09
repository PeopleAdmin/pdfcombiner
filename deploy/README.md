Deploying
=========

These scripts support deployment to an AWS environment using [CloudFormation](http://aws.amazon.com/cloudformation/) to build instances and manage autoscaling.  The app servers sit behind an [Elastic Load Balancer](http://aws.amazon.com/elasticloadbalancing/) instance which terminates SSL.  Instances run Ubuntu 12.04 with the pdfcombiner daemon running as the `ubuntu` user.

Deploying new code without service interruption should be as simple as running `./deploy/release`. The scripts will help you out if you're missing dependencies.  If everything goes well this will build a linux-amd64 copy of the project, copy it to an S3 bucket, and trigger all members of the autoscaling group to pick up the new binary and restart themselves.

Dependencies
------------

The deployment scripts assume:

 1. An S3 bucket is prepared that already contains a linux copy of `cpdf` in its root.
 2. `fog` and `aws_cli` are installed for interacting with amazon, and you have a set of AWS credentials that have permission to write to the deployment bucket and interact with CloudFormation (i.e. pretty much admin access)
 3. The [go source code](https://code.google.com/p/go) is somewhere on your system (not via `brew install`).
 4. You have access to the SSH private key used by the autoscaling group (by default, it's called `pdfcombiner`)
 5. You've already set up an ELB and security group and know what they're called.

You can't bootstrap a stack with these scripts, they're only for redeploying.  To perform an initial deployment, visit the [AWS CFN console](https://console.aws.amazon.com/cloudformation/home?region=us-east-1) and upload [`cloudformation.json`](https://github.com/PeopleAdmin/pdfcombiner/blob/build_deploy/deploy/cloudformation.json).

Scripts
-------

`./deploy/release` invokes the following, which can also be run separately:

 - `./deploy/build` -- set up go and crosscompile for the linux target
 - `./deploy/copy` -- copy the S3 binary to the deploy bucket
 - `./deploy/restart_instances` -- SSH to each running instance, redeploy the binary from the deploy bucket and restart the service daemon.

Other administration scripts:

 - `./deploy/resize` -- adjust the size of the ec2 instances up or down and perform a rolling rebuild.
 - `./deploy/ssh` -- run a command on every running machine a la `knife ssh`

Lower level access to the deployment internals can be got via `deploy/lib/deployer.rb`.

