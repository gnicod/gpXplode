gpXplode
========


#### Why?
If you have one or more big gpx files containing more than one day, and you want 1 file per day.

#### Usage

```shell
$ ls
    bigfile1.gpx bigfile2.
$ gpXplode --output /tmp bigfile1.gpx bigfile2.gpx 
$ ls /tmp
    19-December-2014.gpx 20-December-2014.gpx 21-December-2014.gpx
```
