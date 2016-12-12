

DEPLOYER
===

Is a webhook, that receives a phabricator harbormaster call, and 
execs the script provided by cmd

Every call to execute, needs a phid, param, that is passed as first argument
to the provided script. phid is used to report back to phabricator, 
the status of the current build.

Also is used to attach as an artifact the output log from the shell script.


