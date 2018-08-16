# pwned

## Leaked password checker (in Go)

This is a simple web-frontend for a leaked password checker in Go. It expects a leaked database in SQLite3 format, with two columns per row: The first one containing the SHA1 of the leaked password, and the second the number of times the password was seen in the wild. This project was directly inspired by the great work of Troy Hunt and his site, http://haveibeenpwned.com, but targets those who have concerns about freely typing their passwords on a site they don't control.

## Leaked password file format

For the procedure below, we expect the source file to be formatted with one SHA1 and count per line, pure text, with SHA1 hashes *in uppercase* (sorry about that, but Troy's files are use uppercase SHA1 hashes, and since I use his files, I'm following his standard).

Files should be formatted as:

```
SHA1:count
SHA1:count
...
```

Real life example:
```
000000005AD76BD555C1D6D771DE417A4B87E4B4:4
00000000A8DAE4228F821FB418F59826079BF368:2
00000000DD7F2A1C68A35673713783CA390C9E93:630
00000001E225B908BAC31C56DB04D892E47536E0:5
00000006BAB7FC3113AA73DE3589630FC08218E7:2
```

Also, make sure there are **no repeated hashes** in your file.

## Downloading a password list

You can download a very well maintained file directly from Troy's page at http://haveibeenpwned.com/Passwords. Scroll down to the "Downloading the Pwned Passwords list" section and click in one of the torrent links. These
are large files (~9GB compressed, ~22GB compressed) so make sure you have the required disk space to store and process the files.

## Converting the file to a SQLite3 database

1. Make sure you have even more disk space. A 22GB text file becomes a 57GB database. You have been warned.

1. Make sure your system has sqlite3 installed (on Debian and Debian-based systems like Ubuntu and Mint, use `sudo apt-get install sqlite3`).

1. Unpack the passwords file.

1. Create the sqlite3 database with:

    ```
    $ sqlite3 pwned.db

    create table pwned (
      hash text primary key not null,
      count int
    );
    ```

1. Start the import process with:

    ```
    .mode csv
    .separator :
    .import name-of-the-uncompressed-password-file.txt pwned
    ```

1. This will take a while. Go fetch coffee. Keep an eye for error messages.

## Using the program

1. Download and compile the frontend with `go get -v github.com/marcopaganini/pwned`. This will install the `pwned` binary under `$GOPATH/bin`.

1. Run the program with `$GOPATH/bin/pwned --dbfile=<path_to_your_database_file> &`

1. If everything goes well, access the frontend at http://localhost:8080

1. You can type your password directly, or a SHA1 hash. The frontend uses anything that looks like a SHA1 hash directly. This is useful for those who prefer to enter a SHA1 hash (instead of the password).
