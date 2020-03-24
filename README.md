# I Avatars
generate avatars using name initials

https://iavatar.cc

---------

## Usage

### Use with default setting

```
GET  https://iavatar.cc/image/?name=I+Avatar
```

![default](https://iavatar.cc/image/?name=I+Avatar "Default")


### Params

#### size

size is the size of generated images in pixel, max size is 1024

```
GET https://iavatar.cc/image/?name=秦始皇&size=200
```

![size](https://iavatar.cc/image/?name=秦始皇&size=200 "Size")

#### name

name used to generate initials, split by space or + , max number of initials is 2

```
GET https://iavatar.cc/image/?name=Very+Long+Name&size=100
```

![name](https://iavatar.cc/image/?name=Very+Long+Name&size=100 "Name")

#### color

color will be calculated automatically by names

```
GET https://iavatar.cc/image/?name=Different+Name&size=100
```
![name](https://iavatar.cc/image/?name=Different+Name&size=100 "Color")