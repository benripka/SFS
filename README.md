# ECE_422_Project_1

2 seperate systems, one to run on the client-side (a shell program that can interpret the different commands and
interact with the SFS API over some protocole (http / https / ftp ?)), and the implementation of SFS that will run on a
linux machine in the Cybera cloud The client-side shell will do things like request authentication, ask for users /
groups to be created, ask for a directoy to be created, etc. While the SFS will interpret these requests, and execute
them, using the local server OS where possible. For example if a request comes in from an authenticated user to create a
directory, the SFS will first encrypt the directory name, then use the local os's mkdir function to create it under the
encrypted name.

Q: Should we encrypt the file before sending it to the server (aka should client encrypt it) or shoudl the server
encrypt it? A: We choose which way to do it, "which is more secure?" Just include decision in report

Q: If an external user tries to open a file in vim for example, he shoudln't be able to modify any files or anything.
The internal users shoudl also be notified.

Q: external users shouldn't see file names or contents; should they also not be able to see file permissions? A: See
file names and content in encrypted. Shouldnt' be able to see permissions.

Q: How should we maintain state between the client & server, for example current the current directory? Cookies? Sessions?
A: Could implement a session manager on the server-side like https://astaxie.gitbooks.io/build-web-application-with-golang/content/en/06.2.html

### Supported Operations

1. SFS should allow creating of groups and users like Unix file system.
2. SFS Users (i.e., after authentication) can create, delete, read, write, and rename files.
3. The file system should support directories, including home directory, like Unix file system.
4. Internal users should be able to set permissions on files and directories.
5. File names (and directory names) should be treated as confidential.
6. External users should not be able to modify files or directories without being detected.
7. External users should not be able to read file content, file names, or directory names in plain mode.

### Your SFS should meet the following requirements:

1. SFS should allow creating of groups and users like Unix file system.
2. SFS Users (i.e., after authentication) can create, delete, read, write, and rename files.
3. The file system should support directories, including home directory, like Unix file
system.
4. Internal users should be able to set permissions on files and directories.
5. File names (and directory names) should be treated as confidential.
6. External users should not be able to modify files or directories without being detected.
7. External users should not be able to read file content, file names, or directory names in
plain mode.

### Shell Commands

1. ls - List the contents of the current directory
2. pwd - show the current directory path
3. mkdir <directory_name> - Create a new directory in current directory
4. cd - Change the current directory (support ~/./..)
5. cat <file_name> - Show contents of file, line by line.
6. touch <file_name> - create a new file with provided name in current directory
7. mv <old_path> <new_path> - move a file from one location to another
8. chmod <octal_permissions> <file_name> - Set / reset permissions of file
9. rm <file_name> - Delete file
10. rename <file_name> <new_name> - Rename the file
11. vim <file_name> - Open the file in vim to edit 

## Design Choice Log

1. Decided to pass username and password as json encoded http POST body to avoid having it visible in the server log files
2. Should probably pass request arguments in json http body aswell for the same reason (file names unencrypted in http log files)
3. Use HTTPS to encode all in-transit data between client and server
4. Use HTTP header set cookie to HttpOnly mode to avoiding session hijacking.

## Setup Notes
- GCC will need to be installed on the server to ensure that the go-sqlite3 library c dependencies can be properly built

## Resources

- TLS Implementation https://betterprogramming.pub/create-secure-clients-and-servers-in-golang-using-https-aa970ba36a13
  - TLS Implementation https://medium.com/rungo/secure-https-servers-in-go-a783008b36da
- Session ID Implementation https://astaxie.gitbooks.io/build-web-application-with-golang/content/en/06.2.html
- Rate limiting + channels: https://eli.thegreenplace.net/2019/on-concurrency-in-go-http-servers#:~:text=Go's%20built%2Din%20net%2Fhttp,also%20lead%20to%20some%20gotchas.

