# Guigoes - A serverless blog using GO and HTMX
 
This is the source code for my personal blog https://guigoes.com, it was built using Go and HTMX and deployed on Fly.io.

## Table of Contents

- [Introduction](#introduction)
- [Features](#features)
- [Technologies Used](#technologies-used)
- [Getting Started](#getting-started)
  - [Installation](#installation)
  - [Usage](#usage)
  - [Deploy](#deploy)
- [Contributing](#contributing)
- [License](#license)

## Introduction

I wanted to create a reasonably performant blog that is cheap to run. Go and HTMX were a perfect fit for this, Go programs tend to be lightweight and has great support for text templating which helps a lot to rendering HTML fragments for the HTMX web app.

## Features

- Responsive design
- Full-text search
- Markdown to HTML parsing for blog posts
- Easy deploy

## Technologies Used

- Go - https://github.com/golang/go
- TailwindCSS (standalone CLI) - https://github.com/tailwindlabs/tailwindcss 
- Gin - https://github.com/gin-gonic/gin
- HTMX - https://github.com/bigskysoftware/htmx
- Templ - https://github.com/a-h/templ
- Goldmark - https://github.com/yuin/goldmark
- Bleve - https://github.com/blevesearch/bleve
- Fly.io - https://fly.io

## Getting Started

You will need 
1. Go 1.1.8 or later
2. Make

### Installation

```bash
$ git clone hhttps://github.com/guilycst/guigoes
$ cd guigoes
$ make install
```

```make install``` will run the shell script [install.sh](./install.sh) to install some of the dependencies needed to run and build locally:

- [TailwindCSS](https://tailwindcss.com/) - CSS Framework
- [Air](https://github.com/cosmtrek/air) - Live reloading
- [Templ](https://templ.guide/) - HTML templating language for Go

This step will also run ```go mod tidy``` and ```go mod download``` so all Go's dependencies are available to run the project.

### Usage

First you'll need to change the POSTS_PATH entry on the .env file so it points correctly to the posts dir on this repository

```environment
POSTS_PATH="PATH_TO_REPO_HERE/posts/" 
DIST_PATH="./web/dist/"
BLEVE_IDX_PATH="blog.bleve"
```

Then you can run using ```$ make run```, this will start a local web server at port :8080 with live reloading using air, so changes made to the code should reflect automatically.


![Home page](./docs/home_page.png)

### Deploy

This project was initially meant to be deployed on AWS Lambdas, which do reach my goal of being cheap to run, but the sheer number of artifacts that AWS tooling (CDK) generated is kinda off-putting.

So to keep things as simple as possible i`ve opted for Fly.io, but realistic, any hosting solution that can run binaries would do.

### License
[GNU GENERAL PUBLIC LICENSE Version 3](./LICENSE)

