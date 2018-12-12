# get simple directory statistics (metrics) for prometheus

- all done without prometheus libraries
  - don't know if this is a good decesion or a bad one.
- super simple, the only things provided so far:
  - file item (directory or file) of directory, without subdirectories
  - age of oldest file in directoy

the purpose of this exporter is to generate an alert within prometheus / grafana if there are files waiting longer than a certain time.

exports:
- `dirstat_files_count`: number of files in directory
- `dirstat_oldest_file_age`: age of oldest file in seconds