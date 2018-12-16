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
- `dirstat_files_count`: number of files in directory
- `dirstat_oldest_file_age`: age of oldest file in seconds

## todos
- make sure only files are counted (done)
- implement recursive file walking (done)
- make information gathering concurrent, so more directories can be handled in the same time
- better logging
- better error handling
- test handling of unc paths in windows (yes, it's targeted for windows.)
- *important* stack items correctly

## problems
- it can only handle one directory at the moment because it does not stack the output correctly according to prometheus.
  - metrics that have the same metric name must be bundled and must not have repeating help and type information.
- large directories might not be handled well
  - might use lot of memory, because whole directory is read once (untested)
