document.addEventListener("DOMContentLoaded", () => {
  console.log("DOM Loaded");
  const root = document.getElementById("root");

  const title = document.createElement("h1");
  const node = document.createTextNode("Welcome to Imperium");
  title.appendChild(node);
  root.appendChild(title);
});
