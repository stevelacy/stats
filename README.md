# stats

> Page count and voting counter for static or dynamic systems


### Install

Build from source as a stand-alone binary: `make`

##### Docker (recommended)

Copy the `docker-compose.yaml` file and change the `VIRTUAL_HOST` to match your domain.

```
docker-compose up
```

### Usage

##### JavaScript
To increment the page count:

`GET <stats host>/<page name>`

`POST <stats host>/<page name>`

```js

fetch('http://stats.example.com/page_one').then(() => {})

// Or use POST:
fetch('http://stats.example.com/page_one', { method: 'POST' }).then(() => {})

```

##### Image

You can use an image if you do not want to use JavaScript.

`<img src="<stats host>/<page name>/stat.gif" />`

```html

<body>
  <img src="http://stats.example.com/page_three/stat.gif" height="0" width="0" />
</body>

```


To retrieve the page count statistics make a request to the `page name/stats`

`GET <stats host>/<page name>/stats`

```js

fetch('http://stats.example.com/page_one/stats')
  .then((res) => res.json())
  .then((res) => console.log(res))
  // {count: 448}

```

For fully static-only sites you can use an `iframe` to show the output:

```html
<style>
  .stats {
    height: 20px;
    width: 50px;
    overflow: hidden;
  }
</style>

<iframe class="stats" src="http://stats.example.com/page_three/stats?output=text" scrolling="no"></iframe>
```
