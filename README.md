# directory statistics for prometheus
get simple directory statistics (metrics) for prometheus. if you want to know if a batch job did not process all the files in a directory or a certain folder contains files with error information, this is the right place. the directory_stat_exporter provides minimal metrics to see how many documents are in a folder and the timestamp (unix time in seconds) of the oldest one. with this information, you can decide, depending on the specs of an interface, if everything is ok or not. the analysis of a directory can be generated from only the directory itself, without the subdirectories, or also including all the subdirectories. depending on your needs.

the purpose of this exporter is to provide metrics so prometheus can generate an alert if there are files waiting longer than a certain time.

## features
- all done without prometheus libraries
  - don't know if this is a good decesion or a bad one.
  - that's not really a feature, is it. oh, well, it's a fun project anyway.
- super simple, the only things provided so far:
  - number of files in directory, with or without subdirectories
  - last modified timestamp of oldest file in directoy, with or without subdirs.
  - for calculation reasons the current timestamp

## exports
- `dirstat_files_in_dir`: number of files in directory
- `dirstat_oldest_file_time`: timestamp (unix time) of oldest file in dir
- `dirstat_current_timestamp`: the current timestamp. because it's not provieded by prometheus (or I was not able to find it.)

## todos
- test handling of unc paths in windows (yes, it's targeted for windows.)
- better logging (really? do I need logging here?)
- better error handling
  - it should be fault tolerant and rather give useful metrics if an error occurs
- make information gathering concurrent, so more directories can be handled in the same time
- add the disk usage of the current directory as a metric
  - think of usefulness first. do we really need this? -> so far: no.
- add performance measurements:
  - how fast is the gathering
  - where are the limits on local drives
  - what are the limits on remote directories 

## notes to self
- *important* stack items correctly (types and help text must only appear once in a metric export / per request)
- labels must not contain a single backslash... I replaced all backslashes now with forward slashes. -> there must be a better solution
  - e.g. add labels to the configuration and give them meaningful names.

## problems
- large directories might not be handled well
  - might use lot of memory, because whole directory is read once (untested)

## configuration
### service port
`serviceport: "9999"`

this does not really need more explanation

### cache time
`cachetime: 5`

time of interval in minutes. must be an integer value.

this is the time interval which is used to poll the updates from the directories. it does not make sense to me to poll the directories every 15 seconds. if it is set to `0` then caching is disabled.

### directories
```
directories:
     - path: \tmp
       name: "tmp_dir"
       recursive: true
```
|key|value|
|---|-----|
|path|a path to the directory to be monitored. this can be a relative or an absolute path. for windows it can also be a UNC path.|
|name|a useful name. must be unique. not sure if special characters and spaces are a good idea. I have to check that.|
|recursive|should subdirectories also be analysed. can be `true` or `false`|

## usage
in your GOPATH type

`go get github.com/codestoke/directory_stat_exporter`

you'll find the binary in the `bin` directory of your GOPATH. the `config.yml` must be in the same directory of the executable (there are no other options yet)

currently, this generates the version in the master branch. this might not be very useful and can be buggy.

you can always download the source from a release, unpack it in your GOPATH and then compile it with

`go build`

have fun!

cheers, Oli