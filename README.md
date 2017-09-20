# mysqldumpsplit
[![CircleCI](https://circleci.com/gh/afrase/mysqldumpsplit.svg?style=svg)](https://circleci.com/gh/afrase/mysqldumpsplit)

Split a mysqldump into separate files for each table.

```
Usage of mysqldumpsplit:
  -combine
    	Combine all tables into a single file, deletes individual table files
  -combineFile string
    	The path to output a single SQL file
	Only used if combine flag is set (default "dumpfile.sql")
  -i string
    	The file to read from, can be a gzip file
  -o string
    	The output path  (default "output")
  -skipData value
    	Comma separated list of tables you want to skip outputing the data for.
	Use '*' to skip all.
```

# Performance

Using a sort of pipelined approach, mysqldumpsplit is very memory efficient.
The amount of memory used will mostly be based on how large each line in the file is.

I was able to split a **130 table 16GB** sql file in less than **30 seconds** and use less than **20MB of RAM**.

# Usage

Split a SQL dump into individual tables.

`./mysqldumpsplit -i ~/Downloads/dump.sql -o ~/Downloads/output`

Skip the data for specific tables.

`./mysqldumpsplit -i ~/Downloads/dump.sql -o ~/Downloads/output -skipData table1,table2,table3`

Combine all tables back into a single file. Deletes the output directory once finished unless the combined 
file is in the same directory.

`./mysqldumpsplit -i ~/Downloads/dump.sql.gz -o ~/Downloads/output -skipData table1,table2,table3 -combine -combineFile ~/Downloads/dumpfile.sql`

Skip data for all tables.

`./mysqldumpsplit -i ~/Downloads/dump.sql.gz -o ~/Downloads/output -skipData *`

# Know issues
- If everything in the file is on a single line this will not work.
- Even if skipping the data for a table, it still must be read by the program.
