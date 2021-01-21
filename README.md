# Concurrent file downloader

This project, is a CLI tool, for downloading files via HTTP, in a concurrent manner. It leverages modern CPU architectures in order to manage download at a faster pace.

##

In order to build make sure that you have `make` tool installed, and once cloned this repository run:

```
$ make build
```

Then you can run the app by using following command:

```
$ ./main
```

Keep in mind, that by default it will expect to have `urls.json` file for URLs, that are expected to be downloaded. You will need to create that file for a default functionality, otherwise, use `file` flag.

The full options list you can see here.

```
Usage of ./main:
  -file string
        Path to files that contains the URLs (default "urls.json")
  -format string
        Set the downloadable files format. Note, that all files have to have same MIME type (default "mp4")
  -poolsize int
        Number of concurrent downloading threads (default 3)
```

## Known issues and requirements list

- A need for a flag defined output directory, in which the files will be downloaded
- A pattern for output file names, defined by a flag, with a specific syntax
