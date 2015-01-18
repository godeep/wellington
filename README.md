[![Circle CI](https://circleci.com/gh/wellington/wellington/tree/master.svg?style=svg)](https://circleci.com/gh/wellington/wellington/tree/master)
[![Coverage Status](https://coveralls.io/repos/wellington/wellington/badge.png?branch=master)](https://coveralls.io/r/wellington/wellington?branch=master)

Wellington
===========

Wellington adds spriting to the lightning fast [libsass](http://libsass.org/). No need to learn a new tool, this all happens right in your Sass!

### Speed Matters

Some benchmarks

```
# Vanilla Sass from sass-spec
wt vanilla_css_huge.scss        0.10s user 0.01s system 98% cpu 0.120 total
compass vanilla_css_huge.scss   0.27s user 0.04s system 98% cpu 0.315 total
# 2.7x speedup
```


```
# 40,000 line of code Sass project
wt        14.935s
compass   73.800s
# 5x speedup!
```


#### Example

file.scss
```
$images: sprite-map("sprites/*.png");
div {
  width: image-width(sprite-file($images, "cat"));
  height: image-height(sprite-file($images, "cat"));
  background: sprite($images, "cat");
}
```

file.css
```
div {
  width: 140px;
  height: 79px;
  background: url("genimg/sprites-wehqi.png");
}
```
#### Try before you buy

Fork your own! [http://codepen.io/pen/def?fork=KwggLx](Wellington Playground)

Or, check out the collection of Wellington [http://codepen.io/collection/DbNZQJ/](demos)

#### Installation
Wellington can be installed via brew

	brew install wellington
	wt -h

Run Wellington in docker

	docker run -v $(pwd):/data -it drewwells/wellington wt proj.scss

## Documentation

### Mixins

#### @include sprite-dimensions($map, "file")

Sprite-dimensions outputs the height and width properties of the specified image.

```
div {
  @include sprite-dimensions($spritemap, "file");
}
```

*Output*

```css
div {
	width: 100px;
	height: 50px;
}
```

### Functions

Don't see a function you want?  Check out [handlers](http://godoc.org/github.com/wellington/wellington/context/handlers) and submit a pull request!

#### sprite-map("glob/pattern"[, $spacing: 10px])

sprite-map generates a sprite from the matched images optinally with spacing between the images.  No output is generated by this function, instead the return is used in other functions.

```
$spritemap: sprite-map("*.png");
```

*Output*

```css

```

#### sprite($map, $name: "image"[, $offsetX: 0px, $offsetY: 0px])|

sprite generates a background url with background position to the position of the specified `"image"` in the spritesheet.  Optionally, offsets can be used to slightly modify the background position.

```
div {
	background: sprite($spritemap, "image");
}
```

*Output*

```css
div {
	background: url("spritegen.png") -0px -149px;
}
```

#### sprite-file($map, $name: "image")

Sprite-file returns an encoded string only useful for passing to image-width or image-height.

```
div {
	background: sprite-file($spritemap, "image");
}
```

*Output*

```css
div {
	background: {encodedstring};
}
```

#### image-height($path)

image-height returns the height of the image specified by `$path`.

```
div {
	height: image-height(sprite-file($spritemap, "image"));
}
div {
	height: image-height("path/to/image.png");
}
```

*Output*

```css
div {
	height: 50px;
}
div {
	height: 50px;
}
```

#### image-width($path)

image-width returns the width of the image specified by `$path`.

```
.first {
	width: image-width(sprite-file($spritemap, "image"));
}
.second {
	width: image-width("path/to/image.png");
}
```

*Output*

```css
.first {
	width: 50px;
}
.second {
	width: 50px;
}
```

#### inline-image($path[, $encode: false])

inline-image base64 encodes binary images (png, jpg, gif are currently supported). SVG images are by default url escaped. Optionally SVG can be base64 encoded by specifying `$encode: true`. Base64 encoding incurs a (10-30%) file size penalty.

```
.png {
	background: inline-image("path/to/image.png");
}
.svg {
	background: inline-image("path/to/image.svg", $encode: false);
}
```

*Output*

```css
.png {
	background: inline-image("data:image/png;base64,iVBOR...");
}
.svg {
	background: inline-image("data:image/svg+xml;utf8,%3C%3F...");
}
```

#### image-url($path)

image-url returns a relative path to an image in the image directory from the built css directory.

```
div {
	background: image-url("path/to/image.png");
}
```

*Output*

```css
div {
	background: url('../imgdirectory/path/to/image.png");
}
```

### font-url($path, [$raw:false])

font-url returns a relative path to fonts in your font directory.  You must set the font path to use this function.  By default, font-url will return `url('path/to/font')`, set `$raw: true` to only return the path

```
div {
	$path: font-url("arial.eot", true);
	@font-face {
		src: font-url("arial.eot");
		src: url("#{$path}");
	}
}
```

*Output*

```css
div {
	@font-face {
		src: url("../font/arial.eot");
		src: url("../font/arial.eot");
	}
}
```

### Why?

For the life of Sass, there has been only one tool for doing spriting with Sass. Compass has become to the daily workflow at work, that we really couldn't venture into new tools. As the website grew, Compass and Ruby Sass started to become a real drag on build times. A typical build including transpiling Sass to CSS, RequireJS JavaScript, and minfication of CSS, JS, and images would spend half the time in Compass.

There had to be a better way. I took a look at libsass but besides [image-url](https://github.com/sass/libsass/issues/489) there was no support for the spriting functions we loved from Compass. So I wrote Wellington to be a drop in replacement for all spriting functions found in Compass. This makes it super simple to swap out Compass with Wellington in your Sass projects.

### See how the sausage is made

#### Building from source
Install Go and add $GOPATH/bin to your $PATH. [Detailed instructions](https://golang.org/doc/install)

```
go get -u github.com/wellington/wellington
cd $GOPATH/src/github.com/wellington/wellington
#install libsass
make deps 

PKG_CONFIG_PATH=$(pwd)/libsass/lib/pkgconfig go get -u github.com/wellington/wellington/wt
wt -h
```

It's a good idea to export `PKG_CONFIG_PATH` so that pkg-config can find `libsass.pc`.

Set your fork as the origin.

    cd $GOPATH/src/github.com/wellington/wellington
	git remote rm origin
	git remote add origin git@github.com:username/wellington.git

Testing

    make test

Profiling

	make profile

Docker Container

	make build
	make docker #launch a container

Please use pull requests for contributing code.  [CircleCI](https://circleci.com/gh/wellington/wellington) will automatically test and lint your contributions.  Thanks for helping!

### License

Wellington is licensed under MIT.
