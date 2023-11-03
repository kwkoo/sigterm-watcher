# SIGTERM Watcher

Waits for a `SIGTERM` and terminates as soon as it receives the signal.

This is meant to be deployed on Kubernetes.


## Invoking a subprocess

You can also pass arguments to `sigterm-watcher` - `sigterm-watcher` will invoke the arguments as a subprocess.

`stdout` and `stderr` of the subprocess will be streamed to `sigterm-watcher`'s `stdout` and `stderr`.

If `sigterm-watcher` receives a `SIGTERM` while the command is still running, the subprocess will be killed.


## Tailing a file

You can also tell `sigterm-watcher` to tail the contents of a file to `stdout` by setting the path to the file in the `LOG` environment variable.

If the file is removed, `sigterm-watcher` will wait till a new file appears and tail the contents of the new file.
