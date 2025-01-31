# wp-gcs

upload the wordpress images to gcs.

- read the config to get database and directory informations.
- migration database table. be carefull the length should be 255 for varchar.
- one productor to read directory and push the filename to channel buffer, the buffer should be had cache.
- 200 workers to consume the file from channel buffer.
    - check file if is exists.
    - check file if is in mysql database.
    - upload to gcs
    - insert into mysql database to record.
    - done

