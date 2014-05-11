resizeImages
============

Simple utility to resize images within a directory.

## Parameters:
 * p, containing folder of the images to be resized
 * w, new width
 * h, new height
 * rw, ratio width, auto calculates height
 * rh, ratio height, auto calculates width
 * precent, scales the image by a percentage maintaining its aspect ratio

## Usage:
* Ratio width and height will use their coorisponding height and width values provided in the commandline arguments.  Both ratio height and width cannot be provided at the same time.
* Percent scaling will ignore given width and height, and is incompatable with ratio width and height

## Example:
> ./resizeImages -p /home/user/Photos -w 800 -h 600
> ./resizeImages -p /home/user/Photos -w 800 -rw
> ./resizeImages -p /home/user/Photos -precent 50

[![baby-gopher](https://raw2.github.com/drnic/babygopher-site/gh-pages/images/babygopher-badge.png)](http://www.babygopher.org)