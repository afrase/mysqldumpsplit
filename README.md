# mysqldumpsplit
Split a mysqldump into separate files for each table.

```
Usage of mysqldumpsplit:
  -combine
        Combine all tables into a single file
  -combineFile string
        The path to output a single SQL file
        Only used if combine flag is set (default "dumpfile.sql")
  -i string
        The file to read from
  -o string
        The output path  (default "output")
  -skipData value
        Comma separated list of tables you don't want the data for
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

Combine all tables back into a single file. Deletes the output directory once finished.

`./mysqldumpsplit -i ~/Downloads/dump.sql -o ~/Downloads/output -skipData table1,table2,table3 -combine -combineFile ~/Downloads/dumpfile.sql`

# Know issues
- Don't save the combined file into the output directory or the file grow until you're out of drive space.
- If everything in the file is on a single line this will not work.
