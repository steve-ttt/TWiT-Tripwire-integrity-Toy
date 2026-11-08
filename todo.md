[x] Hash an empty file and match the known SHA-256 of an empty input
[x] Hash a small file with the known digest for "abc"
[x] Return an error for a missing file
[x] Hash a directory tree and return a map of relative path to digest
[x] read a json config file of files/directories to monitor
[x] create a database to perform CRUD on the hashes for the system files
[x] Save a baseline to sqllite and load it back unchanged
[ ] scan files/directories in the config and store them in the db
[ ] Report a file as modified when its digest changes
[ ] Report a file as added or removed
[ ] Report a permission or owner change
