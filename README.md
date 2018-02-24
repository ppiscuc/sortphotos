# Sortphotos

## Sortphotos

## Reason

Sortphotos organizez photos in a subdirectory format using the EXIF information.
This becomes really handy when apps suddenlly delete or clear the EXIF information of photos moved from one device to another.

Sortphotos uses a the xiam/exif library that reads EXIF tags using [libexif](https://libexif.github.io) and CGO.

## Usecase

Sortphotos is a small utility programs that takes photos from this format:

```
DSC04660.JPG
```

to the RFC3339 format

```
2017/06-24/2017:06:24T12:44:04Z.JPG
```

while keeping all the Exif data

## How to run

To run the utility, you need to have `go` installed, and you can run:

```
go get -u github.com/xiam/exif
go run sortphotos.go -src photos/ -dst export/ -sort
```

## Contributing

Any help is greatly apreaciated. If you have feature requests, please use the issue tracker.

