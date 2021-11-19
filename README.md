## About

This tutorial will walk you through main concepts of a Mantil project.

## Prerequisites

 * Mantil command line tool - [install instructions](https://github.com/mantil-io/mantil-docs#installation)
 * [Go](https://golang.org) 
 * access to an AWS account
 
### Mantil node on the AWS account

Before we start with the project we need to setup Mantil node on AWS. Setting up
node is only time you will need your AWS credentials. Node is set of Lambda
functions and supported resources which will be used by Mantil to
deploy/upgrade/remove project Lambda functions. Node is installed in the
specific region of an AWS Account.

To setup new node run something like:

``` 
➜ mantil aws install try-ping --aws-profile org5
```

This will setup node named try-ping in AWS account and region defined in my
org5 AWS profile. There are various options how you can provide AWS credentials
for setting up new node to Mantil. Use `mantil aws install --help` to view them.

For in depth explanation of setting up Mantil node see [this](todo) instructions.

## New Project

Create your Mantil project from the ping template:

``` 
➜ mantil new my-ping --from https://github.com/mantil-io/template-ping

Creating my-ping in /tmp/my-ping...
Replacing import paths with my-ping...

Your project is ready in /tmp/my-ping
```

This command clones [template-ping](https://github.com/mantil-io/template-ping)
repository. _my-ping_ is the folder and the name of the Mantil project.

For all project command ensure that you are positioned somewhere in the project
folder tree.

``` 
cd my-ping
```

## Project stage

Next step is to create first stage for the project. Stages are actual
installations of the project in AWS. A project can have multiple stages.A stage
for each developer, integration stage, production...

Stage is created on the Mantil node. So you need to specify node when creating
new stage.

```
➜ mantil stage new development --node try-ping
Using node try-ping for new stage

Creating stage development and deploying project my-ping
==> Building...
ping

==> Uploading...
ping

==> Setting up AWS infrastructure...
	Initializing, done.
	Planning changes, done.
	Creating infrastructure 100% (35/35), done.

Deploy successful!
Build time: 2.537s, upload: 830ms (5.3 MiB), update: 57.681s

Stage development is ready!
Endpoint: https://qd3tidvbuf.execute-api.eu-central-1.amazonaws.com
```

This creates resources on AWS. Most importantly Lambda function for each of you
API's and an API Gateway to expose those functions on an URL. That URL is shown
at the end of command and could be found any time by:

``` 
➜ mantil env --url
https://qd3tidvbuf.execute-api.eu-central-1.amazonaws.com
```

## Invoke API method

This project has only one API. All API's are in project _/api_ folder. Ping API is
in _/api/ping_ folder. It is exposed at _[endpoint]/ping_ URL. That URL leads to the
Default method. All other exported methods have URL _[endpoint]/ping/[method]_.

You can use curl to reach API methods:

``` 
➜ curl -X POST $(mantil env --url)/ping
pong

➜ curl -X POST $(mantil env --url)/ping/hello
Hello, 

➜ curl -X POST $(mantil env --url)/ping/hello --data Foo
Hello, Foo
```

Easier and with added features way is to use `mantil invoke` command:

``` 
➜ mantil invoke ping
200 OK
pong

mantil invoke ping/hello
200 OK
Hello,

➜ mantil invoke ping/hello --data Bar
200 OK
Hello, Bar
```

If the response is JSON invoke will pretty print that:

``` 
➜ mantil invoke ping/reqrsp --data '{"name":"Bar"}'
200 OK
{
   "Response": "Hello, Bar"
}
```

Invoke will show Lambda function logs during function execution. 

``` 
➜ mantil invoke ping/logs --data '{"name":"Bar"}'
λ start Logs method
λ req.Name: 'Bar'
λ end
200 OK
{
   "Response": "Hello, Bar"
}
```

## Deploy 

Make some change in the _api/ping/ping.go_. For example change return of the
Default method to something other than _"ping"__ string. Execute `mantil deploy` to
update stage.

``` 
➜ mantil deploy
==> Building...
ping

==> Uploading...
ping

==> Updating...
ping

Deploy successful!
Build time: 697ms, upload: 2.599s (5.3 MiB), update: 1.531s
```

Deploy consists of three parts. First it builds Lambda function binary from
API's code. Second it uploads every changed binary to S3. And in third part
updates Lambda function with new binary.

To support iterative build/test cycle there is `mantil watch` command. It will
watch project folder for changes. On each file save it will start new deploy.
You can configure watch to execute invoke or run test after deploy. Check
`mantil watch --help` for details.

Here is example of a watch cycle where I changed response of ReqRsp method two
times. Every file save triggered deploy and invoke after that.

```
➜ mantil watch --method ping/reqrsp --data '{"name": "Foo"}'

Watching changes in /tmp/my-ping

Changes detected! Starting deploy
==> Building...
ping

==> Uploading...
ping

==> Updating...
ping

Deploy successful!
Build time: 694ms, upload: 900ms (5.3 MiB), update: 1.466s

==> Invoking function
200 OK
{
   "Response": "Hello from Lambda, Foo"
}

Watching changes in /tmp/my-ping

Changes detected! Starting deploy
==> Building...
ping

==> Uploading...
ping

==> Updating...
ping

Deploy successful!
Build time: 1.333s, upload: 1.743s (10.7 MiB), update: 2.336s

==> Invoking function
200 OK
{
   "Response": "Hello from Mantil, Foo"
}

Watching changes in /tmp/my-ping
```

## Test

_/test_ folder contains example of end to end tests for Mantil project. Those
tests are using current stage endpoint to execute methods and examine returned
results.

If you were following this tutorial you probable made some changes in the
original method output so the tests will fail. Run `mantil test` to check that.

Running tests after each change in watch cycle is supported by `mantil watch --test`.

## Cleanup

We created a node and a stage in your AWS account. To remove that and leave
account in original state destroy stage and node. Commands for that are:

``` 
➜ mantil stage destroy development

➜ mantil aws uninstall try-ping --aws-profile org5
```

