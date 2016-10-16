dusage
======

Simple du(1) wrapper that lists directory sizes. It will try to get the disk
usage for all directories in the given directory ("." if not given), sort the
result by size in bytes, and print it in a human readable way.

Example to get top ten directories in "Library":

    > dusage Library | tail   
    868.00K	Library/Keychains
      2.89M	Library/Preferences
      4.39M	Library/Logs
      5.31M	Library/Safari
      7.62M	Library/Saved Application State
     26.24M	Library/Calendars
     55.94M	Library/Containers
     95.50M	Library/Application Support
      1.19G	Library/Caches
      3.25G	Library/Mail

Install
-------

go get github.com/sstark/dusage

