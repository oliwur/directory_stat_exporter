# directory statistics for prometheus
get simple directory statistics (metrics) for prometheus

## features
- all done without prometheus libraries
  - don't know if this is a good decesion or a bad one.
- super simple, the only things provided so far:
  - file item (directory or file) of directory, without subdirectories
  - age of oldest file in directoy

the purpose of this exporter is to generate an alert within prometheus / grafana if there are files waiting longer than a certain time.

## exports
- `dirstat_files_in_dir`: number of files in directory
- `dirstat_oldest_file_time`: timestamp (unix time) of oldest file in dir

## todos
- make sure only files are counted (done)
- implement recursive file walking (done)
- test handling of unc paths in windows (yes, it's targeted for windows.)
- better logging
- better error handling
- make information gathering concurrent, so more directories can be handled in the same time

## notes
- *important* stack items correctly (types and help text must only appear once in a metric export / per request)
- note to self: labels must not contain a single backslash... I replaced all backslashes now with forward slashes. -> there must be a better solution
  - e.g. add labels to the configuration and give them meaningful names.

## problems
- large directories might not be handled well
  - might use lot of memory, because whole directory is read once (untested)
