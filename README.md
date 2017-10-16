# twidel

A command line tool to delete your tweets.

## Usage

	# go run twidel.go [-options]

```
> twidel.exe --help
Usage of twidel:
  -dbg
        Debug mode on if dbg=true
  -limit int
        Limit of number to delete tweets (default 3200)
  -minfav int
        Delete tweet less than minfav (default 114514)
  -minrt int
        Delete tweet less than minrt (default 114514)
```

## Example

	$ go run twidel.go -limit 1000 -minfav 10 -minrt 5

In this case, up to 1000 tweets that is favorited less than 10 and Retweeted less than 5 will be deleted.

If you don't use any options, up to 3200 tweets will be deleted.