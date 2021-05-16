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

For ponies use this fragment:

```javascript
names = [];
Array.from(document.getElementsByClassName("listofponies")[0]
  .children[1]
  .children
).forEach(row => {
  let name = row.children[0]
    .textContent
    .toLowerCase()
    .replaceAll(" ", "-")
    .replaceAll(".", "")
    .replaceAll("ö", "o");
  if (name.includes("unnamed")) { return; }
  if (name.includes("[")) { return; }
  if (name.includes("/")) { return; }
  if (name.includes("alt")) { return; }
  if (name.includes("pony")) { return; }
  if (name.includes("mare")) { return; }
  if (name.includes("student")) { return; }
  if (name.includes("'")) { return; }
  if (name.includes('"')) { return; }
  if (name.length > 10) { return; }
  console.log([name, name.length]);
  names.push(name);
});
console.log(JSON.stringify(names));
```
