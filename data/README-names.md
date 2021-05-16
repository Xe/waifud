# How to generate names.json

Open https://xenoblade.github.io/xb2/bdat/common/BLD_NameList.html and paste
this into the browser inspector:

```js
let names = [];
Array.from(document.getElementsByClassName("sortable")[0].children[1].children)
  .forEach(row => names.push(row.children[2]
    .innerHTML
    .toLowerCase()
    .replaceAll(" ", "-")));
console.log(JSON.stringify(names));
```

Then format it with jq.
